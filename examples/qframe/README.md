## Development

```bash
$ docker run -ti --name qframe --rm -e SKIP_ENTRYPOINTS=1 \
            -v ${GOPATH}/src/github.com/qnib/qframe/examples/qframe:/usr/local/src/github.com/qnib/qframe/examples/qframe \
            -v ${GOPATH}/src/github.com/qnib/qframe-collector-tcp:/usr/local/src/github.com/qnib/qframe-collector-tcp \
            -v ${GOPATH}/src/github.com/qnib/qframe-collector-internal:/usr/local/src/github.com/qnib/qframe-collector-internal \
            -v ${GOPATH}/src/github.com/qnib/qframe-collector-docker-events:/usr/local/src/github.com/qnib/qframe-collector-docker-events \
            -v ${GOPATH}/src/github.com/qnib/qframe-collector-docker-log:/usr/local/src/github.com/qnib/qframe-collector-docker-log \
            -v ${GOPATH}/src/github.com/qnib/qframe-collector-docker-stats:/usr/local/src/github.com/qnib/qframe-collector-docker-stats \
            -v ${GOPATH}/src/github.com/qnib/qframe-filter-docker-stats/lib:/usr/local/src/github.com/qnib/qframe-filter-docker-stats/lib \
            -v ${GOPATH}/src/github.com/qnib/qframe-filter-grok/lib:/usr/local/src/github.com/qnib/qframe-filter-grok/lib \
            -v ${GOPATH}/src/github.com/qnib/qframe-filter-inventory/lib:/usr/local/src/github.com/qnib/qframe-filter-inventory/lib \
            -v ${GOPATH}/src/github.com/qnib/qframe-filter-metrics/lib:/usr/local/src/github.com/qnib/qframe-filter-metrics/lib \
            -v ${GOPATH}/src/github.com/qnib/qframe-filter-statsd/lib:/usr/local/src/github.com/qnib/qframe-filter-statsd/lib \
            -v ${GOPATH}/src/github.com/qnib/qframe-inventory/lib:/usr/local/src/github.com/qnib/qframe-inventory/lib \
            -v ${GOPATH}/src/github.com/qnib/qframe-handler-influxdb/lib:/usr/local/src/github.com/qnib/qframe-handler-influxdb/lib \
            -v ${GOPATH}/src/github.com/qnib/qframe-handler-elasticsearch/lib:/usr/local/src/github.com/qnib/qframe-handler-elasticsearch/lib \
            -v ${GOPATH}/src/github.com/qnib/qframe-types:/usr/local/src/github.com/qnib/qframe-types \
            -v ${GOPATH}/src/github.com/qnib/qframe-utils:/usr/local/src/github.com/qnib/qframe-utils \
            -v ${GOPATH}/src/github.com/qnib/statsq/lib:/usr/local/src/github.com/qnib/statsq/lib \
            -v /var/run/docker.sock:/var/run/docker.sock \
            -v $(pwd)/resources/patterns/:/etc/gcollect/patterns/ \
            -w /usr/local/src/github.com/qnib/qframe/examples/qframe \
            qnib/uplain-golang bash
$ govendor update github.com/qnib/qframe-collector-docker-events/lib \
                  github.com/qnib/qframe-collector-docker-log/lib \
                  github.com/qnib/qframe-collector-docker-stats/lib \
                  github.com/qnib/qframe-collector-tcp/lib \
                  github.com/qnib/qframe-collector-internal/lib \
                  github.com/qnib/qframe-filter-docker-stats/lib \
                  github.com/qnib/qframe-filter-grok/lib \
                  github.com/qnib/qframe-filter-metrics/lib \
                  github.com/qnib/qframe-filter-inventory/lib \
                  github.com/qnib/qframe-filter-statsq/lib \
                  github.com/qnib/qframe-inventory/lib \
                  github.com/qnib/qframe-handler-influxdb/lib \
                  github.com/qnib/qframe-handler-elasticsearch/lib \
                  github.com/qnib/statsq/lib \
                  github.com/qnib/qframe-types \
                  github.com/qnib/qframe-utils
$ govendor fetch -v +m
```
