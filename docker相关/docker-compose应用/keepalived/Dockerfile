FROM keepalived:2.0.16

RUN apt-get update && DEBIAN_FRONTEND=noninteractive apt-get install -y locales

RUN sed -i -e 's/# en_US.UTF-8 UTF-8/en_US.UTF-8 UTF-8/' /etc/locale.gen && \
    dpkg-reconfigure --frontend=noninteractive locales && \
    update-locale LANG=en_US.UTF-8

ENV LANG en_US.UTF-8 

RUN wget https://www.keepalived.org/software/keepalived-2.0.16.tar.gz -o keepalived.tar.gz \
    && tar -zxf  keepalived.tar.gz \
    && rm -rf keepalived.tar.gz \
	&& cd keepalived \
	&& ./configure --disable-fwmark --prefix=/usr/local/keepalived \
	&& make \
    && make install \