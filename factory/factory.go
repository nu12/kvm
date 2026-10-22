package factory

import (
	"context"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
)

type AccountClient interface {
	GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error)
}

type SubscriptionClient interface {
	NewListPager(options *armsubscriptions.ClientListOptions) *runtime.Pager[armsubscriptions.ClientListResponse]
}

type VaultClient interface {
	NewListBySubscriptionPager(options *armkeyvault.VaultsClientListBySubscriptionOptions) *runtime.Pager[armkeyvault.VaultsClientListBySubscriptionResponse]
}

type SecretClient interface {
	NewListSecretPropertiesPager(options *azsecrets.ListSecretPropertiesOptions) *runtime.Pager[azsecrets.ListSecretPropertiesResponse]
	GetSecret(ctx context.Context, name string, version string, options *azsecrets.GetSecretOptions) (azsecrets.GetSecretResponse, error)
}

type Factory struct {
	CredentialFactory   func() (AccountClient, error)
	SubscriptionFactory func() (SubscriptionClient, error)
	VaultFactory        func(string) (VaultClient, error)
	SecretFactory       func(string) (SecretClient, error)
}

var instance *Factory
var once sync.Once

func GetAccountClient() (AccountClient, error) {
	return GetInstance().CredentialFactory()
}

func GetSubscriptionClient() (SubscriptionClient, error) {
	return GetInstance().SubscriptionFactory()
}
func GetVaultClient(sub string) (VaultClient, error) {
	return GetInstance().VaultFactory(sub)
}
func GetSecretClient(vault string) (SecretClient, error) {
	return GetInstance().SecretFactory(vault)
}

func GetInstance() *Factory {
	once.Do(func() {
		instance = &Factory{}
	})
	return instance
}
