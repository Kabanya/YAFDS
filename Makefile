# ============================================================================
# YAFDS Project Orchestrator Makefile
# ============================================================================

.DEFAULT_GOAL := help

# Цвета
CYAN   := \033[36m
GREEN  := \033[32m
YELLOW := \033[33m
RED    := \033[31m
RESET  := \033[0m

# Сервисы
SERVICES         := customer courier restaurant front
BACKEND_SERVICES := customer courier restaurant
GO_MODULES       := pkg customer courier restaurant

# Тестовые доступы
TEST_CUSTOMER_LOGIN   := "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
TEST_CUSTOMER_PASS    := "StrongPass001"
TEST_COURIER_LOGIN    := "0x5647445a3e564b189fc35a75079fc924"
TEST_COURIER_PASS     := "aad9D6C7e9R7"
TEST_RESTAURANT_LOGIN := "0xbd2f054c1faa44b2a561c47dd1a0a367"
TEST_RESTAURANT_PASS  := "JO88M0Rtf0BG"

.PHONY: help \
        start start-dev run stop restart \
        start-customer stop-customer restart-customer logs-customer build-customer \
        start-courier  stop-courier  restart-courier  logs-courier  build-courier \
        start-restaurant stop-restaurant restart-restaurant logs-restaurant build-restaurant \
        start-front stop-front restart-front logs-front build-front \
        migrate-up migrate-down migrate-status db-clean db-seed \
        build build-backend install-local build-local run-local stop-local clean-local \
        clean clean-images clean-all \
        status ps logs logs-follow health \
        check-docker check-deps \
        init-env init update-deps test quick-restart show-env \
        setup-python-env prepare-python send-test-data generate-db-html \
        go-tidy go-fmt go-vet go-test go-build go-clean go-deps-check go-upgrade-deps

# ============================================================================
# HELP
# ============================================================================

help:
	@printf '%b\n' "$(CYAN)    ╔══════════════════════════════════════════════════════╗$(RESET)"
	@printf '%b\n' "$(CYAN)    ║              YAFDS Project Orchestrator              ║$(RESET)"
	@printf '%b\n' "$(CYAN)    ╚══════════════════════════════════════════════════════╝$(RESET)"
	@printf '\n'
	@printf '%b\n' "$(GREEN)📦 Основные команды:$(RESET)"
	@printf '%b\n' "  $(YELLOW)make start$(RESET)              - Запустить все сервисы (полный деплой)"
	@printf '%b\n' "  $(YELLOW)make start-dev$(RESET)          - Запустить все сервисы с тестовыми данными"
	@printf '%b\n' "  $(YELLOW)make run-local$(RESET)          - Очистить, собрать и запустить всё локально"
	@printf '%b\n' "  $(YELLOW)make stop-local$(RESET)         - Остановить все локальные сервисы"
	@printf '%b\n' "  $(YELLOW)make init-env$(RESET)           - Инициализировать .env файлы из .env.example"
	@printf '%b\n' "  $(YELLOW)make stop$(RESET)               - Остановить все сервисы"
	@printf '%b\n' "  $(YELLOW)make restart$(RESET)            - Перезапустить все сервисы"
	@printf '%b\n' "  $(YELLOW)make status$(RESET)             - Показать статус всех контейнеров"
	@printf '\n'
	@printf '%b\n' "$(GREEN)🔧 Сервисы (customer, courier, restaurant, front):$(RESET)"
	@printf '%b\n' "  $(YELLOW)make start-<service>$(RESET)    - Запустить конкретный сервис"
	@printf '%b\n' "  $(YELLOW)make stop-<service>$(RESET)     - Остановить конкретный сервис"
	@printf '%b\n' "  $(YELLOW)make restart-<service>$(RESET)  - Перезапустить конкретный сервис"
	@printf '%b\n' "  $(YELLOW)make logs-<service>$(RESET)     - Показать логи конкретного сервиса"
	@printf '%b\n' "  $(YELLOW)make build-<service>$(RESET)    - Пересобрать конкретный сервис"
	@printf '\n'
	@printf '%b\n' "$(GREEN)🗄️  База данных и миграции:$(RESET)"
	@printf '%b\n' "  $(YELLOW)make migrate-up$(RESET)         - Применить все миграции"
	@printf '%b\n' "  $(YELLOW)make migrate-down$(RESET)       - Откатить все миграции"
	@printf '%b\n' "  $(YELLOW)make migrate-status$(RESET)     - Статус миграций всех сервисов"
	@printf '%b\n' "  $(YELLOW)make db-clean$(RESET)           - Очистить все данные в БД"
	@printf '%b\n' "  $(YELLOW)make db-seed$(RESET)            - Заполнить БД тестовыми данными"
	@printf '\n'
	@printf '%b\n' "$(GREEN)📊 Мониторинг и логи:$(RESET)"
	@printf '%b\n' "  $(YELLOW)make logs$(RESET)               - Показать логи всех сервисов"
	@printf '%b\n' "  $(YELLOW)make logs-follow$(RESET)        - Следить за логами в реальном времени"
	@printf '%b\n' "  $(YELLOW)make ps$(RESET)                 - Список всех контейнеров проекта"
	@printf '\n'
	@printf '%b\n' "$(GREEN)🧹 Очистка:$(RESET)"
	@printf '%b\n' "  $(YELLOW)make clean$(RESET)              - Остановить и удалить все контейнеры"
	@printf '%b\n' "  $(YELLOW)make clean-all$(RESET)          - Полная очистка (контейнеры + образы + volumes)"
	@printf '%b\n' "  $(YELLOW)make clean-images$(RESET)       - Удалить Docker образы проекта"
	@printf '%b\n' "  $(YELLOW)make clean-local$(RESET)        - Очистить локальные артефакты (bin, dist, node_modules)"
	@printf '\n'
	@printf '%b\n' "$(GREEN)🏗️  Сборка:$(RESET)"
	@printf '%b\n' "  $(YELLOW)make build$(RESET)              - Пересобрать все сервисы (Docker)"
	@printf '%b\n' "  $(YELLOW)make build-backend$(RESET)      - Пересобрать только backend сервисы (Docker)"
	@printf '%b\n' "  $(YELLOW)make build-local$(RESET)        - Собрать всё локально (Go + Frontend)"
	@printf '%b\n' "  $(YELLOW)make install-local$(RESET)      - Установить зависимости локально"
	@printf '\n'
	@printf '%b\n' "$(GREEN)🔍 Диагностика:$(RESET)"
	@printf '%b\n' "  $(YELLOW)make health$(RESET)             - Проверить здоровье всех сервисов"
	@printf '%b\n' "  $(YELLOW)make check-deps$(RESET)         - Проверить зависимости (Docker, Go, etc.)"
	@printf '\n'
	@printf '%b\n' "$(GREEN)🐍 Python утилиты:$(RESET)"
	@printf '%b\n' "  $(YELLOW)make prepare-python$(RESET)     - Создать Python venv с зависимостями"
	@printf '%b\n' "  $(YELLOW)make send-test-data$(RESET)     - Запустить загрузку тестовых данных"
	@printf '%b\n' "  $(YELLOW)make generate-db-html$(RESET)   - Генерировать HTML таблиц БД"
	@printf '\n'
	@printf '%b\n' "$(GREEN)🔧 Go утилиты:$(RESET)"
	@printf '%b\n' "  $(YELLOW)make go-tidy$(RESET)            - go mod tidy для всех модулей"
	@printf '%b\n' "  $(YELLOW)make go-fmt$(RESET)             - go fmt для всех модулей"
	@printf '%b\n' "  $(YELLOW)make go-vet$(RESET)             - go vet для всех модулей"
	@printf '%b\n' "  $(YELLOW)make go-test$(RESET)            - go test для всех модулей"
	@printf '%b\n' "  $(YELLOW)make go-build$(RESET)           - go build для всех backend сервисов"
	@printf '%b\n' "  $(YELLOW)make go-clean$(RESET)           - Очистить Go кэш и бинарники"
	@printf '\n'
	@printf '%b\n' "$(GREEN)🔑 Тестовые доступы:$(RESET)"
	@printf '%b\n' "  $(CYAN)Customer:$(RESET)   $(YELLOW)$(TEST_CUSTOMER_LOGIN)$(RESET) / $(YELLOW)$(TEST_CUSTOMER_PASS)$(RESET)"
	@printf '%b\n' "  $(CYAN)Courier:$(RESET)    $(YELLOW)$(TEST_COURIER_LOGIN)$(RESET) / $(YELLOW)$(TEST_COURIER_PASS)$(RESET)"
	@printf '%b\n' "  $(CYAN)Restaurant:$(RESET) $(YELLOW)$(TEST_RESTAURANT_LOGIN)$(RESET) / $(YELLOW)$(TEST_RESTAURANT_PASS)$(RESET)"
	@printf '\n'
	@printf '%b\n' "$(GREEN)🌐 Frontend:$(RESET) http://localhost:5174/"
	@printf '\n'

# ============================================================================
# ОРКЕСТРАЦИЯ МИКРОСЕРВИСОВ
# ============================================================================

# --- Запуск всей системы ---

start: check-docker
	@echo -e "$(GREEN)🚀 Запуск всех сервисов (полный деплой)...$(RESET)"
	@cd restaurant && $(MAKE) launch-db
	@cd customer && $(MAKE) launch-redis
	@cd customer && $(MAKE) launch-minio
	@cd customer && $(MAKE) launch-rabbitmq
	@cd migrations && $(MAKE) up
	@cd restaurant && $(MAKE) migrate-up-orders-db
	@cd customer && $(MAKE) build-customer && $(MAKE) launch-customer
	@cd courier && $(MAKE) build-courier && $(MAKE) launch-courier
	@cd restaurant && $(MAKE) build-restaurant && $(MAKE) launch-restaurant
	@cd front && $(MAKE) run
	@echo -e "$(GREEN)✅ Все сервисы запущены!$(RESET)"
	@$(MAKE) status

start-dev: check-docker init-env setup-python-env
	@echo -e "$(GREEN)🚀 Запуск всех сервисов с тестовыми данными...$(RESET)"
	@cd restaurant && $(MAKE) launch-db
	@cd customer && $(MAKE) launch-redis
	@cd customer && $(MAKE) launch-minio
	@cd customer && $(MAKE) launch-rabbitmq
	@cd migrations && $(MAKE) up
	@cd restaurant && $(MAKE) migrate-up-orders-db
	@cd customer && $(MAKE) build-customer && $(MAKE) launch-customer
	@cd courier && $(MAKE) build-courier && $(MAKE) launch-courier
	@cd restaurant && $(MAKE) build-restaurant && $(MAKE) launch-restaurant
	@cd front && $(MAKE) run
	@echo -e "$(GREEN)✅ Все сервисы запущены в dev режиме!$(RESET)"
	@sleep 3
	@$(MAKE) status
	@if [ -f "send_requests.py" ]; then \
		echo -e "$(GREEN)🌱 Загрузка тестовых данных...$(RESET)"; \
		. .venv-yafds/bin/activate && python send_requests.py; \
		echo -e "$(GREEN)✅ Тестовые данные загружены!$(RESET)"; \
	else \
		echo -e "$(YELLOW)⚠️  send_requests.py не найден, пропускаем$(RESET)"; \
	fi

run: init-env clean-all start-dev

stop:
	@echo -e "$(YELLOW)🛑 Остановка всех сервисов...$(RESET)"
	@for service in $(SERVICES); do \
		echo -e "$(CYAN)Останавливаем $$service...$(RESET)"; \
		cd $$service && $(MAKE) clear-all 2>/dev/null || true; \
		cd ..; \
	done
	@echo -e "$(GREEN)✅ Все сервисы остановлены$(RESET)"

restart: stop
	@sleep 2
	@$(MAKE) start

# --- Customer ---

start-customer:
	@cd customer && $(MAKE) run

stop-customer:
	@cd customer && $(MAKE) clear-all

restart-customer: stop-customer
	@sleep 1
	@$(MAKE) start-customer

logs-customer:
	@cd customer && $(MAKE) logs-customer

build-customer:
	@cd customer && $(MAKE) build-customer

# --- Courier ---

start-courier:
	@cd courier && $(MAKE) run

stop-courier:
	@cd courier && $(MAKE) clear-all

restart-courier: stop-courier
	@sleep 1
	@$(MAKE) start-courier

logs-courier:
	@cd courier && $(MAKE) logs-courier

build-courier:
	@cd courier && $(MAKE) build-courier

# --- Restaurant ---

start-restaurant:
	@cd restaurant && $(MAKE) run

stop-restaurant:
	@cd restaurant && $(MAKE) clear-all

restart-restaurant: stop-restaurant
	@sleep 1
	@$(MAKE) start-restaurant

logs-restaurant:
	@cd restaurant && $(MAKE) logs-restaurant

build-restaurant:
	@cd restaurant && $(MAKE) build-restaurant

# --- Front ---

start-front:
	@cd front && $(MAKE) run

stop-front:
	@cd front && $(MAKE) clean

restart-front: stop-front
	@sleep 1
	@$(MAKE) start-front

logs-front:
	@cd front && $(MAKE) logs

build-front:
	@cd front && $(MAKE) build

# --- Миграции и БД ---

migrate-up:
	@echo -e "$(GREEN)📊 Применение всех миграций...$(RESET)"
	@cd migrations && $(MAKE) up

migrate-down:
	@echo -e "$(YELLOW)📊 Откат всех миграций...$(RESET)"
	@cd migrations && $(MAKE) down

migrate-status:
	@echo -e "$(CYAN)📊 Статус миграций всех сервисов:$(RESET)"
	@cd migrations && $(MAKE) migrate-status

db-clean:
	@echo -e "$(YELLOW)🗑️  Очистка базы данных...$(RESET)"
	@cd migrations && $(MAKE) clean
	@echo -e "$(GREEN)✅ База данных очищена$(RESET)"

db-seed:
	@echo -e "$(GREEN)🌱 Заполнение БД тестовыми данными...$(RESET)"
	@cd customer && $(MAKE) seed-up
	@cd courier && $(MAKE) seed-up
	@cd restaurant && $(MAKE) seed-up
	@echo -e "$(GREEN)✅ Тестовые данные загружены$(RESET)"

# --- Сборка ---

build:
	@echo -e "$(CYAN)🔨 Пересборка всех сервисов...$(RESET)"
	@for service in $(SERVICES); do \
		echo -e "$(YELLOW)Сборка $$service...$(RESET)"; \
		cd $$service && $(MAKE) build-$$service 2>/dev/null || $(MAKE) build || true; \
		cd ..; \
	done
	@echo -e "$(GREEN)✅ Все сервисы собраны$(RESET)"

build-backend:
	@echo -e "$(CYAN)🔨 Пересборка backend сервисов...$(RESET)"
	@for service in $(BACKEND_SERVICES); do \
		echo -e "$(YELLOW)Сборка $$service...$(RESET)"; \
		cd $$service && $(MAKE) build-$$service; \
		cd ..; \
	done
	@echo -e "$(GREEN)✅ Backend сервисы собраны$(RESET)"

# --- Очистка ---

clean: stop
	@echo -e "$(RED)🧹 Удаление всех контейнеров...$(RESET)"
	@docker ps -a --filter "name=yafds" -q | xargs -r docker rm -f 2>/dev/null || true
	@echo -e "$(GREEN)✅ Контейнеры удалены$(RESET)"

clean-images:
	@echo -e "$(RED)🗑️  Удаление Docker образов...$(RESET)"
	@docker images --filter "reference=yafds*" -q | xargs -r docker rmi -f 2>/dev/null || true
	@echo -e "$(GREEN)✅ Образы удалены$(RESET)"

clean-all: clean clean-images
	@echo -e "$(RED)🗑️  Удаление volumes...$(RESET)"
	@docker volume ls --filter "name=yafds" -q | xargs -r docker volume rm 2>/dev/null || true
	@echo -e "$(GREEN)✅ Полная очистка завершена$(RESET)"

# --- Мониторинг ---

status:
	@echo -e "$(CYAN)📊 Статус контейнеров YAFDS:$(RESET)"
	@echo -e ""
	@docker ps -a --filter "name=yafds" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" || echo "Нет запущенных контейнеров"

ps:
	@docker ps -a --filter "name=yafds"

logs:
	@echo -e "$(CYAN)📜 Логи всех сервисов:$(RESET)"
	@for service in customer courier restaurant; do \
		echo -e "$(YELLOW)=== $$service ===$(RESET)"; \
		docker logs yafds-$$service-service 2>/dev/null | tail -n 20 || echo "Контейнер не запущен"; \
		echo ""; \
	done

logs-follow:
	@echo -e "$(CYAN)📜 Следим за логами (Ctrl+C для выхода)...$(RESET)"
	@docker logs -f yafds-customer-service 2>/dev/null || echo "Customer не запущен"

health:
	@echo -e "$(CYAN)🏥 Проверка здоровья сервисов:$(RESET)"
	@echo -e "$(YELLOW)Customer:$(RESET)"
	@curl -s http://localhost:8090/health 2>/dev/null && echo "$(GREEN)✅ OK$(RESET)" || echo "$(RED)❌ DOWN$(RESET)"
	@echo -e "$(YELLOW)Courier:$(RESET)"
	@curl -s http://localhost:8091/health 2>/dev/null && echo "$(GREEN)✅ OK$(RESET)" || echo "$(RED)❌ DOWN$(RESET)"
	@echo -e "$(YELLOW)Restaurant:$(RESET)"
	@curl -s http://localhost:8092/health 2>/dev/null && echo "$(GREEN)✅ OK$(RESET)" || echo "$(RED)❌ DOWN$(RESET)"
	@echo -e "$(YELLOW)Frontend:$(RESET)"
	@curl -s http://localhost:5173 2>/dev/null > /dev/null && echo "$(GREEN)✅ OK$(RESET)" || echo "$(RED)❌ DOWN$(RESET)"

# --- Диагностика ---

check-docker:
	@docker info > /dev/null 2>&1 || (echo -e "$(RED)❌ Docker не запущен!$(RESET)" && exit 1)
	@echo -e "$(GREEN)✅ Docker запущен$(RESET)"

check-deps:
	@echo -e "$(CYAN)🔍 Проверка зависимостей:$(RESET)"
	@echo -n "Docker: ";      docker --version 2>/dev/null       && echo -e "$(GREEN)✅$(RESET)" || echo -e "$(RED)❌$(RESET)"
	@echo -n "Docker Compose: "; docker compose version 2>/dev/null && echo -e "$(GREEN)✅$(RESET)" || echo -e "$(RED)❌$(RESET)"
	@echo -n "Go: ";          go version 2>/dev/null             && echo -e "$(GREEN)✅$(RESET)" || echo -e "$(RED)❌$(RESET)"
	@echo -n "Node.js: ";     node --version 2>/dev/null         && echo -e "$(GREEN)✅$(RESET)" || echo -e "$(RED)❌$(RESET)"
	@echo -n "Goose: ";       goose --version 2>/dev/null        && echo -e "$(GREEN)✅$(RESET)" || echo -e "$(RED)❌$(RESET)"

# --- Утилиты ---

init-env:
	@echo -e "$(CYAN)📝 Копирование .env.example → .env для всех сервисов...$(RESET)"
	@for service in $(BACKEND_SERVICES); do \
		(cd $$service && $(MAKE) copy-env); \
	done
	@echo -e "$(GREEN)✅ Все .env файлы созданы$(RESET)"

init: init-env

update-deps:
	@./update_deps.sh

test:
	@echo -e "$(CYAN)🧪 Запуск тестов...$(RESET)"
	@for dir in $(GO_MODULES); do \
		echo -e "$(YELLOW)Тестирование $$dir...$(RESET)"; \
		cd $$dir && go test ./... -v || true; \
		cd ..; \
	done

quick-restart: stop-customer stop-courier stop-restaurant
	@sleep 1
	@$(MAKE) start-customer &
	@$(MAKE) start-courier &
	@$(MAKE) start-restaurant &
	@wait
	@echo -e "$(GREEN)✅ Backend сервисы перезапущены$(RESET)"

show-env:
	@echo -e "$(CYAN)📋 Переменные окружения сервисов:$(RESET)"
	@for service in $(BACKEND_SERVICES); do \
		echo -e "$(YELLOW)=== $$service ===$(RESET)"; \
		cat $$service/config/.env 2>/dev/null || echo "Файл .env не найден"; \
		echo ""; \
	done

# ============================================================================
# ЛОКАЛЬНЫЙ ЗАПУСК
# ============================================================================

install-local:
	@echo -e "$(CYAN)📦 Установка зависимостей локально...$(RESET)"
	@for dir in $(GO_MODULES); do \
		echo -e "  $$dir..."; \
		cd $$dir && go mod download > /dev/null 2>&1 || true; \
		cd ..; \
	done
	@cd front && npm install --silent
	@echo -e "$(GREEN)✅ Все зависимости установлены$(RESET)"

build-local: install-local
	@echo -e "$(CYAN)🔨 Локальная сборка всех сервисов...$(RESET)"
	@for service in $(BACKEND_SERVICES); do \
		echo -e "  $$service..."; \
		cd $$service && go build -o bin/$$service ./cmd/server && echo -e "    $(GREEN)✅$(RESET)" || echo -e "    $(RED)❌$(RESET)"; \
		cd ..; \
	done
	@cd front && npm run build > /dev/null 2>&1 && echo -e "  $(GREEN)✅ front/dist$(RESET)" || echo -e "  $(RED)❌$(RESET)"
	@echo -e "$(GREEN)✅ Локальная сборка завершена!$(RESET)"

run-local: clean-local install-local
	@echo -e "$(GREEN)🚀 Запуск всех сервисов локально (без Docker)...$(RESET)"
	@cd customer && $(MAKE) run-local &
	@cd courier && $(MAKE) run-local &
	@cd restaurant && $(MAKE) run-local &
	@sleep 2
	@cd front && $(MAKE) dev &
	@echo -e "$(GREEN)✅ Все сервисы запущены локально!$(RESET)"
	@echo -e "$(CYAN)📍 Frontend: http://localhost:5173 | Customer: :8091 | Courier: :8090 | Restaurant: :8092$(RESET)"
	@wait

stop-local:
	@echo -e "$(YELLOW)🛑 Остановка всех локальных сервисов...$(RESET)"
	@pkill -f "customer/cmd/server" 2>/dev/null || true
	@pkill -f "courier/cmd/server" 2>/dev/null || true
	@pkill -f "restaurant/cmd/server" 2>/dev/null || true
	@pkill -f "vite" 2>/dev/null || true
	@echo -e "$(GREEN)✅ Все локальные сервисы остановлены$(RESET)"

clean-local:
	@echo -e "$(RED)🧹 Очистка локальных артефактов...$(RESET)"
	@for service in $(BACKEND_SERVICES); do rm -rf $$service/bin 2>/dev/null || true; done
	@rm -rf front/dist front/node_modules 2>/dev/null || true
	@echo -e "$(GREEN)✅ Локальные артефакты очищены$(RESET)"

# ============================================================================
# PYTHON УТИЛИТЫ
# ============================================================================

setup-python-env:
	@echo -e "$(CYAN)🐍 Подготовка Python окружения...$(RESET)"
	@if [ ! -d ".venv-yafds" ]; then \
		python3 -m venv .venv-yafds; \
		echo -e "$(GREEN)✅ Виртуальное окружение создано$(RESET)"; \
	fi
	@. .venv-yafds/bin/activate && pip install --upgrade pip setuptools wheel > /dev/null 2>&1
	@. .venv-yafds/bin/activate && pip install tqdm colorama pandas psycopg2-binary requests > /dev/null 2>&1
	@echo -e "$(GREEN)✅ Зависимости Python установлены$(RESET)"

prepare-python: setup-python-env
	@echo -e "$(GREEN)✅ Python окружение готово$(RESET)"
	@echo -e "$(CYAN)Активация: $(YELLOW)source .venv-yafds/bin/activate$(RESET)"

send-test-data:
	@[ -d ".venv-yafds" ] || (echo -e "$(RED)❌ Запустите make prepare-python$(RESET)" && exit 1)
	@[ -f "send_requests.py" ] || (echo -e "$(RED)❌ send_requests.py не найден$(RESET)" && exit 1)
	@echo -e "$(GREEN)🌱 Загрузка тестовых данных...$(RESET)"
	@. .venv-yafds/bin/activate && python send_requests.py
	@echo -e "$(GREEN)✅ Тестовые данные загружены!$(RESET)"

generate-db-html:
	@[ -d ".venv-yafds" ] || (echo -e "$(RED)❌ Запустите make prepare-python$(RESET)" && exit 1)
	@[ -f "make_nice_sql_html.py" ] || (echo -e "$(RED)❌ make_nice_sql_html.py не найден$(RESET)" && exit 1)
	@echo -e "$(GREEN)📊 Генерация HTML с содержимым БД...$(RESET)"
	@. .venv-yafds/bin/activate && python make_nice_sql_html.py
	@echo -e "$(GREEN)✅ HTML файл сгенерирован$(RESET)"

# ============================================================================
# GO УТИЛИТЫ
# ============================================================================

go-tidy:
	@echo -e "$(CYAN)🔨 go mod tidy для всех Go модулей...$(RESET)"
	@for dir in $(GO_MODULES); do \
		cd $$dir && go mod tidy && echo -e "$(GREEN)✅ $$dir$(RESET)" || echo -e "$(RED)❌ $$dir$(RESET)"; \
		cd ..; \
	done

go-fmt:
	@echo -e "$(CYAN)🎨 go fmt для всех Go модулей...$(RESET)"
	@for dir in $(GO_MODULES); do \
		cd $$dir && go fmt ./... && echo -e "$(GREEN)✅ $$dir$(RESET)" || echo -e "$(RED)❌ $$dir$(RESET)"; \
		cd ..; \
	done

go-vet:
	@echo -e "$(CYAN)🔍 go vet для всех Go модулей...$(RESET)"
	@for dir in $(GO_MODULES); do \
		cd $$dir && go vet ./... && echo -e "$(GREEN)✅ $$dir$(RESET)" || echo -e "$(YELLOW)⚠️  $$dir$(RESET)"; \
		cd ..; \
	done

go-test:
	@echo -e "$(CYAN)🧪 go test для всех Go модулей...$(RESET)"
	@for dir in $(GO_MODULES); do \
		cd $$dir && go test -v -cover ./... || echo -e "$(YELLOW)⚠️  $$dir$(RESET)"; \
		cd ..; \
	done

go-build:
	@echo -e "$(CYAN)🏗️  go build для всех backend сервисов...$(RESET)"
	@for service in $(BACKEND_SERVICES); do \
		cd $$service && go build -o bin/$$service ./cmd/server && echo -e "$(GREEN)✅ $$service$(RESET)" || echo -e "$(RED)❌ $$service$(RESET)"; \
		cd ..; \
	done

go-clean:
	@echo -e "$(YELLOW)🧹 Очистка Go кэша и бинарников...$(RESET)"
	@go clean -cache
	@for service in $(BACKEND_SERVICES); do \
		cd $$service && go clean -v && rm -rf bin/ || true; \
		cd ..; \
	done
	@echo -e "$(GREEN)✅ Очистка завершена$(RESET)"

go-deps-check:
	@echo -e "$(CYAN)📦 Проверка зависимостей Go модулей...$(RESET)"
	@for dir in $(GO_MODULES); do \
		cd $$dir && go list -u -m all | grep -v indirect || true; \
		cd ..; \
	done

go-upgrade-deps:
	@echo -e "$(CYAN)📦 Обновление зависимостей Go модулей...$(RESET)"
	@for dir in $(GO_MODULES); do \
		cd $$dir && go get -u ./... && go mod tidy && echo -e "$(GREEN)✅ $$dir$(RESET)" || echo -e "$(RED)❌ $$dir$(RESET)"; \
		cd ..; \
	done
