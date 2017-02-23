FROM alpine:3.5

ADD bin/tjmsync /usr/local/bin/

EXPOSE 8080

CMD ["tjmsync", "-conf", "/etc/tjmsync.toml"]