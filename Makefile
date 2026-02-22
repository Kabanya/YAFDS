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
# Все директории с Go кодом (включая общие пакеты)
GO_MODULES := pkg customer courier restaurant

# Тестовые доступы
TEST_CUSTOMER_LOGIN := "0x742d35Cc6634C0532925a3b844Bc454e4438f44e"
TEST_CUSTOMER_PASS  := "StrongPass001"
TEST_COURIER_LOGIN  := "0x5647445a3e564b189fc35a75079fc924"
TEST_COURIER_PASS   := "aad9D6C7e9R7"
TEST_RESTAURANT_LOGIN := "0xbd2f054c1faa44b2a561c47dd1a0a367"
TEST_RESTAURANT_PASS  := "JO88M0Rtf0BG"

# ============================================================================
# HELP - Главное меню помощи
# ============================================================================

help:
	@echo -e "$(CYAN)    ╔══════════════════════════════════════════════════════╗$(RESET)"
	@echo -e "$(CYAN)    ║              YAFDS Project Orchestrator              ║$(RESET)"
	@echo -e "$(CYAN)    ╚══════════════════════════════════════════════════════╝$(RESET)"
	@echo -e ""
	@echo -e "$(GREEN)📦 Основные команды:$(RESET)"
	@echo -e "  $(YELLOW)make start$(RESET)              - Запустить все сервисы (полный деплой)"
	@echo -e "  $(YELLOW)make start-dev$(RESET)          - Запустить все сервисы с тестовыми данными"
	@echo -e "  $(YELLOW)make init-env$(RESET)           - Инициализировать .env файлы из .env.example"
	@echo -e "  $(YELLOW)make init$(RESET)               - Инициализировать .env файлы (в будущем расширить)"
	@echo -e "  $(YELLOW)make stop$(RESET)               - Остановить все сервисы"
	@echo -e "  $(YELLOW)make restart$(RESET)            - Перезапустить все сервисы"
	@echo -e "  $(YELLOW)make status$(RESET)             - Показать статус всех контейнеров"
	@echo -e ""
	@echo -e "$(GREEN)🔧 Сервисы (customer, courier, restaurant, front):$(RESET)"
	@echo -e "  $(YELLOW)make start-<service>$(RESET)    - Запустить конкретный сервис"
	@echo -e "  $(YELLOW)make stop-<service>$(RESET)     - Остановить конкретный сервис"
	@echo -e "  $(YELLOW)make restart-<service>$(RESET)  - Перезапустить конкретный сервис"
	@echo -e "  $(YELLOW)make logs-<service>$(RESET)     - Показать логи конкретного сервиса"
	@echo -e "  $(YELLOW)make build-<service>$(RESET)    - Пересобрать конкретный сервис"
	@echo -e ""
	@echo -e "$(GREEN)🗄️  База данных и миграции:$(RESET)"
	@echo -e "  $(YELLOW)make migrate-up$(RESET)         - Применить все миграции"
	@echo -e "  $(YELLOW)make migrate-down$(RESET)       - Откатить все миграции"
	@echo -e "  $(YELLOW)make migrate-status$(RESET)     - Статус миграций всех сервисов"
	@echo -e "  $(YELLOW)make db-clean$(RESET)           - Очистить все данные в БД"
	@echo -e "  $(YELLOW)make db-seed$(RESET)            - Заполнить БД тестовыми данными"
	@echo -e ""
	@echo -e "$(GREEN)📊 Мониторинг и логи:$(RESET)"
	@echo -e "  $(YELLOW)make logs$(RESET)               - Показать логи всех сервисов"
	@echo -e "  $(YELLOW)make logs-follow$(RESET)        - Следить за логами в реальном времени"
	@echo -e "  $(YELLOW)make ps$(RESET)                 - Список всех контейнеров проекта"
	@echo -e ""
	@echo -e "$(GREEN)🧹 Очистка:$(RESET)"
	@echo -e "  $(YELLOW)make clean$(RESET)              - Остановить и удалить все контейнеры"
	@echo -e "  $(YELLOW)make clean-all$(RESET)          - Полная очистка (контейнеры + образы + volumes)"
	@echo -e "  $(YELLOW)make clean-images$(RESET)       - Удалить Docker образы проекта"
	@echo -e ""
	@echo -e "$(GREEN)🏗️  Сборка:$(RESET)"
	@echo -e "  $(YELLOW)make build$(RESET)              - Пересобрать все сервисы"
	@echo -e "  $(YELLOW)make build-backend$(RESET)      - Пересобрать только backend сервисы"
	@echo -e ""
	@echo -e "$(GREEN)🔍 Диагностика:$(RESET)"
	@echo -e "  $(YELLOW)make health$(RESET)             - Проверить здоровье всех сервисов"
	@echo -e "  $(YELLOW)make check-deps$(RESET)         - Проверить зависимости (Docker, Go, etc.)"
	@echo -e ""
	@echo -e "$(GREEN)🐍 Python утилиты:$(RESET)"
	@echo -e "  $(YELLOW)make prepare-python$(RESET)     - Создать Python venv с зависимостями"
	@echo -e "  $(YELLOW)make send-test-data$(RESET)     - Запустить загрузку тестовых данных"
	@echo -e "  $(YELLOW)make generate-db-html$(RESET)   - Генерировать HTML таблиц БД"
	@echo -e ""
	@echo -e "$(GREEN)🔧 Go утилиты:$(RESET)"
	@echo -e "  $(YELLOW)make go-tidy$(RESET)            - Выполнить go mod tidy для всех backend сервисов"
	@echo -e "  $(YELLOW)make go-fmt$(RESET)             - Отформатировать код (go fmt) во всех backend сервисах"
	@echo -e "  $(YELLOW)make go-vet$(RESET)             - Проверить код (go vet) во всех backend сервисах"
	@echo -e "  $(YELLOW)make go-test$(RESET)            - Запустить тесты Go (go test) для всех backend сервисов"
	@echo -e "  $(YELLOW)make go-build$(RESET)           - Собрать бинарники (go build) для всех backend сервисов"
	@echo -e "  $(YELLOW)make go-clean$(RESET)           - Очистить Go кэш и бинарники"
	@echo -e ""
	@echo -e "$(GREEN)🔑 Тестовые доступы (Customer, Courier, Restaurant):$(RESET)"
	@echo -e "  $(CYAN)Customer:$(RESET)   Wallet: $(YELLOW)$(TEST_CUSTOMER_LOGIN)$(RESET), Password: $(YELLOW)$(TEST_CUSTOMER_PASS)$(RESET)"
	@echo -e "  $(CYAN)Courier:$(RESET)    Wallet: $(YELLOW)$(TEST_COURIER_LOGIN)$(RESET), Password: $(YELLOW)$(TEST_COURIER_PASS)$(RESET)"
	@echo -e "  $(CYAN)Restaurant:$(RESET) Wallet: $(YELLOW)$(TEST_RESTAURANT_LOGIN)$(RESET), Password: $(YELLOW)$(TEST_RESTAURANT_PASS)$(RESET)"
	@echo -e ""

