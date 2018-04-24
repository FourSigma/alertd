package log

import (
	log "github.com/sirupsen/logrus"
)

type FieldKey string

type Logger interface {
}

func init() {
	log.SetFormatter(&log.JSONFormatter{})
}
