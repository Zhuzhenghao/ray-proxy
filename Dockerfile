FROM --platform=$BUILDPLATFORM release-ci.daocloud.io/docker/golang:1.19.10 as build

WORKDIR /root

ENV GO111MODULE=on GOPROXY=https://goproxy.cn,direct

COPY . .

ARG TARGETARCH
ARG LDFLAGS
RUN CGO_ENABLED=0 GOOS=linux GOARCH=$TARGETARCH go build -mod vendor -o ray-proxy -ldflags "-w -s $LDFLAGS" ./cmd/main.go

FROM release-ci.daocloud.io/docker/alpine:3.18

USER root

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && apk add --no-cache curl && mkdir /etc/template

COPY --from=build /root/ray-proxy /bin/

COPY ./artifacts/template /etc/template

RUN chmod +x /bin/ray-proxy

CMD ["/bin/ray-proxy"]
