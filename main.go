package main

import (
	"github.com/nu12/kvm/gui/windows"
)

var Version string

func main() {
	w := windows.MainWindow{
		Logo: resourceLogoPng,
	}
	w.Init().CreateComponents().SetContent().Run()
}
