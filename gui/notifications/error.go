package notifications

import (
	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/dialog"
)

func OnErrorContinue(e error, parent fyne.Window) {
	if e != nil {
		dialog.ShowError(e, parent)
	}
}

func OnErrorCloseWindow(e error, parent fyne.Window) {
	if e != nil {
		d := dialog.NewError(e, parent)
		d.Show()
		d.SetOnClosed(func() {
			parent.Close()
		})
	}
}

func OnErrorCloseApplication(e error, parent fyne.Window) {
	if e != nil {
		d := dialog.NewError(e, parent)
		d.Show()
		d.SetOnClosed(func() {
			fyne.CurrentApp().Quit()
		})
	}
}
