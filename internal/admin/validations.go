package admin

import (
	"context"
	"fmt"

	v8n "github.com/go-ozzo/ozzo-validation/v4"
)

// isString is a validation rule that checks if of type string.
var isString = v8n.By(func(value any) error {
	_, ok := value.(string)
	if !ok {
		return fmt.Errorf("must be a string")
	}

	return nil
})

// AnyRule is a validation rule that checks if any of the given rules passes.
type AnyRule struct {
	rules []v8n.Rule
}

// Any returns a validation rule that checks if any of the given rules passes.
func Any(rules ...v8n.Rule) AnyRule {
	return AnyRule{rules}
}

// Validate checks if the given value passes any of the rules. Returns last failed validation error.
func (r AnyRule) Validate(value any) (err error) {
	for _, rule := range r.rules {
		err = v8n.Validate(value, rule)
		if err == nil {
			return
		}
	}

	return
}

// ValidateWithContext checks if the given value passes any of the rules. Returns last failed validation error.
func (r AnyRule) ValidateWithContext(ctx context.Context, value any) (err error) {
	for _, rule := range r.rules {
		err = v8n.ValidateWithContext(ctx, value, rule)
		if err == nil {
			return
		}
	}

	return
}
