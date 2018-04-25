package repo

import (
	"github.com/FourSigma/alertd/internal/core"

	"github.com/FourSigma/alertd/internal/repo/postgres"
	log "github.com/Sirupsen/logrus"
)

type Datastore struct {
	User    core.UserRepo
	Token   core.TokenRepo
	Message core.MessageRepo
	Topic   core.TopicRepo
}

var glbDS *Datastore

func GetDatastore(l *log.Logger) Datastore {
	if glbDS != nil {
		return *glbDS
	}

	glbDS := &Datastore{
		User:    postgres.NewUserRepo(l),
		Token:   postgres.NewTokenRepo(l),
		Message: postgres.NewMessageRepo(l),
		Topic:   postgres.NewTopicRepo(l),
	}

	return *glbDS
}
