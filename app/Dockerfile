FROM golang:1.22.5-alpine3.20 AS builder

WORKDIR /app

COPY go.mod go.sum ./
RUN go mod download

COPY . .

RUN go build -o app cmd/main.go

FROM gcr.io/distroless/base-debian10

COPY --from=builder /app/app /app

EXPOSE 8080

CMD ["/app"]