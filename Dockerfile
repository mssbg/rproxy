FROM alpine:latest
EXPOSE 8080 2012
COPY config.yml /
COPY rproxy.linux /rproxy
CMD /rproxy
