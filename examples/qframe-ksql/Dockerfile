FROM qnib/uplain-golang:1.9 AS build

ARG LIBRDKAFKA_VER=0.11.0.x
RUN apt-get -qq update \
 && apt-get -qq install -y librdkafka-dev
 #&& cd /usr/local/src/ \
 #&& wget -qO -  https://github.com/edenhill/librdkafka/archive/${LIBRDKAFKA_VER}.zip| bsdtar xzf - -C . \
 #&& cd librdkafka-${LIBRDKAFKA_VER} \
 #&& chmod +x configure lds-gen.py \
 #&& ./configure \
 #&& make \
 #&& make install

WORKDIR /usr/local/src/github.com/qnib/qframe/examples/qframe-ksql
COPY main.go ./
COPY vendor ./vendor
RUN go build

FROM qnib/uplain-init

RUN apt-get -qq update \
 && apt-get -qq install -y wget \
 && echo "deb [arch=amd64] http://packages.confluent.io/deb/3.3 stable main" > /etc/apt/sources.list.d/confluent.list \
 && wget -qO - http://packages.confluent.io/deb/3.3/archive.key | apt-key add - \
 && apt-get update \
 && apt-get install -y librdkafka1
RUN apt-get install -y -qq vim

ENV ENTRYPOINTS_DIR=/opt/qframe/entry/ \
    LOG_LEVEL=info \
    DOCKER_HOST=unix:///var/run/docker.sock \
    KAFKA_BROKER_HOSTS=tasks.broker \
    KAFKA_BROKER_PORT=9092 \
    HEALTH_BIND_PORT=8123 \
    HEALTH_BIND_HOST=0.0.0.0
COPY --from=build \
     /usr/local/src/github.com/qnib/qframe/examples/qframe-ksql/qframe-ksql \
     /usr/local/bin/
COPY resources/opt/qnib/qframe/template.yml /opt/qnib/qframe/template.yml
COPY resources/opt/qframe/entry/*.sh /opt/qframe/entry/
CMD ["qframe-ksql", "--config", "/etc/qframe.yml"]
