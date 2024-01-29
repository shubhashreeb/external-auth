FROM golang:alpine as builder

WORKDIR /build

COPY . ./

# Build the binary
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o authsvc 

FROM alpine:latest

WORKDIR /app

COPY --from=builder /build/authsvc /app/authsvc

EXPOSE 8090
EXPOSE 9090
EXPOSE 8080

CMD ["/app/authsvc"]
