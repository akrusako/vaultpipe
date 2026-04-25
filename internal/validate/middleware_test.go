package validate_test

import (
	"errors"
	"testing"

	"github.com/yourusername/vaultpipe/internal/validate"
)

func okFetcher(secrets map[string]string) validate.SecretFetcher {
	return func(_ []string) (map[string]string, error) {
		return secrets, nil
	}
}

func errFetcher(err error) validate.SecretFetcher {
	return func(_ []string) (map[string]string, error) {
		return nil, err
	}
}

func TestMiddleware_PassesThroughOnValid(t *testing.T) {
	v := validate.New(validate.NonEmpty())
	mw := validate.NewMiddleware(v, okFetcher(map[string]string{"K": "v"}))
	result, err := mw([]string{"secret/data/app"})
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if result["K"] != "v" {
		t.Errorf("expected K=v, got %q", result["K"])
	}
}

func TestMiddleware_ReturnsErrorOnInvalid(t *testing.T) {
	v := validate.New(validate.NonEmpty())
	mw := validate.NewMiddleware(v, okFetcher(map[string]string{"EMPTY": ""}))
	_, err := mw([]string{"secret/data/app"})
	if err == nil {
		t.Fatal("expected validation error, got nil")
	}
}

func TestMiddleware_PropagatesFetchError(t *testing.T) {
	v := validate.New(validate.NonEmpty())
	sentinel := errors.New("vault unreachable")
	mw := validate.NewMiddleware(v, errFetcher(sentinel))
	_, err := mw([]string{"secret/data/app"})
	if !errors.Is(err, sentinel) {
		t.Fatalf("expected sentinel error, got %v", err)
	}
}

func TestMiddleware_NoRules_AlwaysPasses(t *testing.T) {
	v := validate.New()
	mw := validate.NewMiddleware(v, okFetcher(map[string]string{"K": ""}))
	_, err := mw(nil)
	if err != nil {
		t.Fatalf("unexpected error with no rules: %v", err)
	}
}
