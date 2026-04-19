// Package template provides Go text/template-based rendering for Vault secret
// paths. This allows vaultpipe users to define dynamic secret paths in
// configuration using environment variables as interpolation context.
//
// Example path template:
//
//	"secret/{{.ENVIRONMENT}}/database"
//
// The renderer is seeded with the current process environment at construction
// time, so any variables set before vaultpipe starts are available for use
// in path templates.
package template
