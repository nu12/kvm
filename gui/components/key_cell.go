package components

import (
	"fmt"
	"sort"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
)

type KeyCell struct {
	widget.Label
	DoubleTapFunc func()
	CellID        widget.TableCellID
}

func NewKeyCell() *KeyCell {
	cell := &KeyCell{}
	cell.ExtendBaseWidget(cell)
	return cell
}

func (cell *KeyCell) SetTextFromKey(key *azkeys.KeyProperties) {
	cell.SetText("")
	if key == nil {
		return
	}

	switch cell.CellID.Col {
	case 0:
		if key.KID != nil {
			cell.SetText(key.KID.Name())
		}
	case 1:
		if key.Attributes != nil && key.Attributes.Created != nil {
			cell.SetText(key.Attributes.Created.String())
		} else {
			cell.SetText("N/A")
		}
	case 2:
		if key.Attributes == nil || key.Attributes.Expires == nil {
			cell.SetText("No expiration")
		} else {
			cell.SetText(key.Attributes.Expires.String())
		}
	case 3:
		if key.Attributes == nil || key.Attributes.Enabled == nil {
			cell.SetText("N/A")
		} else if *key.Attributes.Enabled {
			cell.SetText("Enabled")
		} else {
			cell.SetText("Disabled")
		}
	case 4:
		if key.Managed == nil {
			cell.SetText("N/A")
		} else {
			cell.SetText(fmt.Sprint(*key.Managed))
		}
	case 5:
		var tags []string
		for name, value := range key.Tags {
			if value == nil {
				tags = append(tags, name)
			} else {
				tags = append(tags, name+"="+*value)
			}
		}
		sort.Strings(tags)
		cell.SetText(strings.Join(tags, ", "))
	}
}

func (cell *KeyCell) DoubleTapped(_ *fyne.PointEvent) {
	if cell.DoubleTapFunc != nil {
		cell.DoubleTapFunc()
	}
}

func (cell *KeyCell) TappedSecondary(event *fyne.PointEvent) {
	menu := fyne.NewMenu("Context",
		fyne.NewMenuItem("Download public key", func() {
			if cell.DoubleTapFunc != nil {
				cell.DoubleTapFunc()
			}
		}),
	)
	entryPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(cell)
	popupPos := entryPos.Add(event.Position)
	widget.ShowPopUpMenuAtPosition(menu, fyne.CurrentApp().Driver().CanvasForObject(cell), popupPos)
}
