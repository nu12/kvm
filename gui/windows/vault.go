package windows

import (
	"fyne.io/fyne/v2"
	"github.com/nu12/kvm/gui/notifications"
	"github.com/nu12/kvm/gui/tabs"
)

type VaultWindow struct {
	this   fyne.Window
	action *notifications.Activity

	VaultName string
	Tabs      *tabs.Tabs
}

func NewVaultWindow(vaultName string) *VaultWindow {
	return initVaultWindow(vaultName).createComponents().setContent().run()
}

func initVaultWindow(vaultName string) *VaultWindow {
	vaultWindow := &VaultWindow{
		VaultName: vaultName,
		this:      fyne.CurrentApp().NewWindow("Vault: " + vaultName),
		action:    notifications.NewActivity(),
	}
	vaultWindow.this.Resize(fyne.Size{Width: 1100, Height: 400})
	vaultWindow.this.Show()

	return vaultWindow
}

func (window *VaultWindow) createComponents() *VaultWindow {
	window.Tabs = tabs.New(window.VaultName, window.this, window.action)
	return window
}

func (window *VaultWindow) setContent() *VaultWindow {
	window.this.SetContent(window.Tabs.Widget)
	return window
}

func (window *VaultWindow) run() *VaultWindow {
	fyne.Do(func() {
		window.Tabs.Secrets.Load()
	})
	return window
}
