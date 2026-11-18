package dialogs

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
	"fyne.io/fyne/v2/widget"
	"github.com/nu12/kvm/services"
)

type SecretUpdateMetadata struct {
	VaultName  string
	SecretName string
}

type SecretUpdateDialog struct {
	Metadata SecretUpdateMetadata
	parent   fyne.Window
}

func NewSecretUpdateDialog(parent fyne.Window, vaultName string, secretName string) *SecretUpdateDialog {
	return &SecretUpdateDialog{
		Metadata: SecretUpdateMetadata{
			VaultName:  vaultName,
			SecretName: secretName,
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
				if err := d.Update(valueEntry.Text); err != nil {
					dialog.ShowError(err, d.parent)
				}
			}
		},
		d.parent,
	)
	form.Show()
}
