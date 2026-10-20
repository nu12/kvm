package services

import (
	"context"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions/v2"
)

type TokenCredential interface {
	GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error)
}

type CredentialFactory func() (TokenCredential, error)

type SubscriptionListPager interface {
	NewListPager(options *armsubscriptions.ClientListOptions) *runtime.Pager[armsubscriptions.ClientListResponse]
}

type SubscriptionFactory func() (SubscriptionListPager, error)

type VaultBySubscriptionPager interface {
	NewListBySubscriptionPager(options *armkeyvault.VaultsClientListBySubscriptionOptions) *runtime.Pager[armkeyvault.VaultsClientListBySubscriptionResponse]
}

type VaultFactory func(string) (VaultBySubscriptionPager, error)

type Factory struct {
	CredentialFactory   func() (TokenCredential, error)
	SubscriptionFactory func() (SubscriptionListPager, error)
	VaultFactory        func(string) (VaultBySubscriptionPager, error)
}

var instance *Factory
var once sync.Once

func AccountClient() (TokenCredential, error) {
	return GetInstance().CredentialFactory()
}

func GetInstance() *Factory {
	once.Do(func() {
		instance = &Factory{}
	})
	return instance
}
