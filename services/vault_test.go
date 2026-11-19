package services

import (
	"errors"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions/v2"
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

func TestGetKeyVaultsReturnsSubscriptionError(t *testing.T) {
	factory.GetInstance().SubscriptionFactory = func() (factory.SubscriptionClient, error) {
		return nil, errors.New("subscription client unavailable")
	}

	vaults, err := GetKeyVaults()
	require.Error(t, err)
	require.Empty(t, vaults)
}

func TestGetFromSubscriptionReturnsVaultClientError(t *testing.T) {
	factory.GetInstance().VaultFactory = func(string) (factory.VaultClient, error) {
		return nil, errors.New("vault client unavailable")
	}

	subscription := &armsubscriptions.Subscription{SubscriptionID: to.Ptr("test")}
	vaults, err := getFromSubscription(subscription)
	require.Error(t, err)
	require.Empty(t, vaults)
}
