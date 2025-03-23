.ONESHELL:
TAB=echo "\t"
CURRENT_DIR = $(shell pwd)
BUILD_DIR = $(CURRENT_DIR)/build
VERSION = $(shell git describe --tags --always --dirty || echo "dev")
BUILD_DATE = $(shell date -u '+%Y-%m-%d_%H:%M:%S')
SERVER_ADDRESS ?= 127.0.0.1:8085

# Обновленные LDFLAGS с дополнительными параметрами
LDFLAGS = -ldflags "\
-X main.version=$(VERSION) \
-X main.buildDate=$(BUILD_DATE) \
-X main.serverAddress=$(SERVER_ADDRESS)"

help:
	@$(TAB) up-server - запустить сервер
	@$(TAB) migrate-create - создание миграции
	@$(TAB) up-docker - запуск контейнера
	@$(TAB) down-docker - остановка контейнера
	@$(TAB) migrate-up - выполнение миграций в базе данных
	@$(TAB) migrate-down - откат последней миграции в базе данных
	@$(TAB) doc - сгенерировать документацию
	@$(TAB) open-doc - запустить докер, сервер и открыть документацию в браузере
	@$(TAB) make cover - отчет покрытия тестами
	@$(TAB) make cover-save - сохранить отчет покрытия тестами
	@$(TAB) make cover-func - покрытие по функциям
	@$(TAB) make cover-percent - процент покрытия тестами\(читаем из фаила отчета\)
	@$(TAB) build-client - сборка клиента для текущей платформы
	@$(TAB) build-client-all - сборка клиента для всех платформ \(Windows, Linux, macOS\)
	@$(TAB) build-client-windows - сборка клиента для Windows
	@$(TAB) build-client-linux - сборка клиента для Linux
	@$(TAB) build-client-macos - сборка клиента для macOS
	@$(TAB) build-server-linux - сборка сервера для Linux
	@$(TAB) help - вывод справки по командам

up-server:
	go run ./cmd/server/main.go ./cmd/server/config.yaml

up-client:
	go run ./cmd/client/main.go ./cmd/client/config.yaml

up-docker:
	docker-compose up -d

down-docker:
	docker-compose down

migrate-up:
	migrate -path migrations -database "postgresql://developer:developer@localhost:5432/gophkeeper?sslmode=disable" up

migrate-down:
	migrate -path migrations -database "postgresql://developer:developer@localhost:5432/gophkeeper?sslmode=disable" down 1

migrate-create:
	migrate create -ext sql -dir migrations -seq $(name)

doc:
	swag init -g ./cmd/server/main.go

open-doc:
	$(MAKE) up-docker
	$(MAKE) doc
	@echo "Запуск сервера и открытие документации в браузере..."
	@{ \
		$(MAKE) up-server & \
		SERVER_PID=$$!; \
		sleep 3; \
		if command -v xdg-open > /dev/null; then \
			xdg-open http://127.0.0.1:8085/swagger/index.html; \
		elif command -v open > /dev/null; then \
			open http://127.0.0.1:8085/swagger/index.html; \
		elif command -v start > /dev/null; then \
			start http://127.0.0.1:8085/swagger/index.html; \
		else \
			echo "Не удалось открыть браузер автоматически. Пожалуйста, откройте http://127.0.0.1:8085/swagger/index.html вручную"; \
		fi; \
		wait $$SERVER_PID; \
	}

cover:
	go test -cover ./...

cover-save:
	go test -coverprofile=coverage.out ./...

cover-func:
	go tool cover -func=coverage.out

cover-percent:
	go tool cover -func=coverage.out | grep total

# Создание директории для сборки, если она не существует
build-dir:
	mkdir -p $(BUILD_DIR)

# Сборка клиента для текущей платформы
build-client: build-dir
	go build $(LDFLAGS) -o $(BUILD_DIR)/passcli ./cmd/client

# Сборка клиента для Windows
build-client-windows: build-dir
	GOOS=windows GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/passcli-windows-amd64.exe ./cmd/client

# Сборка клиента для Linux
build-client-linux: build-dir
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/passcli-linux-amd64 ./cmd/client

# Сборка клиента для macOS
build-client-macos: build-dir
	GOOS=darwin GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/passcli-darwin-amd64 ./cmd/client
	GOOS=darwin GOARCH=arm64 go build $(LDFLAGS) -o $(BUILD_DIR)/passcli-darwin-arm64 ./cmd/client

# Сборка клиента для всех платформ
build-client-all: build-client-windows build-client-linux build-client-macos

# Сборка сервера для Linux
build-server-linux: build-dir
	GOOS=linux GOARCH=amd64 go build $(LDFLAGS) -o $(BUILD_DIR)/passserver-linux-amd64 ./cmd/server

