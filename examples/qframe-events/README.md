

## Development

First we start the development services.

```bash
$ docker stack deploy -c docker-compose.yml qframe                                                                                          git:(master|●2✚2…
Creating service qframe_kibana
Creating service qframe_influxdb
Creating service qframe_grafana
Creating service qframe_elasticsearch
```

Afterwards a golang container with all libraries we want to work on is fired up.

```bash
$ docker run -ti --name qframe-collector-tcp --rm -e SKIP_ENTRYPOINTS=1 -p 11001:11001 \
           -v ${GOPATH}/src/github.com/qnib/qframe/:/usr/local/src/github.com/qnib/qframe/ \
           -v ${GOPATH}/src/github.com/qnib/qframe-collector-docker-events:/usr/local/src/github.com/qnib/qframe-collector-docker-events \
           -v ${GOPATH}/src/github.com/qnib/qframe-collector-internal/lib:/usr/local/src/github.com/qnib/qframe-collector-internal/lib \
           -v ${GOPATH}/src/github.com/qnib/qframe-collector-tcp/lib:/usr/local/src/github.com/qnib/qframe-collector-tcp/lib \
           -v ${GOPATH}/src/github.com/qnib/qframe-filter-grok/lib:/usr/local/src/github.com/qnib/qframe-filter-grok/lib \
           -v ${GOPATH}/src/github.com/qnib/qframe-filter-inventory/lib:/usr/local/src/github.com/qnib/qframe-filter-inventory/lib \
           -v ${GOPATH}/src/github.com/qnib/qframe-handler-elasticsearch/lib:/usr/local/src/github.com/qnib/qframe-handler-elasticsearch/lib \
           -v ${GOPATH}/src/github.com/qnib/qframe-handler-influxdb/lib:/usr/local/src/github.com/qnib/qframe-handler-influxdb/lib \
           -v ${GOPATH}/src/github.com/qnib/qframe-inventory/lib:/usr/local/src/github.com/qnib/qframe-inventory/lib \
           -v ${GOPATH}/src/github.com/qnib/qframe-types:/usr/local/src/github.com/qnib/qframe-types \
           -v ${GOPATH}/src/github.com/qnib/qframe-utils:/usr/local/src/github.com/qnib/qframe-utils \
           -v /var/run/docker.sock:/var/run/docker.sock \
           -v $(pwd)/patterns:/etc/qframe/patterns \
           -w /usr/local/src/github.com/qnib/qframe/examples/qframe-events \
            qnib/uplain-golang bash
```

```
$ govendor update github.com/qnib/qframe-collector-docker-events/lib \
	            github.com/qnib/qframe-collector-internal/lib \
	            github.com/qnib/qframe-collector-tcp/lib \
	            github.com/qnib/qframe-filter-grok/lib \
	            github.com/qnib/qframe-filter-inventory/lib \
	            github.com/qnib/qframe-handler-elasticsearch/lib \
	            github.com/qnib/qframe-handler-influxdb/lib \
	            github.com/qnib/qframe-inventory/lib \
                github.com/qnib/qframe-types \
                github.com/qnib/qframe-utils
$ govendor fetch +m
```

Start the daemon...

