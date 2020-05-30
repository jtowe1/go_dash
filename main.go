package main

import (
    "fmt"
    "github.com/joho/godotenv"
    "github.com/rivo/tview"
    "jeremiahtowe.com/go_dash/pkg/github"
    "jeremiahtowe.com/go_dash/pkg/systemProperties/cpu"
    "jeremiahtowe.com/go_dash/pkg/weather"
    "log"
    "os"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }

    // Initialize grid
    grid := tview.NewGrid().SetRows(0, 0).SetColumns(0, 0).SetBorders(true)
    app := tview.NewApplication().SetRoot(grid, true).SetFocus(grid)

    cpuTextView := tview.NewTextView().SetDynamicColors(true)
    cpuTextView.SetChangedFunc(func() {
        app.Draw()
    })
    grid.AddItem(cpuTextView, 0, 0, 1, 1, 0, 100, false)
    go populateCpuDisplay(cpuTextView)

    // Weather info to grid
    weatherTextView := tview.NewTextView().SetDynamicColors(true)
    weatherTextView.SetChangedFunc(func() {
        app.Draw()
    })
    grid.AddItem(weatherTextView, 0, 1, 1, 1, 0, 100, false)
    go populateWeatherDisplay(weatherTextView)

    githubTable := tview.NewTable().SetBorders(true)
    grid.AddItem(githubTable, 1, 0, 1, 2, 0, 100, false)
    //go populateGithubDisplay(githubTable, app)

    // Run application
    err = app.Run()
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

func populateGithubDisplay(githubTextView *tview.TextView) {
    githubInfo, gitHubError := github.GetPullRequests()
    if gitHubError != nil {
        file, err := os.OpenFile("error.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
        if err != nil {
            log.Fatal(err)
        }
        defer file.Close()
        log.SetOutput(file)
        log.Print(gitHubError)
        fmt.Fprintf(githubTextView, "Error, check error.log")
        return
    }

    textViewString := "Pull Requests authored by Jeremiah\n"
    for _, element := range githubInfo.Items {
        textViewString += element.Title + "\n"
    }

    fmt.Fprintf(githubTextView, "%s", textViewString)
}

func populateWeatherDisplay(weatherTextView *tview.TextView) {
    // Weather info
    weatherInfo, getWeatherError := weather.GetWeather()
    if getWeatherError != nil {
        file, err := os.OpenFile("error.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
        if err != nil {
            log.Fatal(err)
        }
        defer file.Close()
        log.SetOutput(file)
        log.Print(getWeatherError)
        fmt.Fprintf(weatherTextView, "Error, check error.log")
        return
    }

    go fmt.Fprintf(
        weatherTextView,
        "Weather in: %s\nCurrent temp: [red]%d °F[white]\nFeels like: [red]%d °F[white]\n",
        weatherInfo.Name,
        int(weatherInfo.Main.Temp),
        int(weatherInfo.Main.FeelsLike))
}
