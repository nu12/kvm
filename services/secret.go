package services

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/nu12/kvm/factory"
)

func ListSecrets(vaultName string) ([]*azsecrets.SecretProperties, error) {
	cred, _ := factory.GetAccountClient()
	client, err := azsecrets.NewClient(fmt.Sprintf("https://%s.vault.azure.net/", vaultName), cred, nil)
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

func GetSecret(vaultName string, secretName string, secretVersion string) (string, error) {
	cred, _ := factory.GetAccountClient()
	client, err := azsecrets.NewClient(fmt.Sprintf("https://%s.vault.azure.net/", vaultName), cred, nil)
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
