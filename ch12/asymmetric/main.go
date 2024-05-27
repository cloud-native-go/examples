/*
 * Copyright 2024 Matthew A. Titmus
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *      http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package main

import (
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha256"
	"crypto/x509"
	"encoding/base64"
	"encoding/pem"
	"fmt"
)

// generateKeyPair generates an RSA key pair.
func generateKeyPair(bits int) (*rsa.PrivateKey, *rsa.PublicKey, error) {
	privateKey, err := rsa.GenerateKey(rand.Reader, bits)
	if err != nil {
		return nil, nil, err
	}

	return privateKey, &privateKey.PublicKey, nil
}

// exportKeys exports keys to PEM format for demonstration purposes
func exportKeys(privateKey *rsa.PrivateKey, publicKey *rsa.PublicKey) {
	privBytes := x509.MarshalPKCS1PrivateKey(privateKey)
	privPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PRIVATE KEY",
			Bytes: privBytes,
		})

	pubBytes, _ := x509.MarshalPKIXPublicKey(publicKey)
	pubPEM := pem.EncodeToMemory(
		&pem.Block{
			Type:  "RSA PUBLIC KEY",
			Bytes: pubBytes,
		})

	fmt.Println(string(privPEM))
	fmt.Println(string(pubPEM))
}

// encryptRSA encrypts the given message with the RSA public key.
func encryptRSA(publicKey *rsa.PublicKey, message []byte) ([]byte, error) {
	ciphertext, err := rsa.EncryptOAEP(
		sha256.New(), rand.Reader, publicKey, message, nil)
	if err != nil {
		return nil, err
	}

	return ciphertext, nil
}

// decryptRSA decrypts the given ciphertext with the RSA private key.
func decryptRSA(privateKey *rsa.PrivateKey, ciphertext []byte) ([]byte, error) {
	plaintext, err := rsa.DecryptOAEP(
		sha256.New(), rand.Reader, privateKey, ciphertext, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func main() {
	// Generate 2048-bit RSA keys
	privateKey, publicKey, err := generateKeyPair(2048)
	if err != nil {
		panic(err)
	}

	exportKeys(privateKey, publicKey)

	plaintext := []byte("Hello, Cloud Native Go!")

	// Encrypt message
	encrypted, err := encryptRSA(publicKey, plaintext)
	if err != nil {
		panic(err)
	}

	// Decrypt message
	decrypted, err := decryptRSA(privateKey, encrypted)
	if err != nil {
		panic(err)
	}

	data := base64.StdEncoding.EncodeToString(encrypted)
	fmt.Println("Encrypted:", data)
	fmt.Println("Decrypted:", string(decrypted))
}
