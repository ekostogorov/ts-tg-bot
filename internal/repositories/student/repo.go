package student_repo

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

type student struct {
	ID             int64          `db:"id"`
	Name           string         `db:"name"`
	Login          string         `db:"login"`
	Folder         sql.NullString `db:"folder"`
	TelegramUserID sql.NullString `db:"telegram_user_id"`
	IsActivated    bool           `db:"is_activated"`
	ActivatedAt    sql.NullTime   `db:"activated_at"`
	CreatedAt      time.Time      `db:"created_at"`
}

func (s *student) toModel() *models.Student {
	return &models.Student{
		ID:             s.ID,
		Name:           s.Name,
		Login:          s.Login,
		Folder:         s.Folder.String,
		TelegramUserID: s.TelegramUserID.String,
		IsActivated:    s.IsActivated,
		ActivatedAt:    s.ActivatedAt.Time,
		CreatedAt:      s.CreatedAt,
	}
}

func (s *student) getSelectColumns() []string {
	return []string{
		"id",
		"name",
		"folder",
		"telegram_user_id",
		"is_activated",
		"activated_at",
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

func (r *Repo) GetByTelegramUserID(ctx context.Context, tgUserID string) (*models.Student, error) {
	nilModel := &student{}
	builder := db_helpers.Psql.
		Select(nilModel.getSelectColumns()...).
		From("students").
		Where("telegram_user_id = ?", tgUserID)

	query, bound, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to build SQL")
	}

	var dest student
	if err = r.db.GetContext(ctx, &dest, query, bound...); err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrNotFound
		}

		return nil, errors.Wrapf(err, "failed to exec SQL")
	}

	return dest.toModel(), nil
}

func (r *Repo) Activate(ctx context.Context, login, tgUserID string) error {
	builder := db_helpers.Psql.
		Update("students").
		Set("telegram_user_id", tgUserID).
		Set("is_activated", true).
		Where("login = ?", login)

	query, bound, err := builder.ToSql()
	if err != nil {
		return errors.Wrapf(err, "failed to build SQL")
	}

	if _, err = r.db.ExecContext(ctx, query, bound...); err != nil {
		return errors.Wrapf(err, "failed to exec SQL")
	}

	return nil
}

func (r *Repo) GetByLogin(ctx context.Context, login string) (*models.Student, error) {
	nilModel := &student{}
	builder := db_helpers.Psql.
		Select(nilModel.getSelectColumns()...).
		From("students").
		Where("login = ?", login)

	query, bound, err := builder.ToSql()
	if err != nil {
		return nil, errors.Wrapf(err, "failed to build SQL")
	}

	var dest student
	if err = r.db.GetContext(ctx, &dest, query, bound...); err != nil {
		if err == sql.ErrNoRows {
			return nil, types.ErrNotFound
		}

		return nil, errors.Wrapf(err, "failed to exec SQL")
	}

	return dest.toModel(), nil
}
