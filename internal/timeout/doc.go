// Package timeout provides deadline enforcement for vaultpipe operations.
//
// An Enforcer wraps a Config that specifies independent timeouts for secret
// fetch calls and child-process execution. Callers obtain a context with the
// appropriate deadline via ExecContext or FetchContext and pass it to the
// relevant subsystem. A zero duration disables the timeout for that operation.
//
// Example:
//
//	cfg := timeout.DefaultConfig()
//	cfg.Exec = 30 * time.Second
//	e, err := timeout.New(cfg)
//	if err != nil { ... }
//
//	ctx, cancel := e.ExecContext(parent)
//	defer cancel()
//	// pass ctx to exec.Runner
package timeout
