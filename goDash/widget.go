package goDash

type WidgetInterface interface {
	GetView() interface{}
	GetRow() int
	GetCol() int
	GetRowSpan() int
	GetColSpan() int
	GetMinGridHeight() int
	GetMinGridWidth() int
	GetModule() string
}
