package main

import (
    "fmt"
    "github.com/joho/godotenv"
    "github.com/rivo/tview"
    "jeremiahtowe.com/go_dash/pkg/calendar"
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
    grid := tview.NewGrid().SetRows(0, 0).SetColumns(0, 0, 0).SetBorders(false)
    app := tview.NewApplication().SetRoot(grid, true).SetFocus(grid)

    cpuTextView := tview.NewTextView().SetDynamicColors(true)
    cpuTextView.SetBorder(true).SetTitle("🖥️  Computer Info")
    cpuTextView.SetChangedFunc(func() {
        app.Draw()
    })
    grid.AddItem(cpuTextView, 0, 1, 1, 1, 0, 100, false)
    go populateCpuDisplay(cpuTextView)

    // Weather info to grid
    weatherTextView := tview.NewTextView().SetDynamicColors(true)
    weatherTextView.SetBorder(true).SetTitle("☁️  Weather")
    weatherTextView.SetChangedFunc(func() {
        app.Draw()
    })
    grid.AddItem(weatherTextView, 0, 2, 1, 1, 0, 100, false)
    go populateWeatherDisplay(weatherTextView)

    githubTable := tview.NewTable()
    githubTable.SetBorders(true)
    githubTable.SetSeparator(tview.BoxDrawingsLightVertical)
    grid.AddItem(githubTable, 1, 1, 1, 2, 0, 100, false)
    go populateGithubDisplay(githubTable, app)

    calendarTextView := tview.NewTextView()
    calendarTextView.SetBorder(true).SetTitle("📅  Calendar")
    calendarTextView.SetChangedFunc(func() {
        app.Draw()
    })
    grid.AddItem(calendarTextView, 0, 0, 2, 1, 0, 100, false)
    go populateCalendarDisplay(calendarTextView)


    // Run application
    err = app.Run()
    if err != nil {
        log.Fatal(err)
    }
}

func populateCalendarDisplay(calenderTextView *tview.TextView) {
    events, err := calendar.GetCalendar()
    if err != nil {
        file, err := os.OpenFile("error.log", os.O_CREATE|os.O_APPEND|os.O_WRONLY, 0644)
        if err != nil {
            log.Fatal(err)
        }
        defer file.Close()
        log.SetOutput(file)
        log.Print(err)
        fmt.Fprintf(calenderTextView, "%s","Error, check error.log")
        return
    }

    for _, item := range events.Items {
        date := item.Start.DateTime
        if date == "" {
            date = item.Start.Date
        }
        fmt.Fprintf(calenderTextView, "%v (%v)\n", item.Summary, date)
    }

}

func populateCpuDisplay(cpuTextView *tview.TextView) {
    // Cpu info
    cpuInfo, err := cpu.GetInfo()
    if err != nil {
        log.Fatal(err)
    }

    // Cpu info to grid
    fmt.Fprintf(cpuTextView, "%s", cpuInfo.Brand)
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

    githubTable.SetCell(0, 0, tview.NewTableCell("️[aquamarine]Open Pull Requests authored by Jeremiah[white]"))
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
        githubTable.SetCell(rowCounter, 3, tview.NewTableCell("[green]" + strconv.Itoa(pullRequest.Additions) + "[white]").SetAlign(tview.AlignCenter))
        githubTable.SetCell(rowCounter, 4, tview.NewTableCell("[red]" + strconv.Itoa(pullRequest.Deletions) + "[white]").SetAlign(tview.AlignCenter))
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
