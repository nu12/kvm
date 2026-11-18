package notifications

import (
	"testing"
	"time"

	"fyne.io/fyne/v2/test"
	"github.com/stretchr/testify/require"
)

func TestActivity(t *testing.T) {
	test.NewApp()
	activity := NewActivity()

	require.True(t, activity.Visible())
	require.True(t, activity.activity.Visible())
	require.Equal(t, "Loading...", activity.label.Text)

	activity.Stop("Done")
	require.Eventually(t, func() bool {
		return !activity.activity.Visible() && activity.label.Text == "Done"
	}, time.Second, 10*time.Millisecond)
}
