// Package mask implements secret-value redaction for vaultpipe.
//
// When vaultpipe surfaces errors or diagnostic messages, any secret values
// fetched from Vault must not appear in plain text. A Masker is constructed
// from the resolved secrets map and can be used to sanitise any string before
// it is written to stderr or a log sink.
//
// Usage:
//
//	m := mask.New(secrets)
//	fmt.Fprintln(os.Stderr, m.Redact(someMessage))
package mask
