FROM golang:1.25.3-alpine3.22 as build

WORKDIR /app

COPY . .

RUN go mod download


RUN go build -o migrator ./cmd/migrator/main.go

FROM alpine:3.22.2

RUN apk --no-cache add ca-certificates

COPY --from=build /app/migrator /usr/local/bin/migrator

WORKDIR /app

CMD ["migrator"]