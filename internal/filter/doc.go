// Package filter provides include/exclude prefix-based filtering for secret
// keys resolved from Vault. It allows operators to restrict which secrets are
// injected into a child process environment, reducing the blast radius of
// over-permissioned Vault policies.
//
// Rules are evaluated in order: a key must match at least one include prefix
// (or all keys are included when the include list is empty), and must not
// match any exclude prefix.
package filter
