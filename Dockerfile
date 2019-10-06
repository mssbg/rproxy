FROM alpine:latest
EXPOSE 8080 2012
ENV RPROXY_STRATEGY random
COPY config.yml /
COPY rproxy.linux /rproxy
CMD /rproxy
