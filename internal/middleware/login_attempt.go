package middleware

import (
	"sync"
	"time"

	"github.com/gofiber/fiber/v2"
)

type LoginAttempt struct {
	Count       int
	LastAttempt time.Time
	LockedUntil time.Time
}

type LoginAttemptStore struct {
	mu      sync.RWMutex
	attempts map[string]*LoginAttempt
}

var store = &LoginAttemptStore{
	attempts: make(map[string]*LoginAttempt),
}

const (
	MaxAttempts    = 5
	LockDuration   = 5 * time.Minute
	WindowDuration = 15 * time.Minute
)

func (s *LoginAttemptStore) GetOrCreate(key string) *LoginAttempt {
	s.mu.Lock()
	defer s.mu.Unlock()

	if attempt, exists := s.attempts[key]; exists {
		return attempt
	}

	attempt := &LoginAttempt{
		Count:       0,
		LastAttempt: time.Now(),
	}
	s.attempts[key] = attempt
	return attempt
}

func (s *LoginAttemptStore) RecordFailed(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	attempt, exists := s.attempts[key]
	if !exists {
		attempt = &LoginAttempt{Count: 0}
		s.attempts[key] = attempt
	}

	attempt.Count++
	attempt.LastAttempt = time.Now()

	if attempt.Count >= MaxAttempts {
		attempt.LockedUntil = time.Now().Add(LockDuration)
	}
}

func (s *LoginAttemptStore) Reset(key string) {
	s.mu.Lock()
	defer s.mu.Unlock()

	delete(s.attempts, key)
}

func (s *LoginAttemptStore) IsLocked(key string) (bool, time.Time) {
	s.mu.RLock()
	defer s.mu.RUnlock()

	attempt, exists := s.attempts[key]
	if !exists {
		return false, time.Time{}
	}

	if attempt.LockedUntil.IsZero() {
		return false, time.Time{}
	}

	if time.Now().After(attempt.LockedUntil) {
		return false, time.Time{}
	}

	return true, attempt.LockedUntil
}

func (s *LoginAttemptStore) GetRemainingAttempts(key string) int {
	s.mu.RLock()
	defer s.mu.RUnlock()

	attempt, exists := s.attempts[key]
	if !exists {
		return MaxAttempts
	}

	if time.Since(attempt.LastAttempt) > WindowDuration {
		return MaxAttempts
	}

	remaining := MaxAttempts - attempt.Count
	if remaining < 0 {
		return 0
	}
	return remaining
}

func GetRemainingAttempts(email string) int {
	return store.GetRemainingAttempts("login:" + email)
}

func RecordLoginAttempt(email string, success bool) {
	key := "login:" + email
	if success {
		store.Reset(key)
		return
	}
	store.RecordFailed(key)
}

func LoginRateLimiter() fiber.Handler {
	return func(c *fiber.Ctx) error {
		if c.Path() != "/api/auth/login" {
			return c.Next()
		}

		var req struct {
			Email string `json:"email"`
		}
		if err := c.BodyParser(&req); err != nil {
			return c.Next()
		}

		if req.Email == "" {
			return c.Next()
		}

		key := "login:" + req.Email

		isLocked, lockedUntil := store.IsLocked(key)
		if isLocked {
			remainingTime := time.Until(lockedUntil)
			minutes := int(remainingTime.Minutes())
			seconds := int(remainingTime.Seconds()) % 60

			return c.Status(fiber.StatusTooManyRequests).JSON(fiber.Map{
				"success": false,
				"message": "Terlalu banyak percobaan gagal. Akun terkunci sementara.",
				"locked_until": lockedUntil,
				"remaining_time": fiber.Map{
					"minutes": minutes,
					"seconds": seconds,
				},
				"can_reset": true,
			})
		}

		return c.Next()
	}
}