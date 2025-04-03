package config_test

import (
	"context"
	"encoding/json"
	"errors"
	"net/http"
	"net/http/httptest"
	"reflect"
	"testing"

	"github.com/labstack/echo/v4"

	"github.com/bidon-io/bidon-backend/config"
)

type mockPinger struct {
	shouldFail bool
}

func (m *mockPinger) Ping(_ context.Context) error {
	if m.shouldFail {
		return errors.New("ping failed")
	}
	return nil
}

func TestUseHealthCheckHandler(t *testing.T) {
	services := config.HealthCheckParams{
		"service1": &mockPinger{shouldFail: false},
		"service2": &mockPinger{shouldFail: true},
		"service3": nil,
	}

	req := httptest.NewRequest(http.MethodGet, "/health_checks", nil)
	req.Header.Set(echo.HeaderContentType, echo.MIMEApplicationJSON)
	rec := httptest.NewRecorder()

	e := config.Echo()

	config.UseHealthCheckHandler(e, services)

	e.ServeHTTP(rec, req)

	if rec.Code != http.StatusInternalServerError {
		t.Fatalf("expected status code %d, got %d", http.StatusInternalServerError, rec.Code)
	}
	expected := `{"status":"error","service1":"ok","service2":"error"}`
	got := rec.Body.String()
	var expectedJSONAsInterface, actualJSONAsInterface interface{}
	if err := json.Unmarshal([]byte(expected), &expectedJSONAsInterface); err != nil {
		t.Fatalf("Expected value ('%s') is not valid json.\nJSON parsing error: '%s'", expected, err.Error())
	}

	if err := json.Unmarshal([]byte(got), &actualJSONAsInterface); err != nil {
		t.Fatalf("Input ('%s') needs to be valid json.\nJSON parsing error: '%s'", got, err.Error())
	}

	if !reflect.DeepEqual(expectedJSONAsInterface, actualJSONAsInterface) {
		t.Fatalf("Expected JSON response '%s' does not match actual response '%s'", expected, got)
	}
}
