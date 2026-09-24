package services

import (
	"context"
	"log"

	armkeyvault "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault/v2"
)

func ListKeyVaults(filter string) ([]*armkeyvault.Vault, error) {
	var vaults = []*armkeyvault.Vault{}
	cred, err := GetAzureAccount()
	if err != nil {
		log.Printf("failed to get credential: %v", err)
		return []*armkeyvault.Vault{}, err
	}

	ctx := context.Background()
	subs, err := GetAzureSubscriptions(cred, ctx)
	if err != nil {
		log.Printf("failed to get subscriptions: %v", err)
		return []*armkeyvault.Vault{}, err
	}

	log.Printf("Loading Azure Key Vaults")
	for _, sub := range subs {
		clientFactory, err := armkeyvault.NewClientFactory(*sub.SubscriptionID, cred, nil)
		if err != nil {
			log.Printf("failed to create client: %v", err)
			return []*armkeyvault.Vault{}, err
		}

		pager := clientFactory.NewVaultsClient().NewListBySubscriptionPager(nil)

		for pager.More() {
			page, err := pager.NextPage(ctx)
			if err != nil {
				log.Printf("failed to advance page: %v", err)
				return []*armkeyvault.Vault{}, err
			}

			for _, vault := range page.Value {
				vaults = append(vaults, vault)
			}
		}

	}
	log.Printf("Done. %d KV(s) found", len(vaults))
	return vaults, nil
}
