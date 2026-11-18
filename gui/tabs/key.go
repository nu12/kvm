package tabs

import (
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
)

type KeyTab struct {
	content *container.TabItem
}

func NewKeyTab() *KeyTab {
	return &KeyTab{
		content: container.NewTabItem("Keys", container.NewVBox(
			widget.NewLabel("Loading keys..."),
			widget.NewProgressBarInfinite(),
		)),
	}
}

func (tab *KeyTab) Content() *container.TabItem {
	return tab.content
}
