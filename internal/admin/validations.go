package admin

import (
	"fmt"

	v8n "github.com/go-ozzo/ozzo-validation/v4"
)

var isString = v8n.By(func(value any) error {
	_, ok := value.(string)
	if !ok {
		return fmt.Errorf("must be a string")
	}

	return nil
})
