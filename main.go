package main

import (
    "github.com/joho/godotenv"
    "github.com/rivo/tview"
    "jeremiahtowe.com/go_dash/goDash"
    "jeremiahtowe.com/go_dash/pkg/calendar"
    "jeremiahtowe.com/go_dash/pkg/github"
    "jeremiahtowe.com/go_dash/pkg/systemProperties/cpu"
    "jeremiahtowe.com/go_dash/pkg/weather"
    "log"
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
            go calendar.PopulateDisplay(view)
        }
        if value.GetModule() == "github" {
            var view = value.GetView().(*tview.Table)
            go github.PopulateDisplay(view, app)
        }
        if value.GetModule() == "cpu" {
            var view = value.GetView().(*tview.TextView)
            go cpu.PopulateDisplay(view)
        }
        if value.GetModule() == "weather" {
            var view = value.GetView().(*tview.TextView)
            go weather.PopulateDisplay(view)
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
