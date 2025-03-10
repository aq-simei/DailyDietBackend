package repositories

import (
	"context"
	"time"

	"daily-diet-backend/models"
	"daily-diet-backend/utils/crypt"
	"daily-diet-backend/utils/errors"
	"daily-diet-backend/utils/logger"

	"github.com/google/uuid"
	"gorm.io/gorm"
)

// Remove the local User struct since we'll use models.User
type UserRepository interface {
	CreateUser(c context.Context, data models.CreateUserDTO) (*models.User, error)
	GetUserByEmail(c context.Context, email string) (*models.User, error)
	Login(c context.Context, data models.LoginDTO) (*models.LoginResponse, error)
	CreateRefreshToken(c context.Context, data models.CreateRefreshTokenDTO) (*models.RefreshToken, error)
	ValidateRefreshToken(c context.Context, refreshToken string) (*models.RefreshToken, error)
	UpdateRefreshToken(c context.Context, refreshToken string, userId string) (*models.RefreshToken, error)
	GetUserByID(c context.Context, id string) (*models.User, error)
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

func (repo *userRepository) Login(c context.Context, data models.LoginDTO) (*models.LoginResponse, error) {
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

	var refreshToken models.RefreshToken
	var finalToken *models.RefreshToken

	existingRefreshToken := repo.db.Where("user_id = ?", user.ID).First(&refreshToken)
	if existingRefreshToken.Error != nil {
		if existingRefreshToken.Error == gorm.ErrRecordNotFound {
			// Create new refresh token
			newRefreshToken := models.CreateRefreshTokenDTO{
				UserID: user.ID,
			}
			if data.DeviceID != nil {
				newRefreshToken.DeviceID = data.DeviceID
			}
			finalToken, err = repo.CreateRefreshToken(c, newRefreshToken)
			if err != nil {
				return nil, errors.NewError(errors.Internal, "error creating refresh token", err)
			}
		} else {
			return nil, errors.NewError(errors.Internal, "error finding refresh token", existingRefreshToken.Error)
		}
	} else {
		// Update existing token
		finalToken, err = repo.UpdateRefreshToken(c, refreshToken.Token, user.ID.String())
		if err != nil {
			return nil, errors.NewError(errors.Internal, "error updating refresh token", err)
		}
	}

	response := &models.LoginResponse{
		User:         *user,
		RefreshToken: finalToken.Token,
	}

	return response, nil
}

func (repo *userRepository) CreateRefreshToken(c context.Context, data models.CreateRefreshTokenDTO) (*models.RefreshToken, error) {
	var token models.RefreshToken

	if data.DeviceID != nil {
		token.DeviceID = data.DeviceID
	}
	token.UserID = data.UserID

	token.Token = uuid.NewString()
	token.CreatedAt = time.Now()
	token.ExpireAt = time.Now().Add(time.Hour * 24 * 7) // 1 week
	token.Revoked = false
	if err := repo.db.Create(&token).Error; err != nil {
		return nil, errors.NewError(errors.Internal, "error creating refresh token", err)
	}
	return &token, nil
}

func (repo *userRepository) ValidateRefreshToken(c context.Context, refreshToken string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	result := repo.db.WithContext(c).Where("token = ?", refreshToken).First(&token)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, result.Error
		}
		return nil, errors.NewError(errors.Internal, "error finding refresh token in database", result.Error)
	}

	if token.Revoked {
		return nil, errors.NewError(errors.Unauthorized, "token revoked", nil)
	}
	if token.ExpireAt.Before(time.Now()) {
		return nil, errors.NewError(errors.Unauthorized, "token expired", nil)
	}
	// token is valid
	// update token and expire time
	return &token, nil
}

func (repo *userRepository) RevokeRefreshToken(c context.Context, refreshToken string) error {
	var token models.RefreshToken
	result := repo.db.WithContext(c).Where("token = ?", refreshToken).First(&token)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return result.Error
		}
		return errors.NewError(errors.Internal, "error finding refresh token in database", result.Error)
	}
	token.Revoked = true
	token.UpdatedAt = time.Now()
	if err := repo.db.Save(&token).Error; err != nil {
		return errors.NewError(errors.Internal, "error updating refresh token", err)
	}
	return nil
}

func (repo *userRepository) UpdateRefreshToken(c context.Context, refreshToken string, userId string) (*models.RefreshToken, error) {
	var token models.RefreshToken
	result := repo.db.WithContext(c).Where("token = ? AND user_id = ?", refreshToken, userId).First(&token)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, result.Error
		}
		return nil, errors.NewError(errors.Internal, "error finding refresh token in database", result.Error)
	}
	patches := map[string]interface{}{
		"expire_at":  time.Now().Add(time.Hour * 24 * 7),
		"updated_at": time.Now(),
		"token":      uuid.New().String(),
	}
	if err := repo.db.Model(&token).Updates(patches).Error; err != nil {
		return nil, errors.NewError(errors.Internal, "error updating refresh token", err)
	}
	// return updated token
	logger.Log(logger.DEBUG, "Refresh token updated")
	logger.Log(logger.DEBUG, token.Token)
	return &token, nil
}

func (repo *userRepository) GetUserByID(c context.Context, id string) (*models.User, error) {
	user := &models.User{}
	result := repo.db.WithContext(c).Where("id = ?", id).First(user)
	if result.Error != nil {
		if result.Error == gorm.ErrRecordNotFound {
			return nil, result.Error
		}
		return nil, errors.NewError(errors.Internal, "error finding user in database", result.Error)
	}
	return user, nil
}
