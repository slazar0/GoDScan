package main

import (
    "log"
    "flag"
	"errors"
	"github.com/godscan/core"
)

func parseArgs() (string, string, string, int, error){
	var inputFile, outputFile, portList string
	//var itemId int64
    var err error
    threads := flag.Int("threads", 100, "Number of threads (go rouitines) to use.")
	flag.StringVar(&inputFile, "i", "", "Input file with IP addresses/CIDR network range to scan.")
	flag.StringVar(&outputFile, "o", "", "Output file name to store the scan results.")
    flag.StringVar(&portList, "p", "", "Comma-separated list of ports to scan.")
	flag.Parse()
	if inputFile == "" || outputFile  == "" || portList == "" {
		flag.PrintDefaults()
		err = errors.New("[!!] Wrong parameters.")
	}

	return inputFile, outputFile, portList, *threads, err
}

func main() {
	//Parse the arguments
    inputFile, outputFile, portList, threads, err := parseArgs()
	if err != nil {
		log.Fatal(err)
	}

    //Prepare the output file
    outputFileHandler, err := utils.CreateOutputFile(outputFile)
    if err != nil {
        log.Fatal(err)
    }

    //Prepare a list of ports
    ports, err := utils.ParsePorts(portList)
    if err !=nil {
        log.Fatal(err)
    }

    //Read the hosts to scan
    hostList, err := utils.ReadInputFile(inputFile, ports)
    if err != nil {
        log.Fatal(err)
    }
    //Scan the targets
    utils.ScanListOfHosts(hostList, outputFileHandler, threads)
}
