package utils

import (
    "fmt"
    "net"
)

type Host struct {
    IPAddress net.IP
    IPNet *net.IPNet
    Ports []Port
}

func (h Host) String() string {
    return fmt.Sprintf("%s:%s:%s", h.IPAddress, h.IPNet, h.Ports)
}
