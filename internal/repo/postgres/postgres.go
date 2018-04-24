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
	gen := sqlhelpers.NewStmtGenerator(PostgresSchemaPrefix, &core.User{}, core.UserKey{})
	return userRepo{
		crud: sqlhelpers.NewCRUD(gen, HandlePSQLError),
	}
}

func NewTopicRepo() topicRepo {
	gen := sqlhelpers.NewStmtGenerator(PostgresSchemaPrefix, &core.Topic{}, core.TopicKey{})
	return topicRepo{
		crud: sqlhelpers.NewCRUD(gen, HandlePSQLError),
	}
}

func NewMessageRepo() messageRepo {
	gen := sqlhelpers.NewStmtGenerator(PostgresSchemaPrefix, &core.Message{}, core.MessageKey{})
	return messageRepo{
		crud: sqlhelpers.NewCRUD(gen, HandlePSQLError),
	}
}

func NewTokenRepo() tokenRepo {
	gen := sqlhelpers.NewStmtGenerator(PostgresSchemaPrefix, &core.Token{}, core.TokenKey{})
	return tokenRepo{
		crud: sqlhelpers.NewCRUD(gen, HandlePSQLError),
	}
}
