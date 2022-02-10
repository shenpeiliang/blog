#### 负载均衡方式

1. round-robin 对于访问请求来说，将请求循环分发到应用服务器上，nginx默认使用此方式

2. least-connect 下一个请求被分发到当前活动连接数最少的应用服务器上

3. ip-hash 通过hash方法决定将下一个请求分发到哪一个应用服务器上


proxy_pass指向的是upstream的配置服务器上

#### 配置说明

> round-robin轮询方式

weight按权重（比例分配）按顺序轮询访问72.24.0.4服务器10次，之后再到72.24.0.4服务器1次，按此规则不断轮询

默认情况下每台应用的权重为1

```
http{
    upstream backserver { 
        server 172.24.0.4:9000 weight=10; 
        server 172.24.0.5:9000 weight=1; 
    } 
    server {
        listen 80;
        location / {
            proxy_pass http://www.my.com;
        }
    }
}

```

> least-connect方式

它会把新的请求分发到负载量较小的应用服务器上

```
http{
    upstream backserver { 
        least_conn;
        server 172.24.0.4:9000; 
        server 172.24.0.5:9000; 
    } 
    server {
        listen 80;
        location / {
            proxy_pass http://www.my.com;
        }
    }
}

```

> ip-hash会话持久性

通过将客户端的ip地址作为hash键去决定将客户端的请求分发到哪一台应用服务器上，这种方式保证了来自同一客户端的请求总是会被分发到特定的应用服务器上去，除非这个特定的服务器停止了工作

```
http{
    upstream backserver { 
        ip_hash;
        server 172.24.0.4:9000; 
        server 172.24.0.5:9000; 
    } 
    server {
        listen 80;
        location / {
            proxy_pass http://www.my.com;
        }
    }
}

```

#### 临时移除某台应用服务器

如果某一台应用服务器需要临时移除，不允许请求访问，我们可以在其后面使用 down 来标记此台服务器，这时nginx不会将请求分发到这台应用服务器上面。当先前由这台服务器响应的客户端再次发起请求的时候，nginx会自动将其分发到其他的应用服务器上

```
http{
    upstream backserver { 
        ip_hash;
        server 172.24.0.4:9000; 
        server 172.24.0.5:9000 down; 
    } 
    server {
        listen 80;
        location / {
            proxy_pass http://www.my.com;
        }
    }
}

```

#### 备用应用服务器

在实际应用中不是所有的应用服务器都要参与，可以将一台或者几台应用作为备用服务器，当其他的应用服务器出现问题不能访问的时候，Nginx会自动启动备用的应用服务器

```
http{
    upstream backserver { 
        ip_hash;
        server 172.24.0.4:9000; 
        server 172.24.0.5:9000 backup; 
    } 
    server {
        listen 80;
        location / {
            proxy_pass http://www.my.com;
        }
    }
}

```

#### 宕机容错机制

超时未访问成功自动轮询下一台服务器

```
server {
    listen      80;
    server_name  0.sk0.com;

    root /usr/share/nginx/html/sk/sites/0.shikee.com/root;
    include public.conf;

    location ~ \\.php
    {
        root /var/www/html/sk/sites/0.shikee.com/root;
        include        fastcgi.conf;

        proxy_connect_timeout 5; #最大连接时间
        proxy_send_timeout 5; #最大发送时间
        proxy_read_timeout 5; #最大读取时间
    }

}

```