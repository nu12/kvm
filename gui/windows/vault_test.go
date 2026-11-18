package windows

import (
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/nu12/kvm/factory"
	"github.com/nu12/kvm/gui/notifications"
	"github.com/stretchr/testify/require"
)

func TestVaultInit(t *testing.T) {
	vaultWindow := initVaultWindow("test")
	require.NotNil(t, vaultWindow.Components)
	require.NotNil(t, vaultWindow.Secrets)
	require.NotNil(t, vaultWindow.VisibleSecrets)
	require.Equal(t, 0, len(vaultWindow.Secrets))
	require.NotNil(t, vaultWindow.action)
	require.True(t, vaultWindow.action.Visible())

	require.NotNil(t, vaultWindow.Components.SecretsTab)
	require.Nil(t, vaultWindow.Components.SecretsTab.Content)
	require.Nil(t, vaultWindow.Components.SecretsTab.Search)
	require.Nil(t, vaultWindow.Components.SecretsTab.Table)

	require.NotNil(t, vaultWindow.Components.CertificatesTab)
	require.Nil(t, vaultWindow.Components.CertificatesTab.Content)

	require.NotNil(t, vaultWindow.Components.KeysTab)
	require.Nil(t, vaultWindow.Components.KeysTab.Content)
}

func TestVaultCreateComponents(t *testing.T) {
	vaultWindow := &VaultWindow{}

	require.Nil(t, vaultWindow.Components.SecretsTab.Search)
	require.Nil(t, vaultWindow.Components.SecretsTab.Table)
	vaultWindow.createComponents()
	require.NotNil(t, vaultWindow.Components.SecretsTab.Search)
	require.NotNil(t, vaultWindow.Components.SecretsTab.SearchEntry)
	require.NotNil(t, vaultWindow.Components.SecretsTab.Table)
}

func TestSearchFiltersVisibleSecrets(t *testing.T) {
	test.NewApp()
	vaultWindow := &VaultWindow{
		Secrets: []*azsecrets.SecretProperties{
			{
				ID:          to.Ptr(azsecrets.ID("https://test.vault.azure.net/secrets/Alpha")),
				ContentType: to.Ptr("text/plain"),
			},
			{
				ID:          to.Ptr(azsecrets.ID("https://test.vault.azure.net/secrets/Beta")),
				ContentType: to.Ptr("application/json"),
			},
		},
	}
	vaultWindow.VisibleSecrets = append([]*azsecrets.SecretProperties(nil), vaultWindow.Secrets...)
	vaultWindow.createComponents()

	vaultWindow.Components.SecretsTab.SearchEntry.SetText("APPLICATION")
	require.Len(t, vaultWindow.VisibleSecrets, 1)
	require.Equal(t, "Beta", vaultWindow.VisibleSecrets[0].ID.Name())

	vaultWindow.Components.SecretsTab.SearchEntry.SetText("alpha")
	require.Len(t, vaultWindow.VisibleSecrets, 1)
	require.Equal(t, "Alpha", vaultWindow.VisibleSecrets[0].ID.Name())

	vaultWindow.Components.SecretsTab.SearchEntry.SetText("")
	require.Equal(t, vaultWindow.Secrets, vaultWindow.VisibleSecrets)
}

func TestVaultSetContent(t *testing.T) {
	test.NewApp()
	vaultWindow := &VaultWindow{
		this:   fyne.CurrentApp().NewWindow("Vault: " + "test"),
		action: notifications.NewActivity(),
	}

	vaultWindow.Components.SecretsTab.Search = container.NewCenter(widget.NewEntry())
	vaultWindow.Components.SecretsTab.Table = widget.NewTable(
		func() (rows int, cols int) { return 0, 0 },
		func() fyne.CanvasObject { return widget.NewLabel("test") },
		func(tci widget.TableCellID, co fyne.CanvasObject) {},
	)
	vaultWindow.setContent()
	require.NotNil(t, vaultWindow.Components.SecretsTab.Content)
	require.NotNil(t, vaultWindow.Components.CertificatesTab.Content)
	require.NotNil(t, vaultWindow.Components.KeysTab.Content)
}

func TestVaultRun(t *testing.T) {
	f := factory.GetInstance()
	f.SecretFactory = factory.GetSecretClientMock(2, nil)

	vaultWindow := &VaultWindow{
		action:  notifications.NewActivity(),
		Secrets: []*azsecrets.SecretProperties{},
		Components: struct {
			SecretsTab struct {
				Content     *container.TabItem
				Search      *fyne.Container
				SearchEntry *widget.Entry
				Table       *widget.Table
			}
			CertificatesTab struct{ Content *container.TabItem }
			KeysTab         struct{ Content *container.TabItem }
		}{
			SecretsTab: struct {
				Content     *container.TabItem
				Search      *fyne.Container
				SearchEntry *widget.Entry
				Table       *widget.Table
			}{
				SearchEntry: widget.NewEntry(),
				Table: widget.NewTable(
					func() (rows int, cols int) { return 0, 0 },
					func() fyne.CanvasObject { return widget.NewLabel("Test") },
					func(tci widget.TableCellID, co fyne.CanvasObject) {},
				),
			},
		},
	}

	vaultWindow.run()
	require.Eventually(t, func() bool {
		return len(vaultWindow.Secrets) == 2
	}, time.Second, 10*time.Millisecond)
}
