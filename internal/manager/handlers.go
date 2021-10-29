package manager

import (
	"context"
	"fmt"
	regexp_helpers "ts-tg-bot/internal/helpers/regexp"
	"ts-tg-bot/internal/types"
	"log"
	"strconv"
	"strings"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	"github.com/pkg/errors"
)

func (m *Manager) handleUpdate(ctx context.Context, update *tgbotapi.Update) {
	if update == nil {
		log.Printf("update is nil")
		return
	}
	if update.Message == nil {
		log.Printf("update.Message is nil")
		return
	}

	if isFeedback(update.Message) {
		if err := m.handleFeedback(ctx, update.Message); err != nil {
			log.Printf("failed to handle review: %s", err.Error())
			return
		}
	}

	tgUserID := strconv.Itoa(update.Message.From.ID)
	isRegistered, err := m.studentService.IsRegistered(ctx, tgUserID)
	if err != nil {
		log.Printf("failed to check if student is registered: %s", err.Error())
		return
	}

	if !isRegistered {
		if err := m.handleRegistration(ctx, update); err != nil {
			log.Printf("failed to reply to message: %s", err.Error())
		}

		return
	}

	switch getMessageType(update.Message) {
	case messageTypeLectures:
		if err := m.handleGetLectures(ctx, update.Message.Chat.ID); err != nil {
			log.Printf("failed to handle lectures cmd: %s", err.Error())
		}

	case messageTypeLecture:
		if err := m.handleGetLecture(ctx, update.Message); err != nil {
			log.Printf("failed to handle lecture cmd: %s", err.Error())
		}

	case messageTypeHomework:
		if err := m.handleGetActiveHomework(ctx, update.Message); err != nil {
			log.Printf("failed to handle active homework cmd: %s", err.Error())
		}

	case messageTypeAllHomeworks:
		if err := m.handleGetAllHomeworks(ctx, update.Message); err != nil {
			log.Printf("failed to handle active homework cmd: %s", err.Error())
		}

	case messageTypeHomeworkPass:
		if err := m.handlePassHomework(ctx, update.Message); err != nil {
			log.Printf("failed to handle pass homework cmd: %s", err.Error())
		}

	case messageTypeHelp:
		if err := m.handleHelpCommand(update.Message); err != nil {
			log.Printf("failed to handle help cmd: %s", err.Error())
		}
	}
}

func (m *Manager) handleRegistration(ctx context.Context, update *tgbotapi.Update) error {
	text := update.Message.Text
	tgUserID := update.Message.From.ID

	if getMessageType(update.Message) != messageTypeRegistration {
		return m.sendTextMessage(update.Message.Chat.ID, msgIDontKnowU)
	}

	emailString := regexp_helpers.StripSpaces(strings.ReplaceAll(text, cmdRegistration, ""))
	isEmail, err := regexp_helpers.IsHSEEmail(strings.ReplaceAll(emailString, " ", ""))
	if err != nil {
		return errors.Wrap(err, "failed to validate email")
	}
	if !isEmail {
		return m.sendTextMessage(update.Message.Chat.ID, msgWrongEmailFormat)
	}

	student, err := m.studentService.Register(ctx, strconv.Itoa(tgUserID), emailString)
	if err != nil {
		return errors.Wrap(err, "failed to register student")
	}

	greatingText := fmt.Sprintf(msgWelcomeUsername, student.Name)

	if err := m.sendTextMessage(update.Message.Chat.ID, greatingText); err != nil {
		return errors.Wrap(err, "failed to send greating text")
	}

	return m.sendTextMessage(update.Message.Chat.ID, msgDescribeCommands)
}

func (m *Manager) handleGetLectures(ctx context.Context, chatID int64) error {
	lectures, err := m.lectureService.List(ctx)
	if err != nil {
		return err
	}

	msgText := fmt.Sprintf("%s \n\n", msgLecturesList)
	for _, lecture := range lectures {
		msgText += fmt.Sprintf("\n\xE2\x9C\x85 Лекция %d. %s \n%s", lecture.ID, lecture.Name, separator)
	}

	return m.sendTextMessage(chatID, msgText)
}

