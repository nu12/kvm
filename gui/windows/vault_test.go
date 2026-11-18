package windows

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"github.com/nu12/kvm/gui/notifications"
	"github.com/stretchr/testify/require"
)

func TestVaultInit(t *testing.T) {
	test.NewApp()
	vaultWindow := initVaultWindow("test")
	require.Equal(t, "test", vaultWindow.VaultName)
	require.Nil(t, vaultWindow.Tabs)
	require.NotNil(t, vaultWindow.action)
	require.True(t, vaultWindow.action.Visible())
}

func TestVaultCreateComponents(t *testing.T) {
	test.NewApp()
	vaultWindow := &VaultWindow{
		this:   fyne.CurrentApp().NewWindow("Vault: test"),
		action: notifications.NewActivity(),
	}

	vaultWindow.createComponents()
	require.NotNil(t, vaultWindow.Tabs)
	require.NotNil(t, vaultWindow.Tabs.Secrets)
	require.NotNil(t, vaultWindow.Tabs.Certificates)
	require.NotNil(t, vaultWindow.Tabs.Keys)
	require.Len(t, vaultWindow.Tabs.Widget.Items, 3)
}

func TestVaultSetContent(t *testing.T) {
	test.NewApp()
	vaultWindow := &VaultWindow{
		this:   fyne.CurrentApp().NewWindow("Vault: test"),
		action: notifications.NewActivity(),
	}

	vaultWindow.createComponents().setContent()
	require.NotNil(t, vaultWindow.this.Content())
	require.Same(t, vaultWindow.Tabs.Widget, vaultWindow.this.Content())
}
