#!/bin/bash


cat /opt/qnib/gcollect/template.yml \
  |sed -e "s/LOG_LEVEL/${LOG_LEVEL}/" \
  |sed -e "s/DOCKER_LOG_SINCE/${DOCKER_LOG_SINCE}/" \
  |sed -e "s#DOCKER_HOST#${DOCKER_HOST}#" \
  |sed -e "s/INFLUXDB_HOST/${INFLUXDB_HOST}/" \
  |sed -e "s/INFLUXDB_DB/${INFLUXDB_DB}/" \
  |sed -e "s/ELASTICSEARCH_HOST/${ELASTICSEARCH_HOST}/" \
  |sed -e "s#GROK_PATTERNS_DIR#${GROK_PATTERNS_DIR}#" \
  |sed -e "s#HEALTH_BIND_PORT#${HEALTH_BIND_PORT}#" \
  |sed -e "s#HEALTH_BIND_HOST#${HEALTH_BIND_HOST}#" \
  > /etc/gcollect.yml

if [[ "X${LOG_ONLY_PLUGINS}" != "X" ]];then
    sed -i'' -e "s/\(#\)only-plugins:.*/only-plugins: \"${LOG_ONLY_PLUGINS}\"/" /etc/gcollect.yml
fi