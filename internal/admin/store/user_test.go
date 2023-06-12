package store_test

import (
	"context"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/admin/store"
	"github.com/google/go-cmp/cmp"
)

func TestUserRepo_List(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := store.NewUserRepo(tx)

	users := []admin.UserAttrs{
		{
			Email: "user1@example.com",
		},
		{
			Email: "user2@example.com",
		},
	}

	want := make([]admin.User, len(users))
	for i, attrs := range users {
		user, err := repo.Create(context.Background(), &attrs)
		if err != nil {
			t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", &attrs, nil, err, user, nil)
		}

		want[i] = *user
	}

	got, err := repo.List(context.Background())
	if err != nil {
		t.Fatalf("repo.List(ctx) = %v, %q; want %+v, %v", got, err, want, nil)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("repo.List(ctx) mismatch (-want, +got):\n%s", diff)
	}
}

func TestUserRepo_Find(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := store.NewUserRepo(tx)

	attrs := &admin.UserAttrs{
		Email: "user1@example.com",
	}

	want, err := repo.Create(context.Background(), attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", attrs, nil, err, want, nil)
	}

	got, err := repo.Find(context.Background(), want.ID)
	if err != nil {
		t.Fatalf("repo.Find(ctx) = %v, %q; want %+v, %v", got, err, want, nil)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("repo.List(ctx) mismatch (-want, +got):\n%s", diff)
	}
}

func TestUserRepo_Update(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := store.NewUserRepo(tx)

	attrs := admin.UserAttrs{
		Email: "user1@example.com",
	}

	user, err := repo.Create(context.Background(), &attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", &attrs, nil, err, user, nil)
	}

	want := user
	want.Email = "user1alt@example.com"

	updateParams := &admin.UserAttrs{
		Email: want.Email,
	}
	got, err := repo.Update(context.Background(), user.ID, updateParams)
	if err != nil {
		t.Fatalf("repo.Update(ctx, %+v) = %v, %q; want %T, %v", updateParams, nil, err, got, nil)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("repo.Find(ctx) mismatch (-want, +got):\n%s", diff)
	}
}

func TestUserRepo_Delete(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := store.NewUserRepo(tx)

	attrs := &admin.UserAttrs{
		Email: "user1@example.com",
	}
	user, err := repo.Create(context.Background(), attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", attrs, nil, err, user, nil)
	}

	err = repo.Delete(context.Background(), user.ID)
	if err != nil {
		t.Fatalf("repo.Delete(ctx, %v) = %q, want %v", user.ID, err, nil)
	}

	got, err := repo.Find(context.Background(), user.ID)
	if got != nil {
		t.Fatalf("repo.Find(ctx, %v) = %+v, %q; want %v, %q", user.ID, got, err, nil, "record not found")
	}
}
