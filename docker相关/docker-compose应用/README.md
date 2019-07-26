### 简介

环境：Centos7


### 安装环境

需要安装docker-compose

官网：
https://docs.docker.com/compose/startup-order/


下载地址：
https://github.com/docker/compose/releases


安装：
curl -L "https://github.com/docker/compose/releases/download/1.23.1/docker-compose-$(uname -s)-$(uname -m)" -o /usr/local/bin/docker-compose


启动命令：

```
docker-compose up -d 

```

如果修改了Dockerfile文件需要重新生效，则启动命令：

```
docker-compose up -d --build

```

查看命令详细说明：

```
docker-compose command --help
```

例如：

```
docker-compose restart --help
```


重启单个服务

```
docker-compose restart elasticsearch6.4.3
```

查看日志，类似tail -f xx.file

```
docker-compose logs -f --tail 10 nginx
```


进入docker容器内部

1.docker-compose exec container_name bash

2.docker-compose run container_name bash

注意是容器的name，不是id

exec回直接进入容器，而run则是在当前容器基础上新建一个一摸一样的容器，相当于clone一个吧。所以exec以后，修改就是对原来容器的修改，而run的修改则与原来的无关。


### 启动后效果图

![image](https://github.com/shenpeiliang/blog/blob/master/docker相关/docker-compose应用/image/docker-compose状态.png)