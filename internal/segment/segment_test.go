package segment

import (
	"testing"
)

func TestMatchCountry(t *testing.T) {
	// Test case 1: IN operator, country in values
	filter := Filter{
		Operator: "IN",
		Values:   []string{"US", "CA"},
	}
	country := "US"
	expected := true
	result := matchCountry(filter, country)
	if result != expected {
		t.Errorf("matchCountry returned unexpected result. Expected: %v, Got: %v", expected, result)
	}

	// Test case 2: IN operator, country not in values
	filter = Filter{
		Operator: "IN",
		Values:   []string{"US", "CA"},
	}
	country = "FR"
	expected = false
	result = matchCountry(filter, country)
	if result != expected {
		t.Errorf("matchCountry returned unexpected result. Expected: %v, Got: %v", expected, result)
	}

	// Test case 3: NOT IN operator, country in values
	filter = Filter{
		Operator: "NOT IN",
		Values:   []string{"UK", "DE"},
	}
	country = "US"
	expected = true
	result = matchCountry(filter, country)
	if result != expected {
		t.Errorf("matchCountry returned unexpected result. Expected: %v, Got: %v", expected, result)
	}

	// Test case 4: NOT IN operator, country not in values
	filter = Filter{
		Operator: "NOT IN",
		Values:   []string{"UK", "DE"},
	}
	country = "DE"
	expected = false
	result = matchCountry(filter, country)
	if result != expected {
		t.Errorf("matchCountry returned unexpected result. Expected: %v, Got: %v", expected, result)
	}
}

func TestMatchCustomString(t *testing.T) {
	// Test case 1: == operator, prop eq value
	filter := Filter{
		Operator: "==",
		Name:     "best_friend",
		Values:   []string{"Winnie Pooh"},
	}

	ext := "{\"custom_attributes\":{\"best_friend\":\"Winnie Pooh\"}}"
	expected := true
	result := matchCustomString(filter, ext)
	if result != expected {
		t.Errorf("matchCustomString returned unexpected result. Expected: %v, Got: %v", expected, result)
	}

	// Test case 2: == operator, prop not eq value
	filter = Filter{
		Operator: "==",
		Name:     "best_friend",
		Values:   []string{"Winnie Pooh"},
	}

	ext = "{\"custom_attributes\":{\"best_friend\":\"Tigger\"}}"
	expected = false
	result = matchCustomString(filter, ext)
	if result != expected {
		t.Errorf("matchCustomString returned unexpected result. Expected: %v, Got: %v", expected, result)
	}

	// Test case 3: != operator, prop not eq value
	filter = Filter{
		Operator: "!=",
		Name:     "best_friend",
		Values:   []string{"Winnie Pooh"},
	}

	ext = "{\"custom_attributes\":{\"best_friend\":\"Tigger\"}}"
	expected = true
	result = matchCustomString(filter, ext)
	if result != expected {
		t.Errorf("matchCustomString returned unexpected result. Expected: %v, Got: %v", expected, result)
	}

	// Test case 4: != operator, prop eq value
	filter = Filter{
		Operator: "!=",
		Name:     "best_friend",
		Values:   []string{"Winnie Pooh"},
	}

	ext = "{\"custom_attributes\":{\"best_friend\":\"Winnie Pooh\"}}"
	expected = false
	result = matchCustomString(filter, ext)
	if result != expected {
		t.Errorf("matchCustomString returned unexpected result. Expected: %v, Got: %v", expected, result)
	}

	// Test case 5: Invalid customProps JSON
	filter = Filter{
		Operator: "==",
		Name:     "key1",
		Values:   []string{"value1"},
	}
	ext = `invalid JSON`
	expected = false
	result = matchCustomString(filter, ext)
	if result != expected {
		t.Errorf("matchCustomString returned unexpected result. Expected: %v, Got: %v", expected, result)
	}
}
