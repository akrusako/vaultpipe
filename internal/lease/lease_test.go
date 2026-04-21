package lease

import (
	"testing"
	"time"
)

func makeInfo(id, path string, dur time.Duration) Info {
	return Info{
		LeaseID:   id,
		Path:      path,
		Duration:  dur,
		Renewable: true,
		Acquired:  time.Now(),
	}
}

func TestAdd_AndGet(t *testing.T) {
	tr := New()
	info := makeInfo("lease-1", "secret/foo", time.Hour)
	tr.Add(info)

	got, ok := tr.Get("lease-1")
	if !ok {
		t.Fatal("expected lease to be found")
	}
	if got.Path != "secret/foo" {
		t.Errorf("path: got %q want %q", got.Path, "secret/foo")
	}
}

func TestGet_Missing(t *testing.T) {
	tr := New()
	_, ok := tr.Get("nonexistent")
	if ok {
		t.Fatal("expected miss for unknown lease ID")
	}
}

func TestRemove_DeletesLease(t *testing.T) {
	tr := New()
	tr.Add(makeInfo("lease-2", "secret/bar", time.Hour))
	tr.Remove("lease-2")
	_, ok := tr.Get("lease-2")
	if ok {
		t.Fatal("expected lease to be removed")
	}
}

func TestRemove_NoOp(t *testing.T) {
	tr := New()
	// should not panic
	tr.Remove("does-not-exist")
}

func TestExpiring_ReturnsNearExpiry(t *testing.T) {
	tr := New()
	short := Info{
		LeaseID:  "short",
		Path:     "secret/short",
		Duration: 5 * time.Second,
		Acquired: time.Now().Add(-4 * time.Second), // 1 s remaining
	}
	long := makeInfo("long", "secret/long", time.Hour)
	tr.Add(short)
	tr.Add(long)

	expiring := tr.Expiring(10 * time.Second)
	if len(expiring) != 1 {
		t.Fatalf("expected 1 expiring lease, got %d", len(expiring))
	}
	if expiring[0].LeaseID != "short" {
		t.Errorf("wrong lease: %q", expiring[0].LeaseID)
	}
}

func TestAll_ReturnsCopy(t *testing.T) {
	tr := New()
	tr.Add(makeInfo("a", "secret/a", time.Hour))
	tr.Add(makeInfo("b", "secret/b", time.Hour))

	all := tr.All()
	if len(all) != 2 {
		t.Fatalf("expected 2 leases, got %d", len(all))
	}
}

func TestAdd_SetsAcquiredIfZero(t *testing.T) {
	tr := New()
	info := Info{LeaseID: "z", Path: "secret/z", Duration: time.Hour}
	before := time.Now()
	tr.Add(info)
	after := time.Now()

	got, _ := tr.Get("z")
	if got.Acquired.Before(before) || got.Acquired.After(after) {
		t.Errorf("Acquired not set correctly: %v", got.Acquired)
	}
}

func TestTTL_DecreaseOverTime(t *testing.T) {
	info := Info{
		LeaseID:  "ttl-test",
		Duration: time.Minute,
		Acquired: time.Now().Add(-30 * time.Second),
	}
	ttl := info.TTL()
	if ttl > 31*time.Second || ttl < 29*time.Second {
		t.Errorf("unexpected TTL: %v", ttl)
	}
}
