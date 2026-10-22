package services

import (
	"context"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/runtime"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/nu12/kvm/factory"
)

type mockSecret struct {
	secrets []*azsecrets.SecretProperties
	secret  azsecrets.GetSecretResponse
}

func (m *mockSecret) GetSecret(ctx context.Context, name string, version string, options *azsecrets.GetSecretOptions) (azsecrets.GetSecretResponse, error) {
	return m.secret, nil
}

func (m *mockSecret) NewListSecretPropertiesPager(options *azsecrets.ListSecretPropertiesOptions) *runtime.Pager[azsecrets.ListSecretPropertiesResponse] {
	return runtime.NewPager(runtime.PagingHandler[azsecrets.ListSecretPropertiesResponse]{
		More: func(lspr azsecrets.ListSecretPropertiesResponse) bool { return false },
		Fetcher: func(ctx context.Context, lspr *azsecrets.ListSecretPropertiesResponse) (azsecrets.ListSecretPropertiesResponse, error) {
			return azsecrets.ListSecretPropertiesResponse{
				SecretPropertiesListResult: azsecrets.SecretPropertiesListResult{
					Value: m.secrets,
				},
			}, nil
		},
	})
}

func TestGetSecret(t *testing.T) {
	expectedSecretValue := "Test"
	f := factory.GetInstance()
	f.SecretFactory = func(s string) (factory.SecretClient, error) {
		return &mockSecret{
			secret: azsecrets.GetSecretResponse{
				Secret: azsecrets.Secret{
					Value: to.Ptr(expectedSecretValue),
				},
			},
		}, nil
	}

	secret, err := GetSecret("test", "test", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if secret != expectedSecretValue {
		t.Fatalf("expected value to be %s, got %s", expectedSecretValue, secret)
	}
}

func TestGetSecrets(t *testing.T) {
	f := factory.GetInstance()
	f.SecretFactory = func(s string) (factory.SecretClient, error) {
		return &mockSecret{
			secrets: []*azsecrets.SecretProperties{
				{
					ID: (*azsecrets.ID)(&s),
				},
			},
		}, nil
	}

	secrets, err := GetSecrets("test")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(secrets) != 1 {
		t.Fatalf("expected 1 secret, got %d", len(secrets))
	}
}
