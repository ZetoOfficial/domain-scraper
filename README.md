# Domain Scraper

**Domain Scraper** — это проект на Go для парсинга и обработки внутренних ссылок веб-страниц. Система состоит из двух приложений:

1. **Parser**: парсит веб-страницу, извлекает ссылки и отправляет их в RabbitMQ.
2. **Consumer**: потребляет ссылки из RabbitMQ, обрабатывает их и извлекает данные с указанных страниц.

---

## Требования

- [Go](https://golang.org/dl/) версии 1.22 или выше.
- [RabbitMQ](https://www.rabbitmq.com/) (локально или в облаке).

---

## Установка

1. Клонируйте репозиторий:

   ```bash
   git clone https://github.com/ZetoOfficial/domain-scraper.git
   cd domain-scraper
   ```

2. Установите зависимости:

   ```bash
   go mod download
   ```

---

## Настройка

1. Создайте файл конфигурации `configs/config.env` (если его нет):

   ```bash
   touch configs/config.env
   ```

2. Заполните `config.env` следующими параметрами:

   ```env
   RABBITMQ_HOST=localhost
   RABBITMQ_PORT=5672
   RABBITMQ_USER=guest
   RABBITMQ_PASSWORD=guest
   RABBITMQ_QUEUE=links_queue
   ```

   **Примечание**: Укажите реальные значения, если RabbitMQ настроен в другой среде.

---

## Сборка и запуск

### Запуск Parser

1. Соберите бинарный файл:

   ```bash
   go build -o parser cmd/parser/main.go
   ```

2. Запустите парсер, указав URL страницы для парсинга:

   ```bash
   ./parser https://example.com
   ```

---

### Запуск Consumer

1. Соберите бинарный файл:

   ```bash
   go build -o consumer cmd/consumer/main.go
   ```

2. Запустите консьюмера с опциональным указанием таймаута (в секундах):

   ```bash
   # По умолчанию используется таймаут 30 секунд
   ./consumer

   # С таймаутом 60 секунд
   ./consumer -timeout 60
   ```
