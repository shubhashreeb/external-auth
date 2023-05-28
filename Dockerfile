FROM golang:alpine as builder

RUN mkdir /build
WORKDIR /build

COPY go.mod .
COPY go.sum .

# Get dependancies - will also be cached if we won't change mod/sum
RUN go mod download

# COPY the source code as the last step
COPY . .

# Build the binary
RUN CGO_ENABLED=0 go build -o envoy-auth-server main.go

# Copy into scratch
FROM scratch
COPY --from=builder /build/envoy-auth-server /bin/envoy-auth-server
CMD ["/bin/envoy-auth-server"]
