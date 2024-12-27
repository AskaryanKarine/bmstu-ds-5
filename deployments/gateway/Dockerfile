FROM golang:1.22.5 AS builder

LABEL org.opencontainers.image.source=https://github.com/askaryankarine/bmstu-ds-4
LABEL org.opencontainers.image.description="DS labs image"
LABEL org.opencontainers.image.licenses=MIT

WORKDIR /app

COPY ["go.mod", "go.sum", "./"]
RUN go mod download

COPY ../.. ./

RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o /go/bin/main ./cmd/gateway/main.go

FROM alpine:latest AS app

COPY --from=builder /go/bin/main ./go/

EXPOSE 80
CMD ./go/main
