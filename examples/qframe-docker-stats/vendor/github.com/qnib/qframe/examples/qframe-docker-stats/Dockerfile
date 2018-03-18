FROM qnib/uplain-golang:1.9 AS build

WORKDIR /usr/local/src/github.com/qnib/qframe/examples/qframe-docker-stats
COPY main.go .
COPY vendor/ vendor
RUN go install

FROM qnib/uplain-config

ENV ENTRYPOINTS_DIR=/opt/qnib/entry/ \
    LOG_LEVEL=info \
    DOCKER_HOST=unix:///var/run/docker.sock \
    ELASTICSEARCH_HOST=172.17.0.1 \
    INFLUXDB_HOST=172.17.0.1 \
    INFLUXDB_HOST=qframe \
    GROK_PATTERNS_DIR=/etc/qframe/patterns \
    DOCKER_LOG_SINCE=1m \
    HEALTH_BIND_PORT=8123 \
    HEALTH_BIND_HOST=0.0.0.0

COPY --from=build /usr/local/bin/qframe-docker-stats /usr/local/bin/
COPY resources/opt/qnib/qframe/template.yml /opt/qnib/qframe/template.yml
COPY resources/patterns /etc/qframe/patterns
COPY resources/opt/qframe/entry/*.sh /opt/qframe/entry/
CMD ["qframe-docker-stats", "--config", "/etc/qframe.yml"]
