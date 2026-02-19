package main

import(
	"flag"
)

var flagReportAddr string
var flagPollInterval int
var flagReportInterval int

func parseFlags() {

	flag.StringVar(&flagReportAddr, "a", "localhost:8080", "address and port to run server")
	flag.IntVar(&flagPollInterval, "p", 2, "poll interval")
	flag.IntVar(&flagReportInterval, "r", 10, "report interval")
	
	flag.Parse()
}
