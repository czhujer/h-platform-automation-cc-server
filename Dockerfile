FROM golang:1.13-alpine3.12 as builder

WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -mod vendor -o cc-server .

FROM golang:1.13-alpine3.12
COPY --from=builder /src/cc-server /opt/cc-server

# because terraform scripts
RUN apk add bash coreutils
# tools
RUN apk add vim

ENTRYPOINT ["/opt/cc-server"]
