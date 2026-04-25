// Package validate provides secret value validation before injection into
// a child process environment. Rules are evaluated in registration order;
// the first failing rule short-circuits evaluation for that key.
package validate

import (
	"fmt"
	"regexp"
)

// Rule is a single validation predicate applied to a secret value.
type Rule struct {
	Name    string
	CheckFn func(value string) error
}

// Validator holds an ordered set of rules and applies them to secret maps.
type Validator struct {
	rules []Rule
}

// New returns a Validator with the provided rules pre-registered.
func New(rules ...Rule) *Validator {
	v := &Validator{}
	for _, r := range rules {
		v.rules = append(v.rules, r)
	}
	return v
}

// Add registers an additional rule.
func (v *Validator) Add(r Rule) {
	v.rules = append(v.rules, r)
}

// Check validates every value in secrets against all registered rules.
// It returns a map of key → first error encountered, or nil when all pass.
func (v *Validator) Check(secrets map[string]string) map[string]error {
	if len(v.rules) == 0 {
		return nil
	}
	errs := make(map[string]error)
	for k, val := range secrets {
		for _, r := range v.rules {
			if err := r.CheckFn(val); err != nil {
				errs[k] = fmt.Errorf("rule %q: %w", r.Name, err)
				break
			}
		}
	}
	if len(errs) == 0 {
		return nil
	}
	return errs
}

// NonEmpty rejects blank secret values.
func NonEmpty() Rule {
	return Rule{
		Name: "non-empty",
		CheckFn: func(v string) error {
			if v == "" {
				return fmt.Errorf("value must not be empty")
			}
			return nil
		},
	}
}

// MatchesPattern rejects values that do not satisfy the supplied regexp.
func MatchesPattern(name, pattern string) Rule {
	re := regexp.MustCompile(pattern)
	return Rule{
		Name: name,
		CheckFn: func(v string) error {
			if !re.MatchString(v) {
				return fmt.Errorf("value does not match pattern %q", pattern)
			}
			return nil
		},
	}
}

// MaxLength rejects values longer than n bytes.
func MaxLength(n int) Rule {
	return Rule{
		Name: "max-length",
		CheckFn: func(v string) error {
			if len(v) > n {
				return fmt.Errorf("value length %d exceeds maximum %d", len(v), n)
			}
			return nil
		},
	}
}
