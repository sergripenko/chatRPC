package mem

import (
	"sync"
)

type Users map[string]*User

type InMemoryRepositoryService struct {
	users  Users            // map of users username:User
	groups map[string]Users // map of groups to Users
	mut    sync.RWMutex     // mutex lock
}

func NewInMemoryRepositoryService() *InMemoryRepositoryService {
	return &InMemoryRepositoryService{
		users:  Users{},
		groups: map[string]Users{},
	}
}
