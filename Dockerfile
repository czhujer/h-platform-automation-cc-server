FROM golang:1.16-alpine3.12 as builder

WORKDIR /src
COPY . .
RUN CGO_ENABLED=0 GOOS=linux go build -ldflags '-extldflags "-static"' -mod vendor -o cc-server .

FROM ubuntu:noble
COPY --from=builder /src/cc-server /opt/cc-server

RUN apt-get update && apt-get install -y \
  libvirt-clients \
  && rm -rf /var/lib/apt/lists/*

ENTRYPOINT ["/opt/cc-server"]
