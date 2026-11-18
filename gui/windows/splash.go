package windows

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/go-gl/glfw/v3.4/glfw"
	"github.com/nu12/kvm/gui/notifications"
	"github.com/nu12/kvm/services"
)

type SplashWindow struct {
	parent *RootWindow
	this   fyne.Window
}

func NewSplash(parent *RootWindow) *SplashWindow {

	var splash fyne.Window
	if drv, ok := fyne.CurrentApp().Driver().(desktop.Driver); ok {
		if glfw.Init() == nil {
			glfw.WindowHint(glfw.TransparentFramebuffer, glfw.True)
		}
		splash = drv.CreateSplashWindow()
		img := canvas.NewImageFromResource(parent.Logo)
		img.FillMode = canvas.ImageFillOriginal
		splash.SetContent(img)
		splash.Show()
		if glfw.Init() == nil {
			glfw.DefaultWindowHints()
		}
	}
	return &SplashWindow{
		parent: parent,
		this:   splash,
	}
}

func (splash *SplashWindow) LoadAsync(table *widget.Table, counter *widget.Label) {
	go func() {
		vaults, err := services.GetKeyVaults()
		notifications.OnErrorCloseApplication(err, splash.parent.this)
		splash.parent.Vaults = vaults
		table.Refresh()
		counter.SetText(fmt.Sprintf("Showing %d out of %d key vaults.", len(splash.parent.Vaults), len(splash.parent.Vaults)))
		fyne.Do(func() {
			splash.this.Close() // ← Hide, not Close (see note below)
			fyne.CurrentApp().Settings().SetTheme(theme.DefaultTheme())
			splash.parent.this.Show() // ← Show, NOT ShowAndRun
		})
	}()
}
