# qframe-logs
Collection framework for a containerised world to subscribe to container logs

## Spin Up Backend

```bash
$ docker stack deploy -c docker-compose.yml qframe
Creating service qframe_elasticsearch
Creating service qframe_kibana
Creating service qframe_qframe-logs
Creating service qframe_influxdb
Creating service qframe_grafana
$ open http://localhost:5601
```


```bash
$ go run main.go
go run main.go --config qframe.yml
2017/05/12 07:47:45 [II] Start Version: 0.0.0
2017/05/12 07:47:45 [II] Use config file: qframe.yml
2017/05/12 07:47:45 [II] Dispatch broadcast for Back, Data and Tick
2017/05/12 07:47:45.621583 [  INFO]   elasticsearch Name:es_logstash >> KV key 'docker-log.log_level' will replace 'Level' when indexing
2017/05/12 07:47:45.635552 [  INFO]   elasticsearch Name:es_logstash >> KV key 'docker-log.log_msg' will replace 'msg' when indexing
2017/05/12 07:47:45.647797 [NOTICE]   elasticsearch Name:es_logstash >> Start elasticsearch handler: es_logstashv0.1.8
2017/05/12 07:47:45.667924 [  INFO]   elasticsearch Name:es_logstash >> Connecting to 172.17.0.1:9200
2017/05/12 07:47:45.692449 [  INFO]            grok Name:app-log    >> Add patterns from directory '/etc/gcollect/patterns'
2017/05/12 07:47:45.705545 [NOTICE]            grok Name:app-log    >> Start grok filter v0.1.9
2017/05/12 07:47:45.736310 [NOTICE]      docker-log Name:docker-log >> Start v0.1.1
2017/05/12 07:47:46.769149 [  INFO]      docker-log Name:docker-log >> Connected to 'moby' / v'17.05.0-ce' (SWARM: active)
2017/05/12 07:47:46.776641 [  INFO]      docker-log Name:docker-log >> Start listeners for already running containers: 5```
```
##### docker-log collector

```yaml
collector:
  docker-log:
    docker-host: "unix:///var/run/docker.sock"
    inputs: "docker-events"
    skip-container-label: "org.qnib.skip-logs"
```

The `docker-log` collector subscribes to stdout/stderr of each container that is running on-top of the local engine.

```go
	logOpts := types.ContainerLogsOptions{ShowStdout: true, ShowStderr: true, Follow: true, Tail: "all"}
	reader, err := cs.cli.ContainerLogs(ctx, cs.CntID, logOpts)
```

The `docker-log` collector will start a listener for already running container - if they do not have a label set to prevent this from happening:

```bash
2017/05/12 07:47:46.786465 [  INFO]      docker-log Name:docker-log >> Skip subscribing to logs of '[/qframe-logs]' as label 'org.qnib.skip-logs' is set
2017/05/12 07:47:46.790583 [  INFO]      docker-log Name:docker-log >> Skip subscribing to logs of '[/qframe_grafana.1.i6mhekwto6aysvod3feaggedt]' as label 'org.qnib.skip-logs' is set
2017/05/12 07:47:46.795457 [  INFO]      docker-log Name:docker-log >> Skip subscribing to logs of '[/qframe_influxdb.1.47gdjb5nuz0me9tj6mnwl3mw1]' as label 'org.qnib.skip-logs' is set
2017/05/12 07:47:46.799130 [  INFO]      docker-log Name:docker-log >> Skip subscribing to logs of '[/qframe_kibana.1.06i4xd1w9yryymw3lqysx3mcw]' as label 'org.qnib.skip-logs' is set
2017/05/12 07:47:46.803635 [  INFO]      docker-log Name:docker-log >> Skip subscribing to logs of '[/qframe_elasticsearch.1.r6bhahze3jw869tpwfle339fi]' as label 'org.qnib.skip-logs' is set
```

###### Example App-Log

A container is started to send an exemplary log line.

```bash
$ docker run --rm -ti --name log-$(date +%s) qnib/plain-qframe-client loop-log.sh 1 WARN
[WARN] Log message No1 at unixepoch 1494574401
```
This freshly started container will trigger a go-routine, which will handle the container output.

```bash
2017/05/12 07:48:02.813383 [  INFO]      docker-log Name:docker-log >> Received: log-1494575282: container.start
2017/05/12 07:48:02.816692 [  INFO]      docker-log Name:docker-log >> Container started: /log-1494575282 | ID:15e403290939ef765dd1cb9838037c7925cc5860cb97636ffed2d5fd08443a2a
2017/05/12 07:48:02.821061 [  INFO]      docker-log Name:docker-log >> Start listener for: 'log-1494575282' [15e403290939ef765dd1cb9838037c7925cc5860cb97636ffed2d5fd08443a2a]
2017/05/12 07:48:02.821976 [  INFO]      docker-log Name:docker-log >> Received: log-1494575282: container.resize
```
Output is processed and forwarded to the internal go-channels.

```bash
2017/05/12 07:48:02.825304 [ DEBUG]      docker-log Name:docker-log >> Container 'docker-log': [WARN] Log message No1 at unixepoch 1494575282
```

##### GROK filter

```yaml
filter:
  app-log:
    pattern-dir: "/etc/qframe/patterns"
    inputs: "docker-log"
    pattern: "%{LOG}"
    overwrite-message-key: "msg"
