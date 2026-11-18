package tabs

import (
	"bytes"
	"encoding/hex"
	"fmt"
	"io"
	"strings"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/widget"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/nu12/kvm/gui/components"
	"github.com/nu12/kvm/gui/notifications"
	"github.com/nu12/kvm/services"
)

type CertificateTab struct {
	vaultName string
	parent    fyne.Window
	action    *notifications.Activity
	content   *container.TabItem

	searchEntry  *widget.Entry
	table        *widget.Table
	certificates []*azcertificates.CertificateProperties
	visible      []*azcertificates.CertificateProperties
	loading      bool
	loaded       bool
}

func NewCertificateTab(vaultName string, parent fyne.Window, action *notifications.Activity) *CertificateTab {
	tab := &CertificateTab{
		vaultName:    vaultName,
		parent:       parent,
		action:       action,
		certificates: []*azcertificates.CertificateProperties{},
		visible:      []*azcertificates.CertificateProperties{},
	}

	tab.searchEntry = widget.NewEntry()
	tab.searchEntry.OnChanged = tab.filter
	tab.searchEntry.PlaceHolder = "Type to filter certificates..."

	tab.table = components.NewTable(components.TableConfig{
		RowCount:    func() int { return len(tab.visible) },
		Headers:     []string{"Name", "Creation date", "Expiration date", "Enabled", "Thumbprint"},
		ColumnWidth: []float32{250, 250, 250, 100, 250},
		NewCell:     func() fyne.CanvasObject { return components.NewCertificateCell() },
		Func: func(id widget.TableCellID, object fyne.CanvasObject) {
			cell := object.(*components.CertificateCell)
			cell.CellID = id
			cell.DoubleTapFunc = nil
			if id.Row < 0 || id.Row >= len(tab.visible) {
				cell.SetText("")
				return
			}

			certificate := tab.visible[id.Row]
			cell.DoubleTapFunc = func() {
				if certificate != nil && certificate.ID != nil {
					tab.download(certificate.ID.Name())
				}
			}
			cell.SetTextFromCertificate(certificate)
		},
	})
	tab.content = container.NewTabItem("Certificates", container.NewBorder(
		container.NewVBox(tab.searchEntry),
		action,
		nil,
		nil,
		tab.table,
	))

	return tab
}

func (tab *CertificateTab) Content() *container.TabItem {
	return tab.content
}

func (tab *CertificateTab) Load() {
	if tab.loaded || tab.loading {
		return
	}
	tab.loading = true
	tab.action.Start("Loading certificates")

	go func() {
		certificates, err := services.GetCertificates(tab.vaultName)
		fyne.Do(func() {
			tab.loading = false
			if err != nil {
				tab.action.Stop("Failed to load certificates")
				notifications.OnErrorContinue(err, tab.parent)
				return
			}

			tab.certificates = certificates
			tab.filter(tab.searchEntry.Text)
			tab.loaded = true
			tab.action.Stop("Done")
		})
	}()
}

func (tab *CertificateTab) filter(query string) {
	query = strings.ToLower(strings.TrimSpace(query))
	tab.visible = make([]*azcertificates.CertificateProperties, 0, len(tab.certificates))

	for _, certificate := range tab.certificates {
		if query == "" || certificateMatchesQuery(certificate, query) {
			tab.visible = append(tab.visible, certificate)
		}
	}

	if tab.table != nil {
		tab.table.Refresh()
	}
}

func certificateMatchesQuery(certificate *azcertificates.CertificateProperties, query string) bool {
	if certificate == nil {
		return false
	}

	var fields []string
	if certificate.ID != nil {
		fields = append(fields, string(*certificate.ID), certificate.ID.Name())
	}
	if certificate.Attributes != nil {
		if certificate.Attributes.Created != nil {
			fields = append(fields, certificate.Attributes.Created.String())
		}
		if certificate.Attributes.Expires != nil {
			fields = append(fields, certificate.Attributes.Expires.String())
		}
		if certificate.Attributes.Enabled != nil {
			fields = append(fields, fmt.Sprint(*certificate.Attributes.Enabled))
		}
	}
	if len(certificate.X509Thumbprint) != 0 {
		fields = append(fields, hex.EncodeToString(certificate.X509Thumbprint))
	}
	for key, value := range certificate.Tags {
		fields = append(fields, key)
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

func (tab *CertificateTab) download(certificateName string) {
	tab.action.Start("Fetching certificate")
	go func() {
		certificate, err := services.GetCertificate(tab.vaultName, certificateName)
		fyne.Do(func() {
			if err != nil {
				tab.action.Stop("Failed to fetch certificate")
				notifications.OnErrorContinue(err, tab.parent)
				return
			}

			tab.action.Stop("Choose where to save certificate")
			tab.newCertificateSaveDialog(certificateName, certificate).Show()
		})
	}()
}

func (tab *CertificateTab) newCertificateSaveDialog(certificateName string, certificate []byte) *dialog.FileDialog {
	saveDialog := dialog.NewFileSave(func(writer fyne.URIWriteCloser, err error) {
		if err != nil {
			tab.action.Stop("Failed to save certificate")
			notifications.OnErrorContinue(err, tab.parent)
			return
		}
		if writer == nil {
			tab.action.Stop("Certificate download cancelled")
			return
		}

		if err := writeCertificate(writer, certificate); err != nil {
			tab.action.Stop("Failed to save certificate")
			notifications.OnErrorContinue(err, tab.parent)
			return
		}
		tab.action.Stop("Certificate downloaded")
	}, tab.parent)
	saveDialog.SetFileName(certificateName + ".cer")
	saveDialog.SetFilter(storage.NewExtensionFileFilter([]string{".cer"}))
	return saveDialog
}

func writeCertificate(writer io.WriteCloser, certificate []byte) error {
	_, writeErr := io.Copy(writer, bytes.NewReader(certificate))
	closeErr := writer.Close()
	if writeErr != nil {
		return writeErr
	}
	return closeErr
}
