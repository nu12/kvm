package windows

import (
	"fmt"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/canvas"
	"fyne.io/fyne/v2/driver/desktop"
	"fyne.io/fyne/v2/theme"
	"fyne.io/fyne/v2/widget"
	"github.com/go-gl/glfw/v3.4/glfw"
)

type SplashWindow struct {
	app    fyne.App
	parent *MainWindow
	this   fyne.Window
}

func (window *MainWindow) NewSplash() (*SplashWindow, error) {
	// 2. Enable transparent framebuffer via GLFW
	if glfw.Init() == nil {
		glfw.WindowHint(glfw.TransparentFramebuffer, glfw.True)
	}

	var splash fyne.Window
	if drv, ok := window.app.Driver().(desktop.Driver); ok {
		splash = drv.CreateSplashWindow()
		img := canvas.NewImageFromResource(window.Logo)
		img.FillMode = canvas.ImageFillOriginal
		splash.SetContent(img)
		splash.Show()
		if glfw.Init() == nil {
			glfw.DefaultWindowHints()
		}
	}
	return &SplashWindow{
		app:    window.app,
		parent: window,
		this:   splash,
	}, nil
}

func (splash *SplashWindow) LoadAsync(table *widget.Table, counter *widget.Label) error {
	go func() {
		splash.parent.VaultService = splash.parent.VaultService.ListKeyVaults()
		//table.Show()
		table.Refresh()
		counter.SetText(fmt.Sprintf("Showing %d out of %d key vaults.", len(splash.parent.VaultService.Vaults), len(splash.parent.VaultService.Vaults)))
		fyne.Do(func() {
			splash.this.Close() // ← Hide, not Close (see note below)

			splash.app.Settings().SetTheme(theme.DefaultTheme())

			splash.parent.this.Show() // ← Show, NOT ShowAndRun
		})
	}()
	return nil
}
