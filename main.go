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
    "strconv"
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
    go populateGithubDisplay(githubTable, app)

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

func populateGithubDisplay(githubTable *tview.Table, app *tview.Application) {
    pullRequests, gitHubError := github.GetPullRequests()
    if gitHubError != nil {
        file, err := os.OpenFile("error.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
        if err != nil {
            log.Fatal(err)
        }
        defer file.Close()
        log.SetOutput(file)
        log.Print(gitHubError)
        githubTable.SetCell(0, 0, tview.NewTableCell("Error, check error.log"))
        app.Draw()
        return
    }

    githubTable.SetCell(0, 0, tview.NewTableCell("[aquamarine]Open Pull Requests authored by Jeremiah[white]"))
    githubTable.SetCell(0, 1, tview.NewTableCell("[aquamarine]Comments[white]"))
    githubTable.SetCell(0, 2, tview.NewTableCell("[aquamarine]Labels[white]"))
    githubTable.SetCell(0, 3, tview.NewTableCell("[aquamarine]Additions[white]"))
    githubTable.SetCell(0, 4, tview.NewTableCell("[aquamarine]Deletions[white]"))

    rowCounter := 1
    for _, pullRequest := range *pullRequests {
        githubTable.SetCell(rowCounter, 0, tview.NewTableCell(pullRequest.Title))
        githubTable.SetCell(rowCounter, 1, tview.NewTableCell(strconv.Itoa(pullRequest.NumberOfComments)).SetAlign(tview.AlignCenter))

        labels := ""
        for _, label := range pullRequest.Labels {
            labels += "[#" + label.Color +"]" + label.Name + " "
        }
        githubTable.SetCell(rowCounter, 2, tview.NewTableCell(labels))
        githubTable.SetCell(rowCounter, 3, tview.NewTableCell("[green]" + strconv.Itoa(pullRequest.Additions) + "[white]"))
        githubTable.SetCell(rowCounter, 4, tview.NewTableCell("[red]" + strconv.Itoa(pullRequest.Deletions) + "[white]"))
        rowCounter++
    }

    app.Draw()
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
