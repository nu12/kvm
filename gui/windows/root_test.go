package windows

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault/v2"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestInit(t *testing.T) {
	window := &RootWindow{}
	test.NewApp()
	window.Init()

	assert.NotNil(t, window.app)
	assert.NotNil(t, window.this)
	assert.Equal(t, fyne.Size(fyne.Size{Width: 820, Height: 400}), window.this.Canvas().Size())
	assert.NotNil(t, window.Vaults)
	assert.NotNil(t, window.VisibleVaults)
	assert.NotNil(t, window.Components)
	assert.Nil(t, window.Components.Counter)
	assert.Nil(t, window.Components.Search)
	assert.Nil(t, window.Components.Table)
}

func TestCreateComponents(t *testing.T) {
	window := &RootWindow{}

	require.Nil(t, window.Components.Counter)
	require.Nil(t, window.Components.Search)
	require.Nil(t, window.Components.Table)
	window.CreateComponents()
	assert.NotNil(t, window.Components.Counter)
	assert.NotNil(t, window.Components.Search)
	assert.NotNil(t, window.Components.Table)
}

func TestSearchFiltersVisibleVaults(t *testing.T) {
	test.NewApp()
	window := &RootWindow{
		Vaults: []*armkeyvault.Vault{
			{
				ID:       to.Ptr("/subscriptions/sub-1/resourceGroups/rg-one/providers/Microsoft.KeyVault/vaults/Alpha"),
				Name:     to.Ptr("Alpha"),
				Location: to.Ptr("West US"),
			},
			{
				ID:       to.Ptr("/subscriptions/sub-2/resourceGroups/rg-two/providers/Microsoft.KeyVault/vaults/Beta"),
				Name:     to.Ptr("Beta"),
				Location: to.Ptr("Canada Central"),
			},
		},
	}
	window.VisibleVaults = append([]*armkeyvault.Vault(nil), window.Vaults...)
	window.CreateComponents()

	window.Components.SearchEntry.SetText("cAnAdA")
	require.Len(t, window.VisibleVaults, 1)
	require.Equal(t, "Beta", *window.VisibleVaults[0].Name)
	require.Equal(t, "Showing 1 out of 2 key vaults.", window.Components.Counter.Text)

	window.Components.SearchEntry.SetText("rg-one")
	require.Len(t, window.VisibleVaults, 1)
	require.Equal(t, "Alpha", *window.VisibleVaults[0].Name)

	window.Components.SearchEntry.SetText("")
	require.Equal(t, window.Vaults, window.VisibleVaults)
	require.Equal(t, "Showing 2 out of 2 key vaults.", window.Components.Counter.Text)
}

func TestSetContent(t *testing.T) {
	window := &RootWindow{
		this: test.NewApp().NewWindow("Test"),
	}
	require.Nil(t, window.this.Content())

	window.Components.Counter = widget.NewLabel("Counter")
	window.Components.Search = container.NewCenter()
	window.Components.Table = widget.NewTable(
		func() (rows int, cols int) { return 0, 0 },
		func() fyne.CanvasObject { return widget.NewLabel("Test") },
		func(tci widget.TableCellID, co fyne.CanvasObject) {},
	)
	window.SetContent()
	assert.NotNil(t, window.this.Content())
}
