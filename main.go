package main

import (
    "fmt"
    "jeremiahtowe.com/go_dash/pkg/systemProperties/cpu"
    "log"
)

func main() {

    cpuInfo, err := cpu.GetInfo()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%s\n", cpuInfo)
}
