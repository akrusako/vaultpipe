// Package env provides the Writer type for constructing process environment
// variable sets from layered sources.
//
// Typical usage:
//
//	import "os"
//
//	base := env.ParseEnvSlice(os.Environ())
//	w   := env.New(base)
//	w.Apply(secrets)          // secrets shadow OS values
//	cmd.Env = w.Build()
//
// The Writer never mutates the base slice passed at construction time,
// keeping the original process environment intact.
package env
