package service

import (
	"crypto/rand"
	"encoding/hex"
	"errors"
	"time"

	"github.com/golang-jwt/jwt/v5"
	"github.com/user/user-management/internal/models"
	"github.com/user/user-management/internal/repository"
	"golang.org/x/crypto/bcrypt"
)

type AuthService interface {
	Register(username, email, password string) (*models.User, error)
	Login(email, password string) (*models.User, string, string, error)
	RefreshToken(refreshToken string) (string, string, error)
	Logout(token string, userID uint) error
	ValidateToken(tokenString string) (uint, error)
}

type authService struct {
	userRepo       repository.UserRepository
	sessionService SessionService
	jwtSecret      string
	tokenExpiry    time.Duration
}

func NewAuthService(userRepo repository.UserRepository, sessionService SessionService, jwtSecret string, tokenExpiry time.Duration) AuthService {
	return &authService{
		userRepo:       userRepo,
		sessionService: sessionService,
		jwtSecret:      jwtSecret,
		tokenExpiry:    tokenExpiry,
	}
}

func (s *authService) Register(username, email, password string) (*models.User, error) {
	// 检查用户是否已存在
	existingUser, _ := s.userRepo.GetByEmail(email)
	if existingUser != nil {
		return nil, errors.New("email already exists")
	}

	existingUser, _ = s.userRepo.GetByUsername(username)
	if existingUser != nil {
		return nil, errors.New("username already exists")
	}

	// 哈希密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		return nil, err
	}

	// 创建用户
	user := &models.User{
		Username:     username,
		Email:        email,
		PasswordHash: string(hashedPassword),
		IsActive:     true,
	}

	if err := s.userRepo.Create(user); err != nil {
		return nil, err
	}

	return user, nil
}

func (s *authService) Login(email, password string) (*models.User, string, string, error) {
	// 查找用户
	user, err := s.userRepo.GetByEmail(email)
	if err != nil {
		return nil, "", "", err
	}
	if user == nil {
		return nil, "", "", errors.New("invalid credentials")
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.PasswordHash), []byte(password)); err != nil {
		return nil, "", "", errors.New("invalid credentials")
	}

	// 检查用户是否激活
	if !user.IsActive {
		return nil, "", "", errors.New("user account is disabled")
	}

	// 生成访问令牌
	accessToken, err := s.generateAccessToken(user.ID)
	if err != nil {
		return nil, "", "", err
	}

	// 在Redis中创建session
	err = s.sessionService.CreateSession(user.ID, accessToken, s.tokenExpiry)
	if err != nil {
		return nil, "", "", err
	}

	// 生成刷新令牌
	refreshToken, err := s.generateRefreshToken(user.ID)
	if err != nil {
		return nil, "", "", err
	}

	return user, accessToken, refreshToken, nil
}

func (s *authService) RefreshToken(refreshToken string) (string, string, error) {
	// 查找刷新令牌
	token, err := s.userRepo.GetRefreshToken(refreshToken)
	if err != nil {
		return "", "", err
	}
	if token == nil {
		return "", "", errors.New("invalid refresh token")
	}

	// 检查是否过期
	if time.Now().After(token.ExpiresAt) {
		s.userRepo.DeleteRefreshToken(refreshToken)
		return "", "", errors.New("refresh token expired")
	}

	// 生成新的访问令牌
	accessToken, err := s.generateAccessToken(token.UserID)
	if err != nil {
		return "", "", err
	}

	// 在Redis中创建新的session
	err = s.sessionService.CreateSession(token.UserID, accessToken, s.tokenExpiry)
	if err != nil {
		return "", "", err
	}

	// 生成新的刷新令牌
	newRefreshToken, err := s.generateRefreshToken(token.UserID)
	if err != nil {
		return "", "", err
	}

	// 删除旧的刷新令牌
	s.userRepo.DeleteRefreshToken(refreshToken)

	return accessToken, newRefreshToken, nil
}

func (s *authService) Logout(token string, userID uint) error {
	// 删除Redis中的session
	err := s.sessionService.DeleteSession(token)
	if err != nil {
		return err
	}
	
	// 删除数据库中的refresh tokens
	return s.userRepo.DeleteUserRefreshTokens(userID)
}

func (s *authService) ValidateToken(tokenString string) (uint, error) {
	// 首先检查Redis中的session
	sessionData, err := s.sessionService.GetSession(tokenString)
	if err != nil {
		return 0, err
	}
	if sessionData == nil {
		return 0, errors.New("session not found")
	}

	// 然后验证JWT token
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		return []byte(s.jwtSecret), nil
	})

	if err != nil {
		return 0, err
	}

	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		userID := uint(claims["user_id"].(float64))
		
		// 验证session中的用户ID与token中的一致
		if userID != sessionData.UserID {
			return 0, errors.New("session user mismatch")
		}
		
		return userID, nil
	}

	return 0, errors.New("invalid token")
}

func (s *authService) generateAccessToken(userID uint) (string, error) {
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(s.tokenExpiry).Unix(),
		"iat":     time.Now().Unix(),
	}

	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(s.jwtSecret))
}

func (s *authService) generateRefreshToken(userID uint) (string, error) {
	// 生成随机令牌
	bytes := make([]byte, 32)
	if _, err := rand.Read(bytes); err != nil {
		return "", err
	}
	tokenString := hex.EncodeToString(bytes)

	// 保存到数据库
	refreshToken := &models.RefreshToken{
		UserID:    userID,
		Token:     tokenString,
		ExpiresAt: time.Now().Add(7 * 24 * time.Hour),
	}

	if err := s.userRepo.SaveRefreshToken(refreshToken); err != nil {
		return "", err
	}

	return tokenString, nil
}