# ============================================================================
# ОСНОВНЫЕ КОМАНДЫ
# ============================================================================

setup-python-env:
	@echo -e "$(CYAN)🐍 Подготовка Python окружения...$(RESET)"
	@if [ ! -d ".venv-yafds" ]; then \
		echo -e "$(YELLOW)Создание виртуального окружения .venv-yafds...$(RESET)"; \
		python3 -m venv .venv-yafds; \
		echo -e "$(GREEN)✅ Виртуальное окружение создано$(RESET)"; \
	fi
	@echo -e "$(YELLOW)Установка зависимостей...$(RESET)"
	@. .venv-yafds/bin/activate && pip install --upgrade pip setuptools wheel > /dev/null 2>&1
	@. .venv-yafds/bin/activate && pip install tqdm colorama pandas psycopg2-binary requests > /dev/null 2>&1
	@echo -e "$(GREEN)✅ Зависимости установлены: tqdm, colorama, pandas, psycopg2, requests$(RESET)"

start: check-docker
	@echo -e "$(GREEN)🚀 Запуск всех сервисов (полный деплой)...$(RESET)"
	@cd customer && $(MAKE) run
	@cd courier && $(MAKE) run
	@cd restaurant && $(MAKE) run
	@cd front && $(MAKE) run
	@echo -e "$(GREEN)✅ Все сервисы запущены!$(RESET)"
	@$(MAKE) status

