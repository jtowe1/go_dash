package main

import (
    "fmt"
    "jeremiahtowe.com/go_dash/pkg/systemProperties/cpu"
    "jeremiahtowe.com/go_dash/pkg/weather"
    "log"
)

func main() {
    cpuInfo, err := cpu.GetInfo()
    if err != nil {
        log.Fatal(err)
    }
    fmt.Printf("%s\n", (*cpuInfo).Brand)

    weatherInfo, err := weather.GetWeather()
    if err != nil {
        log.Fatal(err)
    }

    fmt.Printf("%d °F\n", int((*weatherInfo).Main.Temp))
}
