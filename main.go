package main

import (
    "fmt"
    "github.com/rivo/tview"
    "jeremiahtowe.com/go_dash/pkg/systemProperties/cpu"
    "jeremiahtowe.com/go_dash/pkg/weather"
    "log"
)

func main() {
    // Initialize grid
    grid := tview.NewGrid().SetRows(3, 0, 3).SetColumns(60, 30, 30).SetBorders(true)
    app := tview.NewApplication().SetRoot(grid, true).SetFocus(grid)

    cpuTextView := tview.NewTextView()
    cpuTextView.SetChangedFunc(func() {
        app.Draw()
    })
    grid.AddItem(cpuTextView, 1, 0, 1, 1, 0, 100, false)
    go populateCpuDisplay(cpuTextView)

    // Weather info to grid
    weatherTextView := tview.NewTextView()
    weatherTextView.SetChangedFunc(func() {
        app.Draw()
    })
    grid.AddItem(weatherTextView, 1, 1, 1, 1, 0, 100, false)
    go populateWeatherDisplay(weatherTextView)

    // Run application
    err := app.Run()
    if err != nil {
        log.Fatal(err)
    }
}

func populateCpuDisplay(cpuTextView *tview.TextView) {
    // Cpu info
    cpuInfo, err := cpu.GetInfo()
    if err != nil {
        log.Fatal(err)
    }

    // Cpu info to grid
    go fmt.Fprintf(cpuTextView, "%s", cpuInfo.Brand)
}

func populateWeatherDisplay(weatherTextView *tview.TextView) {
    // Weather info
    weatherInfo, err := weather.GetWeather()
    if err != nil {
        log.Fatal(err)
    }

    go fmt.Fprintf(
        weatherTextView,
        "Weather in: %s\nCurrent temp: %d °F\nFeels like: %d °F\n",
        weatherInfo.Name,
        int(weatherInfo.Main.Temp),
        int(weatherInfo.Main.FeelsLike))
}
