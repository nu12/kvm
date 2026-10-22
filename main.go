package main

import (
	"fmt"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azsecrets"
	"github.com/nu12/kvm/factory"
	"github.com/nu12/kvm/gui/windows"
)

var Version string

func main() {
	cred, _ := azidentity.NewDefaultAzureCredential(nil)
	f := factory.GetInstance()
	f.CredentialFactory = func() (factory.AccountClient, error) { return cred, nil }
	f.SubscriptionFactory = func() (factory.SubscriptionClient, error) { return armsubscriptions.NewClient(cred, nil) }
	f.VaultFactory = func(sub string) (factory.VaultClient, error) {
		return armkeyvault.NewVaultsClient(sub, cred, nil)
	}
	f.SecretFactory = func(vault string) (factory.SecretClient, error) {
		return azsecrets.NewClient(fmt.Sprintf("https://%s.vault.azure.net/", vault), cred, nil)
	}

	w := windows.MainWindow{
		Logo: resourceLogoPng,
	}
	w.Init().CreateComponents().SetContent().Run()
}
