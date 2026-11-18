package services

import (
	"errors"
	"testing"

	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azcertificates"
	"github.com/nu12/kvm/factory"
	"github.com/stretchr/testify/require"
)

func TestGetCertificates(t *testing.T) {
	factory.GetInstance().CertificateFactory = factory.GetCertificateClientMock(2, nil)

	certificates, err := GetCertificates("test")
	require.NoError(t, err)
	require.Len(t, certificates, 2)
	require.Equal(t, azcertificates.ID("https://test.vault.azure.net/certificates/certificate-0"), *certificates[0].ID)
}

func TestGetCertificatesReturnsFactoryError(t *testing.T) {
	factory.GetInstance().CertificateFactory = func(string) (factory.CertificateClient, error) {
		return nil, errors.New("certificate client unavailable")
	}

	certificates, err := GetCertificates("test")
	require.Error(t, err)
	require.Empty(t, certificates)
}

func TestGetCertificate(t *testing.T) {
	factory.GetInstance().CertificateFactory = factory.GetCertificateClientMock(1, nil)

	certificate, err := GetCertificate("test", "certificate-0")
	require.NoError(t, err)
	require.Equal(t, []byte("certificate data"), certificate)
}

func TestGetCertificateRejectsEmptyData(t *testing.T) {
	factory.GetInstance().CertificateFactory = func(vaultName string) (factory.CertificateClient, error) {
		return &factory.MockCertificateClient{
			Certificates: []*azcertificates.CertificateProperties{
				{ID: to.Ptr(azcertificates.ID("https://test.vault.azure.net/certificates/test"))},
			},
			CertificateData: nil,
		}, nil
	}

	_, err := GetCertificate("test", "test")
	require.ErrorContains(t, err, "no downloadable certificate data")
}
