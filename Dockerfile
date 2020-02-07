FROM golang:1.13 as builder

WORKDIR /github.com/videocoin/go-service-manager
COPY . .
RUN make build

FROM bitnami/minideb:jessie

COPY --from=builder /github.com/videocoin/go-service-manager/build/bin/svcd /usr/local/bin/
COPY --from=builder /github.com/videocoin/go-service-manager/scripts/migrations /migrations

RUN install_packages curl
RUN GRPC_HEALTH_PROBE_VERSION=v0.3.0 && curl -L -k https://github.com/grpc-ecosystem/grpc-health-probe/releases/download/${GRPC_HEALTH_PROBE_VERSION}/grpc_health_probe-linux-amd64 --output /usr/local/bin/grpc_health_probe 
RUN curl -L -k https://github.com/golang-migrate/migrate/releases/download/v4.8.0/migrate.linux-amd64.tar.gz -o migrate.linux-amd64.tar.gz && tar -xzvf migrate.linux-amd64.tar.gz && mv migrate.linux-amd64 /usr/local/bin/migrate && rm migrate.linux-amd64.tar.gz

CMD ["svcd"]