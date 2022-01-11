FROM daocloud.io/library/php:5.6-fpm-alpine3.8

ARG timezone

ENV TIMEZONE=${timezone:-"Asia/Shanghai"} \
    PHPIZE_DEPS="autoconf dpkg-dev dpkg file g++ gcc libc-dev make pkgconf re2c pcre-dev pcre2-dev zlib-dev libtool automake libxml2-dev" \
	# php扩展
	PHP_EXTENSION="bcmath gettext opcache pdo pdo_mysql mysqli gd redis xdebug mcrypt"

ENV LANG="en_US.UTF-8"

COPY --from=mlocati/php-extension-installer /usr/bin/install-php-extensions /usr/bin/

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
	&& install-php-extensions ${PHP_EXTENSION} \
    # 修改配置项
    && cd /usr/local/etc/php/ \
    && cp php.ini-production php.ini \
    #安装Composer需要下载证书
    && mkdir -p /usr/local/openssl/ssl/certs/ \
    && wget -P /usr/local/openssl/ssl/certs http://curl.haxx.se/ca/cacert.pem \
    # - config PHP
    && { \
        echo "upload_max_filesize=100M"; \
        echo "post_max_size=108M"; \
        echo "memory_limit=1024M"; \
        echo "date.timezone=${TIMEZONE}"; \
        echo "opcache.enable_cli = 'On'"; \
        echo "openssl.cafile=/usr/local/openssl/ssl/certs/cacert.pem"; \
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