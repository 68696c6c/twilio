IMAGE_NAME = twilio
DCR_SERVICE = docker-compose run --rm app

.DEFAULT:
	@echo 'App targets:'
	@echo
	@echo '    image    build the Docker image for local development'
	@echo '    deps     install dependancies'
	@echo '    test     run unit tests'
	@echo


default: .DEFAULT

image:
	docker build . -f ./Dockerfile --target dev -t $(IMAGE_NAME):dev

deps:
	$(DCR_SERVICE) go mod tidy
	$(DCR_SERVICE) go mod vendor

test:
	$(DCR_SERVICE) go test ./... -cover
