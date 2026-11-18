package dialogs

import (
	"testing"

	"fyne.io/fyne/v2"
	"fyne.io/fyne/v2/test"
	"fyne.io/fyne/v2/widget"
	"github.com/nu12/kvm/factory"
	"github.com/nu12/kvm/gui/notifications"
	"github.com/stretchr/testify/require"
)

func TestNewSecretUpdateDialog(t *testing.T) {
	d := NewSecretUpdateDialog(nil, "test-vault", "test-secret", nil)

	require.Equal(t, "test-vault", d.Metadata.VaultName)
	require.Equal(t, "test-secret", d.Metadata.SecretName)
	require.Nil(t, d.parent)
}

func TestSecretUpdateDialogUpdate(t *testing.T) {
	f := factory.GetInstance()
	f.SecretFactory = factory.GetSecretClientMock(1, nil)
	d := NewSecretUpdateDialog(nil, "test-vault", "test-secret", nil)

	require.NoError(t, d.Update("new-value"))
}

func TestSecretUpdateDialogShowAndConfirm(t *testing.T) {
	test.NewApp()
	factory.GetInstance().SecretFactory = factory.GetSecretClientMock(1, nil)
	parent := test.NewWindow(nil)
	activity := notifications.NewActivity()
	d := NewSecretUpdateDialog(parent, "test-vault", "test-secret", activity)

	d.Show()
	entry := findDialogEntry(t, parent)
	entry.SetText("updated value")
	button := findDialogButton(t, parent, "Update")
	test.Tap(button)
	require.Equal(t, "Secret updated", activityText(activity))
}

func TestSecretUpdateDialogShowAndCancel(t *testing.T) {
	test.NewApp()
	parent := test.NewWindow(nil)
	activity := notifications.NewActivity()
	d := NewSecretUpdateDialog(parent, "test-vault", "test-secret", activity)

	d.Show()
	test.Tap(findDialogButton(t, parent, "Cancel"))
	require.Equal(t, "Loading...", activityText(activity))
}

func activityText(activity *notifications.Activity) string {
	renderer := activity.CreateRenderer()
	defer renderer.Destroy()
	return renderer.Objects()[0].(*fyne.Container).Objects[1].(*widget.Label).Text
}

func findDialogEntry(t *testing.T, parent fyne.Window) *widget.Entry {
	t.Helper()
	for _, object := range test.LaidOutObjects(parent.Canvas().Overlays().Top()) {
		if entry, ok := object.(*widget.Entry); ok {
			return entry
		}
	}
	t.Fatal("dialog entry not found")
	return nil
}

func findDialogButton(t *testing.T, parent fyne.Window, text string) *widget.Button {
	t.Helper()
	for _, object := range test.LaidOutObjects(parent.Canvas().Overlays().Top()) {
		if button, ok := object.(*widget.Button); ok && button.Text == text {
			return button
		}
	}
	t.Fatalf("dialog button %q not found", text)
	return nil
}
