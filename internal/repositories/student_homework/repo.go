package student_homework_repo

import (
	"context"
	"database/sql"
	db_helpers "ts-tg-bot/internal/helpers/db"
	"ts-tg-bot/internal/models"

	"github.com/jmoiron/sqlx"
	"github.com/pkg/errors"
)

type studentHomework struct {
	StudentID  int64         `db:"student_id"`
	HomeworkID int64         `db:"homework_id"`
	PassedAt   sql.NullTime  `db:"passed_at"`
	Filepath   string        `db:"file_path"`
	Grade      sql.NullInt64 `db:"grade"`
}

func newStudentHomework(m models.StudentHomework) *studentHomework {
	return &studentHomework{
		StudentID:  m.StudentID,
		HomeworkID: m.HomeworkID,
		PassedAt: sql.NullTime{
			Time:  m.PassedAt,
			Valid: !m.PassedAt.IsZero(),
		},
		Filepath: m.Filepath,
		Grade: sql.NullInt64{
			Int64: m.Grade,
			Valid: m.Grade != 0,
		},
	}
}

func (h *studentHomework) toModel() *models.StudentHomework {
	return &models.StudentHomework{
		StudentID:  h.StudentID,
		HomeworkID: h.StudentID,
		PassedAt:   h.PassedAt.Time,
		Filepath:   h.Filepath,
		Grade:      h.Grade.Int64,
	}
}

func (h *studentHomework) getSelectColumns() []string {
	return []string{
		"student_id",
		"homework_id",
		"passed_at",
		"file_path",
		"grade",
	}
}

func (h *studentHomework) getInsertColumnsValues() ([]string, []interface{}) {
	columns := []string{
		"student_id",
		"homework_id",
		"file_path",
		"passed_at",
	}

	values := []interface{}{
		h.StudentID,
		h.HomeworkID,
		h.Filepath,
		h.PassedAt,
	}

	return columns, values
}

type Repo struct {
	db *sqlx.DB
}

func New(db *sqlx.DB) *Repo {
	return &Repo{
		db: db,
	}
}

func (r *Repo) Create(ctx context.Context, payload models.StudentHomework) error {
	m := newStudentHomework(payload)
	columns, values := m.getInsertColumnsValues()

	builder := db_helpers.Psql.
		Insert("student_homeworks").
		Columns(columns...).
		Values(values...)

	query, bound, err := builder.ToSql()
	if err != nil {
		return errors.Wrap(err, "failed to build SQL")
	}

	if _, err = r.db.ExecContext(ctx, query, bound...); err != nil {
		return errors.Wrap(err, "failed to exec SQL")
	}

	return nil
}
