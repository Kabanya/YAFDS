# ============================================================================
# YAFDS Project Orchestrator Makefile
# ============================================================================
# Главный оркестратор для управления всеми микросервисами проекта
# ============================================================================

.PHONY: help
.DEFAULT_GOAL := help

# Цвета для вывода
CYAN := \033[36m
GREEN := \033[32m
YELLOW := \033[33m
RED := \033[31m
RESET := \033[0m

# Список сервисов
SERVICES := customer courier restaurant front
BACKEND_SERVICES := customer courier restaurant

# ============================================================================
# HELP - Главное меню помощи
# ============================================================================

help:
	@echo "$(CYAN)    ╔══════════════════════════════════════════════════════╗$(RESET)"
	@echo "$(CYAN)    ║              YAFDS Project Orchestrator              ║$(RESET)"
	@echo "$(CYAN)    ╚══════════════════════════════════════════════════════╝$(RESET)"
	@echo ""
	@echo "$(GREEN)📦 Основные команды:$(RESET)"
	@echo "  $(YELLOW)make start$(RESET)              - Запустить все сервисы (полный деплой)"
	@echo "  $(YELLOW)make start-dev$(RESET)          - Запустить все сервисы с тестовыми данными"
	@echo "  $(YELLOW)make stop$(RESET)               - Остановить все сервисы"
	@echo "  $(YELLOW)make restart$(RESET)            - Перезапустить все сервисы"
	@echo "  $(YELLOW)make status$(RESET)             - Показать статус всех контейнеров"
	@echo ""
	@echo "$(GREEN)🔧 Сервисы (customer, courier, restaurant, front):$(RESET)"
	@echo "  $(YELLOW)make start-<service>$(RESET)    - Запустить конкретный сервис"
	@echo "  $(YELLOW)make stop-<service>$(RESET)     - Остановить конкретный сервис"
	@echo "  $(YELLOW)make restart-<service>$(RESET)  - Перезапустить конкретный сервис"
	@echo "  $(YELLOW)make logs-<service>$(RESET)     - Показать логи конкретного сервиса"
	@echo "  $(YELLOW)make build-<service>$(RESET)    - Пересобрать конкретный сервис"
	@echo ""
	@echo "$(GREEN)🗄️  База данных и миграции:$(RESET)"
	@echo "  $(YELLOW)make migrate-up$(RESET)         - Применить все миграции"
	@echo "  $(YELLOW)make migrate-down$(RESET)       - Откатить все миграции"
	@echo "  $(YELLOW)make migrate-status$(RESET)     - Статус миграций всех сервисов"
	@echo "  $(YELLOW)make db-clean$(RESET)           - Очистить все данные в БД"
	@echo "  $(YELLOW)make db-seed$(RESET)            - Заполнить БД тестовыми данными"
	@echo ""
	@echo "$(GREEN)📊 Мониторинг и логи:$(RESET)"
	@echo "  $(YELLOW)make logs$(RESET)               - Показать логи всех сервисов"
	@echo "  $(YELLOW)make logs-follow$(RESET)        - Следить за логами в реальном времени"
	@echo "  $(YELLOW)make ps$(RESET)                 - Список всех контейнеров проекта"
	@echo ""
	@echo "$(GREEN)🧹 Очистка:$(RESET)"
	@echo "  $(YELLOW)make clean$(RESET)              - Остановить и удалить все контейнеры"
	@echo "  $(YELLOW)make clean-all$(RESET)          - Полная очистка (контейнеры + образы + volumes)"
	@echo "  $(YELLOW)make clean-images$(RESET)       - Удалить Docker образы проекта"
	@echo ""
	@echo "$(GREEN)🏗️  Сборка:$(RESET)"
	@echo "  $(YELLOW)make build$(RESET)              - Пересобрать все сервисы"
	@echo "  $(YELLOW)make build-backend$(RESET)      - Пересобрать только backend сервисы"
	@echo ""
	@echo "$(GREEN)🔍 Диагностика:$(RESET)"
	@echo "  $(YELLOW)make health$(RESET)             - Проверить здоровье всех сервисов"
	@echo "  $(YELLOW)make check-deps$(RESET)         - Проверить зависимости (Docker, Go, etc.)"
	@echo ""
	@echo "$(GREEN)🐍 Python утилиты:$(RESET)"
	@echo "  $(YELLOW)make prepare-python$(RESET)     - Создать Python venv с зависимостями"
	@echo "  $(YELLOW)make send-test-data$(RESET)     - Запустить загрузку тестовых данных"
	@echo "  $(YELLOW)make generate-db-html$(RESET)   - Генерировать HTML таблиц БД"
	@echo ""
	@echo "$(GREEN)🔧 Go утилиты:$(RESET)"
	@echo "  $(YELLOW)make go-tidy$(RESET)            - Выполнить go mod tidy для всех backend сервисов"
	@echo "  $(YELLOW)make go-fmt$(RESET)             - Отформатировать код (go fmt) во всех backend сервисах"
	@echo "  $(YELLOW)make go-vet$(RESET)             - Проверить код (go vet) во всех backend сервисах"
	@echo "  $(YELLOW)make go-test$(RESET)            - Запустить тесты Go (go test) для всех backend сервисов"
	@echo "  $(YELLOW)make go-build$(RESET)           - Собрать бинарники (go build) для всех backend сервисов"
	@echo "  $(YELLOW)make go-clean$(RESET)           - Очистить Go кэш и бинарники"
	@echo ""

# ============================================================================
# ОСНОВНЫЕ КОМАНДЫ
# ============================================================================

setup-python-env:
	@echo "$(CYAN)🐍 Подготовка Python окружения...$(RESET)"
	@if [ ! -d ".venv-yafds" ]; then \
		echo "$(YELLOW)Создание виртуального окружения .venv-yafds...$(RESET)"; \
		python3 -m venv .venv-yafds; \
		echo "$(GREEN)✅ Виртуальное окружение создано$(RESET)"; \
	fi
	@echo "$(YELLOW)Установка зависимостей...$(RESET)"
	@. .venv-yafds/bin/activate && pip install --upgrade pip setuptools wheel > /dev/null 2>&1
	@. .venv-yafds/bin/activate && pip install tqdm colorama pandas psycopg2-binary requests > /dev/null 2>&1
	@echo "$(GREEN)✅ Зависимости установлены: tqdm, colorama, pandas, psycopg2, requests$(RESET)"
	@echo "$(GREEN)🚀 Запуск всех сервисов YAFDS...$(RESET)"
	@echo "$(CYAN)Step 1/4: Запуск Customer сервиса$(RESET)"
	@cd customer && $(MAKE) run
	@echo ""
	@echo "$(CYAN)Step 2/4: Запуск Courier сервиса$(RESET)"
	@cd courier && $(MAKE) run
	@echo ""
	@echo "$(CYAN)Step 3/4: Запуск Restaurant сервиса$(RESET)"
	@cd restaurant && $(MAKE) run
	@echo ""
	@echo "$(CYAN)Step 4/4: Запуск Frontend$(RESET)"
	@cd front && $(MAKE) run
	@echo ""
	@echo "$(GREEN)✅ Все сервисы успешно запущены!$(RESET)"
	@$(MAKE) status

start-dev: check-docker setup-python-env
	@echo "$(GREEN)🚀 Запуск всех сервисов с тестовыми данными...$(RESET)"
	@cd customer && $(MAKE) dev
	@cd courier && $(MAKE) dev
	@cd restaurant && $(MAKE) dev
	@cd front && $(MAKE) run
	@echo "$(GREEN)✅ Все сервисы запущены в dev режиме!$(RESET)"
	@sleep 3
	@$(MAKE) status
	@echo ""
	@if [ -f "send_requests.py" ]; then \
		echo "$(GREEN)🌱 Загрузка тестовых данных через send_requests.py...$(RESET)"; \
		. .venv-yafds/bin/activate && python send_requests.py; \
		echo "$(GREEN)✅ Тестовые данные загружены!$(RESET)"; \
	else \
		echo "$(YELLOW)⚠️  send_requests.py не найден, пропускаем загрузку тестовых данных$(RESET)"; \
	fi

run: start-dev

stop:
	@echo "$(YELLOW)🛑 Остановка всех сервисов...$(RESET)"
	@for service in $(SERVICES); do \
		echo "$(CYAN)Останавливаем $$service...$(RESET)"; \
		cd $$service && $(MAKE) clear-all 2>/dev/null || true; \
		cd ..; \
	done
	@echo "$(GREEN)✅ Все сервисы остановлены$(RESET)"

restart: stop
	@sleep 2
	@$(MAKE) start

# ============================================================================
# УПРАВЛЕНИЕ ОТДЕЛЬНЫМИ СЕРВИСАМИ
# ============================================================================

start-customer:
	@echo "$(GREEN)🚀 Запуск Customer сервиса...$(RESET)"
	@cd customer && $(MAKE) run

stop-customer:
	@echo "$(YELLOW)🛑 Остановка Customer сервиса...$(RESET)"
	@cd customer && $(MAKE) clear-all

restart-customer: stop-customer
	@sleep 1
	@$(MAKE) start-customer

logs-customer:
	@cd customer && $(MAKE) logs-customer

build-customer:
	@echo "$(CYAN)🔨 Сборка Customer сервиса...$(RESET)"
	@cd customer && $(MAKE) build-customer

# Courier
start-courier:
	@echo "$(GREEN)🚀 Запуск Courier сервиса...$(RESET)"
	@cd courier && $(MAKE) run

stop-courier:
	@echo "$(YELLOW)🛑 Остановка Courier сервиса...$(RESET)"
	@cd courier && $(MAKE) clear-all

restart-courier: stop-courier
	@sleep 1
	@$(MAKE) start-courier

logs-courier:
	@cd courier && $(MAKE) logs-courier

build-courier:
	@echo "$(CYAN)🔨 Сборка Courier сервиса...$(RESET)"
	@cd courier && $(MAKE) build-courier

# Restaurant
start-restaurant:
	@echo "$(GREEN)🚀 Запуск Restaurant сервиса...$(RESET)"
	@cd restaurant && $(MAKE) run

stop-restaurant:
	@echo "$(YELLOW)🛑 Остановка Restaurant сервиса...$(RESET)"
	@cd restaurant && $(MAKE) clear-all

restart-restaurant: stop-restaurant
	@sleep 1
	@$(MAKE) start-restaurant

logs-restaurant:
	@cd restaurant && $(MAKE) logs-restaurant

build-restaurant:
	@echo "$(CYAN)🔨 Сборка Restaurant сервиса...$(RESET)"
	@cd restaurant && $(MAKE) build-restaurant

# Front
start-front:
	@echo "$(GREEN)🚀 Запуск Frontend...$(RESET)"
	@cd front && $(MAKE) run

stop-front:
	@echo "$(YELLOW)🛑 Остановка Frontend...$(RESET)"
	@cd front && $(MAKE) clean

restart-front: stop-front
	@sleep 1
	@$(MAKE) start-front

logs-front:
	@cd front && $(MAKE) logs

build-front:
	@echo "$(CYAN)🔨 Сборка Frontend...$(RESET)"
	@cd front && $(MAKE) build

# ============================================================================
# МИГРАЦИИ И БАЗА ДАННЫХ
# ============================================================================

migrate-up:
	@echo "$(GREEN)📊 Применение всех миграций...$(RESET)"
	@cd migrations && $(MAKE) up

migrate-down:
	@echo "$(YELLOW)📊 Откат всех миграций...$(RESET)"
	@cd migrations && $(MAKE) down

migrate-status:
	@echo "$(CYAN)📊 Статус миграций всех сервисов:$(RESET)"
	@cd migrations && $(MAKE) migrate-status

db-clean:
	@echo "$(YELLOW)🗑️  Очистка базы данных...$(RESET)"
	@cd migrations && $(MAKE) clean
	@echo "$(GREEN)✅ База данных очищена$(RESET)"

db-seed:
	@echo "$(GREEN)🌱 Заполнение БД тестовыми данными...$(RESET)"
	@cd customer && $(MAKE) seed-up
	@cd courier && $(MAKE) seed-up
	@cd restaurant && $(MAKE) seed-up
	@echo "$(GREEN)✅ Тестовые данные загружены$(RESET)"

# ============================================================================
# МОНИТОРИНГ И ЛОГИ
# ============================================================================

status:
	@echo "$(CYAN)📊 Статус контейнеров YAFDS:$(RESET)"
	@echo ""
	@docker ps -a --filter "name=yafds" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" || echo "Нет запущенных контейнеров"

ps:
	@docker ps -a --filter "name=yafds"

logs:
	@echo "$(CYAN)📜 Логи всех сервисов:$(RESET)"
	@echo ""
	@for service in customer courier restaurant; do \
		echo "$(YELLOW)=== $$service ===$(RESET)"; \
		docker logs yafds-$$service-service 2>/dev/null | tail -n 20 || echo "Контейнер не запущен"; \
		echo ""; \
	done

logs-follow:
	@echo "$(CYAN)📜 Следим за логами (Ctrl+C для выхода)...$(RESET)"
	@docker logs -f yafds-customer-service 2>/dev/null || echo "Customer не запущен"

health:
	@echo "$(CYAN)🏥 Проверка здоровья сервисов:$(RESET)"
	@echo ""
	@echo "$(YELLOW)Customer:$(RESET)"
	@curl -s http://localhost:8080/health 2>/dev/null && echo "$(GREEN)✅ OK$(RESET)" || echo "$(RED)❌ DOWN$(RESET)"
	@echo ""
	@echo "$(YELLOW)Courier:$(RESET)"
	@curl -s http://localhost:8081/health 2>/dev/null && echo "$(GREEN)✅ OK$(RESET)" || echo "$(RED)❌ DOWN$(RESET)"
	@echo ""
	@echo "$(YELLOW)Restaurant:$(RESET)"
	@curl -s http://localhost:8082/health 2>/dev/null && echo "$(GREEN)✅ OK$(RESET)" || echo "$(RED)❌ DOWN$(RESET)"
	@echo ""
	@echo "$(YELLOW)Frontend:$(RESET)"
	@curl -s http://localhost:5173 2>/dev/null > /dev/null && echo "$(GREEN)✅ OK$(RESET)" || echo "$(RED)❌ DOWN$(RESET)"

# ============================================================================
# СБОРКА
# ============================================================================

build:
	@echo "$(CYAN)🔨 Пересборка всех сервисов...$(RESET)"
	@for service in $(SERVICES); do \
		echo "$(YELLOW)Сборка $$service...$(RESET)"; \
		cd $$service && $(MAKE) build-$$service 2>/dev/null || $(MAKE) build || true; \
		cd ..; \
	done
	@echo "$(GREEN)✅ Все сервисы собраны$(RESET)"

build-backend:
	@echo "$(CYAN)🔨 Пересборка backend сервисов...$(RESET)"
	@for service in $(BACKEND_SERVICES); do \
		echo "$(YELLOW)Сборка $$service...$(RESET)"; \
		cd $$service && $(MAKE) build-$$service; \
		cd ..; \
	done
	@echo "$(GREEN)✅ Backend сервисы собраны$(RESET)"

# ============================================================================
# ОЧИСТКА
# ============================================================================

clean: stop
	@echo "$(RED)🧹 Удаление всех контейнеров...$(RESET)"
	@docker ps -a --filter "name=yafds" -q | xargs -r docker rm -f 2>/dev/null || true
	@echo "$(GREEN)✅ Контейнеры удалены$(RESET)"

clean-images:
	@echo "$(RED)🗑️  Удаление Docker образов...$(RESET)"
	@docker images --filter "reference=yafds*" -q | xargs -r docker rmi -f 2>/dev/null || true
	@echo "$(GREEN)✅ Образы удалены$(RESET)"

clean-all: clean clean-images
	@echo "$(RED)🗑️  Удаление volumes...$(RESET)"
	@docker volume ls --filter "name=yafds" -q | xargs -r docker volume rm 2>/dev/null || true
	@echo "$(GREEN)✅ Полная очистка завершена$(RESET)"

# ============================================================================
# ДИАГНОСТИКА И ПРОВЕРКИ
# ============================================================================

check-docker:
	@docker info > /dev/null 2>&1 || (echo "$(RED)❌ Docker не запущен! Запустите Docker Desktop.$(RESET)" && exit 1)
	@echo "$(GREEN)✅ Docker запущен$(RESET)"

check-deps:
	@echo "$(CYAN)🔍 Проверка зависимостей:$(RESET)"
	@echo ""
	@echo -n "Docker: "
	@docker --version 2>/dev/null && echo "$(GREEN)✅$(RESET)" || echo "$(RED)❌$(RESET)"
	@echo -n "Docker Compose: "
	@docker compose version 2>/dev/null && echo "$(GREEN)✅$(RESET)" || echo "$(RED)❌$(RESET)"
	@echo -n "Go: "
	@go version 2>/dev/null && echo "$(GREEN)✅$(RESET)" || echo "$(RED)❌$(RESET)"
	@echo -n "Node.js: "
	@node --version 2>/dev/null && echo "$(GREEN)✅$(RESET)" || echo "$(RED)❌$(RESET)"
	@echo -n "Goose: "
	@goose --version 2>/dev/null && echo "$(GREEN)✅$(RESET)" || echo "$(RED)❌$(RESET)"
	@echo ""

# ============================================================================
# УТИЛИТЫ
# ============================================================================

update-deps:
	@echo "$(CYAN)📦 Обновление зависимостей...$(RESET)"
	@./update_deps.sh

test:
	@echo "$(CYAN)🧪 Запуск тестов...$(RESET)"
	@for service in $(BACKEND_SERVICES); do \
		echo "$(YELLOW)Тестирование $$service...$(RESET)"; \
		cd $$service && go test ./... -v || true; \
		cd ..; \
	done

# Быстрый перезапуск для разработки
quick-restart: stop-customer stop-courier stop-restaurant
	@sleep 1
	@$(MAKE) start-customer &
	@$(MAKE) start-courier &
	@$(MAKE) start-restaurant &
	@wait
	@echo "$(GREEN)✅ Backend сервисы перезапущены$(RESET)"

# Показать все переменные окружения
show-env:
	@echo "$(CYAN)📋 Переменные окружения сервисов:$(RESET)"
	@for service in $(BACKEND_SERVICES); do \
		echo ""; \
		echo "$(YELLOW)=== $$service ===$(RESET)"; \
		cat $$service/config/.env 2>/dev/null || echo "Файл .env не найден"; \
	done

# ============================================================================
# PYTHON УТИЛИТЫ
# ============================================================================

prepare-python: setup-python-env
	@echo "$(GREEN)✅ Python окружение готово$(RESET)"
	@echo "$(CYAN)📝 Используйте для активации:$(RESET)"
	@echo "  $(YELLOW)source .venv-yafds/bin/activate$(RESET)"

send-test-data:
	@if [ ! -d ".venv-yafds" ]; then \
		echo "$(RED)❌ Python окружение не найдено! Запустите $(YELLOW)make prepare-python$(RED)$(RESET)"; \
		exit 1; \
	fi
	@if [ ! -f "send_requests.py" ]; then \
		echo "$(RED)❌ send_requests.py не найден!$(RESET)"; \
		exit 1; \
	fi
	@echo "$(GREEN)🌱 Запуск загрузки тестовых данных...$(RESET)"
	@. .venv-yafds/bin/activate && python send_requests.py
	@echo "$(GREEN)✅ Тестовые данные загружены!$(RESET)"

generate-db-html:
	@if [ ! -d ".venv-yafds" ]; then \
		echo "$(RED)❌ Python окружение не найдено! Запустите $(YELLOW)make prepare-python$(RED)$(RESET)"; \
		exit 1; \
	fi
	@if [ ! -f "make_nice_sql_html.py" ]; then \
		echo "$(RED)❌ make_nice_sql_html.py не найден!$(RESET)"; \
		exit 1; \
	fi
	@echo "$(GREEN)📊 Генерация HTML с содержимым БД...$(RESET)"
	@. .venv-yafds/bin/activate && python make_nice_sql_html.py
	@echo "$(GREEN)✅ HTML файл сгенерирован: tables.html$(RESET)"

# ============================================================================
# GO УТИЛИТЫ
# ============================================================================

go-tidy:
	@echo "$(CYAN)🔨 Выполнение go mod tidy для всех backend сервисов...$(RESET)"
	@for service in $(BACKEND_SERVICES); do \
		echo "$(YELLOW)go mod tidy в $$service...$(RESET)"; \
		cd $$service && go mod tidy && echo "$(GREEN)✅ $$service$(RESET)" || echo "$(RED)❌ $$service$(RESET)"; \
		cd ..; \
	done
	@echo "$(GREEN)✅ go mod tidy завершено для всех сервисов$(RESET)"

go-fmt:
	@echo "$(CYAN)🎨 Форматирование кода (go fmt) для всех backend сервисов...$(RESET)"
	@for service in $(BACKEND_SERVICES); do \
		echo "$(YELLOW)Форматирование $$service...$(RESET)"; \
		cd $$service && go fmt ./... && echo "$(GREEN)✅ $$service$(RESET)" || echo "$(RED)❌ $$service$(RESET)"; \
		cd ..; \
	done
	@echo "$(GREEN)✅ go fmt завершено для всех сервисов$(RESET)"

go-vet:
	@echo "$(CYAN)🔍 Проверка кода (go vet) для всех backend сервисов...$(RESET)"
	@for service in $(BACKEND_SERVICES); do \
		echo "$(YELLOW)Проверка $$service...$(RESET)"; \
		cd $$service && go vet ./... && echo "$(GREEN)✅ $$service$(RESET)" || echo "$(YELLOW)⚠️  $$service (предупреждения)$(RESET)"; \
		cd ..; \
	done
	@echo "$(GREEN)✅ go vet завершено для всех сервисов$(RESET)"

go-test:
	@echo "$(CYAN)🧪 Запуск тестов (go test) для всех backend сервисов...$(RESET)"
	@for service in $(BACKEND_SERVICES); do \
		echo "$(YELLOW)Тестирование $$service...$(RESET)"; \
		cd $$service && go test -v -cover ./... || echo "$(YELLOW)⚠️  Тесты в $$service не прошли$(RESET)"; \
		cd ..; \
	done
	@echo "$(GREEN)✅ go test завершено для всех сервисов$(RESET)"

go-build:
	@echo "$(CYAN)🏗️  Сборка бинарников (go build) для всех backend сервисов...$(RESET)"
	@for service in $(BACKEND_SERVICES); do \
		echo "$(YELLOW)Сборка $$service...$(RESET)"; \
		cd $$service && go build -o bin/$$service ./cmd/server && echo "$(GREEN)✅ $$service$(RESET)" || echo "$(RED)❌ $$service$(RESET)"; \
		cd ..; \
	done
	@echo "$(GREEN)✅ Сборка завершена для всех сервисов$(RESET)"

go-clean:
	@echo "$(YELLOW)🧹 Очистка Go кэша и бинарников...$(RESET)"
	@go clean -cache
	@for service in $(BACKEND_SERVICES); do \
		echo "$(YELLOW)Очистка $$service...$(RESET)"; \
		cd $$service && go clean -v && rm -rf bin/ || true; \
		cd ..; \
	done
	@echo "$(GREEN)✅ Очистка завершена$(RESET)"

go-deps-check:
	@echo "$(CYAN)📦 Проверка зависимостей для всех backend сервисов...$(RESET)"
	@for service in $(BACKEND_SERVICES); do \
		echo "$(YELLOW)Проверка $$service...$(RESET)"; \
		cd $$service && go list -u -m all | grep -v indirect || echo "$(GREEN)✅ Все зависимости актуальны$(RESET)"; \
		cd ..; \
	done
	@echo "$(GREEN)✅ Проверка завершена$(RESET)"

go-upgrade-deps:
	@echo "$(CYAN)📦 Обновление зависимостей для всех backend сервисов...$(RESET)"
	@for service in $(BACKEND_SERVICES); do \
		echo "$(YELLOW)Обновление $$service...$(RESET)"; \
		cd $$service && go get -u ./... && go mod tidy && echo "$(GREEN)✅ $$service$(RESET)" || echo "$(RED)❌ $$service$(RESET)"; \
		cd ..; \
	done
	@echo "$(GREEN)✅ Обновление завершено для всех сервисов$(RESET)"
