// Package audit provides structured, append-only audit logging for
// vaultpipe operations. Each auditable action — secret reads, process
// execution, and token renewals — is emitted as a newline-delimited
// JSON record to a configurable io.Writer (e.g. stderr or a file).
//
// Audit logs intentionally contain NO secret values; they record only
// paths, commands, outcomes, and timestamps.
package audit
