package iam

import "time"

// service реализует бизнес-логику аутентификации и управления пользователями.
type service struct {
	userRepo    UserRepository
	sessionRepo SessionRepository
	sessionTTL  time.Duration
	bcryptCost  int
}

// New создаёт сервис IAM.
func New(userRepo UserRepository, sessionRepo SessionRepository, sessionTTL time.Duration, bcryptCost int) *service {
	return &service{
		userRepo:    userRepo,
		sessionRepo: sessionRepo,
		sessionTTL:  sessionTTL,
		bcryptCost:  bcryptCost,
	}
}
