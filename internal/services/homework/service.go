package homework

import (
	"context"
	"fmt"
	uid_helpers "ts-tg-bot/internal/helpers/uid"
	"ts-tg-bot/internal/models"
	"ts-tg-bot/internal/types"
	"log"
	"time"

	"github.com/pkg/errors"
)

type IHomeworkRepo interface {
	GetByLectureID(ctx context.Context, lectureID int64) ([]*models.Homework, error)
	GetActive(ctx context.Context) (*models.Homework, error)
	List(ctx context.Context) ([]*models.Homework, error)
}

type IStudentHomeworkRepo interface {
	Create(ctx context.Context, payload models.StudentHomework) error
}

type IFileRepo interface {
	Open(filepath string) ([]byte, error)
	Save(data []byte, filepath string) error
}

type Service struct {
	hwRepo        IHomeworkRepo
	studentHWRepo IStudentHomeworkRepo
	fileRepo      IFileRepo
}

func New(hwRepo IHomeworkRepo, studentHWRepo IStudentHomeworkRepo, fileRepo IFileRepo) *Service {
	return &Service{
		hwRepo:        hwRepo,
		studentHWRepo: studentHWRepo,
		fileRepo:      fileRepo,
	}
}

func (s *Service) GetActive(ctx context.Context) (*models.HomeWorkWithFile, error) {
	hw, err := s.hwRepo.GetActive(ctx)
	if err != nil {
		return nil, err
	}

	file, err := s.fileRepo.Open(hw.Filepath)
	if err != nil {
		return nil, err
	}

	return &models.HomeWorkWithFile{
		Homework: *hw,
		File:     file,
	}, nil
}

func (s *Service) List(ctx context.Context) ([]*models.HomeWorkWithFile, error) {
	hws, err := s.hwRepo.List(ctx)
	if err != nil {
		return nil, err
	}

	output := make([]*models.HomeWorkWithFile, 0, len(hws))
	for _, hw := range hws {
		file, err := s.fileRepo.Open(hw.Filepath)
		if err != nil {
			return nil, err
		}

		cmplx := &models.HomeWorkWithFile{
			Homework: *hw,
			File:     file,
		}

		output = append(output, cmplx)
	}

	return output, nil
}

func (s *Service) Pass(ctx context.Context, file []byte, fileType string, student *models.Student) error {
	active, err := s.GetActive(ctx)
	if err != nil {
		return err
	}

	if active.ExpiresAt.Before(time.Now().UTC()) {
		return types.ErrHWExpired
	}

	fileName, err := uid_helpers.GenerateUUIDv4()
	if err != nil {
		return errors.Wrap(err, "failed to generate filename")
	}

	studentHW := models.StudentHomework{
		StudentID:  student.ID,
		HomeworkID: active.ID,
		Filepath:   fmt.Sprintf("%s%s.ipynb", student.Folder, fileName),
		PassedAt:   time.Now().UTC(),
	}

	log.Printf("Saving file. Name: %s. Type: %s", studentHW.Filepath, fileType)

	if err := s.fileRepo.Save(file, studentHW.Filepath); err != nil {
		return errors.Wrap(err, "failed to save file")
	}

	if err := s.studentHWRepo.Create(ctx, studentHW); err != nil {
		return errors.Wrap(err, "failed to save student homework to db")
	}

	return nil
}
