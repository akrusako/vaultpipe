// Package circuit provides a simple circuit breaker for protecting vaultpipe
// against cascading failures when upstream services such as Vault become
// temporarily unavailable.
//
// The breaker transitions through three states:
//
//   - Closed: normal operation; all calls are allowed through.
//   - Open: the failure threshold has been exceeded; calls are blocked and
//     ErrOpen is returned immediately.
//   - Half-Open: the reset timeout has elapsed; a single probe call is allowed
//     through to test whether the service has recovered.
//
// Usage:
//
//	b, err := circuit.New(circuit.DefaultConfig())
//	if !b.Allow() {
//	    return circuit.ErrOpen
//	}
//	if err := doCall(); err != nil {
//	    b.RecordFailure()
//	} else {
//	    b.RecordSuccess()
//	}
package circuit
