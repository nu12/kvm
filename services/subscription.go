package services

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	armsubscriptions "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions/v2"
)

func GetAzureSubscriptions(cred *azidentity.DefaultAzureCredential, ctx context.Context) ([]*armsubscriptions.Subscription, error) {
	log.Printf("Getting subscriptions information")
	var subs []*armsubscriptions.Subscription

	subscriptionsClient, err := armsubscriptions.NewClient(cred, nil)
	if err != nil {
		log.Printf("failed to get subscriptions client: %v", err)
		return []*armsubscriptions.Subscription{}, err
	}

	subPager := subscriptionsClient.NewListPager(nil)
	for subPager.More() {
		subPage, err := subPager.NextPage(ctx)
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
