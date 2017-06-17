#!/bin/bash


cat /opt/qnib/qframe/template.yml \
  |sed -e "s/LOG_LEVEL/${LOG_LEVEL}/" \
  |sed -e "s/INFLUXDB_HOST/${INFLUXDB_HOST}/" \
  |sed -e "s/DOCKER_HOST/${DOCKER_HOST}/" \
  |sed -e "s/INFLUXDB_DB/${INFLUXDB_DB}/" \
  |sed -e "s/ELASTICSEARCH_HOST/${ELASTICSEARCH_HOST}/" \
  |sed -e "s#GROK_PATTERNS_DIR#${GROK_PATTERNS_DIR}#" \
  > /etc/qframe.yml

if [[ "X${LOG_PLUGINS}" != "X" ]];then
   echo ">> Only plugins plugins: ${LOG_PLUGINS}"
   sed -i '' -e "s/[#]*only-plugins:.*/only-plugins:\"${LOG_PLUGINS}" /etc/qframe.yml
fi