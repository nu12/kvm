package services

import (
	"context"
	"errors"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	armsubscriptions "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions/v2"
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

func TestGetSubscriptionFactory_Success(t *testing.T) {
	_, err := GetSubscriptionsFactory(func() (SubscriptionListPager, error) {
		return &mockPager{}, nil
	})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
}

func TestGetSubscriptionFactory_Error(t *testing.T) {
	wantErr := errors.New("auth failed")

	_, err := GetSubscriptionsFactory(func() (SubscriptionListPager, error) {
		return nil, wantErr
	})

	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
}

func TestGetSubscriptionsFactory(t *testing.T) {

	client, err := GetSubscriptionsFactory(func() (SubscriptionListPager, error) {
		return &mockPager{}, nil
	})
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
