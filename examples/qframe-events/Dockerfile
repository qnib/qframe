FROM qnib/uplain-golang

WORKDIR /usr/local/src/github.com/qnib/qframe-events
COPY main.go ./main.go
COPY vendor/vendor.json ./vendor/vendor.json
RUN govendor fetch +m \
 && govendor build

FROM qnib/uplain-init

COPY --from=0 /usr/local/src/github.com/qnib/qframe-events/qframe-events \
     /usr/local/bin/
COPY qframe.yml /etc/
COPY patterns
ENV SKIP_ENTRYPOINTS=true
CMD ["qframe-events", "--config", "/etc/qframe.yml"]
