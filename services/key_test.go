package services

import (
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"errors"
	"testing"

	"crypto/x509"
	"github.com/Azure/azure-sdk-for-go/sdk/azcore/to"
	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/nu12/kvm/factory"
	"github.com/stretchr/testify/require"
)

func TestGetKeys(t *testing.T) {
	factory.GetInstance().KeyFactory = factory.GetKeyClientMock(2, nil)

	keys, err := GetKeys("test")
	require.NoError(t, err)
	require.Len(t, keys, 2)
	require.Equal(t, azkeys.ID("https://test.vault.azure.net/keys/key-0"), *keys[0].KID)
}

func TestGetKeysReturnsFactoryError(t *testing.T) {
	factory.GetInstance().KeyFactory = func(string) (factory.KeyClient, error) {
		return nil, errors.New("key client unavailable")
	}

	keys, err := GetKeys("test")
	require.Error(t, err)
	require.Empty(t, keys)
}

func TestGetPublicKeyExportsPublicRSAKeyOnly(t *testing.T) {
	factory.GetInstance().KeyFactory = func(string) (factory.KeyClient, error) {
		return &factory.MockKeyClient{
			PublicKey: &azkeys.JSONWebKey{
				Kty: to.Ptr(azkeys.KeyTypeRSA),
				N:   []byte{1, 2, 3},
				E:   []byte{1, 0, 1},
				D:   []byte("private-exponent"),
				P:   []byte("private-prime"),
				K:   []byte("symmetric-secret"),
			},
		}, nil
	}

	export, err := GetPublicKey("test", "example")
	require.NoError(t, err)

	block, rest := pem.Decode(export.Data)
	require.Empty(t, rest)
	require.NotNil(t, block)
	require.Equal(t, "RSA PUBLIC KEY", block.Type)
	require.Equal(t, ".pem", export.Extension)
	publicKey, err := x509.ParsePKCS1PublicKey(block.Bytes)
	require.NoError(t, err)
	require.Equal(t, "10203", publicKey.N.Text(16))
	require.Equal(t, 65537, publicKey.E)
}

func TestGetPublicKeyExportsPublicECJWK(t *testing.T) {
	factory.GetInstance().KeyFactory = func(string) (factory.KeyClient, error) {
		return &factory.MockKeyClient{
			PublicKey: &azkeys.JSONWebKey{
				Kty: to.Ptr(azkeys.KeyTypeEC),
				Crv: to.Ptr(azkeys.CurveNameP256),
				X:   []byte{4, 5},
				Y:   []byte{6, 7},
				D:   []byte("private-scalar"),
			},
		}, nil
	}

	export, err := GetPublicKey("test", "example")
	require.NoError(t, err)

	var fields map[string]string
	require.NoError(t, json.Unmarshal(export.Data, &fields))
	require.Equal(t, ".jwk", export.Extension)
	require.Equal(t, map[string]string{
		"kty": "EC",
		"crv": "P-256",
		"x":   base64.RawURLEncoding.EncodeToString([]byte{4, 5}),
		"y":   base64.RawURLEncoding.EncodeToString([]byte{6, 7}),
	}, fields)
}

func TestGetPublicKeyRejectsKeyWithoutPublicMaterial(t *testing.T) {
	factory.GetInstance().KeyFactory = func(string) (factory.KeyClient, error) {
		return &factory.MockKeyClient{
			PublicKey: &azkeys.JSONWebKey{
				Kty: to.Ptr(azkeys.KeyTypeOct),
				K:   []byte("symmetric-secret"),
			},
		}, nil
	}

	export, err := GetPublicKey("test", "symmetric")
	require.Error(t, err)
	require.Empty(t, export.Data)
	require.ErrorContains(t, err, "no downloadable public key")
}

func TestGetPublicKeyReturnsClientError(t *testing.T) {
	factory.GetInstance().KeyFactory = func(string) (factory.KeyClient, error) {
		return &factory.MockKeyClient{KeyError: errors.New("key retrieval failed")}, nil
	}

	export, err := GetPublicKey("test", "example")
	require.ErrorContains(t, err, "key retrieval failed")
	require.Empty(t, export.Data)
}
