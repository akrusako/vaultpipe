package validate_test

import (
	"testing"

	"github.com/yourusername/vaultpipe/internal/validate"
)

func TestCheck_NoRules_ReturnsNil(t *testing.T) {
	v := validate.New()
	result := v.Check(map[string]string{"KEY": ""})
	if result != nil {
		t.Fatalf("expected nil, got %v", result)
	}
}

func TestCheck_AllPass(t *testing.T) {
	v := validate.New(validate.NonEmpty())
	result := v.Check(map[string]string{"A": "hello", "B": "world"})
	if result != nil {
		t.Fatalf("expected nil errors, got %v", result)
	}
}

func TestCheck_NonEmpty_Fails(t *testing.T) {
	v := validate.New(validate.NonEmpty())
	result := v.Check(map[string]string{"GOOD": "ok", "BAD": ""})
	if result == nil {
		t.Fatal("expected errors, got nil")
	}
	if _, ok := result["BAD"]; !ok {
		t.Error("expected error for key BAD")
	}
	if _, ok := result["GOOD"]; ok {
		t.Error("unexpected error for key GOOD")
	}
}

func TestCheck_MatchesPattern_Fails(t *testing.T) {
	v := validate.New(validate.MatchesPattern("uuid", `^[0-9a-f-]{36}$`))
	result := v.Check(map[string]string{"ID": "not-a-uuid"})
	if result == nil {
		t.Fatal("expected validation error")
	}
	if _, ok := result["ID"]; !ok {
		t.Error("expected error for key ID")
	}
}

func TestCheck_MatchesPattern_Passes(t *testing.T) {
	v := validate.New(validate.MatchesPattern("uuid", `^[0-9a-f-]{36}$`))
	result := v.Check(map[string]string{"ID": "123e4567-e89b-12d3-a456-426614174000"})
	if result != nil {
		t.Fatalf("unexpected errors: %v", result)
	}
}

func TestCheck_MaxLength_Fails(t *testing.T) {
	v := validate.New(validate.MaxLength(5))
	result := v.Check(map[string]string{"K": "toolongvalue"})
	if result == nil {
		t.Fatal("expected error for oversized value")
	}
}

func TestCheck_MaxLength_Passes(t *testing.T) {
	v := validate.New(validate.MaxLength(20))
	result := v.Check(map[string]string{"K": "short"})
	if result != nil {
		t.Fatalf("unexpected error: %v", result)
	}
}

func TestCheck_FirstRuleShortCircuits(t *testing.T) {
	called := false
	secondRule := validate.Rule{
		Name: "never-called",
		CheckFn: func(v string) error {
			called = true
			return nil
		},
	}
	v := validate.New(validate.NonEmpty(), secondRule)
	v.Check(map[string]string{"K": ""})
	if called {
		t.Error("second rule should not be called after first fails")
	}
}

func TestAdd_RegistersRule(t *testing.T) {
	v := validate.New()
	v.Add(validate.NonEmpty())
	result := v.Check(map[string]string{"K": ""})
	if result == nil {
		t.Fatal("expected error after adding NonEmpty rule")
	}
}
