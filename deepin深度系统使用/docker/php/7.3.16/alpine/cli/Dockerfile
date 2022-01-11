FROM php:7.3.16-cli-alpine3.11

ARG timezone

ENV TIMEZONE=${timezone:-"Asia/Shanghai"} \
	PHPIZE_DEPS="autoconf dpkg-dev dpkg file g++ gcc libc-dev make pkgconf re2c pcre-dev pcre2-dev zlib-dev libtool automake libxml2-dev" 

ENV LANG="en_US.UTF-8"

COPY install-php-extensions /usr/local/bin/install-php-extensions

RUN set -ex \
    && sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories \
    && apk update \
    && apk add --no-cache libstdc++ openssl openssl-dev libaio-dev bash \
    && apk add --no-cache --virtual .build-deps  $PHPIZE_DEPS \
    && apk add \
    # Install base packages ('ca-certificates' will install 'nghttp2-libs')
    ca-certificates \
    curl \
    wget \
    tar \
    gzip \
    xz \
    libressl \
    tzdata \
    pcre \   
	# 使用脚本安装php扩展
	&& chmod +x /usr/local/bin/install-php-extensions && sync \
	&& install-php-extensions bcmath gettext opcache pdo pdo_mysql gd redis \
	&& docker-php-ext-install sockets && docker-php-source extract \
    # 修改配置项
    && cd /usr/local/etc/php/ \
    && cp php.ini-production php.ini \
    # - config PHP
    && { \
        echo "upload_max_filesize=100M"; \
        echo "post_max_size=108M"; \
        echo "memory_limit=1024M"; \
        echo "date.timezone=${TIMEZONE}"; \
        echo "opcache.enable_cli = 'On'"; \
    } | tee /usr/local/etc/php/conf.d/99-overrides.ini \
	# 安装Composer
	&& php -r "copy('https://install.phpcomposer.com/installer', 'composer-setup.php');" \
	# 安装
	&& php composer-setup.php \
	# 删除
	&& php -r "unlink('composer-setup.php');" \
	&& chmod u+x composer.phar \
	&& mv composer.phar /usr/local/bin/composer \
	# 将 Composer 镜像设置为阿里云镜像，加速国内下载速度
	&& composer config -g repo.packagist composer https://mirrors.aliyun.com/composer \
	# ---------- clear works ----------
    && apk del .build-deps \
    && rm -rf /tmp/* /var/cache/apk/* \
    && php -v \
    && php -m \
    # 时区设置
    && ln -sf /usr/share/zoneinfo/${TIMEZONE} /etc/localtime \
    && echo "${TIMEZONE}" > /etc/timezone