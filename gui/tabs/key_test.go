package tabs

import (
	"testing"

	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/require"
)

func TestKeyTabContent(t *testing.T) {
	test.NewApp()
	tab := NewKeyTab()

	require.Equal(t, "Keys", tab.Content().Text)
	require.NotNil(t, tab.Content().Content)
}
