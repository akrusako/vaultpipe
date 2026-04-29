// Package drain provides a graceful-shutdown Drainer that tracks
// in-flight operations and blocks until they all complete or a deadline
// is exceeded.
//
// Typical usage:
//
//	d := drain.New(10 * time.Second)
//
//	// Wrap the vault fetcher so every call is tracked.
//	trackedFetch := drain.NewMiddleware(d, vaultFetcher)
//
//	// On SIGTERM / SIGINT:
//	ok := d.Wait(shutdownCtx)
//	if !ok {
//		log.Println("drain: timed out waiting for in-flight operations")
//	}
package drain
