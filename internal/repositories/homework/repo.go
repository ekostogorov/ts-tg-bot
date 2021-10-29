package homework_repo

import (
	"context"
	"database/sql"
	db_helpers "ts-tg-bot/internal/helpers/db"
	"ts-tg-bot/internal/models"
	"ts-tg-bot/internal/types"
	"time"

	"github.com/pkg/errors"

	"github.com/jmoiron/sqlx"
)

type homework struct {
	ID        int64     `db:"id"`
	Name      string    `db:"name"`
	LectureID int64     `db:"lecture_id"`
	Filepath  string    `db:"file_path"`
	IsActive  bool      `db:"is_active"`
	ExpiresAt time.Time `db:"expires_at"`
	CreatedAt time.Time `db:"created_at"`
}

func (h *homework) toModel() *models.Homework {
	return &models.Homework{
		ID:        h.ID,
		Name:      h.Name,
		LectureID: h.LectureID,
		Filepath:  h.Filepath,
		IsActive:  h.IsActive,
		ExpiresAt: h.ExpiresAt,
		CreatedAt: h.CreatedAt,
	}
}

func (h *homework) getSelectColumns() []string {
	return []string{
		"id",
		"name",
		"lecture_id",
		"file_path",
		"expires_at",
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

func (r *Repo) GetByLectureID(ctx context.Context, lectureID int64) ([]*models.Homework, error) {
	nilModel := &homework{}
	builder := db_helpers.Psql.
		Select(nilModel.getSelectColumns()...).
		From("homeworks").
		Where("lecture_id = ?", lectureID)

	query, bound, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build SQL")
	}

	var dest []homework
	if err = r.db.SelectContext(ctx, &dest, query, bound...); err != nil {
		return nil, errors.Wrap(err, "failed to exec SQL")
	}

	output := make([]*models.Homework, 0, len(dest))
	for _, hw := range dest {
		output = append(output, hw.toModel())
	}

	return output, nil
}

func (r *Repo) GetActive(ctx context.Context) (*models.Homework, error) {
	nilModel := &homework{}
	builder := db_helpers.Psql.
		Select(nilModel.getSelectColumns()...).
		From("homeworks").
		Where("is_active IS TRUE")

	query, bound, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build SQL")
	}

	var dest homework
	if err = r.db.GetContext(ctx, &dest, query, bound...); err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrNotFound
		}

		return nil, errors.Wrap(err, "failed to exec SQL")
	}

	return dest.toModel(), nil
}

func (r *Repo) List(ctx context.Context) ([]*models.Homework, error) {
	nilModel := &homework{}
	builder := db_helpers.Psql.
		Select(nilModel.getSelectColumns()...).
		From("homeworks")

	query, bound, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrap(err, "failed to build SQL")
	}

	var dest []homework
	if err = r.db.SelectContext(ctx, &dest, query, bound...); err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrNotFound
		}

		return nil, errors.Wrap(err, "failed to exec SQL")
	}

	output := make([]*models.Homework, 0, len(dest))
	for _, hw := range dest {
		output = append(output, hw.toModel())
	}

	return output, nil
}
