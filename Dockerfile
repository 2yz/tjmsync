FROM debian:8.5

RUN set -ex \
	&& apt-get update \
	&& apt-get install -y git rsync adduser \
	&& rm -rf /var/lib/apt/lists/* \
	&& ln -s -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
	# Debian Mirror
	&& adduser \
		--system \
		--home=/opt/ftp-master.debian.org/archvsync/ \
		--shell=/bin/bash \
		--no-create-home \
		--group \
		archvsync \
	&& mkdir -p /opt/ftp-master.debian.org/ \
	&& git clone \
		https://ftp-master.debian.org/git/archvsync.git \
		/opt/ftp-master.debian.org/archvsyn \
	&& chown -R archvsync:archvsync /opt/ftp-master.debian.org/archvsyn
	
ENV PATH /opt/ftp-master.debian.org/archvsync/bin:${PATH}

ADD bin/* /usr/local/bin/

EXPOSE 8080

CMD ["tjmsync", "-conf", "/etc/tjmsync.toml"]
