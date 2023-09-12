package adminstore

import (
	"context"
	"errors"

	"github.com/bidon-io/bidon-backend/internal/admin"
	"github.com/bidon-io/bidon-backend/internal/admin/auth"
	"github.com/bidon-io/bidon-backend/internal/db"
	"gorm.io/gorm"
)

type UserRepo struct {
	*resourceRepo[admin.User, admin.UserAttrs, db.User]
}

func NewUserRepo(d *db.DB) *UserRepo {
	return &UserRepo{
		resourceRepo: &resourceRepo[admin.User, admin.UserAttrs, db.User]{
			db:           d,
			mapper:       userMapper{},
			associations: []string{},
		},
	}
}

func (r *UserRepo) FindByEmailAndPassword(ctx context.Context, email, password string) (auth.User, error) {
	var user db.User
	if err := r.db.WithContext(ctx).Where("email = ?", email).First(&user).Error; err != nil {
		if errors.Is(err, gorm.ErrRecordNotFound) {
			return auth.User{}, auth.ErrInvalidCredentials
		}
		return auth.User{}, err
	}

	ok, err := db.ComparePassword(user.PasswordHash, password)
	if err != nil {
		return auth.User{}, err
	}
	if !ok {
		return auth.User{}, auth.ErrInvalidCredentials
	}

	return auth.User{
		ID:      user.ID,
		Email:   user.Email,
		IsAdmin: *user.IsAdmin,
	}, nil
}

type userMapper struct{}

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
