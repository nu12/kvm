package services

import (
	"context"
	"log"

	armkeyvault "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault/v2"
	armsubscriptions "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions/v2"
	"github.com/nu12/kvm/factory"
)

func GetKeyVaults() ([]*armkeyvault.Vault, error) {

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
	vaults, _ := getFromSubscriptions(subs)

	log.Printf("Done. %d KV(s) found", len(vaults))

	return vaults, nil
}

func getFromSubscriptions(subs []*armsubscriptions.Subscription) ([]*armkeyvault.Vault, error) {
	var vaults = []*armkeyvault.Vault{}

	log.Printf("Loading Azure Key Vaults")
	for _, sub := range subs {
		vs, _ := getFromSubscription(sub)
		for _, v := range vs {
			vaults = append(vaults, v)
		}
	}
	return vaults, nil
}

func getFromSubscription(sub *armsubscriptions.Subscription) ([]*armkeyvault.Vault, error) {
	var vaults = []*armkeyvault.Vault{}
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
	return vaults, nil
}
