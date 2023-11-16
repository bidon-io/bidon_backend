package dbtest

import (
	"fmt"
	"sync/atomic"
	"testing"

	"github.com/bidon-io/bidon-backend/internal/db"
)

var userCounter atomic.Uint64

func userDefaults(n uint64) func(*db.User) {
	return func(user *db.User) {
		if user.Email == "" {
			user.Email = fmt.Sprintf("test%d@email.com", n)
		}
	}
}

func BuildUser(opts ...func(*db.User)) db.User {
	var user db.User

	n := userCounter.Add(1)

	opts = append(opts, userDefaults(n))
	for _, opt := range opts {
		opt(&user)
	}

	return user
}

func CreateUser(t *testing.T, tx *db.DB, opts ...func(*db.User)) db.User {
	t.Helper()

	user := BuildUser(opts...)
	if err := tx.Create(&user).Error; err != nil {
		t.Fatalf("Failed to create user: %v", err)
	}

	return user
}
