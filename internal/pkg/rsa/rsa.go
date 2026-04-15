package rsa

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/hex"
	"encoding/pem"
	"fmt"
	"io"

	"asona/config"
)

var GlobalRSAKeyPair *KeyPair

type KeyPair struct {
	config     *config.EnvConfig
	PrivateKey *rsa.PrivateKey
	PublicKey  *rsa.PublicKey
}

func NewKeyPair() *KeyPair {
	return &KeyPair{
		config: config.GetConfig(),
	}
}

func (r *KeyPair) GetPemKey() ([]byte, []byte, error) {
	privBytes, err := x509.MarshalPKCS8PrivateKey(r.PrivateKey)
	if err != nil {
		return nil, nil, err
	}
	privBlock := &pem.Block{Type: "RSA PRIVATE KEY", Bytes: privBytes}
	privateKeyPEM := pem.EncodeToMemory(privBlock)

	pubBytes, err := x509.MarshalPKIXPublicKey(r.PublicKey)
	if err != nil {
		return nil, nil, err
	}
	pubBlock := &pem.Block{Type: "PUBLIC KEY", Bytes: pubBytes}
	publicKeyPEM := pem.EncodeToMemory(pubBlock)

	return privateKeyPEM, publicKeyPEM, nil
}

func (r *KeyPair) GenerateRSAKeyPair(bits int) error {
	private, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return err
	}

	public := &private.PublicKey

	r.PublicKey = public
	r.PrivateKey = private

	return nil
}

func (r *KeyPair) Decrypt(encryptedData []byte) ([]byte, error) {
	if r.config.AppEnv == "dev" {
		return encryptedData, nil
	}

	// Giải mã RSA
	decryptedData, err := rsa.DecryptOAEP(sha256.New(), rand.Reader, r.PrivateKey, encryptedData, nil)
	if err != nil {
		return nil, err
	}

	return decryptedData, nil
}

func (r *KeyPair) Encrypt(data []byte) ([]byte, error) {
	if r.config.AppEnv == "dev" {
		return data, nil
	}

	// Mã hóa dữ liệu sử dụng RSA-OAEP
	encryptedData, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, r.PublicKey, data, nil)
	if err != nil {
		return nil, err
	}

	return encryptedData, nil
}

func (r *KeyPair) EncryptToString(data []byte) (string, error) {
	if r.config.AppEnv == "dev" {
		return string(data), nil
	}

	// Mã hóa dữ liệu sử dụng RSA-OAEP
	encryptedData, err := rsa.EncryptOAEP(sha256.New(), rand.Reader, r.PublicKey, data, nil)
	if err != nil {
		return "", err
	}

	return base64.StdEncoding.EncodeToString(encryptedData), nil
}

func (r *KeyPair) ParsePEMPrivateKey(keyData []byte) error {
	block, _ := pem.Decode(keyData)
	if block == nil {
		return fmt.Errorf("failed to parse PEM block containing the key")
	}

	result, err := x509.ParsePKCS8PrivateKey(block.Bytes)
	if err != nil {
		// Fallback to PKCS1
		result, err = x509.ParsePKCS1PrivateKey(block.Bytes)
		if err != nil {
			return err
		}
	}

	switch private := result.(type) {
	case *rsa.PrivateKey:
		r.PrivateKey = private

		return nil
	default:
		return fmt.Errorf("unexpected type of public key: %T", result)
	}
}

func (r *KeyPair) ParsePEMPublicKey(keyData []byte) error {
	block, _ := pem.Decode(keyData)
	if block == nil {
		return fmt.Errorf("failed to parse PEM block containing the key")
	}

	result, err := x509.ParsePKIXPublicKey(block.Bytes)
	if err != nil {
		return err
	}

	rsaPublicKey, ok := result.(*rsa.PublicKey)
	if !ok {
		return fmt.Errorf("not an RSA public key")
	}

	r.PublicKey = rsaPublicKey

	return nil
}

func (r *KeyPair) EncryptAES(key, plaintext []byte) (string, error) {
    // Generate a new AES cipher block
    block, err := aes.NewCipher(key)
    if err != nil {
        return "", err
    }

    // Use GCM mode for the block cipher
    aesGCM, err := cipher.NewGCM(block)
    if err != nil {
        return "", err
    }

    // Create a nonce. Nonce size is recommended by the GCM standard (12 bytes)
    nonce := make([]byte, aesGCM.NonceSize())
    if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
        return "", err
    }

    // Encrypt the data using Seal
    ciphertext := aesGCM.Seal(nonce, nonce, plaintext, nil)

    // Return the result as a hex string
    return hex.EncodeToString(ciphertext), nil
}



func (r *KeyPair) DecryptAES(encryptedBase64, key string) (string, error) {
    // Decode the base64 encoded ciphertext
    ciphertext, err := base64.StdEncoding.DecodeString(encryptedBase64)
    if err != nil {
        return "", err
    }

    // Ensure the key is 32 bytes (for AES-256)
    keyBytes := []byte(key)
    if len(keyBytes) != 32 {
        return "", err
    }

    // Extract the IV from the first 16 bytes of the ciphertext
    iv := ciphertext[:aes.BlockSize]
    ciphertext = ciphertext[aes.BlockSize:]

    // Create a new AES cipher block
    block, err := aes.NewCipher(keyBytes)
    if err != nil {
        return "", err
    }

    // Create a CBC decrypter
    mode := cipher.NewCBCDecrypter(block, iv)

    // Decrypt the data
    decrypted := make([]byte, len(ciphertext))
    mode.CryptBlocks(decrypted, ciphertext)

	decrypted = pkcs5Unpadding(decrypted)

    // Convert decrypted bytes to a string and return
    return string(decrypted), nil
}



// Function to remove PKCS5 padding
func pkcs5Unpadding(data []byte) []byte {
    if len(data) == 0 {
        return nil
    }
    padding := int(data[len(data)-1])
    if padding > aes.BlockSize || padding == 0 {
        return nil
    }
    return data[:len(data)-padding]
}