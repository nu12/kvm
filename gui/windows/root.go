package windows

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault/v2"
	"github.com/nu12/kvm/gui/components"
	"github.com/nu12/kvm/gui/themes"
)

type RootWindow struct {
	Logo   fyne.Resource
	Vaults []*armkeyvault.Vault

	app        fyne.App
	this       fyne.Window
	Components struct {
		Search  *fyne.Container
		Table   *widget.Table
		Counter *widget.Label
	}
}

func (window *RootWindow) Init() *RootWindow {
	window.app = app.New()
	window.app.Settings().SetTheme(&themes.TransparentTheme{})

	window.this = window.app.NewWindow("Key Vault Manager")
	window.this.Resize(fyne.Size{Width: 600, Height: 400})

	window.Vaults = []*armkeyvault.Vault{}
	return window
}

func (window *RootWindow) CreateComponents() *RootWindow {
	window.Components.Search = container.NewBorder(nil, nil, nil, widget.NewButton("Filter", func() {}), widget.NewEntry())

	window.Components.Table = widget.NewTableWithHeaders(
		func() (rows int, cols int) { return len(window.Vaults), 2 },
		func() fyne.CanvasObject { return components.NewDoubleTapCell() },
		func(id widget.TableCellID, o fyne.CanvasObject) {
			b := o.(*components.DoubleTapCell)
			b.CellID = id
			b.DoubleTapFunc = func() { NewVaultWindow(*window.Vaults[id.Row].Name) }
			switch id.Col {
			case 0:
				b.SetText(*window.Vaults[id.Row].Name)
			case 1:
				b.SetText(*window.Vaults[id.Row].ID)
			}
		},
	)
	window.Components.Table.UpdateHeader = func(id widget.TableCellID, template fyne.CanvasObject) {
		b := template.(*widget.Label)
		switch id.Col {
		case 0:
			b.SetText("Name")
		case 1:
			b.SetText("ID")
		}
	}
	window.Components.Table.ShowHeaderColumn = false
	window.Components.Table.SetColumnWidth(0, 150)
	window.Components.Table.SetColumnWidth(1, 150)

	//table.Hide()

	window.Components.Counter = widget.NewLabel(fmt.Sprintf("Showing %d out of %d key vaults.", len(window.Vaults), len(window.Vaults)))

	return window

}
func (window *RootWindow) SetContent() *RootWindow {
	content := container.NewBorder(
		window.Components.Search,
		window.Components.Counter,
		nil,
		nil,
		window.Components.Table,
	)

	window.this.SetContent(content)
	return window
}
func (window *RootWindow) Run() *RootWindow {
	splash, _ := NewSplash(window)
	splash.LoadAsync(window.Components.Table, window.Components.Counter)
	window.app.Run() // blocks here, runs the event loop
	return window
}
