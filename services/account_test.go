package services

import (
	"context"
	"errors"
	"testing"
	"time"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

type mockCredential struct {
	token  string
	err    error
	called bool
}

func (m *mockCredential) GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error) {
	m.called = true
	return azcore.AccessToken{Token: m.token, ExpiresOn: time.Now().Add(time.Hour)}, m.err
}

func TestGetAccountFactory_Success(t *testing.T) {
	mock := &mockCredential{token: "fake-token"}

	cred, err := GetAccountFactory(func() (TokenCredential, error) {
		return mock, nil
	})

	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if cred == nil {
		t.Fatal("expected non-nil credential")
	}
}

func TestGetAccountFactory_Error(t *testing.T) {
	wantErr := errors.New("auth failed")

	_, err := GetAccountFactory(func() (TokenCredential, error) {
		return nil, wantErr
	})

	if !errors.Is(err, wantErr) {
		t.Fatalf("expected %v, got %v", wantErr, err)
	}
}
