PG_DSN=<YOUR-DB-DSN>
TELEGRAM_API_KEY=<YOUR-API-KEY>
BIN_PATH=./bin/bot
CMD_PATH=./cmd/bot

.PHONY: run_local
run_local:
	PG_DSN=$(PG_DSN) TELEGRAM_API_KEY=$(TELEGRAM_API_KEY) go run cmd/bot/*go

.PHONY: build
build:
	go build -o $(BIN_PATH) $(CMD_PATH)

.PHONY: build_for_linux
build_for_linux:
	GOOS=linux GOARCH=amd64 go build -o $(BIN_PATH) $(CMD_PATH)