package cache_test

import (
	"testing"
	"time"

	"github.com/yourusername/vaultpipe/internal/cache"
)

func TestGet_MissingKey(t *testing.T) {
	c := cache.New(time.Minute)
	_, ok := c.Get("secret/missing")
	if ok {
		t.Fatal("expected miss for unknown key")
	}
}

func TestSetAndGet(t *testing.T) {
	c := cache.New(time.Minute)
	secrets := map[string]string{"FOO": "bar"}
	c.Set("secret/app", secrets)

	got, ok := c.Get("secret/app")
	if !ok {
		t.Fatal("expected cache hit")
	}
	if got["FOO"] != "bar" {
		t.Errorf("expected FOO=bar, got %q", got["FOO"])
	}
}

func TestGet_Expired(t *testing.T) {
	c := cache.New(10 * time.Millisecond)
	c.Set("secret/app", map[string]string{"KEY": "val"})
	time.Sleep(20 * time.Millisecond)

	_, ok := c.Get("secret/app")
	if ok {
		t.Fatal("expected cache miss after TTL expiry")
	}
}

func TestInvalidate(t *testing.T) {
	c := cache.New(time.Minute)
	c.Set("secret/app", map[string]string{"KEY": "val"})
	c.Invalidate("secret/app")

	_, ok := c.Get("secret/app")
	if ok {
		t.Fatal("expected miss after invalidation")
	}
}

func TestFlush(t *testing.T) {
	c := cache.New(time.Minute)
	c.Set("secret/a", map[string]string{"A": "1"})
	c.Set("secret/b", map[string]string{"B": "2"})
	c.Flush()

	for _, p := range []string{"secret/a", "secret/b"} {
		if _, ok := c.Get(p); ok {
			t.Errorf("expected miss for %q after flush", p)
		}
	}
}
