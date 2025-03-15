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
	@$(TAB) help - вывод справки по командам

up-server:
	go run ./cmd/server/main.go -a=localhost:41849 -d=postgresql://developer:developer@localhost:5432/postgres?sslmode=disable

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