.PHONY: build run down start-colima swag

build:
	GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o release/app github.com/babadro/forecaster/cmd/server

run: build
	docker-compose down -v && docker-compose build service && docker-compose up

run-test-env: build
	docker-compose down -v && docker-compose build service && START_TELEGRAM_BOT=false docker-compose up

test:
	 (source .env.tests && go test ./... -v)

down:
	docker-compose down -v

start-colima:
	colima start -c 8 -m 8 --arch aarch64 --vm-type=vz --vz-rosetta --mount-type=virtiofs --vz-rosetta

swag:
	swagger generate server --exclude-main --server-package=internal/infra/restapi --model-package=internal/models/swagger -f swagger.yaml