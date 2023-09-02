FROM golang:1.20.6-alpine3.18 AS builder

WORKDIR /app

RUN apk add --no-cache gcc

COPY go.mod go.sum ./
RUN go mod download

COPY ./ ./
RUN CGO_ENABLED=1 GOOS=linux go build -o todolist

FROM alpine:3.18

ARG DB_MIGR=/var/lib/sqlite/migrations

COPY --from=builder /app/todolist /usr/bin/todolist
COPY --from=builder /app/migrations $DB_MIGR

ENV MODE=prod
ENV DB_PATH=/var/lib/sqlite/data
ENV BE_HOST=0.0.0.0
ENV BE_PORT=80
ENV BE_SECRET=secret

RUN curl -fsSL https://raw.githubusercontent.com/pressly/goose/master/install.sh | sh -s v3.15.0
RUN mkdir -p $(basename $DB_PATH)
RUN goose -dir $DB_MIGR sqlite3 $DB_PATH up

EXPOSE $BE_PORT
CMD todolist
