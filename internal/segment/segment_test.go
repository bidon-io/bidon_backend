package segment_test

import (
	"context"
	"reflect"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/segment"
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
		FetchCachedFunc: func(ctx context.Context, appID int64) ([]segment.Segment, error) {
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
	expected := segments[0]
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
		FetchCachedFunc: func(ctx context.Context, appID int64) ([]segment.Segment, error) {
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
		FetchCachedFunc: func(ctx context.Context, appID int64) ([]segment.Segment, error) {
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

func TestSegment_StringID_known(t *testing.T) {
	segment := segment.Segment{ID: 123}

	expected := "123"
	actual := segment.StringID()
	if expected != actual {
		t.Errorf("segmentMatcher.Match returned unexpected segment. Expected: %v, Got: %v", expected, actual)
	}
}

func TestSegment_StringID_nil(t *testing.T) {
	segment := segment.Segment{}

	expected := ""
	actual := segment.StringID()
	if expected != actual {
		t.Errorf("segmentMatcher.Match returned unexpected segment. Expected: %v, Got: %v", expected, actual)
	}
}