start-dev: check-docker init-env setup-python-env
	@echo -e "$(GREEN)🚀 Запуск всех сервисов с тестовыми данными...$(RESET)"
	@cd customer && $(MAKE) dev
	@cd courier && $(MAKE) dev
	@cd restaurant && $(MAKE) dev
	@cd front && $(MAKE) run
	@echo -e "$(GREEN)✅ Все сервисы запущены в dev режиме!$(RESET)"
	@sleep 3
	@$(MAKE) status
	@echo -e ""
	@if [ -f "send_requests.py" ]; then \
		echo -e "$(GREEN)🌱 Загрузка тестовых данных через send_requests.py...$(RESET)"; \
		. .venv-yafds/bin/activate && python send_requests.py; \
		echo -e "$(GREEN)✅ Тестовые данные загружены!$(RESET)"; \
	else \
		echo -e "$(YELLOW)⚠️  send_requests.py не найден, пропускаем загрузку тестовых данных$(RESET)"; \
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

# ============================================================================
# УПРАВЛЕНИЕ ОТДЕЛЬНЫМИ СЕРВИСАМИ
# ============================================================================

start-customer:
	@echo -e "$(GREEN)🚀 Запуск Customer сервиса...$(RESET)"
	@cd customer && $(MAKE) run

stop-customer:
	@echo -e "$(YELLOW)🛑 Остановка Customer сервиса...$(RESET)"
	@cd customer && $(MAKE) clear-all

restart-customer: stop-customer
	@sleep 1
	@$(MAKE) start-customer

logs-customer:
	@cd customer && $(MAKE) logs-customer

build-customer:
	@echo -e "$(CYAN)🔨 Сборка Customer сервиса...$(RESET)"
	@cd customer && $(MAKE) build-customer

# Courier
start-courier:
	@echo -e "$(GREEN)🚀 Запуск Courier сервиса...$(RESET)"
	@cd courier && $(MAKE) run

stop-courier:
	@echo -e "$(YELLOW)🛑 Остановка Courier сервиса...$(RESET)"
	@cd courier && $(MAKE) clear-all

restart-courier: stop-courier
	@sleep 1
	@$(MAKE) start-courier

logs-courier:
	@cd courier && $(MAKE) logs-courier

build-courier:
	@echo -e "$(CYAN)🔨 Сборка Courier сервиса...$(RESET)"
	@cd courier && $(MAKE) build-courier

# Restaurant
start-restaurant:
	@echo -e "$(GREEN)🚀 Запуск Restaurant сервиса...$(RESET)"
	@cd restaurant && $(MAKE) run

stop-restaurant:
	@echo -e "$(YELLOW)🛑 Остановка Restaurant сервиса...$(RESET)"
	@cd restaurant && $(MAKE) clear-all

restart-restaurant: stop-restaurant
	@sleep 1
	@$(MAKE) start-restaurant

logs-restaurant:
	@cd restaurant && $(MAKE) logs-restaurant

build-restaurant:
	@echo -e "$(CYAN)🔨 Сборка Restaurant сервиса...$(RESET)"
	@cd restaurant && $(MAKE) build-restaurant

# Front
start-front:
	@echo -e "$(GREEN)🚀 Запуск Frontend...$(RESET)"
	@cd front && $(MAKE) run

