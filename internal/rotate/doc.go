// Package rotate implements periodic secret rotation for vaultpipe.
//
// A Rotator fetches secrets from Vault on a configurable interval and
// delivers updated values to a caller-supplied Handler function, allowing
// long-running processes to receive fresh credentials without restarting.
package rotate
