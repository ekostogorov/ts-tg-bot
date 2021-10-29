package lectures_repo

import (
	"context"
	db_helpers "ts-tg-bot/internal/helpers/db"
	"ts-tg-bot/internal/models"
	"time"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
)

type lecture struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	Filepath  string    `db:"file_path"`
	CreatedAt time.Time `db:"created_at"`
}

func (l *lecture) toModel() *models.Lecture {
	return &models.Lecture{
		ID:        l.ID,
		Name:      l.Name,
		Filepath:  l.Filepath,
		CreatedAt: l.CreatedAt,
	}
}

func (l *lecture) getSelectColumns() []string {
	return []string{
		"id",
		"name",
		"file_path",
		"created_at",
	}
}

type Repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) GetByID(ctx context.Context, id int64) (*models.Lecture, error) {
	nilModel := &lecture{}
	builder := db_helpers.Psql.
		Select(nilModel.getSelectColumns()...).
		From("lectures").
		Where("id = ?", id)

	query, bound, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build SQL")
	}

	var dest lecture
	if err = r.db.GetContext(ctx, &dest, query, bound...); err != nil {
		return nil, errors.Wrap(err, "failed to exec SQL")
	}

	return dest.toModel(), nil
}

func (r *Repo) List(ctx context.Context) ([]*models.Lecture, error) {
	nilModel := &lecture{}
	builder := db_helpers.Psql.
		Select(nilModel.getSelectColumns()...).
		From("lectures").
		OrderBy("id")

	query, bound, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build SQL")
	}

	var dest []lecture
	if err = r.db.SelectContext(ctx, &dest, query, bound...); err != nil {
		return nil, errors.Wrap(err, "failed to exec SQL")
	}

	output := make([]*models.Lecture, 0, len(dest))
	for _, lect := range dest {
		output = append(output, lect.toModel())
	}

	return output, nil
}
