package chat

import (
	"context"

	"github.com/sergripenko/chatRPC/internal/domain"
)

type UserRepositoryProvider interface {
	IsUserExist(ctx context.Context, username string) (bool, error)
	AddUser(ctx context.Context, user *domain.User) (*domain.User, error)
	GetUser(ctx context.Context, username string) (*domain.User, error)
	GetUserMessages(ctx context.Context, user *domain.User) (<-chan *domain.Message, error)
	GetUserGroups(ctx context.Context, user *domain.User) ([]*domain.Group, error)
	IsUserInGroup(ctx context.Context, user *domain.User, group *domain.Group) (bool, error)
	DeleteUser(ctx context.Context, user *domain.User) error
}

type GroupRepositoryProvider interface {
	IsGroupExist(ctx context.Context, name string) (bool, error)
	CreateGroup(ctx context.Context, user *domain.User, group *domain.Group) (*domain.Group, error)
	JoinGroup(ctx context.Context, user *domain.User, group *domain.Group) (*domain.Group, error)
	LeaveGroup(ctx context.Context, user *domain.User, group *domain.Group) error
	GetGroupUsers(ctx context.Context, group *domain.Group) ([]*domain.User, error)
	GetAllGroups(ctx context.Context) ([]*domain.Group, error)
	DeleteGroup(ctx context.Context, group *domain.Group) error
}

type MessageRepositoryProvider interface {
	AddMessage(ctx context.Context, mess *domain.Message) (*domain.Message, error)
}
