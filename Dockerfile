FROM golang:1.19.6-alpine AS builder
WORKDIR /app
COPY . .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o main main.go


FROM alpine:3.18
WORKDIR /app
COPY --from=builder /app/main .
RUN apk add --no-cache tzdata
ENV TZ=Asia/Jakarta
COPY config.yaml .
CMD ["./main"]
