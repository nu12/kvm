package factory

import (
	"sync"
	"testing"

	"github.com/stretchr/testify/assert"
)

func resetSingleton() {
	instance = nil
	once = sync.Once{}
}

func TestGetInstance(t *testing.T) {
	resetSingleton()
	defer resetSingleton()

	f := GetInstance()
	s := GetInstance()
	assert.NotNil(t, f)
	assert.NotNil(t, s)
	assert.Equal(t, f, s)
}

func TestGetAccountClient(t *testing.T) {
	resetSingleton()
	defer resetSingleton()

	mock := &MockAccountClient{}
	GetInstance().CredentialFactory = func() (AccountClient, error) {
		return mock, nil
	}

	client, err := GetAccountClient()
	assert.NoError(t, err)
	assert.Equal(t, client, mock)
}

func TestGetSubscriptionClient(t *testing.T) {
	resetSingleton()
	defer resetSingleton()

	mock := &MockSubscriptionClient{}
	GetInstance().SubscriptionFactory = func() (SubscriptionClient, error) {
		return mock, nil
	}

	client, err := GetSubscriptionClient()
	assert.NoError(t, err)
	assert.Equal(t, client, mock)
}

func TestGetVaultClient(t *testing.T) {
	resetSingleton()
	defer resetSingleton()

	mock := &MockVaultClient{}
	GetInstance().VaultFactory = func(sub string) (VaultClient, error) {
		if sub != "test-sub" {
			t.Errorf("unexpected subscription: %s", sub)
		}
		return mock, nil
	}

	client, err := GetVaultClient("test-sub")
	assert.NoError(t, err)
	assert.Equal(t, client, mock)
}

func TestGetSecretClient(t *testing.T) {
	resetSingleton()
	defer resetSingleton()

	mock := &MockSecretClient{}
	GetInstance().SecretFactory = func(vault string) (SecretClient, error) {
		if vault != "test-vault" {
			t.Errorf("unexpected vault: %s", vault)
		}
		return mock, nil
	}

	client, err := GetSecretClient("test-vault")
	assert.NoError(t, err)
	assert.Equal(t, client, mock)
}

func TestGetCertificateClient(t *testing.T) {
	resetSingleton()
	defer resetSingleton()

	mock := &MockCertificateClient{}
	GetInstance().CertificateFactory = func(vault string) (CertificateClient, error) {
		if vault != "test-vault" {
			t.Errorf("unexpected vault: %s", vault)
		}
		return mock, nil
	}

	client, err := GetCertificateClient("test-vault")
	assert.NoError(t, err)
	assert.Equal(t, client, mock)
}
