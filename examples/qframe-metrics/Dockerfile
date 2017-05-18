FROM qnib/uplain-golang

WORKDIR /usr/local/src/github.com/qnib/qframe/examples/qframe-metrics
COPY main.go .
COPY vendor/vendor.json vendor/vendor.json
RUN govendor fetch -v +m
RUN govendor install

FROM qnib/uplain-init

ENV ENTRYPOINTS_DIR=/opt/qnib/entry/ \
    LOG_LEVEL=info \
    DOCKER_HOST=unix:///var/run/docker.sock \
    ELASTICSEARCH_HOST=172.17.0.1 \
    INFLUXDB_HOST=172.17.0.1 \
    GROK_PATTERNS_DIR=/etc/qframe/patterns

COPY --from=0 /usr/local/bin/qframe-metrics /usr/local/bin/
COPY resources/opt/qnib/qframe/template.yml /opt/qnib/qframe/template.yml
COPY resources/patterns /etc/qframe/patterns
COPY resources/opt/qnib/entry/*.sh /opt/qnib/entry/
CMD ["qframe-metrics", "--config", "/etc/qframe.yml"]
