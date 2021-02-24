package utils

import (
    "fmt"
    "log"
    "sync"
    "net"
    "time"
)

/**
 Function that returns true if an element read from the input file is a network
 range. This condition is validated by checking the Host IPNet parameter.
 Otherwise, it returns false.
*/
func isNetwork(host Host) bool {
    //Check if we are targeting a CIDR network
    if host.IPNet == nil {
        return false
    }
    return true
}

/**
 Function that connects to a port on a certain IP address. In case of 
 establishing a connection, it returns true. Otherwise, it returns false.
 The timeout time it is set to 10 seconds.
*/
func isPortOpen(ipAddress net.IP, port Port) bool {
    destination := fmt.Sprintf("%s:%d", ipAddress.String(), port.Port)
    //TCP connection. 10 seconds timeout.
    conn, err := net.DialTimeout("tcp", destination, 10 * time.Second)
    if err != nil{
        return false
    }
    defer conn.Close()
    return true
}

/**
 Function that performs the port scans. It sends the string to a channel used by
 the program to handle the routines output. 
*/
func scanHost(channel chan string, wg *sync.WaitGroup, host Host) {
    defer wg.Done()
    portStatusStr := "Close"
    for _, port := range(host.Ports) {
        status := isPortOpen(host.IPAddress, port)
        if status {
            portStatusStr := "Open"
            fmt.Printf("[*] %s:%d:%s\n", host.IPAddress.String(), port.Port, portStatusStr)
        }
        //We send the next line to write into the file.
        channel <- fmt.Sprintf("%s:%d:%s", host.IPAddress.String(), port.Port, portStatusStr)
    }
}

/**
 Function used to iterate over a network range.
 - https://gist.github.com/kotakanbe/d3059af990252ba89a82
 - http://play.golang.org/p/m8TNTtygK0
*/
func inc(ip net.IP) {
	for j := len(ip) - 1; j >= 0; j-- {
		ip[j]++
		if ip[j] > 0 {
			break
		}
	}
}

/**
 Function that returns the IP addresses within the same network range
 established in the Hosts' IPNet parameter.
 Return a list of strings with the IP addresses that are part of the network
 range.
*/
func hostsFromCIDR(host Host) []string{
    var hosts []string
    for ip := host.IPAddress.Mask(host.IPNet.Mask); host.IPNet.Contains(ip); inc(ip) {
        hosts = append(hosts, ip.String())
    }
    //We return the list of IP address excluding the network address and the
    //broadcast address
    return hosts[1 : len(hosts) - 1]
}

/**
 If a file line has been identified as a network range, we iterate over the IP
 addresses that belong to that range. This function depends on the status of the
 main WaitGroup. It will launch a go routine if the WaitGroup is not blocked.
 Otherwise, it will wait until one of the routines finishes. It appends 1 more
 routine to the WaitGroup.
 The function returns the number of runningThreads to continue the execution
 with the next file lines.
*/
func handleNetwork(channel chan string, wg *sync.WaitGroup, threads int, runningThreads int, host Host) int {
    ipAddresses := hostsFromCIDR(host)
    for _, ip := range(ipAddresses) {
        ip := net.ParseIP(ip)
        h := Host{IPAddress: ip, IPNet: host.IPNet, Ports: host.Ports}
        //Validate if we should wait for the current routines to finish
        runningThreads = checkCurrentRoutines(runningThreads, threads, wg)
        //Add a new routine to the created WaitGroup
        wg.Add(1)
        //Start a new routine to scan the current host
        go scanHost(channel, wg, h)
    }
    return runningThreads
}

/**
 Function that manages the output produced by the routines asynchronously. It
 waits the routines to send data that will be appended to the outputFile
 provided in the programs parameters. outputFile is a log.Logger struct that
 ensures atomic write operations.
 There is a special condition that it is used to notice if the execution of the
 program has finished. The channel will receive a timeout after 30 seconds if
 routines did not write into the channel which will mean there are not new 
 scans and thus, the execution can finish.
*/
func hostReader(channel chan string, outputFile *log.Logger) {
    finish := false
    for !finish {
        select {
        case scanOutput := <-channel:
            outputFile.Printf("%s\n", scanOutput)
        //Set a timeout on the channel
        //After 30 seconds without receiving data, the execution will finish
        case <- time.After(30 * time.Second):
            finish = true
            fmt.Println("[!] 30 seconds time out without receiving new messages. Exiting...")
            close(channel)
        }
    }
}

/**
 If there are the maximun number of routines running, we wait for the routines 
 execution to finish. Otherwise, we increment runningThreads and continue.
 The function returns the number of running threads (routines) after being
 incremented or restarted.
*/
func checkCurrentRoutines(runningThreads int, threads int, wg *sync.WaitGroup) int {
    if runningThreads == threads {
        wg.Wait()
        //Change the value to -1 so it is later incremented to return 0
        runningThreads = -1
    }
    runningThreads += 1
    //Simply increment the current number of threads running
    return runningThreads
}

/**
 Main function on this file that is called from the main program.
 The function starts the hostReader in a go routine that reads the input
 received from the routines that will scan the hosts provided in the input file.
 A go routine will be launched for each host that will be scanned, independently
 if we are dealing with a host within a network range or a single IP address.
 After looping over all the hosts retrieved from the input file provided to the
 program, it waits for all the routines to finish.
*/
func ScanListOfHosts (hostList []Host, outputFile *log.Logger, threads int) {
    //Variables
    var wg sync.WaitGroup
    runningThreads := 0
    channel := make(chan string, threads)
    fmt.Printf("[*] Threads: %d\n", threads)

    //Start our reader that will receive data in the channel
    go hostReader(channel, outputFile)

    //Handle concurrency here
    for _, host := range(hostList) {
        if isNetwork(host) { //We have a network range
            //Now we can start a new go routine
            runningThreads = handleNetwork(channel, &wg, threads, runningThreads, host)
        } else { //It is a single IP address
            runningThreads = checkCurrentRoutines(runningThreads, threads, &wg)
            wg.Add(1)
            go scanHost(channel, &wg, host)
        }
    }
    wg.Wait()
}
