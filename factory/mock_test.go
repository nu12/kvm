package factory

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestGetAccountClientMock(t *testing.T) {
	f := GetInstance()
	f.CredentialFactory = GetAccountClientMock(nil)

	cli, err := GetAccountClient()
	assert.NoError(t, err)
	assert.IsType(t, &MockAccountClient{}, cli)

	f.CredentialFactory = GetAccountClientMock(errors.New("Test"))
	_, err = GetAccountClient()
	assert.Error(t, err)
}

func TestGetToken(t *testing.T) {
	f := GetInstance()
	f.CredentialFactory = GetAccountClientMock(nil)

	cli, err := GetAccountClient()
	require.NoError(t, err)
	resp, err := cli.GetToken(context.TODO(), policy.TokenRequestOptions{})
	assert.NoError(t, err)
	assert.Equal(t, time.Time(time.Date(1, time.January, 1, 0, 0, 0, 0, time.UTC)), resp.ExpiresOn)
}

func TestGetSubscriptionClientMock(t *testing.T) {
	f := GetInstance()
	f.SubscriptionFactory = GetSubscriptionClientMock(2, nil)

	cli, err := GetSubscriptionClient()
	assert.NoError(t, err)
	assert.IsType(t, &MockSubscriptionClient{}, cli)

	f.SubscriptionFactory = GetSubscriptionClientMock(2, errors.New("Test"))
	_, err = GetSubscriptionClient()
	assert.Error(t, err)
}

func TestNewListPager(t *testing.T) {
	f := GetInstance()
	f.SubscriptionFactory = GetSubscriptionClientMock(2, nil)

	cli, err := GetSubscriptionClient()
	require.NoError(t, err)
	resp := cli.NewListPager(nil)
	assert.Equal(t, true, resp.More())
	page, err := resp.NextPage(context.TODO())
	assert.NoError(t, err)
	assert.IsType(t, []*armsubscriptions.Subscription{}, page.Value)
	assert.Equal(t, "/subscriptions/sub-0", *page.Value[0].ID)
}

func TestGetVaultClientMock(t *testing.T) {
	f := GetInstance()
	f.VaultFactory = GetVaultClientMock(2, nil)

	cli, err := GetVaultClient("t")
	assert.NoError(t, err)
	assert.IsType(t, &MockVaultClient{}, cli)

	f.VaultFactory = GetVaultClientMock(2, errors.New("Test"))
	_, err = GetVaultClient("t")
	assert.Error(t, err)
}

func TestNewListBySubscriptionPager(t *testing.T) {
	f := GetInstance()
	f.VaultFactory = GetVaultClientMock(2, nil)

	cli, err := GetVaultClient("t")
	require.NoError(t, err)
	resp := cli.NewListBySubscriptionPager(nil)
	assert.Equal(t, true, resp.More())
	page, err := resp.NextPage(context.TODO())
	assert.NoError(t, err)
	assert.IsType(t, []*armkeyvault.Vault{}, page.Value)
	assert.Equal(t, "/subscriptions/t/resourceGroups/test/providers/Microsoft.KeyVault/vaults/test-0", *page.Value[0].ID)
}

func TestGetSecretClientMock(t *testing.T) {
	f := GetInstance()
	f.SecretFactory = GetSecretClientMock(2, nil)

	cli, err := GetSecretClient("t")
	assert.NoError(t, err)
	assert.IsType(t, &MockSecretClient{}, cli)

	f.SecretFactory = GetSecretClientMock(2, errors.New("Test"))
	_, err = GetSecretClient("t")
	assert.Error(t, err)
}

func TestGetSecret(t *testing.T) {
	f := GetInstance()
	f.SecretFactory = GetSecretClientMock(2, nil)

	cli, err := GetSecretClient("vault")
	require.NoError(t, err)
	secret, err := cli.GetSecret(context.TODO(), "test", "", nil)
	assert.NoError(t, err)
	assert.Equal(t, "SecretValue", *secret.Value)
	assert.Equal(t, azsecrets.ID("/secrets/test/"), *secret.ID)
	_, err = cli.SetSecret(context.TODO(), "test", azsecrets.SetSecretParameters{
		Value: to.Ptr("value"),
	}, nil)
	assert.NoError(t, err)
}

func TestNewListSecretPropertiesPager(t *testing.T) {
	f := GetInstance()
	f.SecretFactory = GetSecretClientMock(2, nil)

	cli, err := GetSecretClient("vault")
	require.NoError(t, err)
	resp := cli.NewListSecretPropertiesPager(nil)
	assert.Equal(t, true, resp.More())
	page, err := resp.NextPage(context.TODO())
	assert.NoError(t, err)
	assert.IsType(t, []*azsecrets.SecretProperties{}, page.Value)
	assert.Equal(t, azsecrets.ID("https://vault.vault.azure.net/secrets/secret-0"), *page.SecretPropertiesListResult.Value[0].ID)
}

func TestCertificateClientMock(t *testing.T) {
	f := GetInstance()
	f.CertificateFactory = GetCertificateClientMock(2, nil)

	client, err := GetCertificateClient("vault")
	require.NoError(t, err)

	pager := client.NewListCertificatePropertiesPager(nil)
	require.True(t, pager.More())
	page, err := pager.NextPage(context.Background())
	require.NoError(t, err)
	require.Len(t, page.Value, 2)
	require.Equal(t, azcertificates.ID("https://vault.vault.azure.net/certificates/certificate-0"), *page.Value[0].ID)

	certificate, err := client.GetCertificate(context.Background(), "certificate-0", "", nil)
	require.NoError(t, err)
	require.Equal(t, []byte("certificate data"), certificate.CER)
}

func TestKeyClientMock(t *testing.T) {
	f := GetInstance()
	f.KeyFactory = GetKeyClientMock(2, nil)

	client, err := GetKeyClient("vault")
	require.NoError(t, err)
	require.IsType(t, &MockKeyClient{}, client)

	pager := client.NewListKeyPropertiesPager(nil)
	require.True(t, pager.More())
	page, err := pager.NextPage(context.Background())
	require.NoError(t, err)
	require.Len(t, page.Value, 2)
	require.Equal(t, azkeys.ID("https://vault.vault.azure.net/keys/key-0"), *page.Value[0].KID)

	key, err := client.GetKey(context.Background(), "key-0", "", nil)
	require.NoError(t, err)
	require.Equal(t, azkeys.KeyTypeRSA, *key.Key.Kty)
	require.Equal(t, []byte{1, 2, 3}, key.Key.N)
}
