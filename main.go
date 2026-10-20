package main

import (
	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/keyvault/armkeyvault/v2"
	"github.com/Azure/azure-sdk-for-go/sdk/resourcemanager/resources/armsubscriptions/v2"
	"github.com/nu12/kvm/gui/windows"
	"github.com/nu12/kvm/services"
)

var Version string

func main() {
	cred, _ := azidentity.NewDefaultAzureCredential(nil)
	w := windows.MainWindow{
		Logo: resourceLogoPng,
		CredentialFactory: func() (services.TokenCredential, error) {
			return cred, nil
		},
		SubscriptionFactory: func() (services.SubscriptionListPager, error) {
			return armsubscriptions.NewClient(cred, nil)
		},
		VaultFactory: func(s string) (services.VaultBySubscriptionPager, error) {
			return armkeyvault.NewVaultsClient(s, cred, nil)
		},
	}
	w.Init().CreateComponents().SetContent().Run()
}
