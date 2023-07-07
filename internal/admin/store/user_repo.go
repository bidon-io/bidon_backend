package store

import (
	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/db"
)

type UserRepo = resourceRepo[admin.User, admin.UserAttrs, db.User]

func NewUserRepo(db *db.DB) *UserRepo {
	return &UserRepo{
		db:           db,
		mapper:       userMapper{},
		associations: []string{},
	}
}

type userMapper struct{}

//lint:ignore U1000 this method is used by generic struct
func (m userMapper) dbModel(u *admin.UserAttrs, id int64) *db.User {
	return &db.User{
		Model: db.Model{ID: id},
		Email: u.Email,
	}
}

//lint:ignore U1000 this method is used by generic struct
func (m userMapper) resource(u *db.User) admin.User {
	return admin.User{
		ID:        u.ID,
		UserAttrs: m.resourceAttrs(u),
	}
}

func (m userMapper) resourceAttrs(u *db.User) admin.UserAttrs {
	return admin.UserAttrs{
		Email: u.Email,
	}
}
