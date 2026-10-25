package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
)

type DoubleTapCell struct {
	widget.Label
	CellID        widget.TableCellID
	DoubleTapFunc func()
}

func NewDoubleTapCell() *DoubleTapCell {
	c := &DoubleTapCell{}
	c.ExtendBaseWidget(c)
	return c
}

func (c *DoubleTapCell) DoubleTapped(_ *fyne.PointEvent) {
	c.DoubleTapFunc()
}
