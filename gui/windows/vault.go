package windows

import (
	"context"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/nu12/kvm/gui/components"
	"github.com/nu12/kvm/services"

	"golang.design/x/clipboard"
)

type VaultWindow struct {
	//parent    *RootWindow
	this      fyne.Window
	VaultName string
	Secrets   []*azsecrets.SecretProperties

	Components struct {
		Search *fyne.Container
		Table  *widget.Table
	}
}

func NewVaultWindow(vaultName string) *VaultWindow {
	vaultWindow := &VaultWindow{
		Secrets:   []*azsecrets.SecretProperties{},
		VaultName: vaultName,
		this:      fyne.CurrentApp().NewWindow("Vault: " + vaultName),
	}
	vaultWindow.this.Resize(fyne.Size{Width: 600, Height: 400})
	vaultWindow.this.Show()

	return vaultWindow.createComponents().setContent().run()
}

func (window *VaultWindow) createComponents() *VaultWindow {
	window.Components.Table = widget.NewTableWithHeaders(
		func() (rows int, cols int) { return len(window.Secrets), 2 },
		func() fyne.CanvasObject { return components.NewDoubleTapCell() },
		func(id widget.TableCellID, o fyne.CanvasObject) {
			b := o.(*components.DoubleTapCell)
			b.CellID = id
			b.DoubleTapFunc = func() {
				v, _ := services.GetSecret(window.VaultName, window.Secrets[id.Row].ID.Name(), "")
				if _, err := clipboard.Write(context.Background(), clipboard.FmtText, []byte(v), clipboard.Loops(1)); err != nil {
					log.Fatal(err)
				}
				log.Println("Secret's value copied to clipboard")

			}
			b.SetText(window.Secrets[id.Row].ID.Name())
		},
	)

	window.Components.Search = container.NewBorder(nil, nil, nil, widget.NewButton("Filter", func() {}), widget.NewEntry())

	return window
}

func (window *VaultWindow) setContent() *VaultWindow {
	content := container.NewBorder(
		window.Components.Search,
		nil,
		nil,
		nil,
		window.Components.Table,
	)
	window.this.SetContent(content)

	return window
}

func (window *VaultWindow) run() *VaultWindow {
	window.Secrets, _ = services.GetSecrets(window.VaultName)
	window.Components.Table.Show()

	return window
}
