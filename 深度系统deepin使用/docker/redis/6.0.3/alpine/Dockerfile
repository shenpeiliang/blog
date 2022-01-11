FROM daocloud.io/library/redis:6.0.3-alpine3.11
  
ARG timezone

ENV TIMEZONE=${timezone:-"Asia/Shanghai"}

ENV LANG="en_US.UTF-8"

RUN set -ex \
    && echo "https://mirrors.tuna.tsinghua.edu.cn/alpine/v3.11/main/" > /etc/apk/repositories \
    && echo "https://mirrors.tuna.tsinghua.edu.cn/alpine/v3.11/community/" >> /etc/apk/repositories \
    && apk update \
        && ln -sf /usr/share/zoneinfo/${TIMEZONE} /etc/localtime \
    && echo "${TIMEZONE}" > /etc/timezone \
    && echo -e "\033[42;37m Build Completed :).\033[0m\n"
