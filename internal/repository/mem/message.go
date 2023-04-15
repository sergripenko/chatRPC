package mem

import (
	"context"
	"strings"

	"github.com/sergripenko/chatRPC/internal/domain"
)

func (r *InMemoryRepositoryService) AddMessage(ctx context.Context, mess *domain.Message) (*domain.Message, error) {
	if strings.EqualFold(mess.Type, domain.GroupMessageType) {
		//	publish the message to group.
		r.mut.Lock()
		bTopics := r.groups[mess.Receiver]
		r.mut.Unlock()

		for _, s := range bTopics {
			if !s.active {
				continue
			}
			go (func(s *User) {
				if r.users[s.username].active {
					r.users[s.username].messages <- mess
				}
			})(s)
		}
	} else {
		//	direct message
		if r.users[mess.Receiver].active {
			r.users[mess.Receiver].messages <- mess
		}
	}
	return mess, nil
}
