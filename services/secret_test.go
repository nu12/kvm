package services

import (
	"errors"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/nu12/kvm/factory"
	"github.com/stretchr/testify/require"
)

func TestGetSecret(t *testing.T) {
	f := factory.GetInstance()
	f.SecretFactory = factory.GetSecretClientMock(1, nil)

	secret, err := GetSecret("test", "test", "")
	require.NoError(t, err)
	require.Equal(t, "SecretValue", secret)
}

func TestGetSecrets(t *testing.T) {
	f := factory.GetInstance()
	f.SecretFactory = factory.GetSecretClientMock(2, nil)

	secrets, err := GetSecrets("test")
	require.NoError(t, err)
	require.Equal(t, 2, len(secrets))
	require.Equal(t, azsecrets.ID("https://test.vault.azure.net/secrets/secret-0"), *secrets[0].ID)
}

func TestSetSecret(t *testing.T) {
	f := factory.GetInstance()
	f.SecretFactory = factory.GetSecretClientMock(1, nil)

	err := SetSecret("test", "test", "new-value")
	require.NoError(t, err)
}

func TestSecretServicesReturnFactoryErrors(t *testing.T) {
	factory.GetInstance().SecretFactory = func(string) (factory.SecretClient, error) {
		return nil, errors.New("secret client unavailable")
	}

	value, err := GetSecret("test", "secret", "")
	require.Error(t, err)
	require.Empty(t, value)
	require.Error(t, SetSecret("test", "secret", "value"))
}
