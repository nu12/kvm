package windows

import (
	"testing"
	"time"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/nu12/kvm/factory"
	"github.com/nu12/kvm/services"
	"github.com/stretchr/testify/require"
)

func TestNewSplash(t *testing.T) {
	// Set up the Fyne test app (replaces GLFW driver)
	app := test.NewApp()
	_ = app

	// Create a minimal RootWindow for the test
	root := &RootWindow{
		Logo: test.NewApp().Icon(), // any fyne.Resource
		this: test.NewWindow(nil),
	}

	splash, err := NewSplash(root)
	if err != nil {
		t.Fatalf("NewSplash returned error: %v", err)
	}
	if splash == nil {
		t.Fatal("NewSplash returned nil")
	}
	// With the test driver (not desktop.Driver), splash.this will be nil
	// because the type assertion to desktop.Driver fails.
	// That's expected — the splash window only exists on real desktop.
}

func TestLoadAsync(t *testing.T) {
	f := factory.GetInstance()
	f.SubscriptionFactory = factory.GetSubscriptionClientMock(1, nil)
	f.VaultFactory = factory.GetVaultClientMock(3, nil)

	app := test.NewApp()
	_ = app

	root := &RootWindow{
		this: test.NewWindow(nil),
	}
	root.Vaults, _ = services.GetKeyVaults()

	splash := &SplashWindow{
		parent: root,
		this:   test.NewWindow(nil),
	}

	table := widget.NewTable(
		func() (int, int) { return 3, 1 },
		func() fyne.CanvasObject { return widget.NewLabel("") },
		func(id widget.TableCellID, o fyne.CanvasObject) {
			o.(*widget.Label).SetText(*root.Vaults[id.Row].Name)
		},
	)
	counter := widget.NewLabel("")

	// Synchronize: wrap LoadAsync's goroutine completion
	done := make(chan struct{})
	originalClose := splash.this.Close
	splash.this = &closeTrackingWindow{
		Window: test.NewWindow(nil),
		onClose: func() {
			originalClose()
			close(done)
		},
	}

	err := splash.LoadAsync(table, counter)
	require.NoError(t, err)

	// Wait for the goroutine to finish (with timeout)
	select {
	case <-done:
		// success
	case <-time.After(2 * time.Second):
		t.Fatal("timed out waiting for LoadAsync goroutine to complete")
	}

	require.Equal(t, 3, len(root.Vaults))
	require.Equal(t, "Showing 3 out of 3 key vaults.", counter.Text)
}

// Helper: a window that signals when Close is called
type closeTrackingWindow struct {
	fyne.Window
	onClose func()
}

func (w *closeTrackingWindow) Close() {
	if w.onClose != nil {
		w.onClose()
	}
}
