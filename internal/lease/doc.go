// Package lease provides Vault secret lease tracking and expiry monitoring.
//
// A Tracker records active leases keyed by their Vault lease ID, allowing
// callers to query remaining TTL and enumerate leases that are approaching
// expiry.
//
// A Watcher polls the Tracker on a configurable interval and invokes a
// user-supplied callback for each lease whose TTL falls below a threshold,
// enabling proactive renewal or alerting before secrets expire.
//
// Typical usage:
//
//	tr := lease.New()
//	tr.Add(lease.Info{
//		LeaseID:   resp.LeaseID,
//		Path:      path,
//		Duration:  time.Duration(resp.LeaseDuration) * time.Second,
//		Renewable: resp.Renewable,
//	})
//
//	w := lease.NewWatcher(tr, lease.WatcherConfig{
//		Interval:   30 * time.Second,
//		Threshold:  5 * time.Minute,
//		OnExpiring: func(i lease.Info) { log.Warnf("lease expiring: %s", i.LeaseID) },
//	})
//	go w.Watch(ctx)
package lease
