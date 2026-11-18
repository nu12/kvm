package components

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNewTable_BasicConfig(t *testing.T) {
	// Initialize the headless test driver (no display required)
	_ = test.NewApp()

	config := TableConfig{
		RowCount:    func() int { return 5 },
		Headers:     []string{"Name", "Age", "City"},
		NewCell:     func() fyne.CanvasObject { return widget.NewLabel("") },
		ColumnWidth: []float32{100, 80, 120},
	}

	table := NewTable(config)
	assert.NotNil(t, table)

	rows, cols := table.Length()
	assert.Equal(t, 5, rows)
	assert.Equal(t, 3, cols)
}

func TestNewTable_UpdateHeader(t *testing.T) {
	_ = test.NewApp()

	config := TableConfig{
		RowCount: func() int { return 3 },
		Headers:  []string{"ID", "Name"},
		NewCell:  func() fyne.CanvasObject { return widget.NewLabel("") },
	}

	table := NewTable(config)
	require.NotNil(t, table)

	// Assert row header
	label := widget.NewLabel("")
	table.UpdateHeader(widget.TableCellID{Row: 0, Col: 1}, label)
	assert.Equal(t, "Name", label.Text)

	// Test the row-number column (Col == -1)
	rowLabel := widget.NewLabel("")
	table.UpdateHeader(widget.TableCellID{Row: 2, Col: -1}, rowLabel)
	assert.Equal(t, "3", rowLabel.Text)
}