func (m *Manager) handleGetLecture(ctx context.Context, msg *tgbotapi.Message) error {
	lectureNumber, err := regexp_helpers.GetLectureNumber(msg.Text)
	if err != nil {
		if err == types.ErrGetLectrureNumber {
			return m.sendTextMessage(msg.Chat.ID, msgCantGetLectureNumber)
		}

		return err
	}

	lecture, err := m.lectureService.GetLectureWithFile(ctx, lectureNumber)
	if err != nil {
		return err
	}

	return m.sendFile(msg.Chat.ID, lecture.File, lecture.Name+".ipynb")
}

func (m *Manager) handleGetActiveHomework(ctx context.Context, msg *tgbotapi.Message) error {
	hw, err := m.hwService.GetActive(ctx)
	if err != nil {
		if err == types.ErrNotFound {
			return m.sendTextMessage(msg.Chat.ID, msgNoActiveHomework)
		}

		return err
	}

	return m.sendFile(msg.Chat.ID, hw.File, hw.Name+".ipynb")
}

func (m *Manager) handleGetAllHomeworks(ctx context.Context, msg *tgbotapi.Message) error {
	hws, err := m.hwService.List(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to load homeworks")
	}

	for _, hw := range hws {
		if err = m.sendFile(msg.Chat.ID, hw.File, hw.Name+".ipynb"); err != nil {
			return errors.Wrap(err, "failed to send file")
		}
	}

	return nil
}

func (m *Manager) handlePassHomework(ctx context.Context, msg *tgbotapi.Message) (err error) {
	defer func() {
		if err != nil {
			if sendErr := m.sendTextMessage(msg.Chat.ID, msgSmthWrong); sendErr != nil {
				log.Printf("failed to send message: %s", err.Error())
			}
		}
	}()

	if msg.Document == nil {
		return m.sendTextMessage(msg.Chat.ID, msgFileNotSent)
	}

	file, fileType, err := m.getFile(msg)
	if err != nil {
		return err
	}

	student, err := m.studentService.GetByTelegramUserID(ctx, strconv.Itoa(msg.From.ID))
	if err != nil {
		return err
	}

	if err := m.hwService.Pass(ctx, file, fileType, student); err != nil {
		if err == types.ErrHWExpired {
			return m.sendTextMessage(msg.Chat.ID, msgHWExpired)
		}
		if err == types.ErrNotFound {
			return m.sendTextMessage(msg.Chat.ID, msgNoActiveHomework)
		}

		return err
	}

	return m.sendTextMessage(msg.Chat.ID, msgHomeworkPassed)
}

func (m *Manager) handleHelpCommand(msg *tgbotapi.Message) (err error) {
	return m.sendTextMessage(msg.Chat.ID, msgDescribeCommands)
}

func (m *Manager) handleFeedback(ctx context.Context, msg *tgbotapi.Message) (err error) {
	text := msg.Text

	defer func() {
		replyMsg := msgThankU
		if err != nil {
			replyMsg = msgSmthWrong
		}

		if sendErr := m.sendTextMessage(msg.Chat.ID, replyMsg); sendErr != nil {
			log.Printf("failed to reply to feedback: %s", err.Error())
		}
	}()

	if text == "" {
		return errors.New("feedback text is empty")
	}

	if err := m.feedbackService.Create(ctx, text); err != nil {
		return errors.Wrap(err, "failed to save feedback to DB")
	}

	return err
}

func getMessageType(msg *tgbotapi.Message) messageType {
	text := strings.ToLower(msg.Text)

	if msg.Document != nil {
		return messageTypeHomeworkPass
	}

	if strings.Contains(text, cmdRegistration) {
		return messageTypeRegistration
	}
	if strings.Contains(text, cmdLectures) {
		return messageTypeLectures
	}
	if strings.Contains(text, cmdLecture) {
		return messageTypeLecture
	}
	if strings.Contains(text, cmdHomework) {
		return messageTypeHomework
	}
	if strings.Contains(text, cmdAllHomeworks) {
		return messageTypeAllHomeworks
	}

	if msg.IsCommand() && text == "/help" {
		return messageTypeHelp
	}

	return messageTypeUnknown
}

func isFeedback(msg *tgbotapi.Message) bool {
	text := strings.ToLower(msg.Text)

	return strings.Contains(text, cmdReview)
}
