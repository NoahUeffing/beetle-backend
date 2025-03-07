package postgres

import (
	"beetle/internal/domain"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	ReadDB  *gorm.DB
	WriteDB *gorm.DB
}

// CreateUser implements the domain.IUserService interface
func (s *UserService) CreateUser(input *domain.UserCreateInput) (*domain.User, error) {
	// Hash the password
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(input.Password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// Create a new user from the input
	user := &domain.User{
		ID:       uuid.New().String(),
		Username: input.Username,
		Email:    input.Email,
		Password: string(hashedPassword), // Store the hashed password
	}

	// Save the user to the database
	if err := s.WriteDB.Create(user).Error; err != nil {
		return nil, err
	}

	return user, nil
}

// CheckPassword verifies if the provided password matches the stored hash
func (s *UserService) CheckPassword(user *domain.User, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}
