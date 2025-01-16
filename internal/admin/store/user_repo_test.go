package adminstore_test

import (
	"context"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/admin/resource"

	"github.com/bidon-io/bidon-backend/internal/db"

	"github.com/bidon-io/bidon-backend/internal/admin"
	adminstore "github.com/bidon-io/bidon-backend/internal/admin/store"
	"github.com/google/go-cmp/cmp"
)

func TestUserRepo_List(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := adminstore.NewUserRepo(tx)

	users := []admin.UserAttrs{
		{
			Email:    "user1@example.com",
			IsAdmin:  ptr(true),
			Password: "password",
		},
		{
			Email:    "user2@example.com",
			IsAdmin:  ptr(false),
			Password: "password",
		},
	}

	items := make([]admin.User, len(users))
	for i, attrs := range users {
		user, err := repo.Create(context.Background(), &attrs)
		if err != nil {
			t.Fatalf("repo.Create(ctx, %+v) = %v, %q; items %T, %v", &attrs, nil, err, user, nil)
		}

		items[i] = *user
	}

	want := &resource.Collection[admin.User]{
		Items: items,
		Meta:  resource.CollectionMeta{TotalCount: int64(len(items))},
	}

	got, err := repo.List(context.Background(), nil)
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

	repo := adminstore.NewUserRepo(tx)

	attrs := &admin.UserAttrs{
		Email:    "user1@example.com",
		IsAdmin:  ptr(true),
		Password: "password",
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

	repo := adminstore.NewUserRepo(tx)

	attrs := admin.UserAttrs{
		Email:    "user1@example.com",
		IsAdmin:  ptr(true),
		Password: "password",
	}

	user, err := repo.Create(context.Background(), &attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", &attrs, nil, err, user, nil)
	}

	want := user
	want.Email = "user1alt@example.com"
	want.IsAdmin = ptr(false)

	updateParams := &admin.UserAttrs{
		Email:    want.Email,
		IsAdmin:  want.IsAdmin,
		Password: "passwordalt",
	}
	got, err := repo.Update(context.Background(), user.ID, updateParams)
	if err != nil {
		t.Fatalf("repo.Update(ctx, %+v) = %v, %q; want %T, %v", updateParams, nil, err, got, nil)
	}
	dbModel := &db.User{}
	if err := tx.First(dbModel, user.ID).Error; err != nil {
		t.Fatalf("tx.First(dbModel, %v) = %q, want %v", user.ID, err, nil)
	}
	result, err := db.ComparePassword(dbModel.PasswordHash, updateParams.Password)
	if err != nil {
		t.Fatalf("db.ComparePassword(dbModel.PasswordHash, %v) = %q, want %v", updateParams.Password, err, nil)
	}
	if !result {
		t.Fatalf("db.ComparePassword(dbModel.PasswordHash, %v) = %v, want %v", updateParams.Password, result, true)
	}

	if diff := cmp.Diff(want, got); diff != "" {
		t.Fatalf("repo.Find(ctx) mismatch (-want, +got):\n%s", diff)
	}
}

func TestUserRepo_Delete(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := adminstore.NewUserRepo(tx)

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

func TestUserRepo_UpdatePassword(t *testing.T) {
	tx := testDB.Begin()
	defer tx.Rollback()

	repo := adminstore.NewUserRepo(tx)

	attrs := &admin.UserAttrs{
		Email:    "user1@example.com",
		IsAdmin:  ptr(true),
		Password: "oldpassword",
	}

	user, err := repo.Create(context.Background(), attrs)
	if err != nil {
		t.Fatalf("repo.Create(ctx, %+v) = %v, %q; want %T, %v", attrs, nil, err, user, nil)
	}

	tests := []struct {
		name            string
		userID          int64
		currentPassword string
		newPassword     string
		expectError     bool
		expectErrorMsg  string
	}{
		{
			name:            "successful password update",
			userID:          user.ID,
			currentPassword: "oldpassword",
			newPassword:     "newpassword",
			expectError:     false,
		},
		{
			name:            "incorrect current password",
			userID:          user.ID,
			currentPassword: "wrongpassword",
			newPassword:     "newpassword",
			expectError:     true,
			expectErrorMsg:  "current password is incorrect",
		},
		{
			name:            "user not found",
			userID:          99999,
			currentPassword: "oldpassword",
			newPassword:     "newpassword",
			expectError:     true,
			expectErrorMsg:  "user not found",
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			err := repo.UpdatePassword(context.Background(), tt.userID, tt.currentPassword, tt.newPassword)
			if tt.expectError {
				if err == nil {
					t.Fatalf("expected error but got nil")
				}
				if err.Error() != tt.expectErrorMsg {
					t.Fatalf("expected error %q but got %q", tt.expectErrorMsg, err.Error())
				}
				return
			}

			if err != nil {
				t.Fatalf("unexpected error: %q", err)
			}

			dbUser := &db.User{}
			if err := tx.First(dbUser, tt.userID).Error; err != nil {
				t.Fatalf("tx.First(dbUser, %v) = %q; want %v", tt.userID, err, nil)
			}

			result, err := db.ComparePassword(dbUser.PasswordHash, tt.newPassword)
			if err != nil {
				t.Fatalf("db.ComparePassword(dbUser.PasswordHash, newPassword) = %q; want %v", err, nil)
			}
			if !result {
				t.Fatalf("db.ComparePassword(dbUser.PasswordHash, newPassword) = %v; want %v", result, true)
			}
		})
	}
}
