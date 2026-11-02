package services

import (
	"testing"

	"github.com/nu12/kvm/factory"
	"github.com/stretchr/testify/require"
)

func TestGetKeyVaults(t *testing.T) {
	f := factory.GetInstance()
	f.SubscriptionFactory = factory.GetSubscriptionClientMock(1, nil)
	f.VaultFactory = factory.GetVaultClientMock(1, nil)

	vaults, err := GetKeyVaults()
	require.NoError(t, err)
	require.NotNil(t, vaults)
	require.Equal(t, 1, len(vaults))
	require.Equal(t, "/subscriptions/sub-0/resourceGroups/test/providers/Microsoft.KeyVault/vaults/test-0", *vaults[0].ID)
}
