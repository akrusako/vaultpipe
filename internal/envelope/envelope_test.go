package envelope_test

import (
	"strings"
	"testing"

	"github.com/your-org/vaultpipe/internal/envelope"
)

func TestNew_ValidKeySizes(t *testing.T) {
	for _, size := range []int{16, 24, 32} {
		key := make([]byte, size)
		_, err := envelope.New(key)
		if err != nil {
			t.Errorf("expected no error for key size %d, got %v", size, err)
		}
	}
}

func TestNew_InvalidKeySize(t *testing.T) {
	_, err := envelope.New([]byte("short"))
	if err == nil {
		t.Fatal("expected error for invalid key size")
	}
	if err != envelope.ErrInvalidKey {
		t.Errorf("expected ErrInvalidKey, got %v", err)
	}
}

func TestSealAndOpen_RoundTrip(t *testing.T) {
	key := []byte("0123456789abcdef") // 16 bytes
	e, err := envelope.New(key)
	if err != nil {
		t.Fatalf("New: %v", err)
	}
	plaintext := "super-secret-value"
	sealed, err := e.Seal(plaintext)
	if err != nil {
		t.Fatalf("Seal: %v", err)
	}
	if sealed == plaintext {
		t.Error("sealed text should not equal plaintext")
	}
	opened, err := e.Open(sealed)
	if err != nil {
		t.Fatalf("Open: %v", err)
	}
	if opened != plaintext {
		t.Errorf("expected %q, got %q", plaintext, opened)
	}
}

func TestSeal_ProducesUniqueOutputs(t *testing.T) {
	key := []byte("0123456789abcdef")
	e, _ := envelope.New(key)
	a, _ := e.Seal("value")
	b, _ := e.Seal("value")
	if a == b {
		t.Error("expected different ciphertexts due to random nonce")
	}
}

func TestOpen_InvalidBase64(t *testing.T) {
	key := []byte("0123456789abcdef")
	e, _ := envelope.New(key)
	_, err := e.Open("not-valid-base64!!!")
	if err == nil {
		t.Fatal("expected error for invalid base64")
	}
}

func TestOpen_TooShort(t *testing.T) {
	key := []byte("0123456789abcdef")
	e, _ := envelope.New(key)
	// base64 of a very short byte slice
	_, err := e.Open("dGlueQ==")
	if err == nil {
		t.Fatal("expected error for short ciphertext")
	}
	if err != envelope.ErrInvalidCiphertext {
		t.Errorf("expected ErrInvalidCiphertext, got %v", err)
	}
}

func TestSealMap_AndOpenMap(t *testing.T) {
	key := []byte("0123456789abcdef")
	e, _ := envelope.New(key)
	secrets := map[string]string{
		"DB_PASS": "hunter2",
		"API_KEY": "abc123",
	}
	sealed, err := e.SealMap(secrets)
	if err != nil {
		t.Fatalf("SealMap: %v", err)
	}
	for k, v := range sealed {
		if strings.Contains(v, secrets[k]) {
			t.Errorf("sealed value for %q should not contain plaintext", k)
		}
	}
	opened, err := e.OpenMap(sealed)
	if err != nil {
		t.Fatalf("OpenMap: %v", err)
	}
	for k, want := range secrets {
		if got := opened[k]; got != want {
			t.Errorf("key %q: expected %q, got %q", k, want, got)
		}
	}
}
