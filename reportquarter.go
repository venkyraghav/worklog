package main

import (
	"fmt"
	"math"
	"net/mail"
	"os"
	"strings"
	"time"

	"github.com/akamensky/argparse"
)

type ReportQuarter struct {
	NextQuarter  bool
	CurrentTime  *time.Time
	Quarter      int
	Year         int
	LastWeek     int
	FirstWeek    int
	LastDay      time.Time
	FirstDay     time.Time
	CustomerName string
	FileName     string
	Name         string
	Email        string
	Format       string
}

func (rq *ReportQuarter) validateFormat() error {
	// TODO introduce support for md and rst formats
	if strings.Compare(rq.Format, "adoc") != 0 {
		return fmt.Errorf("unsupported format %s", rq.Format)
	}
	return nil
}

func (rq *ReportQuarter) validateEmailId() error {
	_, err := mail.ParseAddress(rq.Email)
	return err
}

func (rq *ReportQuarter) validateName() error {
	// TODO Space between Names
	if len(rq.CustomerName) == 0 {
		return fmt.Errorf("customername is invalid")
	}
	return nil
}

func (rq *ReportQuarter) constructFilename() error {
	if err := rq.validateFormat(); err != nil {
		return err
	}
	rq.FileName = fmt.Sprintf("WorkLog_%dQ%d.%s", rq.Year, rq.Quarter, rq.Format)
	return nil
}

func (rq *ReportQuarter) isDayPrintable(printDay time.Time) bool {
	return !(printDay.Weekday() == time.Saturday || printDay.Weekday() == time.Sunday) && printDay.Before(rq.LastDay) && printDay.After(rq.FirstDay)
}

func (rq *ReportQuarter) Summarize() {
	fmt.Printf("Generating report for customer %s year %d and Quarter %d in file %s. Using ...\n", rq.CustomerName, rq.Year, rq.Quarter, rq.FileName)
	fmt.Printf("  Name: %s\n", rq.Name)
	fmt.Printf("  Email: %s\n", rq.Email)
	fmt.Printf("  Between weeks: %d and %d\n", rq.FirstWeek, rq.LastWeek)
	fmt.Printf("  Between days: %s and %s\n", rq.FirstDay.Format("Mon Jan 2 2006"), rq.LastDay.Format("Mon Jan 2 2006"))
}

func (rq *ReportQuarter) GenerateWorklog() error {
	if len(rq.FileName) == 0 {
		return fmt.Errorf("filename is empty")
	}

	rq.Summarize()
	file, err := os.OpenFile(rq.FileName, os.O_WRONLY|os.O_TRUNC|os.O_CREATE, 0600)
	if err != nil {
		return err
	}
	defer file.Close()

	var f ReportFormat = FormatAdoc{rq: rq}
	file.WriteString(f.Directives())
	file.WriteString(f.HeaderMain())
	file.WriteString(f.HeaderTopic("WorkLog"))
	startFromDay := rq.LastDay
	for weekNumber := rq.LastWeek; weekNumber >= rq.FirstWeek; weekNumber-- {
		sunday := time.Date(startFromDay.Year(), time.Month(startFromDay.Month()), startFromDay.Day()-int(startFromDay.Weekday()), 10, 10, 10, 10, time.Local)
		file.WriteString(f.HeaderWeek(weekNumber, sunday))
		for day := 0; day < int(startFromDay.Weekday()); day++ {
			addTo := time.Hour * 24 * -1 * time.Duration(day)
			printDay := startFromDay.Add(addTo)
			if rq.isDayPrintable(printDay) {
				file.WriteString(f.HeaderDaily(printDay))
				file.WriteString(f.ItemDaily())
			}
		}
		addTo := time.Hour * 24 * -1 * time.Duration(int(startFromDay.Weekday())+1)
		startFromDay = startFromDay.Add(addTo)
	}
	// DEPRECATED
	// file.WriteString(f.HeaderPTO())
	// file.WriteString(f.HeaderOKR())
	// file.WriteString(f.HeaderCheckIn())

	return nil
}

func (rq *ReportQuarter) String() string {
	return fmt.Sprintf("Name: %s, Email: %s, Year: %d, NextQuarter: %v, Quarter: %d, CustomerName: %s, FirstWeek: %d, LastWeek: %d, FirstDay: %v, LastDay: %v, FileName: %s, Format: %s",
		rq.Name, rq.Email,
		rq.Year, rq.NextQuarter, rq.Quarter,
		rq.CustomerName,
		rq.FirstWeek, rq.LastWeek,
		rq.FirstDay, rq.LastDay,
		rq.FileName, rq.Format)
}

func (rq *ReportQuarter) quarterOf(month int) (int, error) {
	if month > 12 || month < 1 {
		return 0, fmt.Errorf("invalid month %d", month)
	}
	quarter := math.Ceil(float64(month) / 3)
	return int(quarter), nil
}

