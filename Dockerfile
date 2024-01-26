# FROM golang:alpine as builder

# RUN mkdir /build
# WORKDIR /build

# COPY . .

# RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o authsvc

# # Copy into scratch
# FROM scratch
# COPY --from=builder /build/authsvc /bin/authsvc
# CMD ["/bin/authsvc"]

FROM golang:alpine as builder

WORKDIR /build

COPY . ./

# Build the binary
RUN GOOS=linux GOARCH=amd64 CGO_ENABLED=0 go build -o authsvc 

FROM alpine

WORKDIR /app

COPY --from=builder /build/authsvc /app/authsvc

EXPOSE 9002

CMD ["/app/authsvc"]
