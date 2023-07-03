FROM ghcr.io/guoyk93/minit:latest AS minit

FROM golang:1.19 AS minit-demo
ENV CGO_ENABLED 0
WORKDIR /go/src/app
ADD . .
RUN go build -o /minit-demo ./cmd/minit-demo


FROM alpine:3.16

ENV LANG zh_CN.UTF-8
ENV TZ Asia/Shanghai

RUN sed -i "s/dl-cdn.alpinelinux.org/mirrors.cloud.tencent.com/g" /etc/apk/repositories && \
    apk upgrade --no-cache && \
    apk add --no-cache bash coreutils tzdata ca-certificates && \
    ln -snf /usr/share/zoneinfo/$TZ /etc/localtime && \
    echo $TZ > /etc/timezone && \
    echo $TZ > /etc/TZ

RUN mkdir -p /opt/bin
ENV PATH "/opt/bin:${PATH}"

COPY --from=minit /minit /minit
COPY --from=minit-demo /minit-demo /minit-demo

ADD banner.txt /etc/banner.minit.txt
ADD bin     /opt/bin
ADD minit.d /etc/minit.d
ADD templates /opt/templates

ENV MINIT_MAIN          demo-env-main arg1 arg2
ENV MINIT_MAIN_DIR      /tmp
ENV MINIT_MAIN_NAME     demo-env-main
ENV MINIT_MAIN_GROUP    demo-main
ENV MINIT_MAIN_KIND     cron
ENV MINIT_MAIN_IMMEDIATE true
ENV MINIT_MAIN_CRON     "@every 5s"
ENV MINIT_MAIN_CHARSET  gbk

ENV MINIT_DISABLE @demo-daemon-a,demo-daemon-b-2

ENV MINIT_ENV_EXTRA_ENV_A '{{stringsToUpper "aaa"}}'
ENV MINIT_ENV_EXTRA_ENV_C '{{stringsToUpper "ccc"}}'

ENTRYPOINT [ "/minit" ]

CMD ["demo-arg-main", "arg1", "arg2"]