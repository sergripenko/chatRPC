package chat

import (
	"context"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/sergripenko/chatRPC/internal/domain"
)

type ChatService struct {
	userRepo  UserRepositoryProvider
	groupRepo GroupRepositoryProvider
	messRepo  MessageRepositoryProvider
}

func NewChatService(
	userRepo UserRepositoryProvider,
	groupRepo GroupRepositoryProvider,
	messRepo MessageRepositoryProvider,
) *ChatService {
	return &ChatService{
		userRepo:  userRepo,
		groupRepo: groupRepo,
		messRepo:  messRepo,
	}
}

var (
	errUserAlreadyExist            = errors.New("user with this username already exist")
	errGroupDoesNotExist           = errors.New("group does not exist")
	errGroupAlreadyExist           = errors.New("group already exist")
	errMessageReceiverDoesNotExist = errors.New("message receiver does not exist")
	errUserNotInGroup              = errors.New("user not in group")
)

func (s *ChatService) Connect(ctx context.Context, user *domain.User) (<-chan *domain.Message, error) {
	// check if user already exist
	exist, err := s.userRepo.IsUserExist(ctx, user.Username)
	if err != nil {
		return nil, err
	}
	if exist {
		return nil, errUserAlreadyExist
	}
	user, err = s.userRepo.AddUser(ctx, user)
	if err != nil {
		return nil, err
	}
	logrus.Infof("username: %s connected", user.Username)
	return s.userRepo.GetUserMessages(ctx, user)
}

func (s *ChatService) Disconnect(ctx context.Context, user *domain.User) error {
	logrus.Infof("%s disconnected", user.Username)

	userGroups, err := s.userRepo.GetUserGroups(ctx, user)
	if err != nil {
		return err
	}
	for _, group := range userGroups {
		logrus.Infof("user %s leave group %s", user.Username, group.Name)
		if err = s.groupRepo.LeaveGroup(ctx, user, group); err != nil {
			return err
		}
		users, err := s.groupRepo.GetGroupUsers(ctx, group)
		if err != nil {
			return err
		}
		if len(users) == 0 {
			logrus.Infof("group %s deleted", group.Name)
			if err = s.groupRepo.DeleteGroup(ctx, group); err != nil {
				return err
			}
		}
	}
	return s.userRepo.DeleteUser(ctx, user)
}

func (s *ChatService) JoinGroupChat(ctx context.Context, user *domain.User, group *domain.Group) error {
	exist, err := s.groupRepo.IsGroupExist(ctx, group.Name)
	if err != nil {
		return err
	}
	if !exist {
		return errGroupDoesNotExist
	}
	if _, err = s.groupRepo.JoinGroup(ctx, user, group); err != nil {
		return err
	}
	logrus.Infof("user %s joined group %s", user.Username, group.Name)
	return nil
}

func (s *ChatService) LeaveGroupChat(ctx context.Context, user *domain.User, group *domain.Group) error {
	exist, err := s.groupRepo.IsGroupExist(ctx, group.Name)
	if err != nil {
		return err
	}
	if !exist {
		return errGroupDoesNotExist
	}
	if err = s.groupRepo.LeaveGroup(ctx, user, group); err != nil {
		return err
	}
	logrus.Infof("user %s leave group %s", user.Username, group.Name)
	users, err := s.groupRepo.GetGroupUsers(ctx, group)
	if err != nil {
		return err
	}
	if len(users) == 0 {
		logrus.Infof("group %s deleted", group.Name)
		if err = s.groupRepo.DeleteGroup(ctx, group); err != nil {
			return err
		}
	}
	return nil
}

func (s *ChatService) CreateGroupChat(ctx context.Context, user *domain.User, group *domain.Group) error {
	exist, err := s.groupRepo.IsGroupExist(ctx, group.Name)
	if err != nil {
		return err
	}
	if exist {
		return errGroupAlreadyExist
	}
	group, err = s.groupRepo.CreateGroup(ctx, user, group)
	if err != nil {
		return err
	}
	logrus.Infof("user %s created group %s", user.Username, group.Name)
	group, err = s.groupRepo.JoinGroup(ctx, user, group)
	if err != nil {
		return err
	}
	logrus.Infof("user %s joined group %s", user.Username, group.Name)
	return nil
}

func (s *ChatService) SendMessage(ctx context.Context, user *domain.User, mess *domain.Message) error {
	//	check if receiver exist
	if strings.EqualFold(mess.Type, domain.DirectMessageType) {
		exist, err := s.userRepo.IsUserExist(ctx, mess.Receiver)
		if err != nil {
			return err
		}
		if !exist {
			return errMessageReceiverDoesNotExist
		}
	} else {
		exist, err := s.groupRepo.IsGroupExist(ctx, mess.Group)
		if err != nil {
			return err
		}
		if !exist {
			return errMessageReceiverDoesNotExist
		}
		isUserInGroup, err := s.userRepo.IsUserInGroup(ctx, user, &domain.Group{Name: mess.Group})
		if err != nil {
			return err
		}
		if !isUserInGroup {
			return errUserNotInGroup
		}
	}
	_, err := s.messRepo.AddMessage(ctx, mess)
	if err != nil {
		return err
	}
	return nil
}

func (s *ChatService) ListChannels(ctx context.Context) ([]*domain.Group, error) {
	return s.groupRepo.GetAllGroups(ctx)
}
