package services

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	armsubscriptions "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions/v2"
	"github.com/nu12/kvm/factory"
)

type mockPager struct {
	err error
}

func (m *mockPager) NewListPager(options *armsubscriptions.ClientListOptions) *runtime.Pager[armsubscriptions.ClientListResponse] {
	return runtime.NewPager(
		runtime.PagingHandler[armsubscriptions.ClientListResponse]{
			More: func(r armsubscriptions.ClientListResponse) bool {
				return false // une seule page
			},
			Fetcher: func(ctx context.Context, curPage *armsubscriptions.ClientListResponse) (armsubscriptions.ClientListResponse, error) {
				return armsubscriptions.ClientListResponse{
					SubscriptionListResult: armsubscriptions.SubscriptionListResult{
						Value: []*armsubscriptions.Subscription{
							{
								ID:             to.Ptr("/subscriptions/sub-001"),
								SubscriptionID: to.Ptr("sub-001"),
								DisplayName:    to.Ptr("Test Sub"),
								TenantID:       to.Ptr("tenant-001"),
								State:          to.Ptr(armsubscriptions.SubscriptionStateEnabled),
							},
							{
								ID:             to.Ptr("/subscriptions/sub-002"),
								SubscriptionID: to.Ptr("sub-002"),
								DisplayName:    to.Ptr("Prod Sub"),
								TenantID:       to.Ptr("tenant-001"),
								State:          to.Ptr(armsubscriptions.SubscriptionStateEnabled),
							},
						},
					},
				}, nil
			},
		},
	)

}

func TestGetSubscriptions(t *testing.T) {
	f := factory.GetInstance()
	f.SubscriptionFactory = func() (factory.SubscriptionClient, error) { return &mockPager{}, nil }

	client, err := factory.GetSubscriptionClient()
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	subs, err := GetAzureSubscriptions(client)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	if len(subs) != 2 {
		t.Fatalf("expected 2 subscriptions")
	}
}
