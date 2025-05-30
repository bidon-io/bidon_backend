package store

import (
	"context"
	"testing"
	"time"

	"github.com/google/go-cmp/cmp"

	"github.com/bidon-io/bidon-backend/config"
	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/bidon-io/bidon-backend/internal/db/dbtest"
)

func TestAdUnitLookup_GetInternalIDByUID(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	lineItem1 := dbtest.CreateLineItem(t, tx, func(item *db.LineItem) {
		item.PublicUID.Int64 = 12345
		item.PublicUID.Valid = true
	})

	lineItem2 := dbtest.CreateLineItem(t, tx, func(item *db.LineItem) {
		item.PublicUID.Int64 = 67890
		item.PublicUID.Valid = true
	})

	// Pure database method - no cache needed
	lookup := &AdUnitLookup{DB: tx}

	tests := []struct {
		name       string
		uid        string
		expectedID int64
	}{
		{
			name:       "valid UID returns correct internal ID",
			uid:        "12345",
			expectedID: lineItem1.ID,
		},
		{
			name:       "different UID returns different internal ID",
			uid:        "67890",
			expectedID: lineItem2.ID,
		},
		{
			name:       "invalid UID returns 0",
			uid:        "invalid",
			expectedID: 0,
		},
		{
			name:       "empty UID returns 0",
			uid:        "",
			expectedID: 0,
		},
		{
			name:       "non-existent UID returns 0",
			uid:        "99999",
			expectedID: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := lookup.GetInternalIDByUID(context.Background(), tt.uid)

			if err != nil {
				t.Fatalf("GetInternalIDByUID() error = %v", err)
			}

			if diff := cmp.Diff(tt.expectedID, result); diff != "" {
				t.Errorf("GetInternalIDByUID() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}

func TestAdUnitLookup_GetInternalIDByUIDCached(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	lineItem1 := dbtest.CreateLineItem(t, tx, func(item *db.LineItem) {
		item.PublicUID.Int64 = 12345
		item.PublicUID.Valid = true
	})

	lineItem2 := dbtest.CreateLineItem(t, tx, func(item *db.LineItem) {
		item.PublicUID.Int64 = 67890
		item.PublicUID.Valid = true
	})

	tests := []struct {
		name       string
		uid        string
		expectedID int64
	}{
		{
			name:       "valid UID returns correct internal ID",
			uid:        "12345",
			expectedID: lineItem1.ID,
		},
		{
			name:       "different UID returns different internal ID",
			uid:        "67890",
			expectedID: lineItem2.ID,
		},
		{
			name:       "invalid UID returns 0",
			uid:        "invalid",
			expectedID: 0,
		},
		{
			name:       "empty UID returns 0",
			uid:        "",
			expectedID: 0,
		},
		{
			name:       "non-existent UID returns 0",
			uid:        "99999",
			expectedID: 0,
		},
	}

	cache := config.NewMemoryCacheOf[int64](time.Minute)
	lookup := &AdUnitLookup{
		DB:    tx,
		Cache: cache,
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			result, err := lookup.GetInternalIDByUIDCached(context.Background(), tt.uid)

			if err != nil {
				t.Fatalf("GetInternalIDByUIDCached() error = %v", err)
			}

			if diff := cmp.Diff(tt.expectedID, result); diff != "" {
				t.Errorf("GetInternalIDByUIDCached() mismatch (-want +got):\n%s", diff)
			}
		})
	}
}
