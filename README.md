# Pull Request Assign Service

Сервис назначения ревьюеров для Pull Request’ов

## Архитектура

- **Backend**: Go
- **Database**: PostgreSQL
- **Containerization**: Docker + Docker Compose

## Быстрый старт

### Предварительные требования

- Docker и Docker Compose
- Go 1.21+

### Запуск через Docker Compose

1. Клонируйте репозиторий:
```bash
git clone <repository-url>
cd pr-assign-service
```

2. Запустите все сервисы:
```bash
docker-compose up --build
```

3. Откройте приложение:
- Backend API: http://localhost:8080
- Database: localhost:5432

### Структура сервиса

- **postgres**: База данных PostgreSQL
- **backend**: Go API сервер

### Переменные окружения

- `HTTP_SERVER_ADDRESS`: Адрес сервера 
- `HTTP_SERVER_PORT`: Порт сервера
- `HTTP_SERVER_TIMEOUT`: Таймаут обработки запроса
- `DB_NAME`: Имя базы данных
- `DB_USER`: Имя пользователя БД
- `DB_PASSWORD`: Пароль пользователя БД
- `DB_HOST`: Адрес PostgreSQL
- `DB_PORT`: Порт PostgreSQL

### API Endpoints

#### Пользователи
- `POST /users/setIsActive` - Установить флаг активности пользователя
- `GET /users/getReview` - Получить PR'ы, где пользователь назначен ревьювером

#### Команды
- `POST /team/add` - Создать команду с участниками
- `GET /team/get` - Получить команду с участниками

#### Pull Request'ы
- `POST /pullRequest/create` - Создать PR и автоматически назначить до 2 ревьюверов из команды автора
- `POST /pullRequest/merge` - Пометить PR как MERGED (идемпотентная операция)
- `POST /pullRequest/reassign` - Переназначить конкретного ревьювера на другого из его команды

#### Статистика
- `GET /stats` - Статистика сервиса

### Структура проекта

```
pr-assign_service/
├── cmd/pr-assign-service/ # Точка входа
├── docs/           # Документация OpenAPI
├── internal/           # Внутренние пакеты
│   ├── api/           # HTTP API
│      ├── handlers/     # HTTP handlers
│   ├── app/           # Инициализация приложения
│   ├── config/       # Конфигурация
│   ├── domain/        # Модели данных
│   ├── repository/    # Слой данных
│   └──service/      # Бизнес-логика
├── migrations/        # SQL миграции
├── docker-compose.yml    # Docker Compose конфигурация
├── Dockerfile        # Docker образ
└── README.md            # Документация
```

### Тестирование

```bash
make test
```
