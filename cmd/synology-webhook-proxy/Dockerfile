FROM golang:1.25.6 as builder
WORKDIR /app
COPY ../../go.mod ./
COPY ../../go.sum ./
COPY . .
RUN go mod download \
    && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags="-s -w" -o=main .

FROM scratch
WORKDIR /opt/app
COPY --from=builder /app/main ./main
CMD ["/opt/app/main"]
