package main

import (
    "fmt"
    "github.com/rivo/tview"
    "jeremiahtowe.com/go_dash/pkg/systemProperties/cpu"
    "jeremiahtowe.com/go_dash/pkg/weather"
    "log"
)

func main() {

    // Cpu info
    cpuInfo, err := cpu.GetInfo()
    if err != nil {
        log.Fatal(err)
    }

    // Weather info
    weatherInfo, err := weather.GetWeather()
    if err != nil {
        log.Fatal(err)
    }

    // Initialize grid
    grid := tview.NewGrid().SetRows(3, 0, 3).SetColumns(60, 30, 30).SetBorders(true)

    // Cpu info to grid
    cpuTextView := tview.NewTextView()
    fmt.Fprintf(cpuTextView, "%s", cpuInfo.Brand)
    grid.AddItem(cpuTextView, 1, 0, 1, 1, 0, 100, false)

    // Weather info to grid
    weatherTextView := tview.NewTextView()
    fmt.Fprintf(
        weatherTextView,
        "Weather in: %s\nCurrent temp: %d °F\nFeels like: %d °F\n",
        weatherInfo.Name,
        int(weatherInfo.Main.Temp),
        int(weatherInfo.Main.FeelsLike))
    grid.AddItem(weatherTextView, 1, 1, 1, 1, 0, 100, false)

    // Run tview with grid
    err = tview.NewApplication().SetRoot(grid, true).SetFocus(grid).Run()
    if err != nil {
        log.Fatal(err)
    }

}
