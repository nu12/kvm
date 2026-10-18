package windows

import (
	"context"
	"log"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/container"
	"fyne.io/fyne/v2/widget"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"

	"golang.design/x/clipboard"
)

type VaultWindow struct {
	parent    *MainWindow
	this      fyne.Window
	VaultName string

	Components struct {
		Search *fyne.Container
		Table  *widget.Table
	}
}

type SecretTapCell struct {
	widget.Label
	CellID     widget.TableCellID
	Window     *VaultWindow
	VaultName  string
	SecretName string
}

func newSecretTapCell() *SecretTapCell {
	c := &SecretTapCell{}
	c.ExtendBaseWidget(c)
	return c
}

func (c *SecretTapCell) DoubleTapped(_ *fyne.PointEvent) {
	secretId := c.Text
	// your action here
	log.Println(secretId)

	v, _ := c.Window.parent.VaultService.GetSecret(c.VaultName, c.SecretName, "")
	log.Println(v)

	if err := clipboard.Init(); err != nil {
		log.Fatal(err)
	}

	// Copy text to clipboard
	if _, err := clipboard.Write(context.Background(), clipboard.FmtText, []byte(v), clipboard.Loops(1)); err != nil {
		log.Fatal(err)
	}
}

func (window *MainWindow) NewVaultWindow(vaultName string) {
	vaultWindow := &VaultWindow{
		parent:    window,
		this:      window.app.NewWindow("Vault: " + vaultName),
		VaultName: vaultName,
	}
	secrets := []*azsecrets.SecretProperties{}
	//w := window.app.NewWindow("Vault: " + vaultName)
	vaultWindow.this.Resize(fyne.Size{Width: 600, Height: 400})
	vaultWindow.this.Show()

	vaultWindow.Components.Table = widget.NewTableWithHeaders(
		func() (rows int, cols int) { return len(window.VaultService.Vaults), 2 },
		func() fyne.CanvasObject { return newSecretTapCell() },
		func(id widget.TableCellID, o fyne.CanvasObject) {
			b := o.(*SecretTapCell)
			b.CellID = id
			b.Window = vaultWindow
			b.VaultName = vaultName
			b.SecretName = secrets[id.Row].ID.Name()
			b.SetText(secrets[id.Row].ID.Name())
		},
	)
	vaultWindow.Components.Table.Hide()

	content := container.NewBorder(
		container.NewBorder(nil, nil, nil, widget.NewButton("Filter", func() {}), widget.NewEntry()),
		nil,
		nil,
		nil,
		vaultWindow.Components.Table,
	)
	vaultWindow.this.SetContent(content)
	secrets, _ = window.VaultService.ListSecrets(vaultName)
	vaultWindow.Components.Table.Show()

}
