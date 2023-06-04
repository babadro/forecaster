.PHONY: build

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o release/app github.com/babadro/forecaster/cmd/forecaster_bot