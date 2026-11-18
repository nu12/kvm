package dialogs

import (
	"testing"

	"github.com/nu12/kvm/factory"
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
