# Makefile

.PHONY: build up test down

# Собрать контейнеры
build:
	docker-compose build

# Запустить контейнеры (app, db и тесты)
up:
	docker-compose up -d

# Запустить тесты
test:
	docker-compose run --rm test

# Остановить контейнеры
down:
	docker-compose down
