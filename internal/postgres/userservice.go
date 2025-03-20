package postgres

import (
	"beetle/internal/auth"
	"beetle/internal/domain"

	"strings"

	"github.com/google/uuid"
	"golang.org/x/crypto/bcrypt"
	"gorm.io/gorm"
)

type UserService struct {
	ReadDB      *gorm.DB
	WriteDB     *gorm.DB
	AuthService auth.IAuthService
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
		Username: input.Username,
		Email:    strings.ToLower(input.Email),
		Password: string(hashedPassword), // Store the hashed password
	}

	// Save the user to the database and return all fields
	if err := s.WriteDB.Select("*").Create(user).Error; err != nil {
		// TODO: Check if the error is a duplicate username or email
		return nil, err
	}

	return user, nil
}

func (us *UserService) ReadByEmail(email string) (*domain.User, error) {
	var user domain.User
	if err := us.ReadDB.Where("email = LOWER(?)", email).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}

func (us *UserService) CreateAuthToken(u *domain.User) (*domain.UserAuthToken, error) {
	ts, err := us.AuthService.NewToken(&auth.ClaimsData{
		UserId: u.ID,
	})
	if err != nil {
		return nil, err
	}

	return &domain.UserAuthToken{Token: ts}, err
}

// CheckPassword verifies if the provided password matches the stored hash
func (s *UserService) CheckPassword(user *domain.User, password string) error {
	return bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password))
}

func (us *UserService) ReadByID(id uuid.UUID) (*domain.User, error) {
	var user domain.User
	if err := us.ReadDB.Where("id = ?", id).First(&user).Error; err != nil {
		return nil, err
	}
	return &user, nil
}
