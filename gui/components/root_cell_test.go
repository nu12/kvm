package components

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault/v2"
	"github.com/stretchr/testify/assert"
)

func TestNewRootCell(t *testing.T) {
	c := NewRootCell()
	assert.NotNil(t, c)
	assert.Equal(t, widget.TableCellID{Row: 0, Col: 0}, c.CellID)
	assert.Nil(t, c.DoubleTapFunc)
	assert.Equal(t, "", c.Text)
}

func TestSetTextFromVault(t *testing.T) {
	c := &RootCell{
		CellID: widget.TableCellID{
			Row: 1,
			Col: 0,
		},
	}
	c.SetTextFromVault(&armkeyvault.Vault{
		ID:   to.Ptr("/subscriptions/sub-0/resourceGroups/test/providers/Microsoft.KeyVault/vaults/test-0"),
		Name: to.Ptr("test-0"),
	})
	assert.Equal(t, "test-0", c.Text)
}

func TestSetTextFromVaultColumns(t *testing.T) {
	vault := &armkeyvault.Vault{
		ID:       to.Ptr("/subscriptions/sub-0/resourceGroups/test/providers/Microsoft.KeyVault/vaults/test-0"),
		Name:     to.Ptr("test-0"),
		Location: to.Ptr("West US"),
	}
	cell := &RootCell{}

	for col, expected := range map[int]string{
		1: "test",
		2: "sub-0",
		3: "West US",
	} {
		cell.CellID.Col = col
		cell.SetTextFromVault(vault)
		assert.Equal(t, expected, cell.Text)
	}
}

func TestRootCellDoubleTapped(t *testing.T) {
	c := &RootCell{}
	assert.NotPanics(t, func() { c.DoubleTapped(&fyne.PointEvent{}) })

	called := false
	c.DoubleTapFunc = func() { called = true }
	c.DoubleTapped(&fyne.PointEvent{})
	assert.Equal(t, true, called)
}

func TestRootCellTappedSecondary(t *testing.T) {
	// 1. Initialize the headless test app (sets fyne.CurrentApp() + driver)
	test.NewApp()
	defer test.NewApp() // reset between tests if needed

	// 2. Create your widget with a spy on DoubleTapFunc
	doubleTapCalled := false
	cell := &RootCell{
		DoubleTapFunc: func() { doubleTapCalled = true },
	}

	// 3. Place the widget on a test window so the driver can
	//    resolve AbsolutePositionForObject / CanvasForObject
	w := test.NewWindow(cell)
	w.Resize(fyne.NewSize(200, 200))

	// 4. Simulate a secondary (right) tap at a specific position
	test.TapSecondaryAt(cell, fyne.NewPos(10, 10))

	// 5. Assert the popup menu appeared on the canvas overlays.
	//    ShowPopUpMenuAtPosition adds a PopUpMenu to the canvas overlay.
	overlays := w.Canvas().Overlays().List()
	assert.NotEmpty(t, overlays, "expected a popup overlay after secondary tap")

	// 6. Optionally: find the PopUpMenu in the overlay and click the
	//    "Open Key Vault objects" item to verify DoubleTapFunc fires.
	//    The overlay is a *widget.OverlayContainer wrapping the PopUpMenu.
	//    You can use test.LaidOutObjects to walk the tree:
	objects := test.LaidOutObjects(overlays[0])
	var menu *widget.PopUpMenu
	for _, o := range objects {
		if pm, ok := o.(*widget.PopUpMenu); ok {
			menu = pm
			break
		}
	}
	if assert.NotNil(t, menu, "expected a PopUpMenu in the overlay") {
		// menu.Items is []fyne.CanvasObject (internal *menuItem widgets)
		// They implement fyne.Tappable — Tapped() calls the underlying Action
		if tappable, ok := menu.Items[0].(fyne.Tappable); ok {
			test.Tap(tappable)
		}
	}

	assert.True(t, doubleTapCalled, "DoubleTapFunc should have been called")
}
