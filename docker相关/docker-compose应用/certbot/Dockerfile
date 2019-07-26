FROM phusion/baseimage:latest

MAINTAINER Mahmoud Zalt <mahmoud@zalt.me>

RUN apt-get update
RUN apt-get install -y letsencrypt

COPY run-certbot.sh /root/certbot/run-certbot.sh

ENTRYPOINT bash -c "bash /root/certbot/run-certbot.sh && sleep infinity"