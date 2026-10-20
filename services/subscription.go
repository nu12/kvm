package services

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	armsubscriptions "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions/v2"
)

type SubscriptionListPager interface {
	NewListPager(options *armsubscriptions.ClientListOptions) *runtime.Pager[armsubscriptions.ClientListResponse]
}

type SubscriptionFactory func() (SubscriptionListPager, error)

func GetSubscriptionsFactory(factory SubscriptionFactory) (SubscriptionListPager, error) {
	log.Printf("Getting subscriptions information")
	return factory()
}

func GetAzureSubscriptions(client SubscriptionListPager) ([]*armsubscriptions.Subscription, error) {

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
