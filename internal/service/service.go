package service

import (
	"github.com/FourSigma/alertd/internal/repo/postgres"
	log "github.com/Sirupsen/logrus"
)

func NewUserService(l *log.Logger) *UserService {
	return &UserService{
		usrRepo: postgres.NewUserRepo(l),
	}
}

func NewMessageService(l *log.Logger) *MessageService {
	return &MessageService{
		msgRepo: postgres.NewMessageRepo(l),
	}
}

func NewTopicService(l *log.Logger) *TopicService {
	return &TopicService{
		tpRepo: postgres.NewTopicRepo(l),
	}
}

func NewTokenService(l *log.Logger) *TokenService {
	return &TokenService{
		tknRepo: postgres.NewTokenRepo(l),
	}
}
