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
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"fmt"
	"io"
)

// encryptAES encrypts plaintext using the given key with AES-GCM.
func encryptAES(key, plaintext []byte) ([]byte, error) {
	// Create a new `cipher.Block`, which implements the AES cipher
	// using the given key
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Specify the cipher mode to be GCM (Galois/Counter Mode).
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// GCM requires the use of a nonce (number used once), which is
	// a []byte of (pseudo) random values.
	nonce := make([]byte, gcm.NonceSize())
	if _, err = io.ReadFull(rand.Reader, nonce); err != nil {
		return nil, err
	}

	// We can now encrypt our plaintext. Prepends the nonce value
	// to the ciphertext.
	ciphertext := gcm.Seal(nonce, nonce, plaintext, nil)

	return ciphertext, nil
}

// decryptAES decrypts ciphertext using the given key with AES-GCM.
func decryptAES(key, ciphertext []byte) ([]byte, error) {
	// Retrieve our AES cipher
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, err
	}

	// Specify the cipher mode to be GCM (Galois/Counter Mode).
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return nil, err
	}

	// Retrieve the nonce value, which was prepended to the ciphertext,
	// and the ciphertext proper.
	nonceSize := gcm.NonceSize()

	if len(ciphertext) < nonceSize {
		return nil, fmt.Errorf("invalid input")
	}

	nonce, cipherbytes := ciphertext[:nonceSize], ciphertext[nonceSize:]

	// We can now decrypt our ciphertext!
	plaintext, err := gcm.Open(nil, nonce, cipherbytes, nil)
	if err != nil {
		return nil, err
	}

	return plaintext, nil
}

func main() {
	// 32 bytes for AES-256. Obviously, don't do this.
	key := []byte("example.key.12345678.example.key")

	plaintext := []byte("Hello, Cloud Native Go!")

	encrypted, err := encryptAES(key, plaintext)
	if err != nil {
		panic(err)
	}

	decrypted, err := decryptAES(key, encrypted)
	if err != nil {
		panic(err)
	}

	// Encode the encrypted string bites to Base64
	encoded := base64.StdEncoding.EncodeToString(encrypted)

	fmt.Println("Encrypted:", encoded)
	fmt.Println("Decrypted:", string(decrypted))
}
