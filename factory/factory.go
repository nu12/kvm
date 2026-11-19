package factory

import (
	"context"
	"sync"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
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
	SetSecret(ctx context.Context, name string, parameters azsecrets.SetSecretParameters, options *azsecrets.SetSecretOptions) (azsecrets.SetSecretResponse, error)
}

type CertificateClient interface {
	NewListCertificatePropertiesPager(options *azcertificates.ListCertificatePropertiesOptions) *runtime.Pager[azcertificates.ListCertificatePropertiesResponse]
	GetCertificate(ctx context.Context, name string, version string, options *azcertificates.GetCertificateOptions) (azcertificates.GetCertificateResponse, error)
}

type KeyClient interface {
	NewListKeyPropertiesPager(options *azkeys.ListKeyPropertiesOptions) *runtime.Pager[azkeys.ListKeyPropertiesResponse]
	GetKey(ctx context.Context, name string, version string, options *azkeys.GetKeyOptions) (azkeys.GetKeyResponse, error)
}

type Factory struct {
	CredentialFactory   func() (AccountClient, error)
	SubscriptionFactory func() (SubscriptionClient, error)
	VaultFactory        func(string) (VaultClient, error)
	SecretFactory       func(string) (SecretClient, error)
	CertificateFactory  func(string) (CertificateClient, error)
	KeyFactory          func(string) (KeyClient, error)
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
func GetCertificateClient(vault string) (CertificateClient, error) {
	return GetInstance().CertificateFactory(vault)
}
func GetKeyClient(vault string) (KeyClient, error) {
	return GetInstance().KeyFactory(vault)
}

func GetInstance() *Factory {
	once.Do(func() {
		instance = &Factory{}
	})
	return instance
}
