package utils

import (
    "fmt"
    "strings"
    "strconv"
)

type Port struct {
    Port int
    IsOpen bool
}

func (p Port) String() string {
    return fmt.Sprintf("%d:%t", p.Port, p.IsOpen)
}

/**
 Function that parses a comma separated list of ports and return a list of Port
 structs. If an error is triggered, the error is returned.
*/
func ParsePorts(portList string) ([]Port, error) {
    var ports []Port
    portListSlice := strings.Split(portList, ",")
    for _, port := range(portListSlice) {
        //First, convert the string port to an integer
        intPort, err := strconv.Atoi(port)
        if err != nil {
            return ports, err
        }
        //Create our Port struct and append to our list
        port := Port{Port: intPort, IsOpen: false}
        ports = append(ports, port)
    }
    return ports, nil
}
