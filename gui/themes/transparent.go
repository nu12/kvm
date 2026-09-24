package themes

import (
	"image/color"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/theme"
)

type TransparentTheme struct{}

var _ fyne.Theme = (*TransparentTheme)(nil)

func (t *TransparentTheme) Color(name fyne.ThemeColorName, variant fyne.ThemeVariant) color.Color {
	if name == fyne.ThemeColorName("background") {
		return color.NRGBA{R: 0, G: 0, B: 0, A: 0}
	}
	return theme.DefaultTheme().Color(name, variant)
}

func (t *TransparentTheme) Icon(name fyne.ThemeIconName) fyne.Resource {
	return theme.DefaultTheme().Icon(name)
}

func (t *TransparentTheme) Font(style fyne.TextStyle) fyne.Resource {
	return theme.DefaultTheme().Font(style)
}

func (t *TransparentTheme) Size(name fyne.ThemeSizeName) float32 {
	return theme.DefaultTheme().Size(name)
}