stop-front:
	@echo -e "$(YELLOW)🛑 Остановка Frontend...$(RESET)"
	@cd front && $(MAKE) clean

restart-front: stop-front
	@sleep 1
	@$(MAKE) start-front

logs-front:
	@cd front && $(MAKE) logs

build-front:
	@echo -e "$(CYAN)🔨 Сборка Frontend...$(RESET)"
	@cd front && $(MAKE) build

# ============================================================================
# МИГРАЦИИ И БАЗА ДАННЫХ
# ============================================================================

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

# ============================================================================
# МОНИТОРИНГ И ЛОГИ
# ============================================================================

status:
	@echo -e "$(CYAN)📊 Статус контейнеров YAFDS:$(RESET)"
	@echo -e ""
	@docker ps -a --filter "name=yafds" --format "table {{.Names}}\t{{.Status}}\t{{.Ports}}" || echo "Нет запущенных контейнеров"

ps:
	@docker ps -a --filter "name=yafds"

logs:
	@echo -e "$(CYAN)📜 Логи всех сервисов:$(RESET)"
	@echo -e ""
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
	@echo -e ""
	@echo -e "$(YELLOW)Customer:$(RESET)"
	@curl -s http://localhost:8090/health 2>/dev/null && echo "$(GREEN)✅ OK$(RESET)" || echo "$(RED)❌ DOWN$(RESET)"
	@echo -e ""
	@echo -e "$(YELLOW)Courier:$(RESET)"
	@curl -s http://localhost:8091/health 2>/dev/null && echo "$(GREEN)✅ OK$(RESET)" || echo "$(RED)❌ DOWN$(RESET)"
	@echo -e ""
	@echo -e "$(YELLOW)Restaurant:$(RESET)"
	@curl -s http://localhost:8092/health 2>/dev/null && echo "$(GREEN)✅ OK$(RESET)" || echo "$(RED)❌ DOWN$(RESET)"
	@echo -e ""
	@echo -e "$(YELLOW)Frontend:$(RESET)"
	@curl -s http://localhost:5173 2>/dev/null > /dev/null && echo "$(GREEN)✅ OK$(RESET)" || echo "$(RED)❌ DOWN$(RESET)"

# ============================================================================
# СБОРКА
# ============================================================================

build:
	@echo -e "$(CYAN)🔨 Пересборка всех сервисов...$(RESET)"
	@for service in $(SERVICES); do \
		@echo -e "$(YELLOW)Сборка $$service...$(RESET)"; \
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

# ============================================================================
# ОЧИСТКА
# ============================================================================

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

# ============================================================================
# ДИАГНОСТИКА И ПРОВЕРКИ
# ============================================================================

check-docker:
	@docker info > /dev/null 2>&1 || (echo -e "$(RED)❌ Docker не запущен! Запустите Docker Desktop.$(RESET)" && exit 1)
	@echo -e "$(GREEN)✅ Docker запущен$(RESET)"

check-deps:
	@echo -e "$(CYAN)🔍 Проверка зависимостей:$(RESET)"
	@echo -e ""
	@echo -n "Docker: "
	@docker --version 2>/dev/null && echo -e "$(GREEN)✅$(RESET)" || echo -e "$(RED)❌$(RESET)"
	@echo -n "Docker Compose: "
	@docker compose version 2>/dev/null && echo -e "$(GREEN)✅$(RESET)" || echo -e "$(RED)❌$(RESET)"
	@echo -n "Go: "
	@go version 2>/dev/null && echo -e "$(GREEN)✅$(RESET)" || echo -e "$(RED)❌$(RESET)"
	@echo -n "Node.js: "
	@node --version 2>/dev/null && echo -e "$(GREEN)✅$(RESET)" || echo -e "$(RED)❌$(RESET)"
	@echo -n "Goose: "
	@goose --version 2>/dev/null && echo -e "$(GREEN)✅$(RESET)" || echo -e "$(RED)❌$(RESET)"
	@echo -e ""

