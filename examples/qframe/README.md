## Spin It Up

```bash
$ docker stack deploy -c complete-stack.yml qframe
Creating service qframe_grafana
Creating service qframe_elasticsearch
Creating service qframe_kibana
Creating service qframe_agent
Creating service qframe_info
Creating service qframe_warn
Creating service qframe_influxdb
```
Wait until all is green.

```bash
$ docker service ls
ID                  NAME                   MODE                REPLICAS            IMAGE                             PORTS
7gdh16vyx0c2        qframe_grafana         replicated          1/1                 qnib/plain-grafana4:latest        *:3000->3000/tcp
d0o6eabpprn9        qframe_kibana          replicated          1/1                 qnib/plain-kibana5:latest         *:5601->5601/tcp
gv59lvy3ky98        qframe_info            replicated          0/0                 qnib/plain-qframe-client:latest
jil7pwdzz47m        qframe_influxdb        replicated          1/1                 qnib/plain-influxdb:latest        *:8083->8083/tcp,*:8086->8086/tcp
u1mipv74vzwr        qframe_agent           replicated          1/1                 qnib/qframe:latest                *:11001->11001/tcp
v1lxq6vnrdr8        qframe_warn            replicated          0/0                 qnib/plain-qframe-client:latest
xmsxk5kizl82        qframe_elasticsearch   replicated          1/1                 qnib/plain-elasticsearch:latest   *:9200->9200/tcp,*:9300->9300/tcp 
```

### Add WARN/INFO logger

Now you can increase the replicas for WARN and INFO.

```
$ docker service update --replicas=1 qframe_warn
qframe_warn
$ docker service update --replicas=2 qframe_info
qframe_info
```



Now open [localhost:3000](http://localhost:3000) (admin/admin) to check the dashboards.

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

