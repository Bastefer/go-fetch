FROM golang:1.25-alpine

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

RUN go install github.com/pressly/goose/v3/cmd/goose@latest
RUN apk add --no-cache postgresql-client

COPY . .

RUN chmod +x entrypoint.sh

CMD ["./entrypoint.sh"]