package components

import (
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/stretchr/testify/require"
)

func TestCertificateCellText(t *testing.T) {
	certificate := &azcertificates.CertificateProperties{
		ID:             to.Ptr(azcertificates.ID("https://test.vault.azure.net/certificates/example")),
		X509Thumbprint: []byte{0xab, 0xcd},
		Attributes: &azcertificates.CertificateAttributes{
			Enabled: to.Ptr(true),
		},
	}
	cell := NewCertificateCell()

	cell.CellID = widget.TableCellID{Col: 0}
	cell.SetTextFromCertificate(certificate)
	require.Equal(t, "example", cell.Text)

	cell.CellID = widget.TableCellID{Col: 2}
	cell.SetTextFromCertificate(certificate)
	require.Equal(t, "No expiration", cell.Text)

	cell.CellID = widget.TableCellID{Col: 3}
	cell.SetTextFromCertificate(certificate)
	require.Equal(t, "Enabled", cell.Text)

	cell.CellID = widget.TableCellID{Col: 4}
	cell.SetTextFromCertificate(certificate)
	require.Equal(t, "abcd", cell.Text)
}

func TestCertificateCellTextMetadata(t *testing.T) {
	created := time.Date(2025, time.January, 2, 3, 4, 5, 0, time.UTC)
	expires := created.Add(24 * time.Hour)
	certificate := &azcertificates.CertificateProperties{
		Attributes: &azcertificates.CertificateAttributes{
			Created: &created,
			Expires: &expires,
			Enabled: to.Ptr(false),
		},
	}
	cell := NewCertificateCell()

	cell.CellID.Col = 1
	cell.SetTextFromCertificate(certificate)
	require.Equal(t, created.String(), cell.Text)

	cell.CellID.Col = 2
	cell.SetTextFromCertificate(certificate)
	require.Equal(t, expires.String(), cell.Text)

	cell.CellID.Col = 3
	cell.SetTextFromCertificate(certificate)
	require.Equal(t, "Disabled", cell.Text)

	certificate.Attributes = nil
	cell.CellID.Col = 1
	cell.SetTextFromCertificate(certificate)
	require.Equal(t, "N/A", cell.Text)

	cell.CellID.Col = 3
	cell.SetTextFromCertificate(certificate)
	require.Equal(t, "N/A", cell.Text)

	cell.CellID.Col = 4
	cell.SetTextFromCertificate(certificate)
	require.Equal(t, "N/A", cell.Text)
}

func TestCertificateCellDoubleTapped(t *testing.T) {
	cell := NewCertificateCell()
	called := false
	cell.DoubleTapFunc = func() { called = true }
	cell.DoubleTapped(&fyne.PointEvent{})
	require.True(t, called)
}

func TestCertificateCellTappedSecondaryUsesDoubleTapAction(t *testing.T) {
	test.NewApp()
	called := false
	cell := &CertificateCell{
		DoubleTapFunc: func() { called = true },
	}
	window := test.NewWindow(cell)
	window.Resize(fyne.NewSize(200, 200))

	test.TapSecondaryAt(cell, fyne.NewPos(10, 10))

	overlays := window.Canvas().Overlays().List()
	require.NotEmpty(t, overlays)
	var menu *widget.PopUpMenu
	for _, object := range test.LaidOutObjects(overlays[0]) {
		if popupMenu, ok := object.(*widget.PopUpMenu); ok {
			menu = popupMenu
			break
		}
	}
	require.NotNil(t, menu)
	require.NotEmpty(t, menu.Items)
	tappable, ok := menu.Items[0].(fyne.Tappable)
	require.True(t, ok)

	test.Tap(tappable)
	require.True(t, called)
}
