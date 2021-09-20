FROM golang:1.17.1-alpine
WORKDIR /root/user
COPY ./common ../common
COPY ./user .
RUN CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -o server internal/cmd/server/*.go

FROM alpine:3.13
WORKDIR /app/
COPY --from=0 /root/user/server .
CMD ["/app/server"]