```bash
$ root@0915d2c02d6b:/usr/local/src/github.com/qnib/qframe/examples/qframe-events# go run main.go --config qframe.yml
  2017/05/09 18:50:11 [II] Start Version: 0.0.0.0
  2017/05/09 18:50:11 [II] Use config file: qframe.yml
  2017/05/09 18:50:11 [II] Dispatch broadcast for Back, Data and Tick
  2017/05/09 18:50:11.642470 [  INFO] influxdb >> Start log handler influxdbv0.1.1
  2017/05/09 18:50:11.647722 [  INFO] es_logstash >> Start elasticsearch handler: es_logstashv0.1.6
  2017/05/09 18:50:11.652478 [  INFO] inventory >> Start inventory v0.1.1
  2017/05/09 18:50:11.659806 [  INFO] app-event >> Add patterns from directory '/usr/local/src/github.com/qnib/qframe/examples/qframe-events/patterns'
  2017/05/09 18:50:11.703690 [  INFO] app-event >> Start grok filter v0.1.7
  2017/05/09 18:50:11.709047 [  INFO] influxdb >> Established connection to 'http://172.17.0.1:8086
  2017/05/09 18:50:11.712012 [  INFO] docker-events >> Start docker-events collector v0.2.1
  2017/05/09 18:50:11.717885 [  INFO] internal >> Start internal collector v0.1.0
  2017/05/09 18:50:11.740939 [  INFO] tcp >> Listening on 0.0.0.0:11001
  2017/05/09 18:50:12.442592 [  INFO] docker-events >> Connected to 'moby' / v'17.05.0-ce-rc1'
  2017/05/09 18:56:21.523606 [  INFO] tcp >> Received TCP message '@cee{"msg":"My unixepoch is 1494356181", "event_code":"001.001"}' from '172.17.0.4'
  2017/05/09 18:56:21.526645 [  INFO] tcp >> Got msg from buffer: @cee{"msg":"My unixepoch is 1494356181", "event_code":"001.001"}
  2017/05/09 18:56:21.528872 [  INFO] inventory >> Received InventoryRequest for {2017-05-09 18:56:21.526663903 +0000 UTC   172.17.0.4 0xc4201b7bc0}
  2017/05/09 18:56:51.156649 [ ERROR] docker-events >> Container 00db37030e29744bbed18b7a66ee0e4bf02f7ba7a9bb9247b2be55ee8e9cc14b just 'destroy' without having an entry in the Inventory
  2017/05/09 18:57:37.611993 [ ERROR] docker-events >> Container 60a30e0d98ee9e3448e58b6e4075a26e7425198f29938d4eba02942dedd88036 just 'destroy' without having an entry in the Inventory
  2017/05/09 18:57:45.200780 [ ERROR] docker-events >> Container f7a973db7d4f9d66979dc3b7e2e8dd25f113fe440afdb13f4a52b9d06b002f3c just 'destroy' without having an entry in the Inventory
  2017/05/09 18:57:47.056849 [  INFO] tcp >> Received TCP message '@cee{"msg":"My unixepoch is 1494356266", "event_code":"001.001"}' from '172.17.0.3'
  2017/05/09 18:57:47.061090 [  INFO] tcp >> Got msg from buffer: @cee{"msg":"My unixepoch is 1494356266", "event_code":"001.001"}
  2017/05/09 18:57:47.064359 [  INFO] inventory >> Received InventoryRequest for {2017-05-09 18:57:47.061114811 +0000 UTC   172.17.0.3 0xc4203d6360}
```

Send a message...

```bash
$ docker run --rm --name events-$(date +%s) -ti qnib/plain-qframe-client:latest \
             send-tcp-event.sh 001.001 "My unixepoch is $(date +%s)"
[II] qnib/init-plain script v0.4.20
> execute CMD 'send-tcp-event.sh 001.001 My unixepoch is 1494356266'
Send: @cee{"msg":"My unixepoch is 1494356266", "event_code":"001.001"}
$
```

The daemon reacts..

```bash
2017/05/09 18:57:47.056849 [  INFO] tcp >> Received TCP message '@cee{"msg":"My unixepoch is 1494356266", "event_code":"001.001"}' from '172.17.0.3'
2017/05/09 18:57:47.061090 [  INFO] tcp >> Got msg from buffer: @cee{"msg":"My unixepoch is 1494356266", "event_code":"001.001"}
2017/05/09 18:57:47.064359 [  INFO] inventory >> Received InventoryRequest for {2017-05-09 18:57:47.061114811 +0000 UTC   172.17.0.3 0xc4203d6360}
```
