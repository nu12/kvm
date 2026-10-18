package services

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	armkeyvault "github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault/v2"
	azsecrets "github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
)

type KeyVaultService struct {
	Cred   *azidentity.DefaultAzureCredential
	Vaults []*armkeyvault.Vault
	Error  error
}

func ListKeyVaults() *KeyVaultService {
	var vaults = []*armkeyvault.Vault{}
	cred, err := GetAzureAccount()
	if err != nil {
		log.Printf("failed to get credential: %v", err)
		return &KeyVaultService{
			Error: err,
		}
	}

	ctx := context.Background()
	subs, err := GetAzureSubscriptions(cred, ctx)
	if err != nil {
		log.Printf("failed to get subscriptions: %v", err)
		return &KeyVaultService{
			Error: err,
		}
	}

	log.Printf("Loading Azure Key Vaults")
	for _, sub := range subs {
		clientFactory, err := armkeyvault.NewClientFactory(*sub.SubscriptionID, cred, nil)
		if err != nil {
			log.Printf("failed to create client: %v", err)
			return &KeyVaultService{
				Error: err,
			}
		}

		pager := clientFactory.NewVaultsClient().NewListBySubscriptionPager(nil)

		for pager.More() {
			page, err := pager.NextPage(ctx)
			if err != nil {
				log.Printf("failed to advance page: %v", err)
				return &KeyVaultService{
					Error: err,
				}
			}

			for _, vault := range page.Value {
				vaults = append(vaults, vault)
			}
		}

	}
	log.Printf("Done. %d KV(s) found", len(vaults))
	return &KeyVaultService{
		Cred:   cred,
		Vaults: vaults,
	}
}

func (kv *KeyVaultService) ListSecrets(vaultName string) ([]*azsecrets.SecretProperties, error) {
	client, err := azsecrets.NewClient(fmt.Sprintf("https://%s.vault.azure.net/", vaultName), kv.Cred, nil)
	if err != nil {
		log.Printf("failed to create client: %v", err)
		return []*azsecrets.SecretProperties{}, err
	}

	secrets := []*azsecrets.SecretProperties{}

	// List secrets
	// pager := client.NewListSecretsPager(nil)
	pager := client.NewListSecretPropertiesPager(nil)
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			log.Printf("failed to iterate to next page: %v", err)
			return []*azsecrets.SecretProperties{}, err
		}
		for _, secret := range page.Value {
			fmt.Printf("Secret ID: %s\n", *secret.ID)
			secrets = append(secrets, secret)
		}
	}
	return secrets, nil
}

func (kv *KeyVaultService) GetSecret(vaultName string, secretName string, secretVersion string) (string, error) {
	client, err := azsecrets.NewClient(fmt.Sprintf("https://%s.vault.azure.net/", vaultName), kv.Cred, nil)
	if err != nil {
		log.Printf("failed to create client: %v", err)
		return "", err
	}

	resp, err := client.GetSecret(context.TODO(), secretName, "", nil)
	if err != nil {
		log.Printf("failed to create client: %v", err)
		return "", err
	}

	return *resp.Value, nil
}
