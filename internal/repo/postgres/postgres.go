package postgres

import (
	"context"

	"github.com/FourSigma/alertd/internal/core"
	"github.com/FourSigma/alertd/pkg/sqlhelpers"
	log "github.com/Sirupsen/logrus"
	"github.com/jmoiron/sqlx"
	_ "github.com/lib/pq"
)

const (
	PostgresSchemaPrefix = "alertd"
)

var sqlDB *sqlx.DB

func init() {
	var err error
	if sqlDB, err = sqlx.Connect("postgres", "user=sivawork dbname=sivawork sslmode=disable"); err != nil {
		log.Fatal(err)
	}
}

func AddRepoContext(ctx context.Context) context.Context {
	return context.WithValue(ctx, sqlhelpers.CtxDbKey, sqlDB)
}

func NewUserRepo(log *log.Logger) userRepo {
	gen := sqlhelpers.NewStmtGenerator(PostgresSchemaPrefix, &core.User{}, core.UserKey{})
	return userRepo{
		crud: sqlhelpers.NewCRUD(log, gen, HandlePSQLError),
	}
}

func NewTopicRepo(log *log.Logger) topicRepo {
	gen := sqlhelpers.NewStmtGenerator(PostgresSchemaPrefix, &core.Topic{}, core.TopicKey{})
	return topicRepo{
		crud: sqlhelpers.NewCRUD(log, gen, HandlePSQLError),
	}
}

func NewMessageRepo(log *log.Logger) messageRepo {
	gen := sqlhelpers.NewStmtGenerator(PostgresSchemaPrefix, &core.Message{}, core.MessageKey{})
	return messageRepo{
		crud: sqlhelpers.NewCRUD(log, gen, HandlePSQLError),
	}
}

func NewTokenRepo(log *log.Logger) tokenRepo {
	gen := sqlhelpers.NewStmtGenerator(PostgresSchemaPrefix, &core.Token{}, core.TokenKey{})
	return tokenRepo{
		crud: sqlhelpers.NewCRUD(log, gen, HandlePSQLError),
	}
}
