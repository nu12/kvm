package services

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	armkeyvault "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault/v2"
)

type VaultBySubscriptionPager interface {
	NewListBySubscriptionPager(options *armkeyvault.VaultsClientListBySubscriptionOptions) *runtime.Pager[armkeyvault.VaultsClientListBySubscriptionResponse]
}

type VaultFactory func(string) (VaultBySubscriptionPager, error)

type KeyVaultService struct {
	//Cred                TokenCredential
	CredentialFactory   func() (TokenCredential, error)
	SubscriptionFactory func() (SubscriptionListPager, error)
	VaultFactory        func(string) (VaultBySubscriptionPager, error)
	Vaults              []*armkeyvault.Vault
	Error               error
}

func GetVaultFactory(factory VaultFactory, subscriptionID string) (VaultBySubscriptionPager, error) {
	return factory(subscriptionID)
}

func GetKeyVaults(client VaultBySubscriptionPager) ([]*armkeyvault.Vault, error) {
	var vaults = []*armkeyvault.Vault{}

	pager := client.NewListBySubscriptionPager(nil)

	for pager.More() {
		page, err := pager.NextPage(context.TODO())
		if err != nil {
			return []*armkeyvault.Vault{}, err
		}

		for _, vault := range page.Value {
			vaults = append(vaults, vault)
		}
	}
	return vaults, nil
}

func (kv *KeyVaultService) ListKeyVaults() *KeyVaultService {
	var localVaults = []*armkeyvault.Vault{}

	subClient, err := GetSubscriptionsFactory(kv.SubscriptionFactory)
	if err != nil {
		log.Printf("failed to get subscription client: %v", err)
		return &KeyVaultService{
			Error: err,
		}
	}

	subs, err := GetAzureSubscriptions(subClient)
	if err != nil {
		log.Printf("failed to get subscriptions: %v", err)
		return &KeyVaultService{
			Error: err,
		}
	}

	log.Printf("Loading Azure Key Vaults")
	for _, sub := range subs {
		client, err := GetVaultFactory(kv.VaultFactory, *sub.SubscriptionID)
		//client, err := armkeyvault.NewVaultsClient(*sub.SubscriptionID, cred, nil)
		if err != nil {
			log.Printf("failed to create client: %v", err)
			return &KeyVaultService{
				Error: err,
			}
		}

		vaults, err := GetKeyVaults(client)
		if err != nil {
			return &KeyVaultService{
				Error: err,
			}
		}
		localVaults = concat(localVaults, vaults)

	}
	log.Printf("Done. %d KV(s) found", len(localVaults))
	kv.Vaults = localVaults
	return kv
}

func concat(a []*armkeyvault.Vault, b []*armkeyvault.Vault) []*armkeyvault.Vault {
	for _, el := range b {
		a = append(a, el)
	}
	return a
}
