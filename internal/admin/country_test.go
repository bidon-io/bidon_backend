package admin

import (
	"encoding/json"
	"testing"
)

func TestCountryJSONSerialization(t *testing.T) {
	country := Country{
		ID: 2,
		CountryAttrs: CountryAttrs{
			HumanName:  "Canada",
			Alpha2Code: "CA",
			Alpha3Code: "CAN",
		},
	}

	jsonData, err := json.Marshal(country)
	if err != nil {
		t.Fatalf("Error marshaling JSON: %v", err)
	}

	expectedJSON := `{"id":2,"human_name":"Canada","alpha2_code":"CA","alpha3_code":"CAN"}`
	if string(jsonData) != expectedJSON {
		t.Errorf("Expected JSON: %s, got: %s", expectedJSON, string(jsonData))
	}
}
