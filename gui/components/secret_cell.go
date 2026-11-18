package components

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/nu12/kvm/gui/dialogs"
)

type SecretCell struct {
	widget.Label
	DoubleTapFunc func()
	UpdateDialog  *dialogs.SecretUpdateDialog
	CellID        widget.TableCellID
}

func NewSecretCell() *SecretCell {
	c := &SecretCell{}
	c.ExtendBaseWidget(c)
	return c
}

func (c *SecretCell) SetTextFromSecret(secret *azsecrets.SecretProperties) {
	switch c.CellID.Col {
	case 0:
		c.SetText(secret.ID.Name())
	case 1:
		c.SetText(secret.Attributes.Created.String())
	case 2:
		if secret.Attributes.Expires != nil {
			c.SetText(secret.Attributes.Expires.String())
		} else {
			c.SetText("No expiration")
		}
	case 3:
		if secret.ContentType != nil {
			c.SetText(*secret.ContentType)
		} else {
			c.SetText("N/A")
		}
	case 4:
		if *secret.Attributes.Enabled {
			c.SetText("Enabled")
		} else {
			c.SetText("Disabled")
		}
	}
}

func (c *SecretCell) DoubleTapped(_ *fyne.PointEvent) {
	if c.DoubleTapFunc != nil {
		c.DoubleTapFunc()
	}
}

func (c *SecretCell) TappedSecondary(pe *fyne.PointEvent) {
	menu := fyne.NewMenu("Context",
		fyne.NewMenuItem("Copy value", func() {
			c.DoubleTapFunc()
		}),
		fyne.NewMenuItemSeparator(),
		fyne.NewMenuItem("Update value from input", func() {
			if c.UpdateDialog != nil {
				c.UpdateDialog.Show()
			}
		}),
	)
	entryPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(c)
	// Add the mouse click position relative to the widget
	popUpPos := entryPos.Add(pe.Position)
	widget.ShowPopUpMenuAtPosition(menu, fyne.CurrentApp().Driver().CanvasForObject(c), popUpPos)
}
