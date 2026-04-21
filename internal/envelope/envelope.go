// Package envelope provides secret envelope encryption and decryption
// for secrets in transit between Vault and the child process environment.
package envelope

import (
	"crypto/aes"
	"crypto/cipher"
	"crypto/rand"
	"encoding/base64"
	"errors"
	"fmt"
	"io"
)

// ErrInvalidKey is returned when the provided key is not a valid AES key size.
var ErrInvalidKey = errors.New("envelope: key must be 16, 24, or 32 bytes")

// ErrInvalidCiphertext is returned when the ciphertext is too short to be valid.
var ErrInvalidCiphertext = errors.New("envelope: ciphertext too short")

// Envelope encrypts and decrypts secret values using AES-GCM.
type Envelope struct {
	key []byte
}

// New creates a new Envelope with the provided AES key.
// The key must be 16, 24, or 32 bytes for AES-128, AES-192, or AES-256.
func New(key []byte) (*Envelope, error) {
	switch len(key) {
	case 16, 24, 32:
		// valid
	default:
		return nil, ErrInvalidKey
	}
	buf := make([]byte, len(key))
	copy(buf, key)
	return &Envelope{key: buf}, nil
}

// Seal encrypts plaintext and returns a base64-encoded ciphertext string.
func (e *Envelope) Seal(plaintext string) (string, error) {
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", fmt.Errorf("envelope: create cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("envelope: create gcm: %w", err)
	}
	nonce := make([]byte, gcm.NonceSize())
	if _, err := io.ReadFull(rand.Reader, nonce); err != nil {
		return "", fmt.Errorf("envelope: generate nonce: %w", err)
	}
	ciphertext := gcm.Seal(nonce, nonce, []byte(plaintext), nil)
	return base64.StdEncoding.EncodeToString(ciphertext), nil
}

// Open decrypts a base64-encoded ciphertext string and returns the plaintext.
func (e *Envelope) Open(encoded string) (string, error) {
	ciphertext, err := base64.StdEncoding.DecodeString(encoded)
	if err != nil {
		return "", fmt.Errorf("envelope: base64 decode: %w", err)
	}
	block, err := aes.NewCipher(e.key)
	if err != nil {
		return "", fmt.Errorf("envelope: create cipher: %w", err)
	}
	gcm, err := cipher.NewGCM(block)
	if err != nil {
		return "", fmt.Errorf("envelope: create gcm: %w", err)
	}
	if len(ciphertext) < gcm.NonceSize() {
		return "", ErrInvalidCiphertext
	}
	nonce, ciphertext := ciphertext[:gcm.NonceSize()], ciphertext[gcm.NonceSize():]
	plaintext, err := gcm.Open(nil, nonce, ciphertext, nil)
	if err != nil {
		return "", fmt.Errorf("envelope: decrypt: %w", err)
	}
	return string(plaintext), nil
}

// SealMap encrypts all values in a map of secrets, returning a new map.
func (e *Envelope) SealMap(secrets map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		enc, err := e.Seal(v)
		if err != nil {
			return nil, fmt.Errorf("envelope: seal key %q: %w", k, err)
		}
		out[k] = enc
	}
	return out, nil
}

// OpenMap decrypts all values in a map of encrypted secrets, returning a new map.
func (e *Envelope) OpenMap(secrets map[string]string) (map[string]string, error) {
	out := make(map[string]string, len(secrets))
	for k, v := range secrets {
		dec, err := e.Open(v)
		if err != nil {
			return nil, fmt.Errorf("envelope: open key %q: %w", k, err)
		}
		out[k] = dec
	}
	return out, nil
}
