FROM golang:1.25.3-alpine3.22 as build

WORKDIR /app

COPY . .

RUN go mod download


RUN go build -o app ./cmd/server/main.go

FROM alpine:3.22.2

RUN apk --no-cache add ca-certificates

COPY --from=build /app/app /usr/local/bin/app

WORKDIR /app

CMD ["app"]