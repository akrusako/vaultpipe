package expire_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/expire"
)

func fixedNow(t time.Time) func() time.Time {
	return func() time.Time { return t }
}

func TestSet_AndGet_BeforeExpiry(t *testing.T) {
	now := time.Now()
	tr := expire.New(fixedNow(now))
	tr.Set("key", "val", 5*time.Second)

	e, ok := tr.Get("key")
	if !ok {
		t.Fatal("expected entry to exist")
	}
	if e.Value != "val" {
		t.Fatalf("expected val, got %s", e.Value)
	}
}

func TestGet_MissingKey(t *testing.T) {
	tr := expire.New(nil)
	_, ok := tr.Get("missing")
	if ok {
		t.Fatal("expected missing key to return false")
	}
}

func TestGet_ExpiredEntry(t *testing.T) {
	base := time.Now()
	current := base
	tr := expire.New(func() time.Time { return current })

	tr.Set("token", "secret", 1*time.Second)

	// advance past TTL
	current = base.Add(2 * time.Second)

	_, ok := tr.Get("token")
	if ok {
		t.Fatal("expected expired entry to return false")
	}
}

func TestDelete_RemovesEntry(t *testing.T) {
	tr := expire.New(nil)
	tr.Set("k", "v", time.Minute)
	tr.Delete("k")
	_, ok := tr.Get("k")
	if ok {
		t.Fatal("expected deleted entry to be gone")
	}
}

func TestPurge_RemovesOnlyExpired(t *testing.T) {
	base := time.Now()
	current := base
	tr := expire.New(func() time.Time { return current })

	tr.Set("alive", "yes", 10*time.Second)
	tr.Set("dead", "no", 1*time.Second)

	current = base.Add(2 * time.Second)

	n := tr.Purge()
	if n != 1 {
		t.Fatalf("expected 1 purged, got %d", n)
	}
	if tr.Len() != 1 {
		t.Fatalf("expected 1 remaining, got %d", tr.Len())
	}
	_, ok := tr.Get("alive")
	if !ok {
		t.Fatal("expected alive entry to survive purge")
	}
}

func TestLen_CountsAll(t *testing.T) {
	tr := expire.New(nil)
	tr.Set("a", "1", time.Minute)
	tr.Set("b", "2", time.Minute)
	if tr.Len() != 2 {
		t.Fatalf("expected 2, got %d", tr.Len())
	}
}

func TestIsExpired_FalseBeforeDeadline(t *testing.T) {
	now := time.Now()
	e := expire.Entry{Value: "x", ExpiresAt: now.Add(time.Minute)}
	if e.IsExpired(now) {
		t.Fatal("expected not expired")
	}
}

func TestIsExpired_TrueAfterDeadline(t *testing.T) {
	now := time.Now()
	e := expire.Entry{Value: "x", ExpiresAt: now.Add(-time.Second)}
	if !e.IsExpired(now) {
		t.Fatal("expected expired")
	}
}
