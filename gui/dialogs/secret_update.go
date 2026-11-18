package dialogs

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/nu12/kvm/gui/notifications"
	"github.com/nu12/kvm/services"
)

type SecretUpdateMetadata struct {
	VaultName    string
	SecretName   string
	Notification *notifications.Activity
}

type SecretUpdateDialog struct {
	Metadata SecretUpdateMetadata
	parent   fyne.Window
}

func NewSecretUpdateDialog(parent fyne.Window, vaultName string, secretName string, notification *notifications.Activity) *SecretUpdateDialog {
	return &SecretUpdateDialog{
		Metadata: SecretUpdateMetadata{
			VaultName:    vaultName,
			SecretName:   secretName,
			Notification: notification,
		},
		parent: parent,
	}
}

func (d *SecretUpdateDialog) Update(value string) error {
	return services.SetSecret(d.Metadata.VaultName, d.Metadata.SecretName, value)
}

func (d *SecretUpdateDialog) Show() {
	valueEntry := widget.NewEntry()
	valueEntry.Password = true

	form := dialog.NewForm(
		"Update secret value",
		"Update",
		"Cancel",
		[]*widget.FormItem{
			widget.NewFormItem("Value", valueEntry),
		},
		func(confirmed bool) {
			if confirmed {
				err := d.Update(valueEntry.Text)
				notifications.OnErrorContinue(err, d.parent)
				d.Metadata.Notification.Stop("Secret updated")
			}
		},
		d.parent,
	)
	form.Show()
}
