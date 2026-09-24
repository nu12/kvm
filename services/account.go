package services

import (
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/azidentity"
)

func GetAzureAccount() (*azidentity.DefaultAzureCredential, error) {
	log.Printf("Getting account information")
	return azidentity.NewDefaultAzureCredential(nil)
}
