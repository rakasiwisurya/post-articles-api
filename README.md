# Post Articles API

Article microservice, built with Go, [Fiber v3](https://gofiber.io) and MySQL.

## Tech Stack

- Go 1.25
- Fiber v3 (HTTP framework)
- MySQL (tested against XAMPP/MariaDB on localhost)
- [golang-migrate](https://github.com/golang-migrate/migrate) for schema migrations
- [go-playground/validator](https://github.com/go-playground/validator) for request validation

## Project Structure

```
server/
├── cmd/api/              # application entrypoint
├── internal/
│   ├── article/          # domain: model, repository, service, handler, routes
│   ├── config/           # environment configuration
│   └── database/         # mysql connection + migration runner
├── migrations/           # golang-migrate SQL files (embedded in the binary)
└── docs/                 # manual schema + postman collection
```

The article feature follows a layered flow: **handler → service → repository**. Handlers only translate HTTP, the service owns validation/business rules, and the repository owns SQL.

## Getting Started

1. Make sure MySQL is running (e.g. XAMPP on `localhost:3306`).
2. Copy the environment file and adjust if needed:

   ```bash
   cp .env.example .env
   ```

   | Variable      | Default     | Description    |
   | ------------- | ----------- | -------------- |
   | `APP_PORT`    | `8080`      | HTTP port      |
   | `DB_HOST`     | `127.0.0.1` | MySQL host     |
   | `DB_PORT`     | `3306`      | MySQL port     |
   | `DB_USER`     | `root`      | MySQL user     |
   | `DB_PASSWORD` | _(empty)_   | MySQL password |
   | `DB_NAME`     | `article`   | Database name  |

3. Run the API:

   ```bash
   go run ./cmd/api
   # or: make run
   ```

On startup the service creates the `article` database when missing and applies the embedded migrations (posts table). The manual DDL for the database section of the test is in [docs/manual_schema.sql](docs/manual_schema.sql).

## API Endpoints

| #   | Method               | URL                         | Description                                            |
| --- | -------------------- | --------------------------- | ------------------------------------------------------ |
| 1   | `POST`               | `/article/`                 | Create a new article                                   |
| 2   | `GET`                | `/article/<limit>/<offset>` | List articles with paging (`?status=` filter optional) |
| 3   | `GET`                | `/article/<id>`             | Get one article                                        |
| 4   | `POST`/`PUT`/`PATCH` | `/article/<id>`             | Update an article                                      |
| 5   | `DELETE`             | `/article/<id>`             | Delete an article                                      |

Request body for create/update:

```json
{
  "title": "",
  "content": "",
  "category": "",
  "status": ""
}
```

### Validation Rules

- `title`: required, at least 20 characters
- `content`: required, at least 200 characters
- `category`: required, at least 3 characters
- `status`: required, one of `publish`, `draft`, `thrash`

Validation failures respond `400` with field-keyed messages:

```json
{
  "message": "validation failed",
  "errors": {
    "title": "title must be at least 20 characters"
  }
}
```

## Postman

Import [docs/post-articles.postman_collection.json](docs/post-articles.postman_collection.json) — it contains one request per endpoint with example payloads. `base_url` defaults to `http://localhost:8080`.
