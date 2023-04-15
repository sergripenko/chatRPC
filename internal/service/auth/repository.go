package auth

import (
	"context"

	"github.com/sergripenko/chatRPC/internal/domain"
)

type RepositoryProvider interface {
	GetUser(ctx context.Context, username string) (*domain.User, error)
}
