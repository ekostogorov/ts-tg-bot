package feedback_repo

import (
	"context"
	db_helpers "ts-tg-bot/internal/helpers/db"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type Repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) Create(ctx context.Context, text string) error {
	builder := db_helpers.Psql.
		Insert("feedback").
		Columns("text").
		Values(text)

	query, bound, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build SQL")
	}

	if _, err = r.db.ExecContext(ctx, query, bound...); err != nil {
		return errors.Wrap(err, "failed to exec SQL")
	}

	return nil
}
