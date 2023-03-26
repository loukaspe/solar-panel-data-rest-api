generate-mock:
ifndef FILE
	@echo "FILE parameter is missing"
	@exit 1
endif
	make build-dev
	@docker run --volume "$(PWD)"/../:/app --workdir /app \
	dev-build /bin/bash -c "mockgen -source=${FILE} -destination=mocks/mock_${FILE}"

tests-unit:
	make build-dev
	@docker run \
			--rm \
			--volume "$(PWD)"/../:/app \
			--workdir /app \
			dev-build go test -short -cover -count=1 ./...

tests-benchmark:
	make build-dev
	@docker run \
			--rm \
			--volume "$(PWD)"/../:/app \
			--workdir /app \
			dev-build go test ./... -bench=.

tests-file:
	make build-dev
	@docker run \
			--rm \
			--volume "$(PWD)"/../:/app \
			--workdir /app \
			dev-build go test -count=1 -v ${FILE}

tests-package:
	make build-dev
	@docker run \
			--rm \
			--volume "$(PWD)"/../:/app \
			--workdir /app \
			dev-build go test -count=1  ${PACKAGE}


tests-all:
	make build-dev
	@docker run \
			--rm \
			--volume "$(PWD)"/../:/app \
			--workdir /app \
			dev-build go test ./... -cover -count=1

tests-all-with-coverage:
	make build-dev
	@docker run \
			--rm \
			--volume "$(PWD)"/../:/app \
			--workdir /app \
			dev-build go test -count=1 -v -coverpkg=./... -coverprofile=profile.cov ./... ; go tool cover -func profile.cov

run-linter:
	make build-dev
	@docker run \
			--rm \
			--volume "$(PWD)"/../:/app \
			--workdir /app \
			dev-build golangci-lint run --config ./config/.golangci.yml ./...

start-app:
	docker-compose --file ../deployment/docker-compose.yml down
	docker-compose --file ../deployment/docker-compose.yml up

kill-app:
	docker-compose --file ../deployment/docker-compose.yml down --remove-orphans --volumes

rebuild-app:
	docker-compose --file ../deployment/docker-compose.yml stop solar-panel-data-app
	docker-compose --file ../deployment/docker-compose.yml build solar-panel-data-app
	docker-compose --file ../deployment/docker-compose.yml up -d solar-panel-data-app

build-dev:
	@docker build \
			--tag dev-build \
			-f ../build/Dockerfile.utilities ..