```
The grok filter will try to match `%{LOG}`.

```ruby
LOG ^\[%{LOG_LEVEL:log_level}\]%{GREEDYDATA:log_msg}
LOG_LEVEL (DEBUG|INFO|NOTICE|WARN|ERROR|CRITICAL|PANIC)
```

Both try, but only `app-log` matches the log-line.

```bash
2017/05/12 07:48:02.829482 [ DEBUG]            grok Name:app-log    >> Matched pattern '%{LOG}'
2017/05/12 07:48:02.832106 [ DEBUG]            grok Name:app-log    >>             log_msg:  Log message No1 at unixepoch 1494575282
2017/05/12 07:48:02.834983 [ DEBUG]            grok Name:app-log    >>                 LOG: [WARN] Log message No1 at unixepoch 1494575282
2017/05/12 07:48:02.837479 [ DEBUG]            grok Name:app-log    >>           log_level: WARN
```

##### Elasticsearch

The elasticsearch handler subscribes to the grok filter, but only acts if the previous step was successful.

```yaml
handler:
  es_logstash:
    host: "172.17.0.1"
    inputs: "app-log"
    source-success: "true"
    kv-to-field: "docker-log.log_level:Level,docker-log.log_msg:msg"
    kv-skip: "docker-log.LOG,es-log.ES_LOG_LINE"
