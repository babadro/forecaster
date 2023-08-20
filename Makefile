SHELL := /bin/bash

.PHONY: build run run-test-env test test-sleep down start-colima gen_mocks swag proto

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o release/app github.com/babadro/forecaster/cmd/server

run: build
	docker-compose down -v && docker-compose build service && docker-compose up

# example: make run-test-env start-bot=true
run-test-env: build
	docker-compose down -v && docker-compose build service && START_TELEGRAM_BOT=$(start-bot) docker-compose up

run-test-env-with-bot: build
	docker-compose down -v && docker-compose build service && START_TELEGRAM_BOT=true docker-compose up

# example: make test filter=TestPolls
test:
	 (source .env.tests && go test ./... -p 1 -testify.m=$(filter) -v)

sleep-filter ?= TestPolls_Options
# example: make test-sleep sleep-filter=TestPolls_Options
test-sleep:
	(source .env.tests && SLEEP_MODE=true go test ./... -testify.m=$(sleep-filter) -v)

down:
	docker-compose down -v

start-colima:
	colima start -c 8 -m 8 --arch aarch64 --vm-type=vz --vz-rosetta --mount-type=virtiofs --vz-rosetta

gen_mocks:
	@ mockery --name tgBot --structname TelegramBot --filename telegram_bot_mock.go --dir ./internal/core/forecaster/telegram

swag:
	swagger generate server --exclude-main --server-package=internal/infra/restapi --model-package=internal/models/swagger -f swagger.yaml

proto:
	protoc --go_out=./internal/core/forecaster/telegram/proto ./internal/core/forecaster/telegram/proto/*/*.proto