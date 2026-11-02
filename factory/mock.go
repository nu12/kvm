package factory

import (
	"context"
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
)

type MockAccountClient struct{}

func (m *MockAccountClient) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	return azcore.AccessToken{}, nil
}

func GetAccountClientMock(err error) func() (AccountClient, error) {
	return func() (AccountClient, error) {
		return &MockAccountClient{}, err
	}
}

type MockSubscriptionClient struct {
	Subscriptions []*armsubscriptions.Subscription
}

func (m *MockSubscriptionClient) NewListPager(options *armsubscriptions.ClientListOptions) *runtime.Pager[armsubscriptions.ClientListResponse] {
	return runtime.NewPager(
		runtime.PagingHandler[armsubscriptions.ClientListResponse]{
			More: func(r armsubscriptions.ClientListResponse) bool {
				return false // une seule page
			},
			Fetcher: func(ctx context.Context, curPage *armsubscriptions.ClientListResponse) (armsubscriptions.ClientListResponse, error) {
				return armsubscriptions.ClientListResponse{
					SubscriptionListResult: armsubscriptions.SubscriptionListResult{
						Value: m.Subscriptions,
					},
				}, nil
			},
		},
	)
}

func GetSubscriptionClientMock(nbSubscriptions int, err error) func() (SubscriptionClient, error) {
	var subs []*armsubscriptions.Subscription

	for i := range nbSubscriptions {
		subs = append(subs, &armsubscriptions.Subscription{
			ID:             to.Ptr(fmt.Sprintf("/subscriptions/sub-%d", i)),
			SubscriptionID: to.Ptr(fmt.Sprintf("sub-%d", i)),
			DisplayName:    to.Ptr(fmt.Sprintf("Sub %d", i)),
			TenantID:       to.Ptr("Tenant Test"),
			State:          to.Ptr(armsubscriptions.SubscriptionStateEnabled),
		})
	}
	return func() (SubscriptionClient, error) {
		return &MockSubscriptionClient{
			Subscriptions: subs,
		}, err
	}
}

type MockVaultClient struct {
	Vaults []*armkeyvault.Vault
}

func (m *MockVaultClient) NewListBySubscriptionPager(options *armkeyvault.VaultsClientListBySubscriptionOptions) *runtime.Pager[armkeyvault.VaultsClientListBySubscriptionResponse] {
	return runtime.NewPager(runtime.PagingHandler[armkeyvault.VaultsClientListBySubscriptionResponse]{
		More: func(vclbsr armkeyvault.VaultsClientListBySubscriptionResponse) bool { return false },
		Fetcher: func(ctx context.Context, vclbsr *armkeyvault.VaultsClientListBySubscriptionResponse) (armkeyvault.VaultsClientListBySubscriptionResponse, error) {
			return armkeyvault.VaultsClientListBySubscriptionResponse{
				VaultListResult: armkeyvault.VaultListResult{
					Value: m.Vaults,
				},
			}, nil
		},
	})
}

func GetVaultClientMock(nbVaults int, err error) func(string) (VaultClient, error) {
	return func(sub string) (VaultClient, error) {
		var vaults []*armkeyvault.Vault
		for i := range nbVaults {
			vaults = append(vaults, &armkeyvault.Vault{
				ID: to.Ptr(fmt.Sprintf("/subscriptions/%s/resourceGroups/test/providers/Microsoft.KeyVault/vaults/test-%d", sub, i)),
			})
		}
		return &MockVaultClient{
			Vaults: vaults,
		}, err
	}
}

type MockSecretClient struct {
	Secrets   []*azsecrets.SecretProperties
	vaultName string
}

func (m *MockSecretClient) GetSecret(ctx context.Context, name string, version string, options *azsecrets.GetSecretOptions) (azsecrets.GetSecretResponse, error) {
	return azsecrets.GetSecretResponse{
		Secret: azsecrets.Secret{
			ID:    to.Ptr(azsecrets.ID(fmt.Sprintf("%s/secrets/%s/%s", m.vaultName, name, version))),
			Value: to.Ptr("SecretValue"),
		},
	}, nil
}

func (m *MockSecretClient) NewListSecretPropertiesPager(options *azsecrets.ListSecretPropertiesOptions) *runtime.Pager[azsecrets.ListSecretPropertiesResponse] {
	return runtime.NewPager(runtime.PagingHandler[azsecrets.ListSecretPropertiesResponse]{
		More: func(lspr azsecrets.ListSecretPropertiesResponse) bool { return false },
		Fetcher: func(ctx context.Context, lspr *azsecrets.ListSecretPropertiesResponse) (azsecrets.ListSecretPropertiesResponse, error) {
			return azsecrets.ListSecretPropertiesResponse{
				SecretPropertiesListResult: azsecrets.SecretPropertiesListResult{
					Value: m.Secrets,
				},
			}, nil
		},
	})
}

func GetSecretClientMock(nbSecrets int, err error) func(string) (SecretClient, error) {
	return func(vaultName string) (SecretClient, error) {
		var secrets []*azsecrets.SecretProperties
		for i := range nbSecrets {
			secrets = append(secrets, &azsecrets.SecretProperties{
				ID: to.Ptr(azsecrets.ID(fmt.Sprintf("https://%s.vault.azure.net/secrets/secret-%d", vaultName, i))),
			})
		}
		return &MockSecretClient{
			Secrets: secrets,
		}, err
	}
}
