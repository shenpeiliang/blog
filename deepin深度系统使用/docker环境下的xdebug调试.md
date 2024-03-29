> 配置一

![image](image/xdebug-1.png)

设置 >> 调试

查看Xdebug调试端口，端口号和php.ini中的xdebug配置一致：

```
xdebug3配置：

xdebug.client_port = 9003
xdebug.client_host = host.docker.internal

xdebug2配置：
xdebug.remote_host = host.docker.internal
xdebug.remote_port = 9003
```

注意：
docker-compose.yml 内容
```
php
 ...
 extra_hosts:
        - host.docker.internal:host-gateway
```

宿主机防火墙需要开放端口，否则通知不了phpstorm
```
ufw allow 9003
```

> 配置二

![image](image/xdebug-2.png)

运行/调试配置

注意添加配置为：php远程调试

会话ID填写值与配置一致:

xdebug.idekey = PHPSTORM

> 配置三

![image](image/xdebug-3.png)

服务配置：

主机：访问主域名

勾选使用路径映射

文件/目录：本地站点绝对路径

服务器上的绝对路径： 站点所在docker环境php容器内的绝对路径

> 其他

phpstorm激活，试用期为30天，可以通过插件不断重置

设置 >> 插件

安装IDE Eval Reset

可设置每天开启后就更新，或者手动
Click "Help" menu and select "Eval Reset"