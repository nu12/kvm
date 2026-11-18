package windows

import (
	"context"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/nu12/kvm/gui/components"
	"github.com/nu12/kvm/gui/dialogs"
	"github.com/nu12/kvm/services"

	"golang.design/x/clipboard"
)

type VaultWindow struct {
	//parent    *RootWindow
	this      fyne.Window
	action    *components.Activity
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
		action:    components.NewActivity(),
	}
	vaultWindow.this.Resize(fyne.Size{Width: 1100, Height: 400})
	vaultWindow.this.Show()

	return vaultWindow
}

func (window *VaultWindow) createComponents() *VaultWindow {
	window.Components.SecretsTab.Table = components.NewTable(components.TableConfig{
		RowCount:    func() int { return len(window.Secrets) },
		Headers:     []string{"Name", "Creation date", "Expiration date", "Content-type", "Enabled"},
		ColumnWidth: []float32{250, 250, 250, 100, 100},
		NewCell:     func() fyne.CanvasObject { return components.NewSecretCell() },
		Func: func(id widget.TableCellID, o fyne.CanvasObject) {
			b := o.(*components.SecretCell)
			b.CellID = id
			b.DoubleTapFunc = func() {
				v, _ := services.GetSecret(window.VaultName, window.Secrets[id.Row].ID.Name(), "")
				if _, err := clipboard.Write(context.Background(), clipboard.FmtText, []byte(v)); err != nil {
					log.Fatal(err)
				}
				log.Println("Secret's value copied to clipboard")
			}
			b.UpdateDialog = dialogs.NewSecretUpdateDialog(
				window.this,
				window.VaultName,
				window.Secrets[id.Row].ID.Name(),
			)
			b.SetTextFromSecret(window.Secrets[id.Row])
		},
	})

	window.Components.SecretsTab.Search = container.NewBorder(nil, nil, nil, widget.NewButton("Filter", func() {}), widget.NewEntry())

	return window
}

func (window *VaultWindow) setContent() *VaultWindow {
	window.Components.SecretsTab.Content = container.NewTabItem(
		"Secrets", container.NewBorder(
			window.Components.SecretsTab.Search,
			window.action,
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
	window.action.Start("Loading secrets...")
	go func() {
		window.Secrets, _ = services.GetSecrets(window.VaultName)
		window.Components.SecretsTab.Table.Refresh()
		window.action.Stop("Done")
	}()
	return window
}
