package windows

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/app"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault/v2"
	"github.com/nu12/kvm/gui/themes"
)

type DoubleTapCell struct {
	widget.Label
	CellID widget.TableCellID
	Window *MainWindow
}

func newDoubleTapCell() *DoubleTapCell {
	c := &DoubleTapCell{}
	c.ExtendBaseWidget(c)
	return c
}

func (c *DoubleTapCell) DoubleTapped(_ *fyne.PointEvent) {
	switch c.CellID.Col {
	case 0:
		fmt.Printf("Double-tapped A: row=%d col=%d\n", c.CellID.Row, c.CellID.Col)
		c.Window.NewVaultWindow(c.Text)
	case 1:
		fmt.Printf("Double-tapped B: row=%d col=%d\n", c.CellID.Row, c.CellID.Col)
	default:
		fmt.Printf("Double-tapped: row=%d col=%d\n", c.CellID.Row, c.CellID.Col)
	}
	// your action here
}

type MainWindow struct {
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

func (window *MainWindow) Init() *MainWindow {
	window.app = app.New()
	window.app.Settings().SetTheme(&themes.TransparentTheme{})

	window.this = window.app.NewWindow("Key Vault Manager")
	window.this.Resize(fyne.Size{Width: 600, Height: 400})

	window.Vaults = []*armkeyvault.Vault{}
	return window
}

func (window *MainWindow) CreateComponents() *MainWindow {
	window.Components.Search = container.NewBorder(nil, nil, nil, widget.NewButton("Filter", func() {}), widget.NewEntry())

	window.Components.Table = widget.NewTableWithHeaders(
		func() (rows int, cols int) { return len(window.Vaults), 2 },
		func() fyne.CanvasObject { return newDoubleTapCell() },
		func(id widget.TableCellID, o fyne.CanvasObject) {

			b := o.(*DoubleTapCell)
			b.CellID = id
			b.Window = window
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
func (window *MainWindow) SetContent() *MainWindow {
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
func (window *MainWindow) Run() *MainWindow {
	splash, _ := window.NewSplash()
	splash.LoadAsync(window.Components.Table, window.Components.Counter)
	window.app.Run() // blocks here, runs the event loop
	return window
}
