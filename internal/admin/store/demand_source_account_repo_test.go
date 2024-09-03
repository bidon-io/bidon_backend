package adminstore_test

import (
	"context"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/adapter"
	"github.com/bidon-io/bidon-backend/internal/admin"
	adminstore "github.com/bidon-io/bidon-backend/internal/admin/store"
	"github.com/bidon-io/bidon-backend/internal/db"
	"github.com/bidon-io/bidon-backend/internal/db/dbtest"
	"github.com/google/go-cmp/cmp"
)

func TestDemandSourceAccountRepo_List(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := adminstore.NewDemandSourceAccountRepo(tx)
	demandSources := make([]db.DemandSource, 3)
	demandSources[0] = dbtest.CreateDemandSource(t, tx, func(source *db.DemandSource) {
		source.APIKey = string(adapter.ApplovinKey)
	})
	demandSources[1] = dbtest.CreateDemandSource(t, tx, func(source *db.DemandSource) {
		source.APIKey = string(adapter.BidmachineKey)
	})
	demandSources[2] = dbtest.CreateDemandSource(t, tx, func(source *db.DemandSource) {
		source.APIKey = string(adapter.UnityAdsKey)
	})
	user := dbtest.CreateUser(t, tx)
	accounts := []admin.DemandSourceAccountAttrs{
		{
			UserID:         user.ID,
			Type:           "DemandSourceAccount::Applovin",
			DemandSourceID: demandSources[0].ID,
			IsBidding:      ptr(false),
			Extra:          map[string]any{"key": "value"},
		},
		{
			UserID:         user.ID,
			Type:           "DemandSourceAccount::Bidmachine",
			DemandSourceID: demandSources[1].ID,
			IsBidding:      ptr(true),
			Extra:          map[string]any{"key": "value"},
		},
		{
			UserID:         user.ID,
			Type:           "DemandSourceAccount::UnityAds",
			DemandSourceID: demandSources[2].ID,
			IsBidding:      nil,
			Extra:          map[string]any{"key": "value"},
		},
	}

	want := make([]admin.DemandSourceAccount, len(accounts))
	for i, attrs := range accounts {
		account, err := repo.Create(context.Background(), &attrs)
		if err != nil {
			t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", &attrs, nil, err, account, nil)
		}

		want[i] = *account
		want[i].User = *adminstore.UserResource(&user)
		want[i].DemandSource = *adminstore.DemandSourceResource(&demandSources[i])
	}

	got, err := repo.List(context.Background(), nil)
	if err != nil {
		t.Fatalf("repo.List(ctx) = %v, %q; want %+v, %v", got, err, want, nil)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("repo.List(ctx) mismatch (-want, +got):\n%s", diff)
	}
}

