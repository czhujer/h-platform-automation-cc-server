FROM golang:1.13-alpine3.10 as builder

WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -mod vendor -o cc-server .

FROM golang:1.13-alpine3.10
COPY --from=builder /src/cc-server /opt/cc-server

ENTRYPOINT ["/opt/cc-server"]
