package store

import "github.com/bidon-io/bidon-backend/internal/admin"

type UserRepo = resourceRepo[admin.User, admin.UserAttrs, user, *user]

type user struct {
	Model
	Email string `gorm:"column:email;type:varchar;not null"`
}

//lint:ignore U1000 this method is used by generic struct
func (u *user) initFromResourceAttrs(attrs *admin.UserAttrs) {
	u.Email = attrs.Email
}

//lint:ignore U1000 this method is used by generic struct
func (u *user) toResource() admin.User {
	return admin.User{
		ID: u.ID,
		UserAttrs: admin.UserAttrs{
			Email: u.Email,
		},
	}
}
