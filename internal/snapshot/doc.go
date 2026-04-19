// Package snapshot provides point-in-time capture of secret key/value maps
// and change detection between captures. It is used by the rotation pipeline
// to decide whether a child process should be restarted after a secret refresh.
package snapshot
