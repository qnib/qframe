# qframe
GO Framework to create an input/filter/ouput pipeline.

## Description

The goal is to provide a framework that allows to model generic ETLs inspired by Logstash.

![](/resources/pics/architecture.png)


## Channels
The framework provides a tick-channel, ticking along every once so often (5s by default).

The `Data` channel moves messages from collectors to handlers and allows any number of filters in between.


## Plugins

Each plugin is its own golang project. Thus, it is easily written and interchangeable.
Furthermore the plan is to allow the use of GOLANG plugins, so that each plugin can be build as shared object and dynamically loaded, without the need to compile it into the resulting daemon.

* **collector**: Input plugin producing messages
* **filter**: plugin to refine/alter messages from collectors or other filters
* **handler**: output plugin to send/output the data

### Plugins List

The following plugins are available.

#### Collectors

- [docker-events](https://github.com/qnib/qframe-collector-docker-events) Hooks into moby's `/events` API endpoint and parses incoming events like `contianer.create` or `network.attach. 
 For now SWARM events are not provided, but there is already a PR against moby (former called docker) on github.
- [docker-stats](https://github.com/qnib/qframe-collector-docker-stats) For each incoming `docker-event` about a started container, 
 this collector will spawn a goroutine to stream the /container/<id>/stats` API call. Thus, the collector gets (as close as possible) real-time metrics for a container.
- [GELF](https://github.com/qnib/qframe-collector-gelf) Collector for the GELF log-driver of the docker-engine. Should be replaced by a `docker-logs` collector, which spawns a listener for 
 each container like the `docker-stas` collector does. Supposed to be much nicer, because the logs can still be viewed via `docker logs <container>`.
- [tcp](https://github.com/qnib/qframe-collector-tcp) Opens a TCP port which should be used by a container to send messages like AppMetrics.
 By using the `inventory` filter the metadata will be added according to the remote-IP used by the container.
- [file](https://github.com/qnib/qframe-collector-file) Simple collector to tail a file.

#### Filters

- [id](https://github.com/qnib/qframe-filter-id) Relays the message - might be droped as it was used for reversing events.
- [inventory](https://github.com/qnib/qframe-filter-inventory) Listens to `docker-events` and keeps an inventory of all containers. 
 Can be queried by other plugins sending `ContainerRequests down the `Data` channel.
- [grok](https://github.com/qnib/qframe-filter-grok) Allows for matching `QMsg` with GROK patterns (typed RegEx, much nicer to use then RegExp).
- [docker-stats](https://github.com/qnib/qframe-filter-docker-stats) Potential filter to aggregate or transform metrics comming from the `docker-stats` collector.

#### Handlers

- [log](https://github.com/qnib/qframe-handler-log) Outputs to stdout of the daemon.
- [influxdb](https://github.com/qnib/qframe-handler-influxdb) Forwards metrics to an InfluxDB server. 
- [elasticsearch](https://github.com/qnib/qframe-handler-elasticsearch) FOrwards `QMsg` to Elasticsearch.