func (rq *ReportQuarter) nextQuarter(t time.Time) (int, int, error) {
	month := int(t.Month())
	tQuarter, err := rq.quarterOf(month)
	if err != nil {
		return 0, 0, err
	}
	if tQuarter == 4 { // Next Year 1Q
		return 1, t.Year() + 1, nil
	} else {
		return tQuarter + 1, t.Year(), nil
	}
}

func (rq *ReportQuarter) firstMonthOfQuarter(quarter int) (int, error) {
	if quarter < 1 || quarter > 4 {
		return 0, fmt.Errorf("invalid quarter %d", quarter)
	}
	return (quarter-1)*3 + 1, nil // 1st month of the quarter
}

func (rq *ReportQuarter) computeTimeToUse() (*time.Time, error) {
	var timeToUse time.Time
	// check for rq.CurrentTime
	if rq.CurrentTime == nil {
		return nil, fmt.Errorf("currentTime not set")
	}
	if rq.NextQuarter {
		quarter, year, err := rq.nextQuarter(*rq.CurrentTime)
		if err != nil {
			return nil, err
		}
		newMonth, err := rq.firstMonthOfQuarter(quarter)
		if err != nil {
			return nil, err
		}
		newDay := 7 // use 7th day of the month
		timeToUse = time.Date(year, time.Month(newMonth), newDay, 10, 10, 10, 10, time.Local)
		if int(timeToUse.Weekday()) != int(time.Saturday) {
			newDay = 1
		}
	} else {
		quarter, err := rq.quarterOf(int(rq.CurrentTime.Month()))
		if err != nil {
			return nil, err
		}
		month, err := rq.firstMonthOfQuarter(quarter)
		if err != nil {
			return nil, err
		}
		year := rq.CurrentTime.Year()
		timeToUse = time.Date(year, time.Month(month), 1, 10, 10, 10, 10, time.Local)
	}
	return &timeToUse, nil
}

func (rq *ReportQuarter) computeQuater(timeToUse time.Time) error {
	var err error
	rq.Quarter, err = rq.quarterOf(int(timeToUse.Month()))
	if err != nil {
		return err
	}

	// TODO check for firstWeek with any working days
	rq.Year, rq.FirstWeek = timeToUse.ISOWeek()

	var lastDay int
	if timeToUse.Month() == 1 || timeToUse.Month() == 10 { // Q1 and Q4 has 31 days in last month
		lastDay = 31
	} else {
		lastDay = 30
	}
	rq.LastDay = time.Date(timeToUse.Year(), timeToUse.Month()+2, lastDay, 10, 10, 10, 10, time.Local)

	// TODO check for lastWeek with any working days
	_, rq.LastWeek = rq.LastDay.ISOWeek()
	if rq.LastWeek == 1 { // Special case December some years
		rq.LastWeek = 52
	}

	rq.FirstDay = time.Date(timeToUse.Year(), timeToUse.Month(), 1, 10, 10, 10, 10, time.Local)
	return nil
}

func (rq *ReportQuarter) Compute() error {
	timeToUse, err := rq.computeTimeToUse()
	if err != nil {
		return err
	}

	if err = rq.computeQuater(*timeToUse); err != nil {
		return err
	}
	if err := rq.constructFilename(); err != nil {
		return err
	}

	return nil
}

// parses and validates commandline arguments
func (rq *ReportQuarter) Validate() error {
	parser := argparse.NewParser("v-worklog", "Generates worklog for the current quarter")

	cn := parser.String("c", "customername", &argparse.Options{Help: "Customer Name", Default: "Company"})
	nq := parser.Flag("n", "next-quarter", &argparse.Options{Help: "Generate for next quarter", Default: false})
	name := parser.String("p", "name", &argparse.Options{Help: "Name", Default: "MyName"})
	emailid := parser.String("e", "email", &argparse.Options{Help: "Email Id", Default: "myemail@company.io"})
	format := parser.String("f", "format", &argparse.Options{Help: "Report Format", Default: "adoc"})

	if err := parser.Parse(os.Args); err != nil {
		return fmt.Errorf(parser.Usage(err))
	}
	reportQuarter.NextQuarter = *nq
	reportQuarter.CustomerName = *cn
	reportQuarter.Format = *format
	reportQuarter.Name = *name
	reportQuarter.Email = *emailid
	{
		now := time.Now()
		reportQuarter.CurrentTime = &now
	}

	if err := reportQuarter.validateFormat(); err != nil {
		return fmt.Errorf(parser.Usage(err))
	}

	if err := reportQuarter.validateEmailId(); err != nil {
		return fmt.Errorf(parser.Usage(err))
	}

	if err := reportQuarter.validateName(); err != nil {
		return fmt.Errorf(parser.Usage(err))
	}

	return nil
}
