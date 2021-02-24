package utils

import (
    "fmt"
    "os"
    "log"
    "net"
    "bufio"
)

/**
 Function that opens and reads the lines in inputFile.
 It return a list of Host structs or an error.
*/
func ReadInputFile (inputFile string, ports []Port) ([]Host, error) {
    var hostList []Host

    //Open file
    file, err := os.Open(inputFile)
    if err != nil {
        return nil, err
    }
    defer file.Close()

    //Read file line by line
    scanner := bufio.NewScanner(file)
    for scanner.Scan() {
        line := scanner.Text()
        ipAddress, ipNet, err := net.ParseCIDR(line)
        if err != nil {
            //Line is not in CIDR format
            //If is not a CIDR, ipNet will be nil
            ipAddress = net.ParseIP(line)
        }
        fmt.Printf("[*] %s - %s\n", ipAddress, ipNet)
        host := Host{IPAddress: ipAddress, IPNet: ipNet, Ports: ports}
        hostList = append(hostList, host)
    }
    return hostList, nil
}

/**
 Function to create the output log file. It uses the log package which provides
 the Logger struct that handles concurrency while writing to a file.
 It return a log.Logger struct or an error.
*/
func CreateOutputFile (outputFile string) (*log.Logger, error) {
    file, err := os.Create(outputFile)
    log := log.New(file, "", log.LstdFlags|log.Lshortfile)
    if err != nil {
        return nil, err
    }
    return log, nil
}
