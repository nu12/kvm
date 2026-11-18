package components

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type TableConfig struct {
	RowCount    func() int
	Headers     []string
	ColumnWidth []float32
	NewCell     func() fyne.CanvasObject
	Func        func(id widget.TableCellID, o fyne.CanvasObject)
}

func NewTable(config TableConfig) *widget.Table {
	table := widget.NewTableWithHeaders(
		func() (rows int, cols int) { return config.RowCount(), len(config.Headers) },
		func() fyne.CanvasObject { return config.NewCell() },
		config.Func,
	)
	for index, width := range config.ColumnWidth {
		table.SetColumnWidth(index, width)
	}
	//table.ShowHeaderColumn = false
	table.UpdateHeader = func(id widget.TableCellID, template fyne.CanvasObject) {
		fyne.Do(func() {
			b := template.(*widget.Label)
			if id.Col == -1 {
				b.SetText(fmt.Sprintf("%d", id.Row+1))
			} else {
				b.SetText(config.Headers[id.Col])
			}
		})
	}

	table.ExtendBaseWidget(table)
	return table
}
