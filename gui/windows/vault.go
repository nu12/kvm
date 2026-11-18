package windows

import (
	"context"
	"fmt"
	"log"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/nu12/kvm/gui/components"
	"github.com/nu12/kvm/gui/dialogs"
	"github.com/nu12/kvm/gui/notifications"
	"github.com/nu12/kvm/services"

	"golang.design/x/clipboard"
)

type VaultWindow struct {
	//parent    *RootWindow
	this           fyne.Window
	action         *notifications.Activity
	VaultName      string
	Secrets        []*azsecrets.SecretProperties
	VisibleSecrets []*azsecrets.SecretProperties

	Components struct {
		SecretsTab struct {
			Content     *container.TabItem
			Search      *fyne.Container
			SearchEntry *widget.Entry
			Table       *widget.Table
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
		Secrets:        []*azsecrets.SecretProperties{},
		VisibleSecrets: []*azsecrets.SecretProperties{},
		VaultName:      vaultName,
		this:           fyne.CurrentApp().NewWindow("Vault: " + vaultName),
		action:         notifications.NewActivity(),
	}
	vaultWindow.this.Resize(fyne.Size{Width: 1100, Height: 400})
	vaultWindow.this.Show()

	return vaultWindow
}

func (window *VaultWindow) createComponents() *VaultWindow {
	window.Components.SecretsTab.SearchEntry = widget.NewEntry()
	window.Components.SecretsTab.SearchEntry.OnChanged = window.filterSecrets
	window.Components.SecretsTab.SearchEntry.PlaceHolder = "Type to filter secrets..."
	window.Components.SecretsTab.Search = container.NewVBox(window.Components.SecretsTab.SearchEntry)

	window.Components.SecretsTab.Table = components.NewTable(components.TableConfig{
		RowCount:    func() int { return len(window.VisibleSecrets) },
		Headers:     []string{"Name", "Creation date", "Expiration date", "Content-type", "Enabled"},
		ColumnWidth: []float32{250, 250, 250, 100, 100},
		NewCell:     func() fyne.CanvasObject { return components.NewSecretCell() },
		Func: func(id widget.TableCellID, o fyne.CanvasObject) {
			b := o.(*components.SecretCell)
			b.CellID = id
			b.DoubleTapFunc = func() {
				secretName := window.VisibleSecrets[id.Row].ID.Name()
				window.action.Start("Fetching secret value")
				go func() {
					value, _ := services.GetSecret(window.VaultName, secretName, "")
					_, err := clipboard.Write(context.Background(), clipboard.FmtText, []byte(value))
					if err != nil {
						window.action.Stop("Failed to copy secret value")
					} else {
						window.action.Stop("Secret's value copied to clipboard")
					}
				}()
			}
			b.UpdateDialog = dialogs.NewSecretUpdateDialog(
				window.this,
				window.VaultName,
				window.VisibleSecrets[id.Row].ID.Name(),
				window.action,
			)
			b.SetTextFromSecret(window.VisibleSecrets[id.Row])
		},
	})

	return window
}

func (window *VaultWindow) filterSecrets(query string) {
	query = strings.ToLower(strings.TrimSpace(query))
	window.VisibleSecrets = make([]*azsecrets.SecretProperties, 0, len(window.Secrets))

	for _, secret := range window.Secrets {
		if query == "" || secretMatchesQuery(secret, query) {
			window.VisibleSecrets = append(window.VisibleSecrets, secret)
		}
	}

	if window.Components.SecretsTab.Table != nil {
		window.Components.SecretsTab.Table.Refresh()
	}
}

func secretMatchesQuery(secret *azsecrets.SecretProperties, query string) bool {
	if secret == nil {
		return false
	}

	var fields []string
	if secret.ID != nil {
		fields = append(fields, string(*secret.ID), secret.ID.Name())
	}
	if secret.ContentType != nil {
		fields = append(fields, *secret.ContentType)
	}
	if secret.Attributes != nil {
		if secret.Attributes.Created != nil {
			fields = append(fields, secret.Attributes.Created.String())
		}
		if secret.Attributes.Expires != nil {
			fields = append(fields, secret.Attributes.Expires.String())
		}
		if secret.Attributes.Enabled != nil {
			fields = append(fields, fmt.Sprint(*secret.Attributes.Enabled))
		}
	}

	for _, field := range fields {
		if strings.Contains(strings.ToLower(field), query) {
			return true
		}
	}
	return false
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
	fyne.Do(func() {
		secrets, err := services.GetSecrets(window.VaultName)
		notifications.OnErrorCloseWindow(err, window.this)
		window.Secrets = secrets
		query := ""
		if window.Components.SecretsTab.SearchEntry != nil {
			query = window.Components.SecretsTab.SearchEntry.Text
		}
		window.filterSecrets(query)
		window.action.Stop("Done")
	})
	return window
}
