package postgres

import (
	"context"
	"log"

	"github.com/FourSigma/alertd/internal/core"
	"github.com/FourSigma/alertd/pkg/sqlhelpers"
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

func NewUserRepo() userRepo {
	return userRepo{
		gen: sqlhelpers.NewStmtGenerator(PostgresSchemaPrefix, &core.User{}, core.UserKey{}),
	}
}

func NewTopicRepo() topicRepo {
	return topicRepo{
		gen: sqlhelpers.NewStmtGenerator(PostgresSchemaPrefix, &core.Topic{}, core.TopicKey{}),
	}
}

func NewMessageRepo() messageRepo {
	return messageRepo{
		gen: sqlhelpers.NewStmtGenerator(PostgresSchemaPrefix, &core.Message{}, core.MessageKey{}),
	}
}

func NewTokenRepo() tokenRepo {
	return tokenRepo{
		gen: sqlhelpers.NewStmtGenerator(PostgresSchemaPrefix, &core.Token{}, core.TokenKey{}),
	}
}
