package services

import (
	"context"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/nu12/kvm/factory"
)

func GetKeys(vaultName string) ([]*azkeys.KeyProperties, error) {
	client, err := factory.GetKeyClient(vaultName)
	if err != nil {
		log.Printf("failed to create key client: %v", err)
		return []*azkeys.KeyProperties{}, err
	}

	keys := []*azkeys.KeyProperties{}
	pager := client.NewListKeyPropertiesPager(nil)
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			log.Printf("failed to iterate keys: %v", err)
			return []*azkeys.KeyProperties{}, err
		}
		keys = append(keys, page.Value...)
	}
	return keys, nil
}
