# Song Library API
[![loon-icon.png](https://i.postimg.cc/k486BzSr/loon-icon.png)](https://postimg.cc/wyqTrWFF)

Сервис REST API для управления библиотекой песен. Поддерживает добавление, получение, обновление и удаление песен. Реализованы фильтрация, пагинация, интеграция с внешним API, хранение данных в PostgreSQL и документация с помощью Swagger.

---

# Запуск

Для запуска проекта создайте директорию "logs" в корневой папке и с помощью Makefile выполните команду:

```bash
make up
```
Требования
1. Выставить REST методы

    Получение данных библиотеки с фильтрацией по всем полям и пагинацией
    Получение текста песни с пагинацией по куплетам
    Удаление песни
    Изменение данных песни
    Добавление новой песни 

## Конечные точки API

```go
	r.PathPrefix("/swagger/").Handler(httpSwagger.WrapHandler).Methods("GET")
	r.HandleFunc("/songs", h.Create).Methods("POST")
	r.HandleFunc("/songs", h.Get).Methods("GET")
	r.HandleFunc("/songs/{id}", h.Update).Methods("PUT")
	r.HandleFunc("/songs/{id}", h.Delete).Methods("DELETE")
	r.HandleFunc("/songs/verses", h.GetVerses).Methods("GET")
```
## 1. Получение данных о песнях

### GET /songs

Получение списка песен с возможностью фильтрации и пагинации.

Параметры запроса:

    limit (int) - ограничение на количество песен (по умолчанию 10)
    offset (int) - смещение для пагинации
    id (int) - ID песни для фильтрации
    song (string) - название песни
    group_name (string) - название группы
    releasedate (string) - дата релиза песни в формате 02.01.2006

## 2. Добавление новой песни

### POST /songs

Добавление новой песни.

Тело запроса:

{
  "title": "название песни",
  "group_name": "группа",
  "releasedate": "дата релиза"
}

## 3. Обновление данных песни

### PUT /songs/{id}

Обновление информации о песне.

Параметры запроса:

    id (int) - ID песни для обновления 

Тело запроса:

{
  "title": "новое название песни",
  "group_name": "новая группа",
  "releasedate": "новая дата релиза"
}

## 4. Удаление песни

### DELETE /songs/{id}

Удаление песни по ID.

Параметры запроса:

    id (int) - ID песни для удаления 

## 5. Получение текста песни с пагинацией по куплетам

### GET /songs/verses

Получение текста песни с пагинацией по куплетам.

Параметры запроса:

    song_id (int) - ID песни для получения текста
    limit (int) - ограничение на количество куплетов
    offset (int) - смещение для пагинации 

# Интеграция с внешним API

Реализовал отдельный клиент для запросов в сторонее API по пути internal/client/client.go

# Работа с базой данных

Информация о песнях сохраняется в базе данных PostgreSQL. Структура базы данных создается с помощью миграций при старте сервиса.

Функция подключения к базе данных:
```go
func ConnectDB(connStr string) (*sql.DB, error) {
    db, err := sql.Open("postgres", connStr)
    if err != nil {
        return nil, err
    }

    migrationsDir := "./migrations"

    if err = initMigrations(db, migrationsDir); err != nil {
        return nil, err
    }

    return db, nil
}
```
# Покрыть код debug- и info-логами

В проекте использованы логирования с уровнями debug и info в репозиториях и ключевых местах приложения для отслеживания важных событий и ошибок.

# Конфигурационные данные

Все базовые конфигурационные данные для работы приложения вынесены в .env файл:
```go
DB_HOST=postgres
DB_PORT=5432
DB_USER=postgres
DB_PASS=postgres
DB_NAME=postgres

SRV_HOST=app
SRV_PORT=8000
SRV_TIMEOUT=10s
SRV_IDLE_TIMEOUT=30s
```

# Генерация Swagger документации

Для генерации документации Swagger использованы аннотации в коде. Сгенерировать документацию можно с помощью команды в Makefile:
```bash
make swag init
```
# Тестирование

#### Для запуска тестов выполните команду:
```bash
make run-tests
```
