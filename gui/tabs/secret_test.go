package tabs

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/nu12/kvm/factory"
	"github.com/nu12/kvm/gui/notifications"
	"github.com/stretchr/testify/require"
)

func TestSecretTabFiltersVisibleSecrets(t *testing.T) {
	test.NewApp()
	tab := NewSecretTab("test", nil, notifications.NewActivity())
	tab.secrets = []*azsecrets.SecretProperties{
		{
			ID:          to.Ptr(azsecrets.ID("https://test.vault.azure.net/secrets/Alpha")),
			ContentType: to.Ptr("text/plain"),
		},
		{
			ID:          to.Ptr(azsecrets.ID("https://test.vault.azure.net/secrets/Beta")),
			ContentType: to.Ptr("application/json"),
		},
	}

	tab.searchEntry.SetText("APPLICATION")
	require.Len(t, tab.visible, 1)
	require.Equal(t, "Beta", tab.visible[0].ID.Name())

	tab.searchEntry.SetText("alpha")
	require.Len(t, tab.visible, 1)
	require.Equal(t, "Alpha", tab.visible[0].ID.Name())

	tab.searchEntry.SetText("")
	require.Equal(t, tab.secrets, tab.visible)
}

func TestSecretTabLoadsSecrets(t *testing.T) {
	test.NewApp()
	factory.GetInstance().SecretFactory = factory.GetSecretClientMock(2, nil)
	tab := NewSecretTab("test", nil, notifications.NewActivity())

	tab.Load()

	require.Len(t, tab.secrets, 2)
	require.Equal(t, tab.secrets, tab.visible)
}
