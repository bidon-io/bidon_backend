package admin

import (
	"context"
	"testing"
)

func ptr[T any](t T) *T {
	return &t
}

func Test_userAttrsValidator_ValidateWithContext(t *testing.T) {
	tests := []struct {
		name    string
		attrs   *UserAttrs
		wantErr bool
	}{
		{
			"valid User",
			&UserAttrs{
				Email:    "example@email.com",
				IsAdmin:  ptr(true),
				Password: "password",
			},
			false,
		},
		{
			"invalid when email is incorrect",
			&UserAttrs{
				Email:    "example",
				IsAdmin:  ptr(true),
				Password: "password",
			},
			true,
		},
		{
			"invalid when password is too short",
			&UserAttrs{
				Email:    "example@email.com",
				IsAdmin:  ptr(true),
				Password: "pass",
			},
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			v := &userAttrsValidator{
				attrs: tt.attrs,
			}
			if err := v.ValidateWithContext(context.Background()); (err != nil) != tt.wantErr {
				t.Errorf("ValidateWithContext() error = %v, wantErr %v", err, tt.wantErr)
			}
		})
	}
}
