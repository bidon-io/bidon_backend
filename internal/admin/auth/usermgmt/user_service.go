package usermgmt

import (
	"github.com/bidon-io/bidon-backend/internal/db"
)

type UserService struct {
	db *db.DB
}

func NewUserService(db *db.DB) *UserService {
	return &UserService{db: db}
}

func (s *UserService) GetUserByEmail(email string) (*db.User, error) {
	var user db.User
	query := s.db.Where("email = ?", email)
	if err := query.First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (s *UserService) ComparePassword(storedPasswordHash, password string) bool {
	result, _ := db.ComparePassword(storedPasswordHash, password)
	return result
}
