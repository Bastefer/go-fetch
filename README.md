# go-fetch
для запуска нужно создать `.env` и можно передать свои параметры для postgres

```bash
cp .env.local .env
```

после создания `.env` можно запустить приложение 

```bash
docker compose up --build
```

если нужно чтобы при перезапуске бд полностью была чистой то нужно раскоментировать в entrypoint.sh

```bash
# goose -dir migrations postgres \
# "host=$POSTGRES_HOST port=$POSTGRES_PORT user=$POSTGRES_USER password=$POSTGRES_PASSWORD dbname=$POSTGRES_DB sslmode=disable" down-to 0
```