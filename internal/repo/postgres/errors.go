package postgres

import (
	"database/sql"

	rerr "github.com/FourSigma/alertd/internal/repo/errors"
	"github.com/lib/pq"
)

func HandlePSQLError(e error) error {
	if e == sql.ErrNoRows {
		return rerr.ErrorDoesNotExist
	}

	switch err := e.(type) {
	case *pq.Error:
		switch err.Code.Name() {
		case "unique_violation":
			return rerr.ErrorAlreadyExists
		default:
			return err

		}
	default:
		return err
	}

}
