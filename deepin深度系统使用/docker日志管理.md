#### docker

```
docker ps

查看命令格式
docker logs -h

Usage:	docker logs [OPTIONS] CONTAINER

Fetch the logs of a container

Options:
      --details        Show extra details provided to logs
  -f, --follow         Follow log output
      --since string   Show logs since timestamp (e.g. 2013-01-02T13:23:37) or relative (e.g. 42m for 42 minutes)
      --tail string    Number of lines to show from the end of the logs (default "all")
  -t, --timestamps     Show timestamps
      --until string   Show logs before a timestamp (e.g. 2013-01-02T13:23:37) or relative (e.g. 42m for 42 minutes)

命令说明：
–since : 此参数指定了输出日志开始日期，即只输出指定日期之后的日志。
-f : 查看实时日志
-t : 查看日志产生的日期
-tail=10 : 查看最后的10条日志
```

#### docker-compose

```
docker-compose ps

docker-compose logs -h
Usage: logs [options] [SERVICE...]

Options:
    --no-color          Produce monochrome output.
    -f, --follow        Follow log output.
    -t, --timestamps    Show timestamps.
    --tail="all"        Number of lines to show from the end of the logs
                        for each container.
注：没有开始时间

使用例子：
docker logs -f -t --since=2017-05-31 --tail=10 php
```

#### 日志清理

日志目录：

/var/lib/docker/containers/容器ID/*-json.log

> 方式一

在Linux或者Unix系统中，通过rm -rf或者文件管理器删除文件，将会从文件系统的目录结构上解除链接（unlink）。如果文件是被打开的（有一个进程正在使用），那么进程将仍然可以读取该文件，磁盘空间也一直被占用。正确姿势是cat /dev/null > *-json.log，当然你也可以通过rm -rf删除后重启docker


但是，这样清理之后，随着时间的推移，容器日志会像杂草一样，卷土重来。

> 方式二

设置Docker容器日志大小
```
nginx: 
  image: nginx:1.12.1 
  restart: always 
  logging: 
    driver: “json-file” 
    options: 
      max-size: “5g” 
```

> 方式三

全局设置/etc/docker/daemon.json
```
vim /etc/docker/daemon.json

{
  "registry-mirrors": ["http://f613ce8f.m.daocloud.io"],
  "log-driver":"json-file",
  "log-opts": {"max-size":"500m", "max-file":"3"}
}

max-size=500m，意味着一个容器日志大小上限是500M，
max-file=3，意味着一个容器有三个日志，分别是id+.json、id+1.json、id+2.json。

//重启docker守护进程
#systemctl daemon-reload

#systemctl restart docker

注意：这种方式设置的日志大小，只对新建的容器有效
```

> 方式四
```
vim /etc/profile.d/clean_docker_log.sh

#!/bin/sh
echo "======== start clean docker containers logs ========"
logs=$(find /var/lib/docker/containers/ -name *-json.log)
for log in $logs
do
echo "clean logs : $log"
cat /dev/null > $log
done
echo "======== end clean docker containers logs ========"

权限
chmod +x /etc/profile.d/clean_docker_log.sh

执行
/etc/profile.d/clean_docker_log.sh
```