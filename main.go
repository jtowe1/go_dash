package main

import (
    "fmt"
    "github.com/joho/godotenv"
    "github.com/rivo/tview"
    "jeremiahtowe.com/go_dash/goDash"
    "jeremiahtowe.com/go_dash/pkg/calendar"
    "jeremiahtowe.com/go_dash/pkg/github"
    "jeremiahtowe.com/go_dash/pkg/systemProperties/cpu"
    "jeremiahtowe.com/go_dash/pkg/weather"
    "log"
    "os"
    "strconv"
    "strings"
    "time"
)

func main() {
    err := godotenv.Load()
    if err != nil {
        log.Fatal("Error loading .env file")
    }


    // Initialize grid
    grid := initializeGrid()
    app := initializeApp(grid)

    var widgets []goDash.WidgetInterface

    getWidgets(&widgets, app)

    for _, value := range widgets {
        viewInterface := value.GetView()
        if val, ok := viewInterface.(*tview.TextView); ok {
            grid.AddItem(
                val,
                value.GetRow(),
                value.GetCol(),
                value.GetRowSpan(),
                value.GetColSpan(),
                value.GetMinGridHeight(),
                value.GetMinGridWidth(),
                false)
        }
        if val, ok := viewInterface.(*tview.Table); ok {
            grid.AddItem(
                val,
                value.GetRow(),
                value.GetCol(),
                value.GetRowSpan(),
                value.GetColSpan(),
                value.GetMinGridHeight(),
                value.GetMinGridWidth(),
                false)
        }
        if value.GetModule() == "calendar" {
            var view = value.GetView().(*tview.TextView)
            go populateCalendarDisplay(view)
        }
        if value.GetModule() == "github" {
            var view = value.GetView().(*tview.Table)
            go populateGithubDisplay(view, app)
        }
        if value.GetModule() == "cpu" {
            var view = value.GetView().(*tview.TextView)
            go populateCpuDisplay(view)
        }
        if value.GetModule() == "weather" {
            var view = value.GetView().(*tview.TextView)
            go populateWeatherDisplay(view)
        }

    }

    // Run application
    err = app.Run()
    if err != nil {
        log.Fatal(err)
    }
}

func getWidgets(widgets *[]goDash.WidgetInterface, app *tview.Application) {

    cpuWidget := cpu.GetWidget(app)
    *widgets = append(*widgets, cpuWidget)


    weatherWidget := weather.GetWidget(app)
    *widgets = append(*widgets, weatherWidget)


    githubWidget := github.GetWidget()
    *widgets = append(*widgets, githubWidget)


    calendarWidget := calendar.GetWidget(app)
    *widgets = append(*widgets, calendarWidget)
}

func initializeApp(grid *tview.Grid) *tview.Application {
    app := tview.NewApplication().SetRoot(grid, true).SetFocus(grid)
    return app
}

func initializeGrid() *tview.Grid{
    grid := tview.NewGrid().SetRows(0, 0).SetColumns(0, 0, 0).SetBorders(false)
    return grid
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

    statusIcons := [2]string{"✖️", "✅️"}

    for _, event := range events.Items {
       date, _ := time.Parse(time.RFC3339, event.Start.DateTime)
       statusIcon := statusIcons[0]
       for _, attendee := range event.Attendees {
           if strings.ToLower(attendee.Email) == strings.ToLower(events.Summary) {
                if strings.ToLower(attendee.ResponseStatus) == "accepted" {
                    statusIcon = statusIcons[1]
                }
           }
       }

       fmt.Fprintf(
           calenderTextView,
           "%s ️%v \n\t(%v)\n",
           statusIcon, event.Summary, date.Format(time.ANSIC))
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
