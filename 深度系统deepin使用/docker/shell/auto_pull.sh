#!/bin/bash
GIT=/usr/local/bin/git
WWW_ROOT='/home/docker/html/www'

for dir in $WWW_ROOT/{cms}
do
        cd $dir
        su - docker -c $GIT pull --quiet > /dev/null 2>&1
done
