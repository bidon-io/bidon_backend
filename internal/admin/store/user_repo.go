package adminstore

import (
	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/db"
)

type UserRepo = resourceRepo[admin.User, admin.UserAttrs, db.User]

type userMapper struct{}

func NewUserRepo(d *db.DB) *UserRepo {
	return &UserRepo{
		db:           d,
		mapper:       userMapper{},
		associations: []string{},
	}
}

//lint:ignore U1000 this method is used by generic struct
func (m userMapper) dbModel(u *admin.UserAttrs, id int64) *db.User {
	du := &db.User{
		Model:   db.Model{ID: id},
		Email:   u.Email,
		IsAdmin: u.IsAdmin,
	}

	if u.Password != "" {
		passwordHash, _ := db.HashPassword(u.Password)
		du.PasswordHash = passwordHash
	}

	return du
}

//lint:ignore U1000 this method is used by generic struct
func (m userMapper) resource(u *db.User) admin.User {
	return admin.User{
		ID:      u.ID,
		Email:   u.Email,
		IsAdmin: u.IsAdmin,
	}
}
