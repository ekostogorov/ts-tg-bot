package lecture

import (
	"context"
	"ts-tg-bot/internal/models"

	"github.com/pkg/errors"
)

type ILectureRepo interface {
	List(ctx context.Context) ([]*models.Lecture, error)
	GetByID(ctx context.Context, id int64) (*models.Lecture, error)
}

type IFileRepo interface {
	Open(filepath string) ([]byte, error)
}

type Service struct {
	lectureRepo ILectureRepo
	fileRepo    IFileRepo
}

func New(lectureRepo ILectureRepo, fileRepo IFileRepo) *Service {
	return &Service{
		lectureRepo: lectureRepo,
		fileRepo:    fileRepo,
	}
}

func (s *Service) List(ctx context.Context) ([]*models.Lecture, error) {
	return s.lectureRepo.List(ctx)
}

func (s *Service) GetLectureWithFile(ctx context.Context, id int64) (*models.LectureWithFile, error) {
	lecture, err := s.lectureRepo.GetByID(ctx, id)
	if err != nil {
		return nil, errors.Wrap(err, "failed to load lection")
	}

	file, err := s.fileRepo.Open(lecture.Filepath)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to open file (%s)", lecture.Filepath)
	}

	return &models.LectureWithFile{
		Lecture: *lecture,
		File:    file,
	}, nil
}
