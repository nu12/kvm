package services

import (
	"context"
	"encoding/base64"
	"encoding/json"
	"encoding/pem"
	"fmt"
	"log"
	"math/big"

	"crypto/rsa"
	"crypto/x509"

	"github.com/Azure/azure-sdk-for-go/sdk/security/keyvault/azkeys"
	"github.com/nu12/kvm/factory"
)

type PublicKeyExport struct {
	Data      []byte
	Extension string
}

// publicJWK intentionally excludes private and symmetric JSON Web Key fields.
type publicJWK struct {
	KeyType  string `json:"kty"`
	Curve    string `json:"crv,omitempty"`
	Modulus  string `json:"n,omitempty"`
	Exponent string `json:"e,omitempty"`
	X        string `json:"x,omitempty"`
	Y        string `json:"y,omitempty"`
}

func GetPublicKey(vaultName string, keyName string) (PublicKeyExport, error) {
	client, err := factory.GetKeyClient(vaultName)
	if err != nil {
		log.Printf("failed to create key client: %v", err)
		return PublicKeyExport{}, err
	}

	response, err := client.GetKey(context.Background(), keyName, "", nil)
	if err != nil {
		log.Printf("failed to get public key: %v", err)
		return PublicKeyExport{}, err
	}

	publicKey, err := marshalPublicKey(response.Key)
	if err != nil {
		log.Printf("failed to serialize public key %q: %v", keyName, err)
		return PublicKeyExport{}, err
	}
	return publicKey, nil
}

func marshalPublicKey(key *azkeys.JSONWebKey) (PublicKeyExport, error) {
	if key == nil || key.Kty == nil {
		return PublicKeyExport{}, fmt.Errorf("key has no public key type")
	}

	jwk := publicJWK{KeyType: string(*key.Kty)}
	switch *key.Kty {
	case azkeys.KeyTypeRSA, azkeys.KeyTypeRSAHSM:
		if len(key.N) == 0 || len(key.E) == 0 {
			return PublicKeyExport{}, fmt.Errorf("RSA key has incomplete public material")
		}
		exponent := new(big.Int).SetBytes(key.E)
		maxInt := int64(^uint(0) >> 1)
		if !exponent.IsInt64() {
			return PublicKeyExport{}, fmt.Errorf("RSA key has invalid public exponent")
		}
		exponentValue := exponent.Int64()
		if exponentValue <= 0 || exponentValue > maxInt {
			return PublicKeyExport{}, fmt.Errorf("RSA key has invalid public exponent")
		}
		rsaKey := &rsa.PublicKey{N: new(big.Int).SetBytes(key.N), E: int(exponentValue)}
		if rsaKey.N.Sign() <= 0 {
			return PublicKeyExport{}, fmt.Errorf("RSA key has invalid public modulus")
		}
		der := x509.MarshalPKCS1PublicKey(rsaKey)
		return PublicKeyExport{
			Data:      pem.EncodeToMemory(&pem.Block{Type: "RSA PUBLIC KEY", Bytes: der}),
			Extension: ".pem",
		}, nil
	case azkeys.KeyTypeEC, azkeys.KeyTypeECHSM:
		if key.Crv == nil || len(key.X) == 0 || len(key.Y) == 0 {
			return PublicKeyExport{}, fmt.Errorf("elliptic curve key has incomplete public material")
		}
		jwk.Curve = string(*key.Crv)
		jwk.X = base64.RawURLEncoding.EncodeToString(key.X)
		jwk.Y = base64.RawURLEncoding.EncodeToString(key.Y)
	default:
		return PublicKeyExport{}, fmt.Errorf("key type %q has no downloadable public key", *key.Kty)
	}

	data, err := json.MarshalIndent(jwk, "", "  ")
	if err != nil {
		return PublicKeyExport{}, err
	}
	return PublicKeyExport{Data: data, Extension: ".jwk"}, nil
}
