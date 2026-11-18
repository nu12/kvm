package services

import (
	"context"
	"fmt"
	"log"

	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/nu12/kvm/factory"
)

func GetCertificates(vaultName string) ([]*azcertificates.CertificateProperties, error) {
	client, err := factory.GetCertificateClient(vaultName)
	if err != nil {
		log.Printf("failed to create certificate client: %v", err)
		return []*azcertificates.CertificateProperties{}, err
	}

	certificates := []*azcertificates.CertificateProperties{}
	pager := client.NewListCertificatePropertiesPager(nil)
	for pager.More() {
		page, err := pager.NextPage(context.Background())
		if err != nil {
			log.Printf("failed to iterate certificates: %v", err)
			return []*azcertificates.CertificateProperties{}, err
		}
		certificates = append(certificates, page.Value...)
	}
	return certificates, nil
}

func GetCertificate(vaultName string, certificateName string) ([]byte, error) {
	client, err := factory.GetCertificateClient(vaultName)
	if err != nil {
		log.Printf("failed to create certificate client: %v", err)
		return nil, err
	}

	response, err := client.GetCertificate(context.Background(), certificateName, "", nil)
	if err != nil {
		log.Printf("failed to get certificate: %v", err)
		return nil, err
	}
	if len(response.CER) == 0 {
		err = fmt.Errorf("certificate %q has no downloadable certificate data", certificateName)
		log.Printf("%v", err)
		return nil, err
	}
	return response.CER, nil
}
