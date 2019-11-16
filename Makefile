IMAGE_NAME = twilio-client

DOC_PATH_BASE = docs/swagger.json
DOC_PATH_FINAL = docs/api-spec.json

.PHONY: docs build

.DEFAULT:
	@echo 'App targets:'
	@echo
	@echo '    image          build the Docker image for local development'
	@echo '    build          build docker image and compile the app'
	@echo '    deps           install dependancies'
	@echo '    test           run unit tests'
	@echo


default: .DEFAULT

image:
	docker build . -f ./Dockerfile --target dev -t $(IMAGE_NAME):dev

build:
	docker-compose run --rm app go build -i -o twilio

deps:
	docker-compose run --rm app go mod tidy
	docker-compose run --rm app go mod vendor

test:
	docker-compose run --rm app go test ./... -cover
