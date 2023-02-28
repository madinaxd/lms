# Build stage
FROM golang:1.19-alpine3.16 AS builder
WORKDIR /app
COPY /students_svc .
RUN go build -o main students_svc/cmd/main.go

# Run stage
FROM alpine:3.16
WORKDIR /app
COPY --from=builder /app/main .
COPY students_svc/app.env .
COPY students_svc/db/migrations ./db/migrations

EXPOSE 8081
CMD [ "/app/main" ]

