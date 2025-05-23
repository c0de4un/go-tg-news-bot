# go-tg-news-bot

Telegram bot for news with clients and users.

## Minimal Software requirements
* Go 1.2
* PostgreSQL

## Minimal Hardware Requirements
* 512 MB RAM
* Any CPU x86
* 8 GB file system
* Ethernet

## Build

```shell
go build ./cmd/bot/bot.go
```

## Dev

To start in debug mode, run:
```shell
go run ./cmd/bot/bot.go
```

### Contribution

### Database:
We are using Goose as DB migrations utility `https://github.com/pressly/goose`

To create new migration, use:
```sh
goose -dir db/migrations create new_user_table sql
```

To apply new migration, use:
```sh
goose -dir db/migrations postgres "postgresql://postgres:password@127.0.0.1:5432/newsbot?sslmode=disable" up
```

To rollback, use command:
```sh
goose -dir db/migrations postgres "postgresql://postgres:password@127.0.0.1:5432/newsbot?sslmode=disable" down
```