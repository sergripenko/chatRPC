package mem

import (
	"context"

	"github.com/sergripenko/chatRPC/internal/domain"
)

func (r *InMemoryRepositoryService) IsGroupExist(ctx context.Context, name string) (bool, error) {
	r.mut.RLock()
	defer r.mut.RUnlock()

	_, exist := r.groups[name]
	return exist, nil
}

func (r *InMemoryRepositoryService) CreateGroup(ctx context.Context, user *domain.User, group *domain.Group) (
	*domain.Group, error) {
	r.mut.Lock()
	defer r.mut.Unlock()

	r.groups[group.Name] = Users{}
	return group, nil
}

func (r *InMemoryRepositoryService) JoinGroup(ctx context.Context, user *domain.User, group *domain.Group) (
	*domain.Group, error) {
	r.mut.Lock()
	defer r.mut.Unlock()

	r.users[user.Username].groups[group.Name] = true
	r.groups[group.Name][user.Username] = r.users[user.Username]
	return group, nil
}

func (r *InMemoryRepositoryService) LeaveGroup(ctx context.Context, user *domain.User, group *domain.Group) error {
	r.mut.Lock()
	defer r.mut.Unlock()
	delete(r.groups[group.Name], user.Username)
	delete(r.users[user.Username].groups, group.Name)
	return nil
}

func (r *InMemoryRepositoryService) DeleteGroup(ctx context.Context, group *domain.Group) error {
	r.mut.Lock()
	defer r.mut.Unlock()

	delete(r.groups, group.Name)
	return nil
}

func (r *InMemoryRepositoryService) GetGroupUsers(ctx context.Context, group *domain.Group) ([]*domain.User, error) {
	r.mut.RLock()
	defer r.mut.RUnlock()

	var users []*domain.User

	for k, _ := range r.groups[group.Name] {
		users = append(users, &domain.User{Username: k})
	}
	return users, nil
}

func (r *InMemoryRepositoryService) GetAllGroups(ctx context.Context) ([]*domain.Group, error) {
	r.mut.RLock()
	defer r.mut.RUnlock()

	var groups []*domain.Group

	for groupName, users := range r.groups {
		tmpGroup := &domain.Group{Name: groupName}

		for userName, _ := range users {
			tmpGroup.Users = append(tmpGroup.Users, &domain.User{Username: userName})
		}
		groups = append(groups, tmpGroup)
	}
	return groups, nil
}
