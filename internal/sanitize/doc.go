// Package sanitize cleans and validates environment variable keys before they
// are forwarded to a child process.
//
// Vault secret engines often store keys that contain characters illegal in
// POSIX environment variable names (hyphens, dots, slashes). This package
// normalises such keys so that the exec runner can safely inject them without
// causing shell or kernel errors.
//
// Usage:
//
//	s := sanitize.New(sanitize.WithUppercase())
//	clean, errs := s.Map(rawSecrets)
package sanitize
