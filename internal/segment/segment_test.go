package segment

import (
	"github.com/google/go-cmp/cmp"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/admin"
)

func TestMatch(t *testing.T) {
	// Define test segments and params
	segments := []admin.Segment{
		{
			SegmentAttrs: admin.SegmentAttrs{
				Filters: []admin.SegmentFilter{
					{Type: "country", Operator: "IN", Values: []string{"US", "CA"}},
				}},
		},
		{
			SegmentAttrs: admin.SegmentAttrs{
				Filters: []admin.SegmentFilter{
					{Type: "country", Operator: "NOT IN", Values: []string{"UK", "DE"}},
				}},
		},
	}

	params := Params{
		Country: "US",
	}

	// Test case 1: Matching segment exists
	expected := &segments[0]
	result := Match(segments, params)
	if diff := cmp.Diff(result, expected); diff != "" {
		t.Errorf("Match returned unexpected result. Expected: %v, Got: %v", expected, result)
	}
}

func TestMatchCountry(t *testing.T) {
	// Test case 1: IN operator, country in values
	filter := admin.SegmentFilter{
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
	filter = admin.SegmentFilter{
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
	filter = admin.SegmentFilter{
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
	filter = admin.SegmentFilter{
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
	filter := admin.SegmentFilter{
		Operator: "==",
		Name:     "best_friend",
		Values:   []string{"Winnie Pooh"},
	}

	customProps := `{"best_friend":"Winnie Pooh"}`
	expected := true
	result := matchCustomString(filter, customProps)
	if result != expected {
		t.Errorf("matchCustomString returned unexpected result. Expected: %v, Got: %v", expected, result)
	}

	// Test case 2: == operator, prop not eq value
	filter = admin.SegmentFilter{
		Operator: "==",
		Name:     "best_friend",
		Values:   []string{"Winnie Pooh"},
	}

	customProps = `{"best_friend":"Tigger"}`
	expected = false
	result = matchCustomString(filter, customProps)
	if result != expected {
		t.Errorf("matchCustomString returned unexpected result. Expected: %v, Got: %v", expected, result)
	}

	// Test case 3: != operator, prop not eq value
	filter = admin.SegmentFilter{
		Operator: "!=",
		Name:     "best_friend",
		Values:   []string{"Winnie Pooh"},
	}

	customProps = `{"best_friend":"Tigger"}`
	expected = true
	result = matchCustomString(filter, customProps)
	if result != expected {
		t.Errorf("matchCustomString returned unexpected result. Expected: %v, Got: %v", expected, result)
	}

	// Test case 4: != operator, prop eq value
	filter = admin.SegmentFilter{
		Operator: "!=",
		Name:     "best_friend",
		Values:   []string{"Winnie Pooh"},
	}

	customProps = `{"best_friend":"Winnie Pooh"}`
	expected = false
	result = matchCustomString(filter, customProps)
	if result != expected {
		t.Errorf("matchCustomString returned unexpected result. Expected: %v, Got: %v", expected, result)
	}

	// Test case 5: Invalid customProps JSON
	filter = admin.SegmentFilter{
		Operator: "==",
		Name:     "key1",
		Values:   []string{"value1"},
	}
	customProps = `invalid JSON`
	expected = false
	result = matchCustomString(filter, customProps)
	if result != expected {
		t.Errorf("matchCustomString returned unexpected result. Expected: %v, Got: %v", expected, result)
	}
}
