package service

import (
	"context"
	"encoding/json"
	"fmt"
	"time"

	"github.com/redis/go-redis/v9"
)

type SessionService interface {
	CreateSession(userID uint, token string, expiry time.Duration) error
	GetSession(token string) (*SessionData, error)
	DeleteSession(token string) error
	DeleteUserSessions(userID uint) error
	RefreshSession(token string, expiry time.Duration) error
}

type SessionData struct {
	UserID    uint      `json:"user_id"`
	CreatedAt time.Time `json:"created_at"`
}

type sessionService struct {
	redis *redis.Client
	ctx   context.Context
}

func NewSessionService(redisClient *redis.Client) SessionService {
	return &sessionService{
		redis: redisClient,
		ctx:   context.Background(),
	}
}

func (s *sessionService) CreateSession(userID uint, token string, expiry time.Duration) error {
	sessionData := SessionData{
		UserID:    userID,
		CreatedAt: time.Now(),
	}

	data, err := json.Marshal(sessionData)
	if err != nil {
		return err
	}

	// 存储session
	key := fmt.Sprintf("session:%s", token)
	err = s.redis.Set(s.ctx, key, data, expiry).Err()
	if err != nil {
		return err
	}

	// 将token添加到用户的session集合中
	userKey := fmt.Sprintf("user:sessions:%d", userID)
	err = s.redis.SAdd(s.ctx, userKey, token).Err()
	if err != nil {
		return err
	}

	// 设置用户session集合的过期时间（比token稍长）
	s.redis.Expire(s.ctx, userKey, expiry+time.Hour)

	return nil
}

func (s *sessionService) GetSession(token string) (*SessionData, error) {
	key := fmt.Sprintf("session:%s", token)
	
	data, err := s.redis.Get(s.ctx, key).Result()
	if err == redis.Nil {
		return nil, nil
	}
	if err != nil {
		return nil, err
	}

	var sessionData SessionData
	err = json.Unmarshal([]byte(data), &sessionData)
	if err != nil {
		return nil, err
	}

	return &sessionData, nil
}

func (s *sessionService) DeleteSession(token string) error {
	// 获取session数据以获取用户ID
	sessionData, err := s.GetSession(token)
	if err != nil {
		return err
	}

	if sessionData != nil {
		// 从用户的session集合中移除
		userKey := fmt.Sprintf("user:sessions:%d", sessionData.UserID)
		s.redis.SRem(s.ctx, userKey, token)
	}

	// 删除session
	key := fmt.Sprintf("session:%s", token)
	return s.redis.Del(s.ctx, key).Err()
}

func (s *sessionService) DeleteUserSessions(userID uint) error {
	userKey := fmt.Sprintf("user:sessions:%d", userID)
	
	// 获取用户的所有session
	tokens, err := s.redis.SMembers(s.ctx, userKey).Result()
	if err != nil {
		return err
	}

	// 删除所有session
	for _, token := range tokens {
		key := fmt.Sprintf("session:%s", token)
		s.redis.Del(s.ctx, key)
	}

	// 删除用户的session集合
	return s.redis.Del(s.ctx, userKey).Err()
}

func (s *sessionService) RefreshSession(token string, expiry time.Duration) error {
	key := fmt.Sprintf("session:%s", token)
	return s.redis.Expire(s.ctx, key, expiry).Err()
}