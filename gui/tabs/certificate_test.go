package tabs

import (
	"bytes"
	"errors"
	"os"
	"path/filepath"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/storage"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/nu12/kvm/factory"
	"github.com/nu12/kvm/gui/components"
	"github.com/nu12/kvm/gui/notifications"
	"github.com/stretchr/testify/require"
)

type testCertificateWriter struct {
	bytes.Buffer
	closed   bool
	writeErr error
	closeErr error
}

func (writer *testCertificateWriter) Write(data []byte) (int, error) {
	if writer.writeErr != nil {
		return 0, writer.writeErr
	}
	return writer.Buffer.Write(data)
}

func (writer *testCertificateWriter) Close() error {
	writer.closed = true
	return writer.closeErr
}

func TestCertificateTabFiltersVisibleCertificates(t *testing.T) {
	test.NewApp()
	tab := NewCertificateTab("test", nil, notifications.NewActivity())
	tab.certificates = []*azcertificates.CertificateProperties{
		{
			ID:             to.Ptr(azcertificates.ID("https://test.vault.azure.net/certificates/Alpha")),
			X509Thumbprint: []byte{0xab, 0xcd},
			Attributes: &azcertificates.CertificateAttributes{
				Created: to.Ptr(time.Date(2025, time.January, 2, 3, 4, 5, 0, time.UTC)),
				Enabled: to.Ptr(true),
			},
			Tags: map[string]*string{"environment": to.Ptr("Production")},
		},
		{
			ID: to.Ptr(azcertificates.ID("https://test.vault.azure.net/certificates/Beta")),
		},
	}

	tab.searchEntry.SetText("ABCD")
	require.Len(t, tab.visible, 1)
	require.Equal(t, "Alpha", tab.visible[0].ID.Name())

	tab.searchEntry.SetText("beta")
	require.Len(t, tab.visible, 1)
	require.Equal(t, "Beta", tab.visible[0].ID.Name())

	for _, query := range []string{"production", "environment", "true", "2025-01-02"} {
		tab.searchEntry.SetText(query)
		require.Len(t, tab.visible, 1, query)
		require.Equal(t, "Alpha", tab.visible[0].ID.Name())
	}

	tab.searchEntry.SetText("")
	require.Equal(t, tab.certificates, tab.visible)

	cell := tab.table.CreateCell().(*components.CertificateCell)
	tab.table.UpdateCell(widget.TableCellID{Row: 0, Col: 0}, cell)
	require.Equal(t, "Alpha", cell.Text)
	require.NotNil(t, cell.DoubleTapFunc)
	tab.table.UpdateCell(widget.TableCellID{Row: -1, Col: -1}, cell)
	require.Nil(t, cell.DoubleTapFunc)
}

func TestCertificateTabLoadsOnSelection(t *testing.T) {
	test.NewApp()
	factory.GetInstance().CertificateFactory = factory.GetCertificateClientMock(2, nil)

	tabs := New("test", fyne.CurrentApp().NewWindow("test"), notifications.NewActivity())
	require.False(t, tabs.Certificates.loaded)

	tabs.Widget.Select(tabs.Certificates.Content())
	require.Eventually(t, func() bool {
		return tabs.Certificates.loaded
	}, time.Second, 10*time.Millisecond)
	require.Len(t, tabs.Certificates.certificates, 2)
	require.Equal(t, tabs.Certificates.certificates, tabs.Certificates.visible)
}

func TestCertificateTabLoadCanRetryAfterError(t *testing.T) {
	test.NewApp()
	factory.GetInstance().CertificateFactory = func(string) (factory.CertificateClient, error) {
		return nil, errors.New("certificate client unavailable")
	}

	parent := test.NewWindow(nil)
	tabs := New("test", parent, notifications.NewActivity())
	tabs.Widget.Select(tabs.Certificates.Content())
	require.Eventually(t, func() bool {
		return !tabs.Certificates.loading
	}, time.Second, 10*time.Millisecond)
	require.False(t, tabs.Certificates.loaded)
	require.NotEmpty(t, parent.Canvas().Overlays().List())
}

func TestCertificateTabDownloadOpensSaveDialog(t *testing.T) {
	test.NewApp()
	factory.GetInstance().CertificateFactory = factory.GetCertificateClientMock(1, nil)
	parent := test.NewWindow(nil)
	tab := NewCertificateTab("test", parent, notifications.NewActivity())

	tab.download("certificate-0")

	require.Eventually(t, func() bool {
		return len(parent.Canvas().Overlays().List()) > 0
	}, time.Second, 10*time.Millisecond)
	entry := findFileDialogEntry(t, parent)
	require.Equal(t, "certificate-0.cer", entry.Text)
}

func TestCertificateSaveDialogWritesCertificate(t *testing.T) {
	test.NewApp()
	parent := test.NewWindow(nil)
	tab := NewCertificateTab("test", parent, notifications.NewActivity())
	content := []byte("certificate data")
	directory := t.TempDir()
	location, err := storage.ListerForURI(storage.NewFileURI(directory))
	require.NoError(t, err)
	saveDialog := tab.newCertificateSaveDialog("example", content)
	saveDialog.SetLocation(location)
	saveDialog.Show()

	entry := findFileDialogEntry(t, parent)
	require.Equal(t, "example.cer", entry.Text)
	entry.OnSubmitted(entry.Text)
	saved, err := os.ReadFile(filepath.Join(directory, "example.cer"))
	require.NoError(t, err)
	require.Equal(t, content, saved)
}

func TestWriteCertificate(t *testing.T) {
	writer := &testCertificateWriter{}
	content := []byte("certificate bytes")

	require.NoError(t, writeCertificate(writer, content))
	require.Equal(t, content, writer.Bytes())
	require.True(t, writer.closed)

	writer = &testCertificateWriter{writeErr: errors.New("write failed")}
	require.Error(t, writeCertificate(writer, content))
	require.True(t, writer.closed)

	writer = &testCertificateWriter{closeErr: errors.New("close failed")}
	require.Error(t, writeCertificate(writer, content))
	require.True(t, writer.closed)
}

func findFileDialogEntry(t *testing.T, parent fyne.Window) *widget.Entry {
	t.Helper()
	for _, object := range test.LaidOutObjects(parent.Canvas().Overlays().Top()) {
		if entry, ok := object.(*widget.Entry); ok {
			return entry
		}
	}
	t.Fatal("file dialog entry not found")
	return nil
}
