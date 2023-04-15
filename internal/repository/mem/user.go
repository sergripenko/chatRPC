package mem

import (
	"context"
	"github.com/sergripenko/chatRPC/internal/domain"
	"github.com/sergripenko/chatRPC/internal/repository"
)

type User struct {
	username string
	messages chan *domain.Message // messages channel
	groups   map[string]bool      // groups it is subscribed to.
	active   bool                 // if given subscriber is active
}

func NewUser(username string) *User {
	// returns a new subscriber.
	return &User{
		username: username,
		messages: make(chan *domain.Message),
		groups:   map[string]bool{},
		active:   true,
	}
}

func (r *InMemoryRepositoryService) IsUserExist(ctx context.Context, username string) (bool, error) {
	r.mut.RLock()
	defer r.mut.RUnlock()
	_, exist := r.users[username]
	return exist, nil
}

func (r *InMemoryRepositoryService) GetUser(ctx context.Context, username string) (*domain.User, error) {
	r.mut.RLock()
	defer r.mut.RUnlock()

	userObj, exist := r.users[username]
	if !exist {
		return nil, repository.ErrRecordNotFount
	}
	return &domain.User{Username: userObj.username}, nil
}

func (r *InMemoryRepositoryService) AddUser(ctx context.Context, user *domain.User) (*domain.User, error) {
	r.mut.Lock()
	defer r.mut.Unlock()

	s := NewUser(user.Username)
	r.users[user.Username] = s
	return user, nil
}

func (r *InMemoryRepositoryService) GetUserMessages(ctx context.Context, user *domain.User) (
	<-chan *domain.Message, error) {
	return r.users[user.Username].messages, nil
}

func (r *InMemoryRepositoryService) GetUserGroups(ctx context.Context, user *domain.User) ([]*domain.Group, error) {
	r.mut.RLock()
	defer r.mut.RUnlock()

	var userGroups []*domain.Group
	for name, active := range r.users[user.Username].groups {
		if active {
			userGroups = append(userGroups, &domain.Group{Name: name})
		}
	}
	return userGroups, nil
}

func (r *InMemoryRepositoryService) IsUserInGroup(ctx context.Context, user *domain.User, group *domain.Group) (bool, error) {
	r.mut.RLock()
	defer r.mut.RUnlock()

	_, exist := r.groups[group.Name][user.Username]
	return exist, nil
}

func (r *InMemoryRepositoryService) DeleteUser(ctx context.Context, user *domain.User) error {
	r.mut.Lock()
	defer r.mut.Unlock()

	r.users[user.Username].active = false
	close(r.users[user.Username].messages)
	delete(r.users, user.Username)
	return nil
}
