package tabs

import (
	"bytes"
	"fmt"
	"io"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/nu12/kvm/gui/components"
	"github.com/nu12/kvm/gui/notifications"
	"github.com/nu12/kvm/services"
)

type KeyTab struct {
	vaultName string
	parent    fyne.Window
	action    *notifications.Activity
	content   *container.TabItem

	searchEntry *widget.Entry
	table       *widget.Table
	keys        []*azkeys.KeyProperties
	visible     []*azkeys.KeyProperties
	loading     bool
	loaded      bool
}

func NewKeyTab(vaultName string, parent fyne.Window, action *notifications.Activity) *KeyTab {
	tab := &KeyTab{
		vaultName: vaultName,
		parent:    parent,
		action:    action,
		keys:      []*azkeys.KeyProperties{},
		visible:   []*azkeys.KeyProperties{},
	}

	tab.searchEntry = widget.NewEntry()
	tab.searchEntry.OnChanged = tab.filter
	tab.searchEntry.PlaceHolder = "Type to filter keys..."
	tab.table = components.NewTable(components.TableConfig{
		RowCount:    func() int { return len(tab.visible) },
		Headers:     []string{"Name", "Creation date", "Expiration date", "Enabled", "Managed", "Tags"},
		ColumnWidth: []float32{250, 250, 250, 100, 100, 250},
		NewCell:     func() fyne.CanvasObject { return components.NewKeyCell() },
		Func: func(id widget.TableCellID, object fyne.CanvasObject) {
			cell := object.(*components.KeyCell)
			cell.CellID = id
			cell.DoubleTapFunc = nil
			if id.Row < 0 || id.Row >= len(tab.visible) {
				cell.SetText("")
				return
			}
			key := tab.visible[id.Row]
			cell.DoubleTapFunc = func() {
				if key != nil && key.KID != nil {
					tab.downloadPublicKey(key.KID.Name())
				}
			}
			cell.SetTextFromKey(key)
		},
	})
	tab.content = container.NewTabItem("Keys", container.NewBorder(
		container.NewVBox(tab.searchEntry),
		action,
		nil,
		nil,
		tab.table,
	))

	return tab
}

func (tab *KeyTab) Content() *container.TabItem {
	return tab.content
}

func (tab *KeyTab) Load() {
	if tab.loaded || tab.loading {
		return
	}
	tab.loading = true
	tab.action.Start("Loading keys")

	go func() {
		keys, err := services.GetKeys(tab.vaultName)
		fyne.Do(func() {
			tab.loading = false
			if err != nil {
				tab.action.Stop("Failed to load keys")
				notifications.OnErrorContinue(err, tab.parent)
				return
			}

			tab.keys = keys
			tab.filter(tab.searchEntry.Text)
			tab.loaded = true
			tab.action.Stop("Done")
		})
	}()
}

func (tab *KeyTab) filter(query string) {
	query = strings.ToLower(strings.TrimSpace(query))
	tab.visible = make([]*azkeys.KeyProperties, 0, len(tab.keys))

	for _, key := range tab.keys {
		if query == "" || keyMatchesQuery(key, query) {
			tab.visible = append(tab.visible, key)
		}
	}

	if tab.table != nil {
		tab.table.Refresh()
	}
}

func keyMatchesQuery(key *azkeys.KeyProperties, query string) bool {
	if key == nil {
		return false
	}

	var fields []string
	if key.KID != nil {
		fields = append(fields, string(*key.KID), key.KID.Name())
	}
	if key.Attributes != nil {
		if key.Attributes.Created != nil {
			fields = append(fields, key.Attributes.Created.String())
		}
		if key.Attributes.Expires != nil {
			fields = append(fields, key.Attributes.Expires.String())
		}
		if key.Attributes.Enabled != nil {
			fields = append(fields, fmt.Sprint(*key.Attributes.Enabled))
		}
	}
	if key.Managed != nil {
		fields = append(fields, fmt.Sprint(*key.Managed))
	}
	for name, value := range key.Tags {
		fields = append(fields, name)
		if value != nil {
			fields = append(fields, *value)
		}
	}

	for _, field := range fields {
		if strings.Contains(strings.ToLower(field), query) {
			return true
		}
	}
	return false
}

func (tab *KeyTab) downloadPublicKey(keyName string) {
	tab.action.Start("Fetching public key")
	go func() {
		publicKey, err := services.GetPublicKey(tab.vaultName, keyName)
		fyne.Do(func() {
			if err != nil {
				tab.action.Stop("Failed to fetch public key")
				notifications.OnErrorContinue(err, tab.parent)
				return
			}

			tab.action.Stop("Choose where to save public key")
			tab.newPublicKeySaveDialog(keyName, publicKey).Show()
		})
	}()
}

func (tab *KeyTab) newPublicKeySaveDialog(keyName string, publicKey services.PublicKeyExport) *dialog.FileDialog {
	saveDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			tab.action.Stop("Failed to save public key")
			notifications.OnErrorContinue(err, tab.parent)
			return
		}
		if writer == nil {
			tab.action.Stop("Public key download cancelled")
			return
		}

		if err := writePublicKey(writer, publicKey.Data); err != nil {
			tab.action.Stop("Failed to save public key")
			notifications.OnErrorContinue(err, tab.parent)
			return
		}
		tab.action.Stop("Public key downloaded")
	}, tab.parent)
	saveDialog.SetFileName(keyName + publicKey.Extension)
	saveDialog.SetFilter(storage.NewExtensionFileFilter([]string{publicKey.Extension}))
	return saveDialog
}

func writePublicKey(writer io.WriteCloser, publicKey []byte) error {
	_, writeErr := io.Copy(writer, bytes.NewReader(publicKey))
	closeErr := writer.Close()
	if writeErr != nil {
		return writeErr
	}
	return closeErr
}
