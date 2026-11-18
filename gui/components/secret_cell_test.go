package components

import (
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/stretchr/testify/assert"
)

func TestNewSecretCell(t *testing.T) {
	c := NewSecretCell()
	assert.NotNil(t, c)
	assert.Equal(t, widget.TableCellID{Row: 0, Col: 0}, c.CellID)
	assert.Nil(t, c.DoubleTapFunc)
	assert.Equal(t, "", c.Text)
}

func TestSetTextFromSecret(t *testing.T) {
	c := &SecretCell{
		CellID: widget.TableCellID{
			Row: 0,
			Col: 0,
		},
	}

	c.SetTextFromSecret(&azsecrets.SecretProperties{
		ID: to.Ptr(azsecrets.ID("https://test.vault.azure.net/secrets/secret-0")),
	})
	assert.Equal(t, "secret-0", c.Text)
}

func TestSetTextFromSecretColumns(t *testing.T) {
	created := time.Date(2025, time.January, 2, 3, 4, 5, 0, time.UTC)
	expires := created.Add(24 * time.Hour)
	secret := &azsecrets.SecretProperties{
		Attributes: &azsecrets.SecretAttributes{
			Created: &created,
			Expires: &expires,
			Enabled: to.Ptr(true),
		},
		ContentType: to.Ptr("application/json"),
	}
	cell := &SecretCell{}

	cell.CellID.Col = 1
	cell.SetTextFromSecret(secret)
	assert.Equal(t, created.String(), cell.Text)

	cell.CellID.Col = 2
	cell.SetTextFromSecret(secret)
	assert.Equal(t, expires.String(), cell.Text)
	secret.Attributes.Expires = nil
	cell.SetTextFromSecret(secret)
	assert.Equal(t, "No expiration", cell.Text)

	cell.CellID.Col = 3
	cell.SetTextFromSecret(secret)
	assert.Equal(t, "application/json", cell.Text)
	secret.ContentType = nil
	cell.SetTextFromSecret(secret)
	assert.Equal(t, "N/A", cell.Text)

	cell.CellID.Col = 4
	cell.SetTextFromSecret(secret)
	assert.Equal(t, "Enabled", cell.Text)
	*secret.Attributes.Enabled = false
	cell.SetTextFromSecret(secret)
	assert.Equal(t, "Disabled", cell.Text)
}

func TestDoubleTapped(t *testing.T) {
	c := &SecretCell{}
	assert.NotPanics(t, func() { c.DoubleTapped(&fyne.PointEvent{}) })

	called := false
	c.DoubleTapFunc = func() { called = true }
	c.DoubleTapped(&fyne.PointEvent{})
	assert.True(t, called)
}

func TestTappedSecondary(t *testing.T) {
	// 1. Initialize the headless test app (sets fyne.CurrentApp() + driver)
	test.NewApp()
	defer test.NewApp() // reset between tests if needed

	// 2. Create your widget with a spy on DoubleTapFunc
	doubleTapCalled := false
	cell := &SecretCell{
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
