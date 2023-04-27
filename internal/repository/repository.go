package repository

import (
	"database/sql"

	"github.com/Masterminds/squirrel"
	"go.uber.org/zap"
)

type (
	Repository struct {
		db *sql.DB
		sq squirrel.StatementBuilderType
		logger *zap.SugaredLogger
	}
)

func NewRepository(db *sql.DB, logger *zap.SugaredLogger) *Repository {
	return &Repository{
		db: db,
		sq: squirrel.StatementBuilder.PlaceholderFormat(squirrel.Dollar).RunWith(db),
		logger: logger,
	}
}


