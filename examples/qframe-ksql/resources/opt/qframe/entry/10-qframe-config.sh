#!/bin/bash


cat /opt/qnib/qframe/template.yml \
  |sed -e "s/LOG_LEVEL/${LOG_LEVEL}/" \
  |sed -e "s#DOCKER_HOST#${DOCKER_HOST}#" \
  |sed -e "s/KAFKA_BROKER_HOSTS/${KAFKA_BROKER_HOSTS}/" \
  |sed -e "s/KAFKA_BROKER_PORT/${KAFKA_BROKER_PORT}/" \
  |sed -e "s#HEALTH_BIND_PORT#${HEALTH_BIND_PORT}#" \
  |sed -e "s#HEALTH_BIND_HOST#${HEALTH_BIND_HOST}#" \
  > /etc/qframe.yml

if [[ "X${LOG_ONLY_PLUGINS}" != "X" ]];then
    sed -i'' -e "s/\(#\)only-plugins:.*/only-plugins: \"${LOG_ONLY_PLUGINS}\"/" /etc/qframe.yml
fi
