package services

import (
	"testing"

	"github.com/nu12/kvm/factory"
	"github.com/stretchr/testify/require"
)

func TestGetSubscriptions(t *testing.T) {
	f := factory.GetInstance()
	f.SubscriptionFactory = factory.GetSubscriptionClientMock(2, nil)

	client, err := factory.GetSubscriptionClient()
	require.NoError(t, err)

	subs, err := GetAzureSubscriptions(client)
	require.NoError(t, err)
	require.Equal(t, 2, len(subs))
}
