package cpu

import (
	"github.com/rivo/tview"
	"jeremiahtowe.com/go_dash/goDash"
	"os/exec"
)

type Widget struct {
	goDash.TextViewWidget
	Row int
	Col int
	RowSpan int
	ColSpan int
	MinGridHeight int
	MinGridWidth int
	View *tview.TextView
}

type Info struct {
	Brand string
}

func GetWidget(app *tview.Application) *Widget {
	cpuTextView := tview.NewTextView().SetDynamicColors(true)
	cpuTextView.SetBorder(true).SetTitle("🖥️  Computer Info")
	cpuTextView.SetChangedFunc(func() {
		app.Draw()
	})

	widget := Widget{
		View: cpuTextView,
		Row: 0,
		Col: 1,
		RowSpan: 1,
		ColSpan: 1,
		MinGridHeight: 0,
		MinGridWidth: 100,
	}

	return &widget
}

func GetInfo() (*Info, error) {
	var info Info

	out, err := exec.Command("sysctl","-n",  "machdep.cpu.brand_string").Output()
	if err != nil {
		return nil, err
	}

	info.Brand = string(out)
	return &info, nil
}