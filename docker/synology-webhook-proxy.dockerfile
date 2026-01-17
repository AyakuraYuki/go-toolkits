FROM golang:1.25.6 AS builder
WORKDIR /app
COPY . .
WORKDIR /app/cmd/synology-webhook-proxy
RUN go mod download
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o=/opt/main .

FROM scratch
WORKDIR /opt/app
COPY --from=builder /opt/main ./main
CMD ["/opt/app/main"]
