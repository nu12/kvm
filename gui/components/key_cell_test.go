package components

import (
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/stretchr/testify/require"
)

func TestKeyCellDisplaysMetadata(t *testing.T) {
	created := time.Date(2025, time.January, 2, 3, 4, 5, 0, time.UTC)
	expires := created.Add(24 * time.Hour)
	key := &azkeys.KeyProperties{
		KID: to.Ptr(azkeys.ID("https://test.vault.azure.net/keys/example")),
		Attributes: &azkeys.KeyAttributes{
			Created: &created,
			Expires: &expires,
			Enabled: to.Ptr(true),
		},
		Managed: to.Ptr(true),
		Tags: map[string]*string{
			"environment": to.Ptr("test"),
			"owner":       nil,
		},
	}
	cell := NewKeyCell()

	tests := []struct {
		column   int
		expected string
	}{
		{0, "example"},
		{1, created.String()},
		{2, expires.String()},
		{3, "Enabled"},
		{4, "true"},
		{5, "environment=test, owner"},
	}
	for _, testCase := range tests {
		cell.CellID = widget.TableCellID{Col: testCase.column}
		cell.SetTextFromKey(key)
		require.Equal(t, testCase.expected, cell.Text)
	}
}

func TestKeyCellHandlesMissingMetadata(t *testing.T) {
	cell := NewKeyCell()

	cell.CellID.Col = 0
	cell.SetTextFromKey(&azkeys.KeyProperties{})
	require.Empty(t, cell.Text)

	for col, expected := range map[int]string{
		1: "N/A",
		2: "No expiration",
		3: "N/A",
		4: "N/A",
		5: "",
	} {
		cell.CellID.Col = col
		cell.SetTextFromKey(&azkeys.KeyProperties{})
		require.Equal(t, expected, cell.Text)
	}

	cell.SetTextFromKey(nil)
	require.Empty(t, cell.Text)
}

func TestKeyCellDoubleTapped(t *testing.T) {
	cell := NewKeyCell()
	called := false
	cell.DoubleTapFunc = func() { called = true }

	cell.DoubleTapped(&fyne.PointEvent{})

	require.True(t, called)
}

func TestKeyCellContextMenuUsesDoubleTapAction(t *testing.T) {
	test.NewApp()
	called := false
	cell := &KeyCell{DoubleTapFunc: func() { called = true }}
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
