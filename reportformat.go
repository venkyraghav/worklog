package main

import (
	"fmt"
	"time"
)

type ReportFormat interface {
	Directives() string
	HeaderMain() string
	HeaderTopic(topic string) string
	HeaderWeek(weekNumber int, startDay time.Time) string
	HeaderDaily(day time.Time) string
	ItemDaily() string
	HeaderPTO() string
	HeaderOKR() string
	HeaderCheckIn() string
}

type FormatAdoc struct {
	rq *ReportQuarter
}

func (f FormatAdoc) HeaderPTO() string {
	return `== PTO
.This Quarter
. [todo]

.Next Quarter
. [todo]

`
}

func (f FormatAdoc) HeaderOKR() string {
	return `== OKR
. https://docs.google.com/spreadsheets/d/1Uqe_SAT3x3Sby8pGoZ4Id9dkOUGkKGrm_JcSK8KGRC4/edit#gid=0

. Obective 1: Build a Strong Team
* [ ] KR 1: Assist in hiring new SAs, CEs, and Vendors
* [ ] KR 2: Be a mentor or buddy for new hires
* [ ] KR 3: Conduct bootcamps (Security & Advanced) for new hires -- SAs, CEs, and ACEs
* [ ] KR 4: Collaborate with Sven and Todd on Vendor bootcamp

. Obective 2: Earning Customer Love & Customer Success
* [ ] KR 1: Collaborate with Nikoleta on updating the onboarding experience for new team members

. Objective 3: Personal Development
* [ ] KR 1: Complete Azure Associate Data Engineer certification
* [ ] KR 2: Get beyond CFK deployment (MRC & CL)

`
}

func (f FormatAdoc) HeaderCheckIn() string {
	str := `CheckIn
.Impacts
. [todo]

.Challenges
. [todo]

.Priorities
. [todo]

`
	return fmt.Sprintf("== Q%d %s", f.rq.Quarter, str)
}

func (f FormatAdoc) Directives() string {
	return `// Directives
:toc:
:sectnums:
:sectnumlevels: 2
:hardbreaks:
`
}

func (f FormatAdoc) HeaderMain() string {
	return fmt.Sprintf("= %d Q%d WorkLog\n%s <%s>\n\n", f.rq.Year, f.rq.Quarter, f.rq.Name, f.rq.Email)
}

func (f FormatAdoc) HeaderTopic(topic string) string {
	return fmt.Sprintf("== %s\n", topic)
}

func (f FormatAdoc) HeaderWeek(weekNumber int, startDay time.Time) string {
	return fmt.Sprintf("=== Week %d (TODO) %s\n", weekNumber, startDay.Format("Mon Jan 2 2006"))
}

func (f FormatAdoc) HeaderDaily(day time.Time) string {
	return fmt.Sprintf("==== %s: %s\n", f.rq.CustomerName, day.Format("Mon Jan _2 2006"))
}

func (f FormatAdoc) ItemDaily() string {
	return `.Schedule update
None
	
.Issues/Blockers
None
	
.Progress today
. [todo]

.Plans for tomorrow
. [todo]

.References/Links
None

`
}
