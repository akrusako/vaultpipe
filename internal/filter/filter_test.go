package filter_test

import (
	"testing"

	"github.com/your-org/vaultpipe/internal/filter"
)

func TestAllow_NoRules(t *testing.T) {
	f := filter.New(nil, nil)
	if !f.Allow("ANY_KEY") {
		t.Fatal("expected key to be allowed when no rules set")
	}
}

func TestAllow_IncludePrefix(t *testing.T) {
	f := filter.New([]string{"DB_"}, nil)
	if !f.Allow("DB_PASSWORD") {
		t.Fatal("expected DB_PASSWORD to be allowed")
	}
	if f.Allow("AWS_SECRET") {
		t.Fatal("expected AWS_SECRET to be excluded")
	}
}

func TestAllow_ExcludePrefix(t *testing.T) {
	f := filter.New(nil, []string{"INTERNAL_"})
	if !f.Allow("DB_HOST") {
		t.Fatal("expected DB_HOST to be allowed")
	}
	if f.Allow("INTERNAL_TOKEN") {
		t.Fatal("expected INTERNAL_TOKEN to be excluded")
	}
}

func TestAllow_IncludeAndExclude(t *testing.T) {
	f := filter.New([]string{"DB_"}, []string{"DB_INTERNAL_"})
	if !f.Allow("DB_PASSWORD") {
		t.Fatal("expected DB_PASSWORD allowed")
	}
	if f.Allow("DB_INTERNAL_SECRET") {
		t.Fatal("expected DB_INTERNAL_SECRET excluded")
	}
}

func TestAllow_CaseInsensitive(t *testing.T) {
	f := filter.New([]string{"db_"}, nil)
	if !f.Allow("DB_HOST") {
		t.Fatal("expected case-insensitive match")
	}
}

func TestApply_FiltersMap(t *testing.T) {
	f := filter.New([]string{"APP_"}, nil)
	secrets := map[string]string{
		"APP_SECRET": "s3cr3t",
		"OTHER_KEY":  "value",
	}
	out := f.Apply(secrets)
	if _, ok := out["APP_SECRET"]; !ok {
		t.Fatal("expected APP_SECRET in output")
	}
	if _, ok := out["OTHER_KEY"]; ok {
		t.Fatal("expected OTHER_KEY to be filtered out")
	}
}

func TestApply_EmptyInput(t *testing.T) {
	f := filter.New(nil, nil)
	out := f.Apply(map[string]string{})
	if len(out) != 0 {
		t.Fatalf("expected empty output, got %d keys", len(out))
	}
}