```

```bash
2017/05/12 07:48:02.839926 [  INFO]   elasticsearch Name:es_logstash >> [WARN] Log message No1 at unixepoch 1494575282
2017/05/12 07:48:02.848131 [ DEBUG]   elasticsearch Name:es_logstash >> Overwrite field 'msg' with  Log message No1 at unixepoch 1494575282
2017/05/12 07:48:02.850819 [ DEBUG]   elasticsearch Name:es_logstash >> Skip key docker-log.LOG in qm.KV
2017/05/12 07:48:02.856046 [ DEBUG]   elasticsearch Name:es_logstash >> Overwrite field 'Level' with WARN
2017/05/12 07:48:02.859248 [  INFO]   elasticsearch Name:es_logstash >>                      Timestamp: 2017-05-12T07:48:02.822893+00:00
2017/05/12 07:48:02.862259 [  INFO]   elasticsearch Name:es_logstash >>                            msg:  Log message No1 at unixepoch 1494575282
2017/05/12 07:48:02.865901 [  INFO]   elasticsearch Name:es_logstash >>                    source_path: docker-log,app-log
2017/05/12 07:48:02.869044 [  INFO]   elasticsearch Name:es_logstash >>                 container_name: log-1494575282
2017/05/12 07:48:02.876198 [  INFO]   elasticsearch Name:es_logstash >>                  container_cmd: loop-log.sh 1 WARN
2017/05/12 07:48:02.880393 [  INFO]   elasticsearch Name:es_logstash >>                          image: sha256:5f45663c9e8bec708cff7d4009bbcabbdd2decac10393c04b9fc5890f4e02c31
2017/05/12 07:48:02.883125 [  INFO]   elasticsearch Name:es_logstash >>                    msg_version: 0.5.0
2017/05/12 07:48:02.885667 [  INFO]   elasticsearch Name:es_logstash >>                          Level: WARN
2017/05/12 07:48:02.891021 [  INFO]   elasticsearch Name:es_logstash >>                   container_id: 15e403290939ef765dd1cb9838037c7925cc5860cb97636ffed2d5fd08443a2a
2017/05/12 07:48:02.893713 [  INFO]   elasticsearch Name:es_logstash >>                     image_name: qnib/plain-qframe-client
```

This event will end up in Elasticsearch to be visualised with Kibana.

![](resources/pics/kibana_warn.png)



## Development

```bash
$ docker run -ti --name qframe-logs --rm -e SKIP_ENTRYPOINTS=1 --label org.qnib.skip-logs=true \
            -v ${GOPATH}/src/github.com/qnib/qframe/examples/qframe-logs:/usr/local/src/github.com/qnib/qframe/examples/qframe-logs \
            -v ${GOPATH}/src/github.com/qnib/qframe-collector-docker-events:/usr/local/src/github.com/qnib/qframe-collector-docker-events \
            -v ${GOPATH}/src/github.com/qnib/qframe-collector-docker-log:/usr/local/src/github.com/qnib/qframe-collector-docker-log \
            -v ${GOPATH}/src/github.com/qnib/qframe-collector-internal:/usr/local/src/github.com/qnib/qframe-collector-internal \
            -v ${GOPATH}/src/github.com/qnib/qframe-filter-grok/lib:/usr/local/src/github.com/qnib/qframe-filter-grok/lib \
            -v ${GOPATH}/src/github.com/qnib/qframe-filter-inventory/lib:/usr/local/src/github.com/qnib/qframe-filter-inventory/lib \
            -v ${GOPATH}/src/github.com/qnib/qframe-inventory/lib:/usr/local/src/github.com/qnib/qframe-inventory/lib \
            -v ${GOPATH}/src/github.com/qnib/qframe-handler-elasticsearch/lib:/usr/local/src/github.com/qnib/qframe-handler-elasticsearch/lib \
            -v ${GOPATH}/src/github.com/qnib/qframe-handler-influxdb/lib:/usr/local/src/github.com/qnib/qframe-handler-influxdb/lib \
            -v ${GOPATH}/src/github.com/qnib/qframe-types:/usr/local/src/github.com/qnib/qframe-types \
            -v ${GOPATH}/src/github.com/qnib/qframe-utils:/usr/local/src/github.com/qnib/qframe-utils \
            -v /var/run/docker.sock:/var/run/docker.sock \
            -v $(pwd)/resources//patterns/:/etc/qframe/patterns/ \
            -w /usr/local/src/github.com/qnib/qframe/examples/qframe-logs \
            qnib/uplain-golang bash
$ govendor update github.com/qnib/qframe-collector-docker-events/lib \
                  github.com/qnib/qframe-collector-docker-log/lib \
                  github.com/qnib/qframe-collector-internal/lib \
                  github.com/qnib/qframe-filter-grok/lib \
                  github.com/qnib/qframe-filter-inventory/lib \
                  github.com/qnib/qframe-inventory/lib \
                  github.com/qnib/qframe-handler-elasticsearch/lib \
                  github.com/qnib/qframe-handler-influxdb/lib \
                  github.com/qnib/qframe-types \
                  github.com/qnib/qframe-utils
$ govendor fetch +m
```

