package main

import (
	"flag"
	"fmt"
	"net"
	"sync"
	"time"
)

func main() {
    host := flag.String("host", "", "Host to scan")
    timeout := flag.Int("timeout", 1000, "Timeout in millisseconds")
    maxPort := flag.Int("max", 1024, "The maximum port you'd like to scan (can't be higher than 65535)")
    minPort := flag.Int("min", 0, "The minimum port you'd like to scan (can't be bellow 0)")

    flag.Parse()

    if len(*host) == 0 || *minPort < 0 || *maxPort > 65535 {
        flag.Usage()
        return
    }


    fmt.Printf("[+] Scanning %s...\n", *host)

    start := time.Now()

    var wg sync.WaitGroup
    var mutex sync.Mutex
    var openPorts []int

    for i := *minPort ;i <= *maxPort; i++ {
        wg.Add(1)
        go ScanPort(i, *host, *timeout, &wg, &mutex, &openPorts)
    }
    wg.Wait()

    elapsed := time.Since(start)

    fmt.Println("[+] Open ports:")
    for _, p := range openPorts {
        fmt.Printf("\t- %d\n", p)
    }

    fmt.Printf("[+] Scan done in %s sec", elapsed)

}

func ScanPort(port int, addr string, timeout int, wg *sync.WaitGroup, mutex *sync.Mutex, openPorts* []int) {
    defer wg.Done()
    conn, err := net.DialTimeout("tcp", fmt.Sprintf("%s:%d", addr, port), time.Duration(timeout) * time.Millisecond)
    if err != nil {
        return
    }
    defer conn.Close()

    *openPorts = append(*openPorts, port)

    // fmt.Printf("[+] %d Port is open\n", port)
}
