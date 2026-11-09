package components

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault/v2"
)

type RootCell struct {
	widget.Label
	DoubleTapFunc func()
	CellID        widget.TableCellID
}

func NewRootCell() *RootCell {
	c := &RootCell{}
	c.ExtendBaseWidget(c)
	return c
}

func (c *RootCell) SetTextFromVault(vault *armkeyvault.Vault) {
	parts := strings.Split(*vault.ID, "/")
	switch c.CellID.Col {
	case 0:
		c.SetText(*vault.Name)
	case 1:
		c.SetText(parts[4])
	case 2:
		c.SetText(parts[2])
	case 3:
		c.SetText(*vault.Location)
	}
}

func (c *RootCell) DoubleTapped(_ *fyne.PointEvent) {
	if c.DoubleTapFunc != nil {
		c.DoubleTapFunc()
	}
}

func (c *RootCell) TappedSecondary(pe *fyne.PointEvent) {
	menu := fyne.NewMenu("Context",
		fyne.NewMenuItem("Open Key Vault objects", func() {
			c.DoubleTapFunc()
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Action 2", func() {
			fmt.Println("Action 2 clicked")
		}),
	)
	entryPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(c)
	// Add the mouse click position relative to the widget
	popUpPos := entryPos.Add(pe.Position)
	widget.ShowPopUpMenuAtPosition(menu, fyne.CurrentApp().Driver().CanvasForObject(c), popUpPos)
}
