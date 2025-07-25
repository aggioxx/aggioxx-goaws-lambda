# syntax=docker/dockerfile:1
FROM golang:1.24 as builder
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY .. .
RUN CGO_ENABLED=0 GOOS=linux go build -o main main.go

FROM --platform=linux/amd64 public.ecr.aws/lambda/go:1
COPY --from=builder /app/main /var/task/main
CMD ["main"]
