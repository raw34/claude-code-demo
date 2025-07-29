package service

import (
	"errors"

	"github.com/user/user-management/internal/models"
	"github.com/user/user-management/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type UserService interface {
	GetByID(id uint) (*models.User, error)
	UpdateUser(id uint, updates map[string]interface{}) (*models.User, error)
	DeleteUser(id uint) error
	ListUsers(page, limit int) ([]models.User, int64, error)
}

type userService struct {
	userRepo repository.UserRepository
}

func NewUserService(userRepo repository.UserRepository) UserService {
	return &userService{
		userRepo: userRepo,
	}
}

func (s *userService) GetByID(id uint) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}
	return user, nil
}

func (s *userService) UpdateUser(id uint, updates map[string]interface{}) (*models.User, error) {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return nil, err
	}
	if user == nil {
		return nil, errors.New("user not found")
	}

	// 更新字段
	if username, ok := updates["username"].(string); ok && username != "" {
		// 检查用户名是否已被占用
		existingUser, _ := s.userRepo.GetByUsername(username)
		if existingUser != nil && existingUser.ID != id {
			return nil, errors.New("username already exists")
		}
		user.Username = username
	}

	if email, ok := updates["email"].(string); ok && email != "" {
		// 检查邮箱是否已被占用
		existingUser, _ := s.userRepo.GetByEmail(email)
		if existingUser != nil && existingUser.ID != id {
			return nil, errors.New("email already exists")
		}
		user.Email = email
	}

	if password, ok := updates["password"].(string); ok && password != "" {
		hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
		if err != nil {
			return nil, err
		}
		user.PasswordHash = string(hashedPassword)
	}

	if isActive, ok := updates["is_active"].(bool); ok {
		user.IsActive = isActive
	}

	if err := s.userRepo.Update(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *userService) DeleteUser(id uint) error {
	user, err := s.userRepo.GetByID(id)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("user not found")
	}

	return s.userRepo.Delete(id)
}

func (s *userService) ListUsers(page, limit int) ([]models.User, int64, error) {
	offset := (page - 1) * limit
	return s.userRepo.List(offset, limit)
}