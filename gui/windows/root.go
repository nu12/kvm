package windows

import (
	"fmt"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault/v2"
	"github.com/nu12/kvm/gui/components"
	"github.com/nu12/kvm/gui/themes"
)

type RootWindow struct {
	Logo          fyne.Resource
	Vaults        []*armkeyvault.Vault
	VisibleVaults []*armkeyvault.Vault

	app        fyne.App
	this       fyne.Window
	Components struct {
		Search      *fyne.Container
		SearchEntry *widget.Entry
		Table       *widget.Table
		Counter     *widget.Label
	}
}

func (window *RootWindow) Init() *RootWindow {
	window.app = fyne.CurrentApp()
	window.app.Settings().SetTheme(&themes.TransparentTheme{})

	window.this = window.app.NewWindow("Key Vault Manager")
	window.this.Resize(fyne.Size{Width: 820, Height: 400})

	window.Vaults = []*armkeyvault.Vault{}
	window.VisibleVaults = []*armkeyvault.Vault{}
	return window
}

func (window *RootWindow) CreateComponents() *RootWindow {
	window.Components.SearchEntry = widget.NewEntry()
	window.Components.SearchEntry.OnChanged = window.filterVaults
	window.Components.SearchEntry.PlaceHolder = "Type to filter vaults..."
	window.Components.Search = container.NewVBox(window.Components.SearchEntry)

	window.Components.Table = components.NewTable(components.TableConfig{
		RowCount:    func() int { return len(window.VisibleVaults) },
		Headers:     []string{"Name", "RG", "Subscription", "Location"},
		ColumnWidth: []float32{150, 150, 300, 150},
		NewCell:     func() fyne.CanvasObject { return components.NewRootCell() },
		Func: func(id widget.TableCellID, o fyne.CanvasObject) {
			b := o.(*components.RootCell)
			b.CellID = id
			b.DoubleTapFunc = func() { NewVaultWindow(*window.VisibleVaults[id.Row].Name) }
			b.SetTextFromVault(window.VisibleVaults[id.Row])
		},
	})

	window.Components.Counter = widget.NewLabel(fmt.Sprintf("Showing %d out of %d key vaults.", len(window.Vaults), len(window.Vaults)))
	return window
}

func (window *RootWindow) filterVaults(query string) {
	query = strings.ToLower(strings.TrimSpace(query))
	window.VisibleVaults = make([]*armkeyvault.Vault, 0, len(window.Vaults))

	for _, vault := range window.Vaults {
		if query == "" || vaultMatchesQuery(vault, query) {
			window.VisibleVaults = append(window.VisibleVaults, vault)
		}
	}

	if window.Components.Table != nil {
		window.Components.Table.Refresh()
	}
	if window.Components.Counter != nil {
		window.Components.Counter.SetText(fmt.Sprintf(
			"Showing %d out of %d key vaults.",
			len(window.VisibleVaults),
			len(window.Vaults),
		))
	}
}

func vaultMatchesQuery(vault *armkeyvault.Vault, query string) bool {
	if vault == nil {
		return false
	}

	for _, field := range []*string{vault.Name, vault.ID, vault.Location} {
		if field != nil && strings.Contains(strings.ToLower(*field), query) {
			return true
		}
	}
	return false
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
	splash := NewSplash(window)
	splash.LoadAsync(window.Components.Table, window.Components.Counter)
	window.app.Run() // blocks here, runs the event loop
	return window
}
