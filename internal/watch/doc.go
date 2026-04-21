// Package watch implements periodic polling of Vault secret paths.
//
// A Watcher re-fetches secrets on a configurable interval, compares
// the result with the previous snapshot using the diff package, and
// invokes a ChangeHandler whenever additions, removals, or value
// changes are detected.
//
// Example usage:
//
//	w := watch.New(vaultClient, paths, 30*time.Second, func(added, removed, changed map[string]string) {
//	    log.Println("secrets changed:", changed)
//	})
//	w.Watch(ctx)
package watch
