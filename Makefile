.ONESHELL:
TAB=echo "\t"
CURRENT_DIR = $(shell pwd)

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
	@$(TAB) help - вывод справки по командам

up-server:
	go run ./cmd/server/main.go ./cmd/server/config.yaml

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

