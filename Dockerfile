# build release image
FROM centos:7

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

COPY bin/influx-proxy /opt/influxdb-proxy/

WORKDIR /opt/influxdb-proxy/

CMD ["./influx-proxy"]
