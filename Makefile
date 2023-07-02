.PHONY: build

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o release/app github.com/babadro/forecaster/cmd/server

run: build
	docker-compose down -v && docker-compose build service && docker-compose up

down:
	docker-compose down -v

start-colima:
	colima start -c 8 -m 8 --arch aarch64 --vm-type=vz --vz-rosetta --mount-type=virtiofs --vz-rosetta