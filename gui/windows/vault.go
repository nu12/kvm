package windows

import (
	"context"
	"fmt"
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
		SecretsTab struct {
			Content *container.TabItem
			Search  *fyne.Container
			Table   *widget.Table
		}
		CertificatesTab struct {
			Content *container.TabItem
		}
		KeysTab struct {
			Content *container.TabItem
		}
	}
}

func NewVaultWindow(vaultName string) *VaultWindow {
	return initVaultWindow(vaultName).createComponents().setContent().run()
}

func initVaultWindow(vaultName string) *VaultWindow {
	vaultWindow := &VaultWindow{
		Secrets:   []*azsecrets.SecretProperties{},
		VaultName: vaultName,
		this:      fyne.CurrentApp().NewWindow("Vault: " + vaultName),
	}
	vaultWindow.this.Resize(fyne.Size{Width: 1100, Height: 400})
	vaultWindow.this.Show()

	return vaultWindow
}

func (window *VaultWindow) createComponents() *VaultWindow {
	window.Components.SecretsTab.Table = widget.NewTableWithHeaders(
		func() (rows int, cols int) { return len(window.Secrets), 5 },
		func() fyne.CanvasObject { return components.NewDoubleTapCell() },
		func(id widget.TableCellID, o fyne.CanvasObject) {
			b := o.(*components.DoubleTapCell)
			b.CellID = id
			b.DoubleTapFunc = func() {
				v, _ := services.GetSecret(window.VaultName, window.Secrets[id.Row].ID.Name(), "")
				if _, err := clipboard.Write(context.Background(), clipboard.FmtText, []byte(v)); err != nil {
					log.Fatal(err)
				}
				log.Println("Secret's value copied to clipboard")

			}
			switch id.Col {
			case 0:
				b.SetText(window.Secrets[id.Row].ID.Name())
			case 1:
				b.SetText(window.Secrets[id.Row].Attributes.Created.String())
			case 2:
				if window.Secrets[id.Row].Attributes.Expires != nil {
					b.SetText(window.Secrets[id.Row].Attributes.Expires.String())
				} else {
					b.SetText("No expiration")
				}
			case 3:
				if window.Secrets[id.Row].ContentType != nil {
					b.SetText(*window.Secrets[id.Row].ContentType)
				} else {
					b.SetText("N/A")
				}
			case 4:
				if *window.Secrets[id.Row].Attributes.Enabled {
					b.SetText("Enabled")
				} else {
					b.SetText("Disabled")
				}

			}
		},
	)
	window.Components.SecretsTab.Table.UpdateHeader = func(id widget.TableCellID, template fyne.CanvasObject) {
		b := template.(*widget.Label)
		switch id.Col {
		case -1:
			b.SetText(fmt.Sprintf("%d", id.Row+1))
		case 0:
			b.SetText("Name")
		case 1:
			b.SetText("Creation date")
		case 2:
			b.SetText("Expiration date")
		case 3:
			b.SetText("Content-type")
		case 4:
			b.SetText("Enabled")
		}
	}
	window.Components.SecretsTab.Table.SetColumnWidth(0, 250)
	window.Components.SecretsTab.Table.SetColumnWidth(1, 250)
	window.Components.SecretsTab.Table.SetColumnWidth(2, 250)
	window.Components.SecretsTab.Table.SetColumnWidth(3, 100)
	window.Components.SecretsTab.Table.SetColumnWidth(4, 100)

	window.Components.SecretsTab.Search = container.NewBorder(nil, nil, nil, widget.NewButton("Filter", func() {}), widget.NewEntry())

	return window
}

func (window *VaultWindow) setContent() *VaultWindow {
	window.Components.SecretsTab.Content = container.NewTabItem(
		"Secrets", container.NewBorder(
			window.Components.SecretsTab.Search,
			nil,
			nil,
			nil,
			window.Components.SecretsTab.Table,
		))

	window.Components.CertificatesTab.Content = container.NewTabItem("Certificates", container.NewVBox(widget.NewLabel("Loading certificates..."), widget.NewProgressBarInfinite()))

	window.Components.KeysTab.Content = container.NewTabItem("Keys", container.NewVBox(widget.NewLabel("Loading keys..."), widget.NewProgressBarInfinite()))

	tabs := container.NewAppTabs(
		window.Components.SecretsTab.Content,
		window.Components.CertificatesTab.Content,
		window.Components.KeysTab.Content)

	tabs.OnSelected = func(ti *container.TabItem) {
		log.Printf("Selected tab %s", ti.Text)
	}
	window.this.SetContent(tabs)

	return window
}

func (window *VaultWindow) run() *VaultWindow {
	window.Secrets, _ = services.GetSecrets(window.VaultName)
	window.Components.SecretsTab.Table.Show()
	return window
}
