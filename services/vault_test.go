package services

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	armkeyvault "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault/v2"
	armsubscriptions "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions/v2"
	"github.com/nu12/kvm/factory"
)

type mockVault struct {
	vaults []*armkeyvault.Vault
}

func (m *mockVault) NewListBySubscriptionPager(options *armkeyvault.VaultsClientListBySubscriptionOptions) *runtime.Pager[armkeyvault.VaultsClientListBySubscriptionResponse] {
	return runtime.NewPager(runtime.PagingHandler[armkeyvault.VaultsClientListBySubscriptionResponse]{
		More: func(vclbsr armkeyvault.VaultsClientListBySubscriptionResponse) bool { return false },
		Fetcher: func(ctx context.Context, vclbsr *armkeyvault.VaultsClientListBySubscriptionResponse) (armkeyvault.VaultsClientListBySubscriptionResponse, error) {
			return armkeyvault.VaultsClientListBySubscriptionResponse{
				VaultListResult: armkeyvault.VaultListResult{
					Value: m.vaults,
				},
			}, nil
		},
	})
}

func TestGetKeyVaults(t *testing.T) {
	f := factory.GetInstance()
	f.SubscriptionFactory = func() (factory.SubscriptionClient, error) {
		return &mockSubscription{
			subscriptions: []*armsubscriptions.Subscription{
				{
					SubscriptionID: to.Ptr("sub-001"),
				},
			},
		}, nil
	}
	f.VaultFactory = func(s string) (factory.VaultClient, error) {
		return &mockVault{
			vaults: []*armkeyvault.Vault{
				{
					ID: to.Ptr(""),
				},
			},
		}, nil
	}

	vaults, err := GetKeyVaults()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(vaults) != 1 {
		t.Fatalf("expected 1 key vault, got %d", len(vaults))
	}
}
