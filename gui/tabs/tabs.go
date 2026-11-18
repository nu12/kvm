package tabs

import (
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"github.com/nu12/kvm/gui/notifications"
)

type Tabs struct {
	Secrets      *SecretTab
	Certificates *CertificateTab
	Keys         *KeyTab
	Widget       *container.AppTabs
}

func New(vaultName string, parent fyne.Window, action *notifications.Activity) *Tabs {
	tabSet := &Tabs{
		Secrets:      NewSecretTab(vaultName, parent, action),
		Certificates: NewCertificateTab(vaultName, parent, action),
		Keys:         NewKeyTab(),
	}
	tabSet.Widget = container.NewAppTabs(
		tabSet.Secrets.Content(),
		tabSet.Certificates.Content(),
		tabSet.Keys.Content(),
	)
	tabSet.Widget.OnSelected = func(item *container.TabItem) {
		log.Printf("Selected tab %s", item.Text)
		if item == tabSet.Certificates.Content() {
			tabSet.Certificates.Load()
		}
	}
	return tabSet
}
