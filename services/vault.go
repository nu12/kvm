package services

import (
	"context"
	"log"

	armkeyvault "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault/v2"
	"github.com/nu12/kvm/factory"
)

func GetKeyVaults() ([]*armkeyvault.Vault, error) {
	var vaults = []*armkeyvault.Vault{}

	subClient, err := factory.GetSubscriptionClient()
	if err != nil {
		log.Printf("failed to get subscription client: %v", err)
		return []*armkeyvault.Vault{}, err
	}

	subs, err := GetAzureSubscriptions(subClient)
	if err != nil {
		log.Printf("failed to get subscriptions: %v", err)
		return []*armkeyvault.Vault{}, err
	}

	log.Printf("Loading Azure Key Vaults")
	for _, sub := range subs {

		client, err := factory.GetVaultClient(*sub.SubscriptionID)
		if err != nil {
			log.Printf("failed to create client: %v", err)
			return []*armkeyvault.Vault{}, err
		}

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
	}
	log.Printf("Done. %d KV(s) found", len(vaults))

	return vaults, nil
}
