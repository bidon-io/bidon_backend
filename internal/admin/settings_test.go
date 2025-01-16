package admin

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"github.com/labstack/echo/v4"
	"net/http"
	"net/http/httptest"
	"testing"
)

func TestSettingsService_UpdatePassword(t *testing.T) {
	mockUserRepo := &UserRepoMock{
		UpdatePasswordFunc: func(ctx context.Context, userID int64, currentPassword string, newPassword string) error {
			if userID == 1 && currentPassword == "oldpassword" {
				return nil
			}
			if userID == 1 && currentPassword != "oldpassword" {
				return fmt.Errorf("current password is incorrect")
			}
			if userID != 1 {
				return fmt.Errorf("user not found")
			}
			return nil
		},
	}

	settingsService := &SettingsService{
		UserRepo: mockUserRepo,
	}

	tests := []struct {
		name                string
		authCtx             AuthContext
		requestPayload      PasswordUpdateRequest
		expectedHTTPStatus  int
		expectedErrorString string
	}{
		{
			name: "successful password update",
			authCtx: &AuthContextMock{
				UserIDFunc: func() int64 { return 1 },
				IsAdminFunc: func() bool {
					return false
				},
			},
			requestPayload: PasswordUpdateRequest{
				CurrentPassword:         "oldpassword",
				NewPassword:             "NewPassword1",
				NewPasswordConfirmation: "NewPassword1",
			},
			expectedHTTPStatus: http.StatusNoContent,
		},
		{
			name: "missing uppercase letter",
			authCtx: &AuthContextMock{
				UserIDFunc: func() int64 { return 1 },
				IsAdminFunc: func() bool {
					return false
				},
			},
			requestPayload: PasswordUpdateRequest{
				CurrentPassword:         "oldpassword",
				NewPassword:             "newpassword1",
				NewPasswordConfirmation: "newpassword1",
			},
			expectedHTTPStatus:  http.StatusBadRequest,
			expectedErrorString: "new_password: Password must include at least one uppercase letter.",
		},
		{
			name: "missing lowercase letter",
			authCtx: &AuthContextMock{
				UserIDFunc: func() int64 { return 1 },
				IsAdminFunc: func() bool {
					return false
				},
			},
			requestPayload: PasswordUpdateRequest{
				CurrentPassword:         "oldpassword",
				NewPassword:             "NEWPASSWORD1",
				NewPasswordConfirmation: "NEWPASSWORD1",
			},
			expectedHTTPStatus:  http.StatusBadRequest,
			expectedErrorString: "new_password: Password must include at least one lowercase letter.",
		},
		{
			name: "missing number",
			authCtx: &AuthContextMock{
				UserIDFunc: func() int64 { return 1 },
				IsAdminFunc: func() bool {
					return false
				},
			},
			requestPayload: PasswordUpdateRequest{
				CurrentPassword:         "oldpassword",
				NewPassword:             "NewPassword",
				NewPasswordConfirmation: "NewPassword",
			},
			expectedHTTPStatus:  http.StatusBadRequest,
			expectedErrorString: "new_password: Password must include at least one number.",
		},
		{
			name: "password too short",
			authCtx: &AuthContextMock{
				UserIDFunc: func() int64 { return 1 },
				IsAdminFunc: func() bool {
					return false
				},
			},
			requestPayload: PasswordUpdateRequest{
				CurrentPassword:         "oldpassword",
				NewPassword:             "Short1",
				NewPasswordConfirmation: "Short1",
			},
			expectedHTTPStatus:  http.StatusBadRequest,
			expectedErrorString: "new_password: the length must be between 8 and 50.",
		},
		{
			name: "password confirmation mismatch",
			authCtx: &AuthContextMock{
				UserIDFunc: func() int64 { return 1 },
				IsAdminFunc: func() bool {
					return false
				},
			},
			requestPayload: PasswordUpdateRequest{
				CurrentPassword:         "oldpassword",
				NewPassword:             "NewPassword1",
				NewPasswordConfirmation: "Mismatch1",
			},
			expectedHTTPStatus:  http.StatusBadRequest,
			expectedErrorString: "new_password_confirmation: New password and confirmation do not match.",
		},
		{
			name: "incorrect current password",
			authCtx: &AuthContextMock{
				UserIDFunc: func() int64 { return 1 },
				IsAdminFunc: func() bool {
					return false
				},
			},
			requestPayload: PasswordUpdateRequest{
				CurrentPassword:         "wrongpassword",
				NewPassword:             "NewPassword1",
				NewPasswordConfirmation: "NewPassword1",
			},
			expectedHTTPStatus:  http.StatusForbidden,
			expectedErrorString: "current password is incorrect",
		},
		{
			name: "user not found",
			authCtx: &AuthContextMock{
				UserIDFunc: func() int64 { return 9999 },
				IsAdminFunc: func() bool {
					return false
				},
			},
			requestPayload: PasswordUpdateRequest{
				CurrentPassword:         "oldpassword",
				NewPassword:             "NewPassword1",
				NewPasswordConfirmation: "NewPassword1",
			},
			expectedHTTPStatus:  http.StatusNotFound,
			expectedErrorString: "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			e := echo.New()
			body, _ := json.Marshal(tt.requestPayload)
			req := httptest.NewRequest(http.MethodPut, "/settings/password", bytes.NewReader(body))
			req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
			rec := httptest.NewRecorder()
			c := e.NewContext(req, rec)

			c.Set("authCtx", tt.authCtx)

			err := settingsService.UpdatePassword(c, tt.authCtx)

			if tt.expectedErrorString != "" {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
				httpErr, ok := err.(*echo.HTTPError)
				if !ok {
					t.Fatalf("expected echo.HTTPError but got %T", err)
				}
				if httpErr.Code != tt.expectedHTTPStatus {
					t.Fatalf("expected HTTP status %d but got %d", tt.expectedHTTPStatus, httpErr.Code)
				}
				if httpErr.Message != tt.expectedErrorString {
					t.Fatalf("expected error message %q but got %q", tt.expectedErrorString, httpErr.Message)
				}
			} else {
				if err != nil {
					t.Fatalf("unexpected error: %v", err)
				}
				if rec.Code != tt.expectedHTTPStatus {
					t.Fatalf("expected HTTP status %d but got %d", tt.expectedHTTPStatus, rec.Code)
				}
			}
		})
	}
}