func TestDemandSourceAccountRepo_ListOwnedByUserOrSharedRepo(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	users := make([]db.User, 2)
	for i := range users {
		users[i] = dbtest.CreateUser(t, tx)
	}

	dbFirstUserAccounts := make([]db.DemandSourceAccount, 2)
	for i := range dbFirstUserAccounts {
		dbFirstUserAccounts[i] = dbtest.CreateDemandSourceAccount(t, tx, func(account *db.DemandSourceAccount) {
			account.User = users[0]
		})
	}
	dbSecondUserAccounts := make([]db.DemandSourceAccount, 2)
	for i := range dbSecondUserAccounts {
		dbSecondUserAccounts[i] = dbtest.CreateDemandSourceAccount(t, tx, func(account *db.DemandSourceAccount) {
			account.User = users[1]
		})
	}
	// Need to remove constraint
	//dbSharedAccounts := dbtest.CreateList[db.DemandSourceAccount](t, tx, dbtest.DemandSourceAccountFactory{
	//	User: func(i int) db.User {
	//		return db.User{}
	//	},
	//}, 2)

	firstUserAccounts := make([]admin.DemandSourceAccount, 2)
	secondUserAccounts := make([]admin.DemandSourceAccount, 2)
	//sharedAccounts := make([]admin.DemandSourceAccount, 2)
	for i := 0; i < 2; i++ {
		firstUserAccounts[i] = adminstore.DemandSourceAccountResource(&dbFirstUserAccounts[i])
		secondUserAccounts[i] = adminstore.DemandSourceAccountResource(&dbSecondUserAccounts[i])
		//sharedAccounts[i] = adminstore.DemandSourceAccountResource(&dbSharedAccounts[i])
	}
	//for _, a := range dbSharedAccounts {
	//	firstUserAccounts = append(firstUserAccounts, adminstore.DemandSourceAccountResource(&a))
	//	secondUserAccounts = append(secondUserAccounts, adminstore.DemandSourceAccountResource(&a))
	//}

	repo := adminstore.NewDemandSourceAccountRepo(tx)

	tests := []struct {
		name   string
		userID int64
		want   []admin.DemandSourceAccount
	}{
		{
			"first user",
			users[0].ID,
			firstUserAccounts,
		},
		{
			"second user",
			users[1].ID,
			secondUserAccounts,
		},
		{
			"other user",
			999,
			[]admin.DemandSourceAccount{},
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.ListOwnedByUserOrShared(context.Background(), tt.userID)
			if err != nil {
				t.Fatalf("ListOwnedByUserOrShared() got %v; want %+v", err, tt.want)
			}

			if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Fatalf("ListOwnedByUserOrShared() mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestDemandSourceAccountRepo_Find(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := adminstore.NewDemandSourceAccountRepo(tx)

	user := dbtest.CreateUser(t, tx)
	demandSource := dbtest.CreateDemandSource(t, tx, func(source *db.DemandSource) {
		source.APIKey = string(adapter.BidmachineKey)
	})
	attrs := &admin.DemandSourceAccountAttrs{
		UserID:         user.ID,
		Type:           "DemandSourceAccount::Bidmachine",
		DemandSourceID: demandSource.ID,
		IsBidding:      ptr(true),
		Extra:          map[string]any{"key": "value"},
	}

	want, err := repo.Create(context.Background(), attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", attrs, nil, err, want, nil)
	}
	want.User = *adminstore.UserResource(&user)
	want.DemandSource = *adminstore.DemandSourceResource(&demandSource)

	got, err := repo.Find(context.Background(), want.ID)
	if err != nil {
		t.Fatalf("repo.Find(ctx) = %v, %q; want %+v, %v", got, err, want, nil)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("repo.List(ctx) mismatch (-want, +got):\n%s", diff)
	}
}

func TestDemandSourceAccountRepo_FindOwnedByUserOrSharedRepo(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	users := make([]db.User, 2)
	for i := range users {
		users[i] = dbtest.CreateUser(t, tx)
	}

	dbFirstUserAccounts := make([]db.DemandSourceAccount, 2)
	for i := range dbFirstUserAccounts {
		dbFirstUserAccounts[i] = dbtest.CreateDemandSourceAccount(t, tx, func(account *db.DemandSourceAccount) {
			account.User = users[0]
		})
	}
	dbSecondUserAccounts := make([]db.DemandSourceAccount, 2)
	for i := range dbSecondUserAccounts {
		dbSecondUserAccounts[i] = dbtest.CreateDemandSourceAccount(t, tx, func(account *db.DemandSourceAccount) {
			account.User = users[1]
		})
	}
	// Need to remove constraint
	//dbSharedAccounts := dbtest.CreateList[db.DemandSourceAccount](t, tx, dbtest.DemandSourceAccountFactory{
	//	User: func(i int) db.User {
	//		return db.User{}
	//	},
	//}, 2)

	firstUserAccounts := make([]admin.DemandSourceAccount, 2)
	secondUserAccounts := make([]admin.DemandSourceAccount, 2)
	//sharedAccounts := make([]admin.DemandSourceAccount, 2)
	for i := 0; i < 2; i++ {
		firstUserAccounts[i] = adminstore.DemandSourceAccountResource(&dbFirstUserAccounts[i])
		secondUserAccounts[i] = adminstore.DemandSourceAccountResource(&dbSecondUserAccounts[i])
		//sharedAccounts[i] = adminstore.DemandSourceAccountResource(&dbSharedAccounts[i])
	}
	//for _, a := range dbSharedAccounts {
	//	firstUserAccounts = append(firstUserAccounts, adminstore.DemandSourceAccountResource(&a))
	//	secondUserAccounts = append(secondUserAccounts, adminstore.DemandSourceAccountResource(&a))
	//}

	repo := adminstore.NewDemandSourceAccountRepo(tx)

	tests := []struct {
		name    string
		userID  int64
		id      int64
		want    *admin.DemandSourceAccount
		wantErr bool
	}{
		{
			"first user, first user's account",
			users[0].ID,
			firstUserAccounts[0].ID,
			&firstUserAccounts[0],
			false,
		},
		{
			"first user, second user's account",
			users[0].ID,
			secondUserAccounts[0].ID,
			nil,
			true,
		},
		{
			"second user, second user's account",
			users[1].ID,
			secondUserAccounts[0].ID,
			&secondUserAccounts[0],
			false,
		},
		{
			"second user, first user's account",
			users[1].ID,
			firstUserAccounts[0].ID,
			nil,
			true,
		},
		{
			"non-existent user",
			999,
			999,
			nil,
			true,
		},
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			got, err := repo.FindOwnedByUserOrShared(context.Background(), tt.userID, tt.id)
			if tt.wantErr {
				if err == nil {
					t.Errorf("FindOwnedByUserOrShared() = %+v; want error", got)
				}
			} else if diff := cmp.Diff(tt.want, got); diff != "" {
				t.Errorf("FindOwnedByUserOrShared() mismatch (-want, +got):\n%s", diff)
			}
		})
	}
}

func TestDemandSourceAccountRepo_Update(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := adminstore.NewDemandSourceAccountRepo(tx)

	user := dbtest.CreateUser(t, tx)
	demandSource := dbtest.CreateDemandSource(t, tx, func(source *db.DemandSource) {
		source.APIKey = string(adapter.BidmachineKey)
	})
	attrs := admin.DemandSourceAccountAttrs{
		UserID:         user.ID,
		Type:           "DemandSourceAccount::Bidmachine",
		DemandSourceID: demandSource.ID,
		IsBidding:      ptr(true),
		Extra:          map[string]any{"key": "value"},
	}

	account, err := repo.Create(context.Background(), &attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", &attrs, nil, err, account, nil)
	}

	want := account
	want.Extra = map[string]any{"key": "value2"}
	want.IsBidding = ptr(false)

	updateParams := &admin.DemandSourceAccountAttrs{
		Extra:     want.Extra,
		IsBidding: want.IsBidding,
	}
	got, err := repo.Update(context.Background(), account.ID, updateParams)
	if err != nil {
		t.Fatalf("repo.Update(ctx, %+v) = %v, %q; want %T, %v", updateParams, nil, err, got, nil)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("repo.Find(ctx) mismatch (-want, +got):\n%s", diff)
	}
}

func TestDemandSourceAccountRepo_Delete(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := adminstore.NewDemandSourceAccountRepo(tx)

	user := dbtest.CreateUser(t, tx)
	demandSource := dbtest.CreateDemandSource(t, tx, func(source *db.DemandSource) {
		source.APIKey = string(adapter.BidmachineKey)
	})
	attrs := &admin.DemandSourceAccountAttrs{
		UserID:         user.ID,
		Type:           "DemandSourceAccount::Bidmachine",
		DemandSourceID: demandSource.ID,
		IsBidding:      ptr(true),
		Extra:          map[string]any{"key": "value"},
	}
	account, err := repo.Create(context.Background(), attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", attrs, nil, err, account, nil)
	}

	err = repo.Delete(context.Background(), account.ID)
	if err != nil {
		t.Fatalf("repo.Delete(ctx, %v) = %q, want %v", account.ID, err, nil)
	}

	got, err := repo.Find(context.Background(), account.ID)
	if got != nil {
		t.Fatalf("repo.Find(ctx, %v) = %+v, %q; want %v, %q", account.ID, got, err, nil, "record not found")
	}
}
