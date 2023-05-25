package main

import (
	"log"
)

var (
	reportQuarter ReportQuarter
)

func main() {
	if err := reportQuarter.Validate(); err != nil {
		log.Fatalln(err)
		return
	}
	if err := reportQuarter.Compute(); err != nil {
		log.Fatalln(err)
		return
	}
	if err := reportQuarter.GenerateWorklog(); err != nil {
		log.Fatalln(err)
		return
	}
}
