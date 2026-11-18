package themes

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/theme"
	"github.com/stretchr/testify/require"
)

func TestTransparentTheme(t *testing.T) {
	test.NewApp()
	testTheme := &TransparentTheme{}

	background := testTheme.Color(fyne.ThemeColorName("background"), theme.VariantLight)
	_, _, _, alpha := background.RGBA()
	require.Zero(t, alpha)

	defaultColor := testTheme.Color(theme.ColorNamePrimary, theme.VariantDark)
	require.NotNil(t, defaultColor)
	require.NotNil(t, testTheme.Icon(theme.IconNameSearch))
	require.NotNil(t, testTheme.Font(fyne.TextStyle{Bold: true}))
	require.Positive(t, testTheme.Size(theme.SizeNameText))
}
