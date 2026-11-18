package notifications

import (
	"errors"
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/require"
)

type closeTrackingWindow struct {
	fyne.Window
	closed bool
}

func (window *closeTrackingWindow) Close() {
	window.closed = true
	window.Window.Close()
}

func TestOnErrorContinue(t *testing.T) {
	test.NewApp()
	parent := test.NewWindow(nil)

	OnErrorContinue(nil, parent)
	require.Empty(t, parent.Canvas().Overlays().List())

	OnErrorContinue(errors.New("continue error"), parent)
	require.NotEmpty(t, parent.Canvas().Overlays().List())
}

func TestOnErrorCloseWindow(t *testing.T) {
	test.NewApp()
	parent := &closeTrackingWindow{Window: test.NewWindow(nil)}
	parent.Show()

	OnErrorCloseWindow(errors.New("window error"), parent)
	require.NotEmpty(t, parent.Canvas().Overlays().List())
	dismissFirstDialogButton(t, parent)
	require.True(t, parent.closed)
}

func TestOnErrorCloseApplication(t *testing.T) {
	app := test.NewApp()
	parent := test.NewWindow(nil)

	OnErrorCloseApplication(errors.New("application error"), parent)
	require.NotEmpty(t, parent.Canvas().Overlays().List())
	dismissFirstDialogButton(t, parent)
	require.NotNil(t, app)
}

func dismissFirstDialogButton(t *testing.T, parent fyne.Window) {
	t.Helper()
	overlay := parent.Canvas().Overlays().Top()
	require.NotNil(t, overlay)
	var button *widget.Button
	for _, object := range test.LaidOutObjects(overlay) {
		if candidate, ok := object.(*widget.Button); ok {
			button = candidate
			break
		}
	}
	require.NotNil(t, button)
	test.Tap(button)
}
