package service

import "github.com/FourSigma/alertd/internal/repo/postgres"

func NewUserService() *UserService {
	return &UserService{
		usrRepo: postgres.NewUserRepo(),
	}
}

func NewMessageService() *MessageService {
	return &MessageService{
		msgRepo: postgres.NewMessageRepo(),
	}
}

func NewTopicService() *TopicService {
	return &TopicService{
		tpRepo: postgres.NewTopicRepo(),
	}
}

func NewTokenService() *TokenService {
	return &TokenService{
		tknRepo: postgres.NewTokenRepo(),
	}
}
