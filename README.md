# Онлайн Библиотека Песен 🎶
Тестовое задание для Effective Mobile на позицию Junior Go Developer

Задание: https://docs.google.com/document/d/1DeSr_ogIb802V9RVjm1rJ_t0iY7vDrvf1qYT6hggf5A/edit?tab=t.0
## Описание

Этот проект реализует онлайн-библиотеку песен с возможностью:

    •	Получения информации о песнях с фильтрацией и пагинацией
    •	Получения текста песни с пагинацией по куплетам
	•	Добавления новой песни
	•	Удаления и обновления существующей песни

При добавлении новой песни выполняется запрос к внешнему API для обогащения данных (дата выпуска, текст, ссылка на песню), которые сохраняются в базе данных PostgreSQL.

## Стек технологий

	•	Backend: Go
	•	База данных: PostgreSQL
	•	Тестовая база данных может быть развернута с помощью Docker Compose
	•	Миграции базы данных: Goose
	•	Документация API: Swagger
	•	Логирование: slog

## Структура проекта

Проект реализован с использованием слоистой архитектуры:
```bash
├── Makefile
├── README.md
├── bin
├── cmd
│   └── main.go # Точка входа в приложение
├── docker-compose.yml
├── docs
├── go.mod
├── go.sum
├── internal
│   ├── controller # REST API хендлеры
│   ├── models # Структуры данных
│   ├── repository # Работа с данными
│   ├── service # Бизнес-логика
│   └── utils # Вспомогательные функции
└── migrations # Миграции базы данных

```

## Возможности API

	•	GET /songs: Получение всех песен с возможностью фильтрации и пагинации
	•	POST /songs: Добавление новой песни
	•	PUT /songs/{id}: Обновление данных песни
	•	DELETE /songs/{id}: Удаление песни
	•	GET /songs/{id}/verses: Получение текста песни с пагинацией по куплетам

Полную документацию API можно найти по адресу: http://localhost:8080/swagger/index.html после запуска приложения.

## Установка и запуск

Предварительные требования

	•	Go версии 1.20 или выше
	•	Docker Compose # Для запуска PostgreSQL
	•	Make

Шаги для развертывания

1.	Клонируйте репозиторий:

```bash
git clone https://github.com/Qquiqlerr/test-task-effective-mobile.git
cd test-task-effective-mobile
```

2.	Создайте файл .env с необходимыми переменными окружения(в проекте уже есть файл .env, измените его при необходимости):

```.env
# Database configuration
DB_NAME = test
DB_PASSWORD = 1234
DB_USER = postgres
DB_HOST = localhost
DB_PORT = 5432

# Logger configuration
LOG_LEVEL = debug

# Server configuration
PORT = 8080
```
3. Запустите PostgreSQL с помощью Docker Compose(при желании можно поднять базу данных вручную):
```bash
docker-compose up -d
```

3.	Соберите и запустите проект:
```bash
make all
```



4.	Приложение будет доступно по адресу http://localhost:8080.



Логирование

В проекте используется библиотека slog для логирования. Поддерживается два уровня логов: debug и prod.

## Все команды Makefile

```bash
make all # Сборка и запуск приложения
make init # Инициализация проекта
make build # Сборка приложения
make run # Запуск приложения
make migrate-up # Применение миграций
make migrate-down # Откат миграций
make doc-gen # Генерация документации API
make clean # Очистка бинарников и документации
```

## Контакты
Автор: Алексей Метлушко
- [Telegram](https://t.me/sslowerr)
- [Email](mailto:leha.metlushko@bk.ru)

По любым вопросам связанным с проектом, пожалуйста, свяжитесь со мной.
