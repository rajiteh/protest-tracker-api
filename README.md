# Protest Tracker API

## Development

- Ensure you have a go 1.17 environment (Use [gvm](https://github.com/moovweb/gvm) to procure one easily.)
- Copy `.env.sample` to `.env` at the root of the repository, and set the required values. See [Env File](#env-file) section below.
- Setup the dependencies `go mod tidy`
- Run the program `go run main.go`

## Env File

| Var Name | Description |
|---|---|
|TELEGRAM_APITOKEN | Telegram API token, you can obtain one by messaging [@BotFather](https://t.me/botfather) |
|TELEGRAM_BOT_ADMINS| Comma seperated list of user IDs that will be considered admins |
|GCP_AUTH_JSON_B64| Base64 Encoded GCP service account JSON, this account is used for ingesting watchdog sheet. No permissions are required. |
|DATABASE_DSN | A DSN string to connect to the database, right now only sqlite is supported. |

## Deployment

TBD
