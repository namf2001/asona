package server

import (
	"log"
	"os"
	"path/filepath"

	"asona/config"
	"asona/internal/pkg/rsa"
)

const BitsGlobal = 2048

// InitRSA initializes the global RSA key pair from PEM files
func InitRSA() {
	c := config.GetConfig()

	// Find project root to resolve relative paths
	root := findProjectRoot()

	privateKeyPath := filepath.Join(root, c.RSAPrivateKeyPath)
	publicKeyPath := filepath.Join(root, c.RSAPublicKeyPath)

	globalKey := rsa.NewKeyPair()

	// Check if key files exist
	privateKeyExists := fileExists(privateKeyPath)
	publicKeyExists := fileExists(publicKeyPath)

	if privateKeyExists && publicKeyExists {
		// Load existing keys from files
		privateKeyData, err := os.ReadFile(privateKeyPath)
		if err != nil {
			log.Fatalf("[InitRSA] Failed to read private key file: %v", err)
		}

		if err := globalKey.ParsePEMPrivateKey(privateKeyData); err != nil {
			log.Fatalf("[InitRSA] Failed to parse private key: %v", err)
		}

		publicKeyData, err := os.ReadFile(publicKeyPath)
		if err != nil {
			log.Fatalf("[InitRSA] Failed to read public key file: %v", err)
		}

		if err := globalKey.ParsePEMPublicKey(publicKeyData); err != nil {
			log.Fatalf("[InitRSA] Failed to parse public key: %v", err)
		}

		log.Println("[InitRSA] RSA key pair loaded from files successfully")
	} else {
		// Generate new key pair if files don't exist
		if err := globalKey.GenerateRSAKeyPair(BitsGlobal); err != nil {
			log.Fatalf("[InitRSA] Failed to generate RSA key pair: %v", err)
		}

		privatePem, publicPem, err := globalKey.GetPemKey()
		if err != nil {
			log.Fatalf("[InitRSA] Failed to get PEM keys: %v", err)
		}

		// Ensure directory exists
		if err := os.MkdirAll(filepath.Dir(privateKeyPath), 0755); err != nil {
			log.Fatalf("[InitRSA] Failed to create key directory: %v", err)
		}

		// Save private key
		if err := os.WriteFile(privateKeyPath, privatePem, 0600); err != nil {
			log.Fatalf("[InitRSA] Failed to write private key file: %v", err)
		}

		// Save public key
		if err := os.WriteFile(publicKeyPath, publicPem, 0644); err != nil {
			log.Fatalf("[InitRSA] Failed to write public key file: %v", err)
		}

		log.Println("[InitRSA] RSA key pair generated and saved to files successfully")
	}

	rsa.GlobalRSAKeyPair = globalKey
}

// fileExists checks if a file exists
func fileExists(path string) bool {
	_, err := os.Stat(path)
	return err == nil
}

// findProjectRoot finds the project root directory
func findProjectRoot() string {
	wd, err := os.Getwd()
	if err != nil {
		return "."
	}
	cur := wd
	for {
		if _, err := os.Stat(filepath.Join(cur, "go.mod")); err == nil {
			return cur
		}
		parent := filepath.Dir(cur)
		if parent == cur {
			return wd
		}
		cur = parent
	}
}