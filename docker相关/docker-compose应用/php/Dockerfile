FROM php:7.3-fpm

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y locales

RUN sed -i -e 's/# en_US.UTF-8 UTF-8/en_US.UTF-8 UTF-8/' /etc/locale.gen && \
    dpkg-reconfigure --frontend=noninteractive locales && \
    update-locale LANG=en_US.UTF-8

ENV LANG en_US.UTF-8 

RUN apt-get install -y \
	gcc \
	build-essential \
	curl \
	autoconf \
	automake \
	bzip2 \
	libz-dev \
        libssl-dev \
        libnghttp2-dev \
        libpcre3-dev \
	libfreetype6-dev \
        libmcrypt-dev \
        libjpeg-dev \
        libpng-dev \
	&& apt-get clean \
        && apt-get autoremove \
	&& docker-php-ext-install -j$(nproc) iconv

RUN cd /var/www/html \
 	&& curl -fsSL 'https://sourceforge.net/projects/freetype/files/freetype2/2.7.1/freetype-2.7.1.tar.gz/download' -o freetype-2.7.1.tar.gz \
        && tar -zxf freetype-2.7.1.tar.gz \
	&& cd freetype-2.7.1 \        
	&& ./configure --prefix=/usr/local/freetype --without-harfbuzz \
        && make && make install \
        && cd .. \
	&& rm -rf freetype-2.7.1.tar.gz \
        && rm -rf freetype-2.7.1 \
	&& docker-php-ext-configure gd \
		--enable-gd-native-ttf \
		--with-freetype-dir=/usr/local/freetype \
		--with-png-dir=/usr/include \
		--with-jpeg-dir=/usr/include \ 
	&& docker-php-ext-install -j$(nproc) gd

RUN docker-php-ext-install bcmath \
        && docker-php-ext-enable bcmath

RUN docker-php-ext-install gettext \
        && docker-php-ext-enable gettext
	
RUN docker-php-ext-install opcache \
	&& docker-php-ext-enable opcache

RUN docker-php-ext-install mysqli \
        && docker-php-ext-enable mysqli

RUN docker-php-ext-install sockets \
        && docker-php-ext-enable sockets

RUN docker-php-ext-install pdo_mysql \
	&& docker-php-ext-enable pdo_mysql
	
RUN pecl install redis-4.0.1 \
	&& docker-php-ext-enable redis
	
RUN apt-get install -y libssl-dev openssl

RUN curl -fsSL 'http://pecl.php.net/get/swoole-4.2.9.tgz' -o swoole-4.2.9.tgz \
    && tar -zxf swoole-4.2.9.tgz \
    && rm -rf swoole-4.2.9.tgz \
    && cd swoole-4.2.9 \
    && /usr/local/bin/phpize \
    && ./configure --enable-openssl --enable-sockets --enable-async-redis --enable-mysqlnd --enable-coroutine \
    && make -j$(nproc) \
    && make install \
    && rm -rf swoole-4.2.9 \
    && docker-php-ext-enable swoole
	
RUN curl -fsSL 'http://pecl.php.net/get/mongodb-1.5.3.tgz' -o mongodb-1.5.3.tgz \
	&& pecl install mongodb-1.5.3 \
	&& rm -f mongodb-1.5.3.tgz \
	&& docker-php-ext-enable mongodb

COPY start.sh /usr/local/bin/swchat-start.sh

RUN chmod +x /usr/local/bin/swchat-start.sh 

#RUN ["/bin/bash","-c","swchat-start.sh"]
#docker-compose exec -d php7.3 swchat-start.sh
