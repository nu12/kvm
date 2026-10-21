package services

import (
	"context"
	"log"

	armsubscriptions "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions/v2"
	"github.com/nu12/kvm/factory"
)

func GetAzureSubscriptions(client factory.SubscriptionClient) ([]*armsubscriptions.Subscription, error) {

	var subs []*armsubscriptions.Subscription
	subPager := client.NewListPager(nil)
	for subPager.More() {
		subPage, err := subPager.NextPage(context.Background())
		if err != nil {
			log.Printf("failed to advance page: %v", err)
			return []*armsubscriptions.Subscription{}, err
		}
		for _, sub := range subPage.Value {
			subs = append(subs, sub)
		}
	}
	return subs, nil
}
