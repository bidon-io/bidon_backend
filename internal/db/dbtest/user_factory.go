package dbtest

import (
	"fmt"

	"github.com/bidon-io/bidon-backend/internal/db"
)

type UserFactory struct {
	Email func(int) string
}

func (u UserFactory) Build(i int) db.User {
	user := db.User{}

	if u.Email == nil {
		user.Email = fmt.Sprintf("test%d@email.com", i)
	}

	return user
}
