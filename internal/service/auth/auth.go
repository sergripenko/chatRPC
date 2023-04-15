package auth

import (
	"context"

	"github.com/sergripenko/chatRPC/internal/domain"
)

type AuthService struct {
	repo RepositoryProvider
}

func NewAuthService(repo RepositoryProvider) *AuthService {
	return &AuthService{repo: repo}
}

func (s *AuthService) Authorize(ctx context.Context, user *domain.User) (*domain.User, error) {
	return s.repo.GetUser(ctx, user.Username)
}
