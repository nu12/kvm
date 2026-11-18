package components

import (
	"encoding/hex"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/widget"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
)

type CertificateCell struct {
	widget.Label
	DoubleTapFunc func()
	CellID        widget.TableCellID
}

func NewCertificateCell() *CertificateCell {
	cell := &CertificateCell{}
	cell.ExtendBaseWidget(cell)
	return cell
}

func (cell *CertificateCell) SetTextFromCertificate(certificate *azcertificates.CertificateProperties) {
	cell.SetText("")
	if certificate == nil {
		return
	}

	switch cell.CellID.Col {
	case 0:
		if certificate.ID != nil {
			cell.SetText(certificate.ID.Name())
		}
	case 1:
		if certificate.Attributes != nil && certificate.Attributes.Created != nil {
			cell.SetText(certificate.Attributes.Created.String())
		} else {
			cell.SetText("N/A")
		}
	case 2:
		if certificate.Attributes == nil || certificate.Attributes.Expires == nil {
			cell.SetText("No expiration")
		} else {
			cell.SetText(certificate.Attributes.Expires.String())
		}
	case 3:
		if certificate.Attributes == nil || certificate.Attributes.Enabled == nil {
			cell.SetText("N/A")
		} else if *certificate.Attributes.Enabled {
			cell.SetText("Enabled")
		} else {
			cell.SetText("Disabled")
		}
	case 4:
		if len(certificate.X509Thumbprint) == 0 {
			cell.SetText("N/A")
		} else {
			cell.SetText(hex.EncodeToString(certificate.X509Thumbprint))
		}
	}
}

func (cell *CertificateCell) DoubleTapped(_ *fyne.PointEvent) {
	if cell.DoubleTapFunc != nil {
		cell.DoubleTapFunc()
	}
}

func (cell *CertificateCell) TappedSecondary(event *fyne.PointEvent) {
	menu := fyne.NewMenu("Context",
		fyne.NewMenuItem("Download certificate", func() {
			if cell.DoubleTapFunc != nil {
				cell.DoubleTapFunc()
			}
		}),
	)
	entryPos := fyne.CurrentApp().Driver().AbsolutePositionForObject(cell)
	popupPos := entryPos.Add(event.Position)
	widget.ShowPopUpMenuAtPosition(menu, fyne.CurrentApp().Driver().CanvasForObject(cell), popupPos)
}
