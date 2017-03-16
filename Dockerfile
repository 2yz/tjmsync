FROM debian:8.5

ADD out/* /usr/local/bin/

EXPOSE 3000

CMD ["tjmsync", "-conf", "/etc/tjmsync.toml"]
