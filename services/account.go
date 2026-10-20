package services

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/policy"
)

type TokenCredential interface {
	GetToken(ctx context.Context, opts policy.TokenRequestOptions) (azcore.AccessToken, error)
}

type CredentialFactory func() (TokenCredential, error)

func GetAccountFactory(factory CredentialFactory) (TokenCredential, error) {
	log.Printf("Getting account information")
	return factory()
}
