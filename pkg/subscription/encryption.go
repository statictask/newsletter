package subscription

import (
	"crypto/aes"
	"crypto/cipher"
	"encoding/json"
	"golang.org/x/crypto/scrypt"
	"encoding/hex"
	"fmt"

	"github.com/statictask/newsletter/internal/config"
)

func (s *Subscription) Encrypt() (string, error) {
	// Convert the map to a JSON string
	jsonData, err := json.Marshal(s)
	if err != nil {
		return "", fmt.Errorf("failed marshaling subscription: %v", err)
	}

	// Generate a new AES key from password with
	key, err := scrypt.Key([]byte(config.C.SubscriptionAESPassword), nil, 1<<15, 8, 1, 32)
	if err != nil {
		return "", fmt.Errorf("failed generating derived key: %v", err)
	}

	// Create a new AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return "", fmt.Errorf("failed creating AES cypher block: %v", err)
	}

	// Encrypt the JSON data
	ciphertext := make([]byte, aes.BlockSize+len(jsonData))
	iv := ciphertext[:aes.BlockSize]
	cfb := cipher.NewCFBEncrypter(block, iv)
	cfb.XORKeyStream(ciphertext[aes.BlockSize:], jsonData)

	// Convert to a token that can be shared in different places
	token := hex.EncodeToString(ciphertext)

	return token, nil
}

func Decrypt(token string) (*Subscription, error) {
	// Generate the key from the password
	key, err := scrypt.Key([]byte(config.C.SubscriptionAESPassword), nil, 1<<15, 8, 1, 32)
	if err != nil {
		return nil, fmt.Errorf("failed generating derived key: %v", err)
	}

	// Create the AES cipher block
	block, err := aes.NewCipher(key)
	if err != nil {
		return nil, fmt.Errorf("failed creating AES cypher block: %v", err)
	}

	// Parse b64 token before decrypting
	ciphertext, err := hex.DecodeString(token)
	if err != nil {
		return nil, fmt.Errorf("failed decrypting subscription token: %v", err)
	}

	// Decrypt the data
	iv := ciphertext[:aes.BlockSize]
	cfb := cipher.NewCFBDecrypter(block, iv)
	plaintext := make([]byte, len(ciphertext)-aes.BlockSize)
	cfb.XORKeyStream(plaintext, ciphertext[aes.BlockSize:])

	// Unmarshal the JSON data
	var s = New() 
	err = json.Unmarshal(plaintext, &s)
	if err != nil {
		return nil, err
	}

	return s, nil
}

