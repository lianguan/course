.PHONY:
.SILENT:
.DEFAULT_GOAL := run

build:
	@echo "========================================="
	@echo " UltraThreads API - Building..."
	@echo "========================================="
	@echo "[1/3] Downloading dependencies..."
	go mod download
	@echo "[2/3] Building binary (linux/amd64)..."
	CGO_ENABLED=0 GOOS=linux go build -o ./app ./main.go
	@echo "[3/3] Build complete!"
	@echo " Output: ./app"
	@ls -lh ./app 2>/dev/null || dir /b .\app 2>nul
	@echo "========================================="

run: build
	docker-compose up --remove-orphans app

debug: build
	docker-compose up --remove-orphans debug

test:
	go test --short -coverprofile=cover.out -v ./...
	make test.coverage

# Testing Vars
export TEST_MYSQL_HOST=localhost
export TEST_MYSQL_PORT=3306
export TEST_MYSQL_USER=root
export TEST_MYSQL_PASSWORD=root1234
export TEST_MYSQL_DATABASE=test
export TEST_CONTAINER_NAME=test_mysql

test.integration:
	docker run --rm -d -p 3306:3306 --name $$TEST_CONTAINER_NAME \
		-e MYSQL_ROOT_PASSWORD=$$TEST_MYSQL_PASSWORD \
		-e MYSQL_DATABASE=$$TEST_MYSQL_DATABASE \
		mysql:8.0

	GIN_MODE=release go test -v ./tests/
	docker stop $$TEST_CONTAINER_NAME

test.coverage:
	go tool cover -func=cover.out | grep "total"

swag:
	swag init -g internal/app/app.go

lint:
	golangci-lint run

gen:
	mockgen -source=internal/service/service.go -destination=internal/service/mocks/mock.go
	mockgen -source=internal/repository/repository.go -destination=internal/repository/mocks/mock.go
