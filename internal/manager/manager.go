package manager

import (
	"context"
	"ts-tg-bot/internal/models"
	"io/ioutil"
	"log"
	"net/http"
	"sync"
	"time"

	"github.com/pkg/errors"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
)

const (
	separator = "------------"
)

type IStudentService interface {
	Register(ctx context.Context, tgUserID, email string) (*models.Student, error)
	IsRegistered(ctx context.Context, tgUserID string) (bool, error)
	GetByTelegramUserID(ctx context.Context, tgUserID string) (*models.Student, error)
}

type ILectureService interface {
	List(ctx context.Context) ([]*models.Lecture, error)
	GetLectureWithFile(ctx context.Context, id int64) (*models.LectureWithFile, error)
}

type IHomeworkService interface {
	GetActive(ctx context.Context) (*models.HomeWorkWithFile, error)
	List(ctx context.Context) ([]*models.HomeWorkWithFile, error)
	Pass(ctx context.Context, file []byte, fileType string, student *models.Student) error
}

type IFeedbackService interface {
	Create(ctx context.Context, text string) error
}

type Manager struct {
	ctx    context.Context
	cancel context.CancelFunc

	studentService  IStudentService
	lectureService  ILectureService
	hwService       IHomeworkService
	feedbackService IFeedbackService

	bot       *tgbotapi.BotAPI
	updatesCh tgbotapi.UpdatesChannel

	activeRegistrations sync.Map // map[tgUserID]struct{}
}

func New(
	ctx context.Context,
	studentSrv IStudentService,
	lectureSrv ILectureService,
	hwSrv IHomeworkService,
	feedbackSrv IFeedbackService,
	bot *tgbotapi.BotAPI) *Manager {
	innerCtx, cancel := context.WithCancel(ctx)

	return &Manager{
		ctx:    innerCtx,
		cancel: cancel,

		studentService:  studentSrv,
		lectureService:  lectureSrv,
		hwService:       hwSrv,
		feedbackService: feedbackSrv,

		bot: bot,
	}
}

func (m *Manager) Init() error {
	ch, err := m.bot.GetUpdatesChan(tgbotapi.UpdateConfig{})
	if err != nil {
		return errors.Wrap(err, "failed to get updates channel")
	}

	m.updatesCh = ch
	return nil
}

func (m *Manager) Start() {
	go m.listen()
}

func (m *Manager) Exit() {
	m.cancel()
}

func (m *Manager) listen() {
	log.Print("starting to listen...")
	for {
		select {
		case update := <-m.updatesCh:
			log.Printf("received update: %+v", update)
			go m.handleUpdate(m.ctx, &update)

		case <-m.ctx.Done():
			log.Print("exiting due to context cancellation")
			return
		}
	}
}

func (m *Manager) sendTextMessage(chatID int64, text string) error {
	msg := tgbotapi.NewMessage(chatID, text)
	if _, err := m.bot.Send(msg); err != nil {
		return err
	}

	return nil
}

func (m *Manager) sendFile(chatID int64, file []byte, name string) error {
	f := tgbotapi.FileBytes{
		Name:  name,
		Bytes: file,
	}

	document := tgbotapi.NewDocumentUpload(chatID, f)
	if _, err := m.bot.Send(document); err != nil {
		return err
	}

	return nil
}

func (m *Manager) getFile(msg *tgbotapi.Message) ([]byte, string, error) {
	fileID := msg.Document.FileID
	fileType := msg.Document.MimeType

	url, err := m.bot.GetFileDirectURL(fileID)
	if err != nil {
		return nil, fileType, errors.Wrap(err, "failed to download file")
	}

	req, err := http.NewRequest(http.MethodGet, url, nil)
	if err != nil {
		return nil, fileType, errors.Wrap(err, "failed to create request")
	}

	client := &http.Client{
		Timeout: time.Second * 30,
	}

	resp, err := client.Do(req)
	if err != nil {
		return nil, fileType, errors.Wrap(err, "failed to exec request")
	}
	defer resp.Body.Close()

	buf, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, fileType, errors.Wrap(err, "failed to read file")
	}

	return buf, fileType, nil
}
