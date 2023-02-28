# Build stage
FROM golang:1.19-alpine3.16 AS builder
WORKDIR /app
COPY /courses_svc .
RUN go build -o main courses_svc/cmd/main.go

# Run stage
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .
COPY courses_svc/app.env .
COPY courses_svc/db/migrations ./db/migrations

EXPOSE 7071
CMD [ "/app/main" ]