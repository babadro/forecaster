.PHONY: build

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o release/app github.com/babadro/forecaster/cmd/forecaster_bot

run: build
	docker-compose down -v && docker-compose build service && docker-compose up

down:
	docker-compose down -v