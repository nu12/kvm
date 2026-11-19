package tabs

import (
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
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/nu12/kvm/factory"
	"github.com/nu12/kvm/gui/components"
	"github.com/nu12/kvm/gui/notifications"
	"github.com/nu12/kvm/services"
	"github.com/stretchr/testify/require"
)

func TestKeyTabContent(t *testing.T) {
	test.NewApp()
	tab := NewKeyTab("test", nil, notifications.NewActivity())

	require.Equal(t, "Keys", tab.Content().Text)
	require.NotNil(t, tab.Content().Content)
	require.Equal(t, "Type to filter keys...", tab.searchEntry.PlaceHolder)
}

func TestKeyTabFiltersVisibleKeys(t *testing.T) {
	test.NewApp()
	created := time.Date(2025, time.January, 2, 3, 4, 5, 0, time.UTC)
	tab := NewKeyTab("test", nil, notifications.NewActivity())
	tab.keys = []*azkeys.KeyProperties{
		{
			KID: to.Ptr(azkeys.ID("https://test.vault.azure.net/keys/Alpha")),
			Attributes: &azkeys.KeyAttributes{
				Created: &created,
				Enabled: to.Ptr(true),
			},
			Managed: to.Ptr(true),
			Tags:    map[string]*string{"environment": to.Ptr("Production")},
		},
		{
			KID: to.Ptr(azkeys.ID("https://test.vault.azure.net/keys/Beta")),
		},
	}

	for _, query := range []string{"ALPHA", "production", "environment", "true", "2025-01-02"} {
		tab.searchEntry.SetText(query)
		require.Len(t, tab.visible, 1, query)
		require.Equal(t, "Alpha", tab.visible[0].KID.Name())
	}
	tab.searchEntry.SetText("beta")
	require.Len(t, tab.visible, 1)
	require.Equal(t, "Beta", tab.visible[0].KID.Name())
	tab.searchEntry.SetText("")
	require.Equal(t, tab.keys, tab.visible)

	cell := tab.table.CreateCell().(*components.KeyCell)
	tab.table.UpdateCell(widget.TableCellID{Row: 0, Col: 0}, cell)
	require.Equal(t, "Alpha", cell.Text)
	require.NotNil(t, cell.DoubleTapFunc)
	tab.table.UpdateCell(widget.TableCellID{Row: -1, Col: -1}, cell)
	require.Empty(t, cell.Text)
	require.Nil(t, cell.DoubleTapFunc)
}

func TestKeyTabLoadsOnSelection(t *testing.T) {
	test.NewApp()
	factory.GetInstance().KeyFactory = factory.GetKeyClientMock(2, nil)

	tabSet := New("test", fyne.CurrentApp().NewWindow("test"), notifications.NewActivity())
	require.False(t, tabSet.Keys.loaded)

	tabSet.Widget.Select(tabSet.Keys.Content())
	require.Eventually(t, func() bool {
		return tabSet.Keys.loaded
	}, time.Second, 10*time.Millisecond)
	require.Len(t, tabSet.Keys.keys, 2)
	require.Equal(t, tabSet.Keys.keys, tabSet.Keys.visible)
}

func TestKeyTabShowsLoadError(t *testing.T) {
	test.NewApp()
	factory.GetInstance().KeyFactory = func(string) (factory.KeyClient, error) {
		return nil, errors.New("key client unavailable")
	}
	parent := test.NewWindow(nil)
	tabSet := New("test", parent, notifications.NewActivity())

	tabSet.Widget.Select(tabSet.Keys.Content())
	require.Eventually(t, func() bool {
		return !tabSet.Keys.loading
	}, time.Second, 10*time.Millisecond)
	require.False(t, tabSet.Keys.loaded)
	require.NotEmpty(t, parent.Canvas().Overlays().List())
}

func TestKeyTabDoubleTapOpensPublicKeySaveDialog(t *testing.T) {
	test.NewApp()
	factory.GetInstance().KeyFactory = factory.GetKeyClientMock(1, nil)
	parent := test.NewWindow(nil)
	tab := NewKeyTab("test", parent, notifications.NewActivity())
	tab.keys = []*azkeys.KeyProperties{
		{KID: to.Ptr(azkeys.ID("https://test.vault.azure.net/keys/example"))},
	}
	tab.filter("")

	cell := tab.table.CreateCell().(*components.KeyCell)
	tab.table.UpdateCell(widget.TableCellID{Row: 0, Col: 0}, cell)
	cell.DoubleTapped(&fyne.PointEvent{})

	require.Eventually(t, func() bool {
		return len(parent.Canvas().Overlays().List()) > 0
	}, time.Second, 10*time.Millisecond)
	entry := findFileDialogEntry(t, parent)
	require.Equal(t, "example.pem", entry.Text)
}

func TestKeyTabPublicKeySaveDialogWritesFile(t *testing.T) {
	test.NewApp()
	parent := test.NewWindow(nil)
	tab := NewKeyTab("test", parent, notifications.NewActivity())
	content := []byte("-----BEGIN RSA PUBLIC KEY-----\nAQID\n-----END RSA PUBLIC KEY-----\n")
	directory := t.TempDir()
	location, err := storage.ListerForURI(storage.NewFileURI(directory))
	require.NoError(t, err)
	saveDialog := tab.newPublicKeySaveDialog("example", services.PublicKeyExport{
		Data:      content,
		Extension: ".pem",
	})
	saveDialog.SetLocation(location)
	saveDialog.Show()

	entry := findFileDialogEntry(t, parent)
	require.Equal(t, "example.pem", entry.Text)
	entry.OnSubmitted(entry.Text)
	saved, err := os.ReadFile(filepath.Join(directory, "example.pem"))
	require.NoError(t, err)
	require.Equal(t, content, saved)
}

func TestWritePublicKey(t *testing.T) {
	writer := &testCertificateWriter{}
	content := []byte("public key bytes")

	require.NoError(t, writePublicKey(writer, content))
	require.Equal(t, content, writer.Bytes())
	require.True(t, writer.closed)

	writer = &testCertificateWriter{writeErr: errors.New("write failed")}
	require.Error(t, writePublicKey(writer, content))
	require.True(t, writer.closed)

	writer = &testCertificateWriter{closeErr: errors.New("close failed")}
	require.Error(t, writePublicKey(writer, content))
	require.True(t, writer.closed)
}
