package tabs

import (
	"bytes"
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/nu12/kvm/factory"
	"github.com/nu12/kvm/gui/notifications"
	"github.com/stretchr/testify/require"
)

type testCertificateWriter struct {
	bytes.Buffer
	closed bool
}

func (writer *testCertificateWriter) Close() error {
	writer.closed = true
	return nil
}

func TestCertificateTabFiltersVisibleCertificates(t *testing.T) {
	test.NewApp()
	tab := NewCertificateTab("test", nil, notifications.NewActivity())
	tab.certificates = []*azcertificates.CertificateProperties{
		{
			ID:             to.Ptr(azcertificates.ID("https://test.vault.azure.net/certificates/Alpha")),
			X509Thumbprint: []byte{0xab, 0xcd},
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

	tab.searchEntry.SetText("")
	require.Equal(t, tab.certificates, tab.visible)
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

func TestWriteCertificate(t *testing.T) {
	writer := &testCertificateWriter{}
	content := []byte("certificate bytes")

	require.NoError(t, writeCertificate(writer, content))
	require.Equal(t, content, writer.Bytes())
	require.True(t, writer.closed)
}
