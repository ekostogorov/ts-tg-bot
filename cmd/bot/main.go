package main

import (
	"context"
	"fmt"
	"log"
	"time"
	"ts-tg-bot/internal/config"
	"ts-tg-bot/internal/manager"
	feedback_repo "ts-tg-bot/internal/repositories/feedback"
	file_repo "ts-tg-bot/internal/repositories/file"
	homework_repo "ts-tg-bot/internal/repositories/homework"
	lectures_repo "ts-tg-bot/internal/repositories/lectures"
	student_repo "ts-tg-bot/internal/repositories/student"
	student_homework_repo "ts-tg-bot/internal/repositories/student_homework"
	"ts-tg-bot/internal/runner"
	"ts-tg-bot/internal/services/feedback"
	"ts-tg-bot/internal/services/homework"
	"ts-tg-bot/internal/services/lecture"
	"ts-tg-bot/internal/services/student"

	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api"
	_ "github.com/jackc/pgx/stdlib"
	"github.com/jmoiron/sqlx"
)

func main() {
	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	defer func() {
		if r := recover(); r != nil {
			fmt.Printf("panic occurred, recovered: %+v", r)
		}
	}()

	cfg := config.GetConfig()
	log.Printf("Config: %+v", cfg)

	db, err := sqlx.Connect("pgx", cfg.PostgresDSN)
	if err != nil {
		log.Fatalf("failed to open postgres conn: %s", err.Error())
	}

	studentRepo := student_repo.New(db)
	lecturesRepo := lectures_repo.New(db)
	homeworkRepo := homework_repo.New(db)
	studentHomeworkRepo := student_homework_repo.New(db)
	fileRepo := file_repo.New()
	feedbackRepo := feedback_repo.New(db)

	studentService := student.New(studentRepo)
	lecturesService := lecture.New(lecturesRepo, fileRepo)
	homeworkService := homework.New(homeworkRepo, studentHomeworkRepo, fileRepo)
	feedbackService := feedback.New(feedbackRepo)

	bot, err := tgbotapi.NewBotAPI(cfg.TelegramAPIKey)
	if err != nil {
		log.Fatalf("failed to build tg bot: %s", err.Error())
	}

	mngr := manager.New(ctx, studentService, lecturesService, homeworkService, feedbackService, bot)

	r := runner.New()
	r.SetPreExitSleepDuration(time.Second)
	r.Add(mngr)
	r.Run()
}