# ============================================================================
# УТИЛИТЫ
# ============================================================================
init-env:
	@echo -e "$(CYAN)📝 Копирование .env.example в .env для всех сервисов...$(RESET)"
	@for service in $(BACKEND_SERVICES); do \
		echo -e "$(YELLOW)Копирование для $$service...$(RESET)"; \
		(cd $$service && $(MAKE) copy-env); \
	done
	@echo -e "$(GREEN)✅ Все .env файлы созданы$(RESET)"

init: init-env

update-deps:
	@echo -e "$(CYAN)📦 Обновление зависимостей...$(RESET)"
	@./update_deps.sh

test:
	@echo -e "$(CYAN)🧪 Запуск тестов...$(RESET)"
	@for dir in $(GO_MODULES); do \
		echo -e "$(YELLOW)Тестирование $$dir...$(RESET)"; \
		cd $$dir && go test ./... -v || true; \
		cd ..; \
	done

# Быстрый перезапуск для разработки
quick-restart: stop-customer stop-courier stop-restaurant
	@sleep 1
	@$(MAKE) start-customer &
	@$(MAKE) start-courier &
	@$(MAKE) start-restaurant &
	@wait
	@echo -e "$(GREEN)✅ Backend сервисы перезапущены$(RESET)"

# Показать все переменные окружения
show-env:
	@echo -e "$(CYAN)📋 Переменные окружения сервисов:$(RESET)"
	@for service in $(BACKEND_SERVICES); do \
		echo ""; \
		@echo -e "$(YELLOW)=== $$service ===$(RESET)"; \
		cat $$service/config/.env 2>/dev/null || echo "Файл .env не найден"; \
	done

# ============================================================================
# PYTHON УТИЛИТЫ
# ============================================================================

prepare-python: setup-python-env
	@echo -e "$(GREEN)✅ Python окружение готово$(RESET)"
	@echo -e "$(CYAN)📝 Используйте для активации:$(RESET)"
	@echo -e "  $(YELLOW)source .venv-yafds/bin/activate$(RESET)"

send-test-data:
	@if [ ! -d ".venv-yafds" ]; then \
		echo -e "$(RED)❌ Python окружение не найдено! Запустите $(YELLOW)make prepare-python$(RED)$(RESET)"; \
		exit 1; \
	fi
	@if [ ! -f "send_requests.py" ]; then \
		echo -e "$(RED)❌ send_requests.py не найден!$(RESET)"; \
		exit 1; \
	fi
	@echo -e "$(GREEN)🌱 Запуск загрузки тестовых данных...$(RESET)"
	@. .venv-yafds/bin/activate && python send_requests.py
	@echo -e "$(GREEN)✅ Тестовые данные загружены!$(RESET)"

generate-db-html:
	@if [ ! -d ".venv-yafds" ]; then \
		echo -e "$(RED)❌ Python окружение не найдено! Запустите $(YELLOW)make prepare-python$(RED)$(RESET)"; \
		exit 1; \
	fi
	@if [ ! -f "make_nice_sql_html.py" ]; then \
		echo -e "$(RED)❌ make_nice_sql_html.py не найден!$(RESET)"; \
		exit 1; \
	fi
	@echo -e "$(GREEN)📊 Генерация HTML с содержимым БД...$(RESET)"
	@. .venv-yafds/bin/activate && python make_nice_sql_html.py
	@echo -e "$(GREEN)✅ HTML файл сгенерирован: tables.html$(RESET)"

# ============================================================================
# GO УТИЛИТЫ
# ============================================================================

go-tidy:
	@echo -e "$(CYAN)🔨 Выполнение go mod tidy для всех Go модулей...$(RESET)"
	@for dir in $(GO_MODULES); do \
		echo -e "$(YELLOW)go mod tidy в $$dir...$(RESET)"; \
		cd $$dir && go mod tidy && echo -e "$(GREEN)✅ $$dir$(RESET)" || echo -e "$(RED)❌ $$dir$(RESET)"; \
		cd ..; \
	done
	@echo -e "$(GREEN)✅ go mod tidy завершено для всех модулей$(RESET)"

