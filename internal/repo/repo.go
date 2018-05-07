package repo

import (
	"github.com/FourSigma/alertd/internal/core"

	"github.com/FourSigma/alertd/internal/repo/postgres"
	log "github.com/Sirupsen/logrus"
)

type Datastore struct {
	User struct {
		core.UserRepo
		Reslove struct {
			Tokens func(ls core.UserList, opts ...core.Opts) error
			Topics func(ls core.UserList, opts ...core.Opts) error
		}
	}
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
		Token:   postgres.NewTokenRepo(l),
		Message: postgres.NewMessageRepo(l),
		Topic:   postgres.NewTopicRepo(l),
	}
	glbDS.User.UserRepo = postgres.NewUserRepo(l)

	return *glbDS
}
