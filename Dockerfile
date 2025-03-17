FROM golang:1.24-alpine as builder

WORKDIR /app

COPY . .

RUN go mod tidy
RUN go build -o /app/main ./main.go

FROM alpine:latest

WORKDIR /app

COPY --from=builder /app/main /app/main
COPY --from=builder /app/docs /app/docs
COPY --from=builder /app/templates /app/templates

EXPOSE 8080

CMD ["/app/main"]