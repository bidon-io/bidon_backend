package segment_test

import (
	"context"
	"github.com/bidon-io/bidon-backend/internal/segment"
	"reflect"
	"testing"

	segmentmocks "github.com/bidon-io/bidon-backend/internal/segment/mocks"
)

func TestMatchTwoFilters(t *testing.T) {
	ctx := context.Background()
	segments := []segment.Segment{
		{
			ID: 123,
			Filters: []segment.Filter{
				{Type: "country", Name: "country", Operator: "IN", Values: []string{"US"}},
				{Type: "custom_string", Name: "best_friend", Operator: "==", Values: []string{"Winnie Pooh"}},
			},
		},
	}

	segmentFetcher := &segmentmocks.FetcherMock{
		FetchFunc: func(ctx context.Context, appID int64) ([]segment.Segment, error) {
			return segments, nil
		},
	}

	segmentMatcher := &segment.Matcher{
		Fetcher: segmentFetcher,
	}

	// Test case 1: two filters are matched
	params := &segment.Params{
		Ext:     "{\"custom_attributes\":{\"best_friend\":\"Winnie Pooh\"}}",
		Country: "US",
		AppID:   1,
	}

	result := segmentMatcher.Match(ctx, params)
	expected := segment.Segment{ID: 123}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("segmentMatcher.Match returned unexpected segment. Expected: %v, Got: %v", expected, result)
	}

	// Test case 2: only one filters is matched
	params = &segment.Params{
		Ext:     "{\"custom_attributes\":{\"best_friend\":\"Winnie Pooh\"}}",
		Country: "RU",
		AppID:   1,
	}

	result = segmentMatcher.Match(ctx, params)
	expected = segment.Segment{ID: 0}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("segmentMatcher.Match returned unexpected segment. Expected: %v, Got: %v", expected, result)
	}

	// Test case 3: filters are blank
	segments = []segment.Segment{
		{
			ID:      123,
			Filters: []segment.Filter{},
		},
	}

	segmentFetcher = &segmentmocks.FetcherMock{
		FetchFunc: func(ctx context.Context, appID int64) ([]segment.Segment, error) {
			return segments, nil
		},
	}

	segmentMatcher = &segment.Matcher{
		Fetcher: segmentFetcher,
	}
	params = &segment.Params{
		Ext:     "{\"custom_attributes\":{\"best_friend\":\"Winnie Pooh\"}}",
		Country: "US",
		AppID:   1,
	}

	result = segmentMatcher.Match(ctx, params)
	expected = segment.Segment{ID: 0}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("segmentMatcher.Match returned unexpected segment. Expected: %v, Got: %v", expected, result)
	}

	// Test case 4: 1 filter is unsupported
	segments = []segment.Segment{
		{
			ID: 123,
			Filters: []segment.Filter{
				{Type: "country", Name: "country", Operator: "IN", Values: []string{"US"}},
				{Type: "custom_string", Name: "best_friend", Operator: "==", Values: []string{"Winnie Pooh"}},
				{Type: "unsupported", Name: "unsupported", Operator: "==", Values: []string{"Smth"}},
			},
		},
	}

	segmentFetcher = &segmentmocks.FetcherMock{
		FetchFunc: func(ctx context.Context, appID int64) ([]segment.Segment, error) {
			return segments, nil
		},
	}

	segmentMatcher = &segment.Matcher{
		Fetcher: segmentFetcher,
	}
	params = &segment.Params{
		Ext:     "{\"custom_attributes\":{\"best_friend\":\"Winnie Pooh\"}}",
		Country: "US",
		AppID:   1,
	}

	result = segmentMatcher.Match(ctx, params)
	expected = segment.Segment{ID: 0}
	if !reflect.DeepEqual(result, expected) {
		t.Errorf("segmentMatcher.Match returned unexpected segment. Expected: %v, Got: %v", expected, result)
	}
}
