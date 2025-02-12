package repositories

import (
	"context"

	"daily-diet-backend/models"
	"daily-diet-backend/utils/crypt"
	"daily-diet-backend/utils/errors"
	"daily-diet-backend/utils/logger"

	"gorm.io/gorm"
)

// Remove the local User struct since we'll use models.User
type UserRepository interface {
	CreateUser(c context.Context, data models.CreateUserDTO) (*models.User, error)
	GetUserByEmail(c context.Context, email string) (*models.User, error)
	Login(c context.Context, data models.LoginDTO) (*models.User, error)
}

type userRepository struct {
	db *gorm.DB
}

func NewUserRepository(db *gorm.DB) UserRepository {
	return &userRepository{db: db}
}

func (repo *userRepository) CreateUser(c context.Context, data models.CreateUserDTO) (*models.User, error) {
	// Check if user exists
	existingUser, err := repo.GetUserByEmail(c, data.Email)

	if existingUser != nil {
		return nil, errors.NewError(errors.Invalid, "user already exists", nil)
	}
	if err != gorm.ErrRecordNotFound {
		return nil, errors.NewError(errors.Internal, "database error", err)
	}

	// Hash password
	hashedPassword, err := crypt.HashPassword(data.Password)
	if err != nil {
		return nil, errors.NewError(errors.Internal, "error hashing password", err)
	}

	// Create user
	user := &models.User{
		Email:    data.Email,
		Name:     data.Name,
		Password: hashedPassword,
	}

	if err := repo.db.Create(user).Error; err != nil {
		return nil, errors.NewError(errors.Internal, "error creating user", err)
	}

	return user, nil
}

func (repo *userRepository) GetUserByEmail(c context.Context, email string) (*models.User, error) {
	user := &models.User{}
	result := repo.db.WithContext(c).
		Where("email = ?", email).First(user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, result.Error
		}
		return nil, errors.NewError(errors.Internal, "error finding user in database", result.Error)
	}
	logger.Log(logger.DEBUG, "User found: "+user.Email)
	return user, nil
}

func (repo *userRepository) Login(c context.Context, data models.LoginDTO) (*models.User, error) {
	user, err := repo.GetUserByEmail(c, data.Email)
	if err != nil {
		if err == gorm.ErrRecordNotFound {
			return nil, errors.NewError(errors.NotFound, "user not found", err)
		}
		return nil, err
	}
	if err := crypt.ComparePassword(user.Password, data.Password); err != nil {
		return nil, errors.NewError(errors.Unauthorized, "invalid password", err)
	}
	return user, nil
}
