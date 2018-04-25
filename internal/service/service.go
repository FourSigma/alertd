package service

import (
	"github.com/FourSigma/alertd/internal/repo"
	log "github.com/Sirupsen/logrus"
)

var glbSvc *ServiceManager

func GetService(l *log.Logger) ServiceManager {
	if glbSvc != nil {
		return *glbSvc
	}

	svc := &ServiceManager{
		User:    newUserService(l),
		Topic:   newTopicService(l),
		Message: newMessageService(l),
		Token:   newTokenService(l),
	}

	return *svc
}

type ServiceManager struct {
	User    UserService
	Topic   TopicService
	Message MessageService
	Token   TokenService
}

func newUserService(l *log.Logger) UserService {
	return UserService{
		repo: repo.GetDatastore(l),
	}
}

func newMessageService(l *log.Logger) MessageService {
	return MessageService{
		repo: repo.GetDatastore(l),
	}
}

func newTopicService(l *log.Logger) TopicService {
	return TopicService{
		repo: repo.GetDatastore(l),
	}
}

func newTokenService(l *log.Logger) TokenService {
	return TokenService{
		repo: repo.GetDatastore(l),
	}
}
