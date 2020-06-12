package cpu

import (
	"github.com/rivo/tview"
	"os/exec"
)

type Widget struct {
	Row int
	Col int
	RowSpan int
	ColSpan int
	MinGridHeight int
	MinGridWidth int
	View *tview.TextView
	Module string
}

func (w *Widget) GetView() interface{} {
	return w.View
}

func (w *Widget) GetRow() int {
	return w.Row
}

func (w *Widget) GetCol() int {
	return w.Col
}

func (w *Widget) GetRowSpan() int {
	return w.RowSpan
}

func (w *Widget) GetColSpan() int {
	return w.ColSpan
}

func (w *Widget) GetMinGridHeight() int {
	return w.MinGridHeight
}

func (w *Widget) GetMinGridWidth() int {
	return w.MinGridWidth
}

func (w *Widget) GetModule() string {
	return w.Module
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
		Module: "cpu",
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