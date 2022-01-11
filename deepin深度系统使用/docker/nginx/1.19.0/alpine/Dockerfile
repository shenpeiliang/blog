FROM daocloud.io/library/nginx:1.19.0-alpine

ARG timezone

ENV TIMEZONE=${timezone:-"Asia/Shanghai"}

ENV LANG="en_US.UTF-8"

RUN set -ex \
    && echo "https://mirrors.tuna.tsinghua.edu.cn/alpine/v3.11/main/" > /etc/apk/repositories \
    && echo "https://mirrors.tuna.tsinghua.edu.cn/alpine/v3.11/community/" >> /etc/apk/repositories \
    && apk update \
	&& ln -sf /usr/share/zoneinfo/${TIMEZONE} /etc/localtime \
	# By default the Alphine image automatically streams Nginx logs (access and error logs) to stdout and stderr
	# by creating a symbolic link from stdout to /var/log/nginx/access.log and stderr to /var/log/nginx/error.log
	# Lets remove these symbolic links so that we can setup a volume and persist the logs so that they are available
	# between restarts
    && unlink /var/log/nginx/access.log \
    && unlink /var/log/nginx/error.log \
    && echo "${TIMEZONE}" > /etc/timezone \
    && echo -e "\033[42;37m Build Completed :).\033[0m\n"