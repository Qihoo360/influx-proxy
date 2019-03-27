# build bin
FROM golang:1.10 as builder

COPY . /go/src/github.com/wilhelmguo/influx-proxy/

RUN go build -o bin/influx-proxy github.com/wilhelmguo/influx-proxy/service


# build release image
FROM centos:7

RUN ln -sf /usr/share/zoneinfo/Asia/Shanghai /etc/localtime

COPY --from=builder /go/src/github.com/wilhelmguo/influx-proxy/bin/influx-proxy  /opt/influxdb-proxy/

WORKDIR /opt/influxdb-proxy/

CMD ["./influx-proxy"]
