package student

import (
	"context"
	"ts-tg-bot/internal/models"
	"ts-tg-bot/internal/types"
	"log"

	_ "github.com/jackc/pgx"
	"github.com/pkg/errors"
)

type IStudentRepo interface {
	GetByTelegramUserID(ctx context.Context, tgUserID string) (*models.Student, error)
	Activate(ctx context.Context, login, tgUserID string) error
	GetByLogin(ctx context.Context, login string) (*models.Student, error)
}

type Service struct {
	studentRepo IStudentRepo
}

func New(studentRepo IStudentRepo) *Service {
	return &Service{
		studentRepo: studentRepo,
	}
}

func (s *Service) Register(ctx context.Context, tgUserID, email string) (*models.Student, error) {
	student, err := s.studentRepo.GetByLogin(ctx, email)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get user by login")
	}
	if student.IsActivated {
		return nil, types.ErrUserExists
	}

	if err = s.studentRepo.Activate(ctx, email, tgUserID); err != nil {
		return nil, errors.Wrap(err, "faield to activate user")
	}

	return s.studentRepo.GetByTelegramUserID(ctx, tgUserID)
}

func (s *Service) IsRegistered(ctx context.Context, tgUserID string) (bool, error) {
	student, err := s.studentRepo.GetByTelegramUserID(ctx, tgUserID)
	if err != nil {
		if err == types.ErrNotFound {
			return false, nil
		}

		return false, errors.Wrapf(err, "failed to get user by tg user ID (%v)", tgUserID)
	}

	log.Printf("STUDENT: %+v", student)

	if student.IsActivated {
		return true, nil
	}

	return false, nil
}

func (s *Service) GetByTelegramUserID(ctx context.Context, tgUserID string) (*models.Student, error) {
	return s.studentRepo.GetByTelegramUserID(ctx, tgUserID)
}
