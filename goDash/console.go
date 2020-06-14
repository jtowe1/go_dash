package goDash

import (
	"github.com/rivo/tview"
	"jeremiahtowe.com/go_dash/pkg/calendar"
	"jeremiahtowe.com/go_dash/pkg/github"
	"jeremiahtowe.com/go_dash/pkg/systemProperties/cpu"
	"jeremiahtowe.com/go_dash/pkg/weather"
	"log"
)

type goDash struct {
	widgets *[]WidgetInterface
	app *tview.Application
	grid *tview.Grid
}

func Run() {
	goDash := setup()

	// Initialize grid
	goDash.initializeGrid()
	goDash.initializeApp()

	goDash.getWidgets()
	goDash.populateGrid()

	// Run application
	err := goDash.app.Run()
	if err != nil {
		log.Fatal(err)
	}
}

func setup() *goDash {
	var goDash goDash
	goDash.widgets = new([]WidgetInterface)
	return &goDash
}

func (gd *goDash) getWidgets() {
	cpuWidget := cpu.GetWidget(gd.app)
	*gd.widgets = append(*gd.widgets, cpuWidget)

	weatherWidget := weather.GetWidget(gd.app)
	*gd.widgets = append(*gd.widgets, weatherWidget)

	githubWidget := github.GetWidget()
	*gd.widgets = append(*gd.widgets, githubWidget)

	calendarWidget := calendar.GetWidget(gd.app)
	*gd.widgets = append(*gd.widgets, calendarWidget)
}

func (gd *goDash) initializeApp() {
	gd.app = tview.NewApplication().SetRoot(gd.grid, true).SetFocus(gd.grid)
}

func (gd *goDash) initializeGrid() {
	gd.grid = tview.NewGrid().SetRows(0, 0).SetColumns(0, 0, 0).SetBorders(false)
}

func (gd *goDash) populateGrid() {
	for _, value := range *gd.widgets {
		viewInterface := value.GetView()
		if val, ok := viewInterface.(*tview.TextView); ok {
			gd.grid.AddItem(
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
			gd.grid.AddItem(
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
			go github.PopulateDisplay(view, gd.app)
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
}
