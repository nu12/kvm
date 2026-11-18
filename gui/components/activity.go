package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type Activity struct {
	widget.BaseWidget
	activity *widget.Activity
	label    *widget.Label
}

func NewActivity() *Activity {
	this := &Activity{
		activity: widget.NewActivity(),
		label:    widget.NewLabel("Loading..."),
	}
	this.activity.Start()
	this.activity.Show()
	this.ExtendBaseWidget(this)
	return this
}

func (a *Activity) Start(text string) {
	a.activity.Start()
	a.activity.Show()
	a.label.SetText(text)
}

func (a *Activity) Stop(text string) {
	a.activity.Stop()
	a.activity.Hide()
	a.label.SetText(text)
}

func (a *Activity) CreateRenderer() fyne.WidgetRenderer {
	return widget.NewSimpleRenderer(container.NewHBox(a.activity, a.label))
}
