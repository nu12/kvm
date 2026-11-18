package tabs

import (
	"context"
	"fmt"
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

type SecretTab struct {
	vaultName string
	parent    fyne.Window
	action    *notifications.Activity
	content   *container.TabItem

	searchEntry *widget.Entry
	table       *widget.Table
	secrets     []*azsecrets.SecretProperties
	visible     []*azsecrets.SecretProperties
}

func NewSecretTab(vaultName string, parent fyne.Window, action *notifications.Activity) *SecretTab {
	tab := &SecretTab{
		vaultName: vaultName,
		parent:    parent,
		action:    action,
		secrets:   []*azsecrets.SecretProperties{},
		visible:   []*azsecrets.SecretProperties{},
	}

	tab.searchEntry = widget.NewEntry()
	tab.searchEntry.OnChanged = tab.filter
	tab.searchEntry.PlaceHolder = "Type to filter secrets..."

	table := components.NewTable(components.TableConfig{
		RowCount:    func() int { return len(tab.visible) },
		Headers:     []string{"Name", "Creation date", "Expiration date", "Content-type", "Enabled"},
		ColumnWidth: []float32{250, 250, 250, 100, 100},
		NewCell:     func() fyne.CanvasObject { return components.NewSecretCell() },
		Func: func(id widget.TableCellID, object fyne.CanvasObject) {
			cell := object.(*components.SecretCell)
			cell.CellID = id
			if id.Row < 0 || id.Row >= len(tab.visible) {
				cell.SetText("")
				return
			}

			secret := tab.visible[id.Row]
			cell.DoubleTapFunc = func() {
				secretName := tab.visible[id.Row].ID.Name()
				tab.action.Start("Fetching secret value")
				go func() {
					value, _ := services.GetSecret(tab.vaultName, secretName, "")
					_, err := clipboard.Write(context.Background(), clipboard.FmtText, []byte(value))
					if err != nil {
						tab.action.Stop("Failed to copy secret value")
					} else {
						tab.action.Stop("Secret's value copied to clipboard")
					}
				}()
			}
			cell.UpdateDialog = dialogs.NewSecretUpdateDialog(
				tab.parent,
				tab.vaultName,
				secret.ID.Name(),
				tab.action,
			)
			cell.SetTextFromSecret(secret)
		},
	})
	tab.table = table
	tab.content = container.NewTabItem("Secrets", container.NewBorder(
		container.NewVBox(tab.searchEntry),
		action,
		nil,
		nil,
		tab.table,
	))

	return tab
}

func (tab *SecretTab) Content() *container.TabItem {
	return tab.content
}

func (tab *SecretTab) Load() {
	secrets, err := services.GetSecrets(tab.vaultName)
	notifications.OnErrorCloseWindow(err, tab.parent)
	tab.secrets = secrets
	tab.filter(tab.searchEntry.Text)
	tab.action.Stop("Done")
}

func (tab *SecretTab) filter(query string) {
	query = strings.ToLower(strings.TrimSpace(query))
	tab.visible = make([]*azsecrets.SecretProperties, 0, len(tab.secrets))

	for _, secret := range tab.secrets {
		if query == "" || secretMatchesQuery(secret, query) {
			tab.visible = append(tab.visible, secret)
		}
	}

	if tab.table != nil {
		tab.table.Refresh()
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