go-fmt:
	@echo -e "$(CYAN)🎨 Форматирование кода (go fmt) для всех Go модулей...$(RESET)"
	@for dir in $(GO_MODULES); do \
		echo -e "$(YELLOW)Форматирование $$dir...$(RESET)"; \
		cd $$dir && go fmt ./... && echo -e "$(GREEN)✅ $$dir$(RESET)" || echo -e "$(RED)❌ $$dir$(RESET)"; \
		cd ..; \
	done
	@echo -e "$(GREEN)✅ go fmt завершено для всех модулей$(RESET)"

go-vet:
	@echo -e "$(CYAN)🔍 Проверка кода (go vet) для всех Go модулей...$(RESET)"
	@for dir in $(GO_MODULES); do \
		echo -e "$(YELLOW)Проверка $$dir...$(RESET)"; \
		cd $$dir && go vet ./... && echo -e "$(GREEN)✅ $$dir$(RESET)" || echo -e "$(YELLOW)⚠️  $$dir (предупреждения)$(RESET)"; \
		cd ..; \
	done
	@echo -e "$(GREEN)✅ go vet завершено для всех модулей$(RESET)"

go-test:
	@echo -e "$(CYAN)🧪 Запуск тестов (go test) для всех Go модулей...$(RESET)"
	@for dir in $(GO_MODULES); do \
		echo -e "$(YELLOW)Тестирование $$dir...$(RESET)"; \
		cd $$dir && go test -v -cover ./... || echo -e "$(YELLOW)⚠️  Тесты в $$dir не прошли$(RESET)"; \
		cd ..; \
	done
	@echo -e "$(GREEN)✅ go test завершено для всех модулей$(RESET)"

go-build:
	@echo -e "$(CYAN)🏗️  Сборка бинарников (go build) для всех backend сервисов...$(RESET)"
	@for service in $(BACKEND_SERVICES); do \
		echo -e "$(YELLOW)Сборка $$service...$(RESET)"; \
		cd $$service && go build -o bin/$$service ./cmd/server && echo -e "$(GREEN)✅ $$service$(RESET)" || echo -e "$(RED)❌ $$service$(RESET)"; \
		cd ..; \
	done
	@echo -e "$(GREEN)✅ Сборка завершена для всех сервисов$(RESET)"

go-clean:
	@echo -e "$(YELLOW)🧹 Очистка Go кэша и бинарников...$(RESET)"
	@go clean -cache
	@for service in $(BACKEND_SERVICES); do \
		echo -e "$(YELLOW)Очистка $$service...$(RESET)"; \
		cd $$service && go clean -v && rm -rf bin/ || true; \
		cd ..; \
	done
	@echo -e "$(GREEN)✅ Очистка завершена$(RESET)"

go-deps-check:
	@echo -e "$(CYAN)📦 Проверка зависимостей для всех Go модулей...$(RESET)"
	@for dir in $(GO_MODULES); do \
		echo -e "$(YELLOW)Проверка $$dir...$(RESET)"; \
		cd $$dir && go list -u -m all | grep -v indirect || echo -e "$(GREEN)✅ Все зависимости актуальны$(RESET)"; \
		cd ..; \
	done
	@echo -e "$(GREEN)✅ Проверка завершена$(RESET)"

go-upgrade-deps:
	@echo -e "$(CYAN)📦 Обновление зависимостей для всех Go модулей...$(RESET)"
	@for dir in $(GO_MODULES); do \
		echo -e "$(YELLOW)Обновление $$dir...$(RESET)"; \
		cd $$dir && go get -u ./... && go mod tidy && echo -e "$(GREEN)✅ $$dir$(RESET)" || echo -e "$(RED)❌ $$dir$(RESET)"; \
		cd ..; \
	done
	@echo -e "$(GREEN)✅ Обновление завершено для всех модулей$(RESET)"
