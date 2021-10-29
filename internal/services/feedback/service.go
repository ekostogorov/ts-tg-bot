package feedback

import "context"

type IFeedbackRepo interface {
	Create(ctx context.Context, text string) error
}

type Service struct {
	feedbackRepo IFeedbackRepo
}

func New(feedbackRepo IFeedbackRepo) *Service {
	return &Service{
		feedbackRepo: feedbackRepo,
	}
}

func (s *Service) Create(ctx context.Context, text string) error {
	return s.feedbackRepo.Create(ctx, text)
}
