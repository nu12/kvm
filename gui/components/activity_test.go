package components

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestActivity(t *testing.T) {
	activity := NewActivity()

	require.True(t, activity.Visible())
	require.True(t, activity.activity.Visible())
	require.Equal(t, "Loading...", activity.label.Text)

	activity.Stop("Done")
	require.False(t, activity.activity.Visible())
	require.Equal(t, "Done", activity.label.Text)
}
