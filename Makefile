include .env

.PHONY: lint
lint:
	golangci-lint run ./...

.PHONY: test
test:
	GOEXPERIMENT=loopvar \
	DOCKER_HOST=${TEST_DOCKER_HOST} \
	TEST_MYSQL_DATABASE=${TEST_MYSQL_DATABASE} \
	TEST_MYSQL_ROOT_PASSWORD=${TEST_MYSQL_ROOT_PASSWORD} \
	TEST_MYSQL_USER=${TEST_MYSQL_USER} \
	TEST_MYSQL_PASSWORD=${TEST_MYSQL_PASSWORD} \
	TEST_MYSQL_HOST=${TEST_MYSQL_HOST} \
	TEST_MYSQL_TCP_PORT=${TEST_MYSQL_TCP_PORT} \
	go test -race ./...

.PHONY: dev
dev:
	MYSQL_DATABASE=${MYSQL_DATABASE} \
	MYSQL_ROOT_PASSWORD=${MYSQL_ROOT_PASSWORD} \
	MYSQL_USER=${MYSQL_USER} \
	MYSQL_PASSWORD=${MYSQL_PASSWORD} \
	MYSQL_HOST=${MYSQL_HOST} \
	MYSQL_TCP_PORT=${MYSQL_TCP_PORT} \
	go run ./...

.PHONY: docker-compose-up
docker-compose-up:
	docker compose -f ./docker/compose.yml up -d

.PHONY: docker-compose-down
docker-compose-down:
	docker compose -f ./docker/compose.yml down

.PHONY: mysql
mysql:
	mysql -h ${MYSQL_HOST} -u${MYSQL_USER} -p${MYSQL_PASSWORD} ${MYSQL_DATABASE}

.PHONY: sqlc-generate
sqlc-generate:
	sqlc generate

.PHONY: golang-migrate-up
golang-migrate-up:
	migrate -path db/migrations -database "mysql://${MYSQL_USER}:${MYSQL_PASSWORD}@tcp(${MYSQL_HOST}:${MYSQL_TCP_PORT})/${MYSQL_DATABASE}" --verbose up
	
.PHONY: golang-migrate-down
golang-migrate-down:
	migrate -path db/migrations -database "mysql://${MYSQL_USER}:${MYSQL_PASSWORD}@tcp(${MYSQL_HOST}:${MYSQL_TCP_PORT})/${MYSQL_DATABASE}" --verbose down
	
.PHONY: golang-migrate-drop
golang-migrate-drop:
	migrate -path db/migrations -database "mysql://${MYSQL_USER}:${MYSQL_PASSWORD}@tcp(${MYSQL_HOST}:${MYSQL_TCP_PORT})/${MYSQL_DATABASE}" --verbose drop
	