#### 限流配置

limit_req_zone 用来限制单位时间内的请求数，即速率限制,采用的漏桶算法 

limit_req_conn 用来限制同一时间连接数，即并发限制

参考：

https://www.nginx.cn/doc/standard/httplimitzone.html

https://www.nginx.cn/doc/standard/httplimitrequest.html

> limit_req_zone 按请求速率限流

HttpLimitReqest模块与HttpLimit zone模块，limit_req_zone定义到http区域，limit_req定义到server或location区域，通常limit_req用于限制动态内容的访问频率


```
格式：
limit_req zone=name [burst=number] [nodelay];

```

例如：
```
limit_req_zone $binary_remote_addr zone=one:10m rate=1r/s;
```

1. $binary_remote_addr表示通过remote_addr这个标识来做限制，“binary_”的目的是缩写内存占用量，是限制同一客户端ip地址

2. zone=one:10m表示生成一个大小为10M，名字为one的内存区域，用来存储访问的频次信息

3. rate=1r/s表示允许相同标识的客户端的访问频次，这里限制的是每秒1次，或者其他比如30r/m

注意：

Nginx的限流统计是基于毫秒的，我们设置的速度是2r/s，转换一下就是500ms内单个IP只允许通过1个请求，从501ms开始才允许通过第二个请求

例如：
```
limit_req zone=one burst=5 nodelay;
```

1. zone=one 设置使用哪个配置区域来做限制，与上面limit_req_zone 里的name对应

2. burst=5 设置一个大小为5的缓冲区，当有大量请求（爆发）过来时，超过了访问频次限制的请求可以先放到这个缓冲区内，默认为0

3. nodelay 如果设置了该参数，超过访问频次而且缓冲区也满了的时候就会直接返回503，如果没有设置，则所有请求会等待排队，通常与burst一起使用。如果设置该选项则第一时间处理，反之严格使用平均速率限制请求数


实例：
```
http {
    limit_req_zone $binary_remote_addr zone=one:10m rate=10r/s; //所有访问ip限制每秒10个请求
    ...
    server {
        ...
        location  ~ \.php$ {
            limit_req zone=one burst=5 nodelay;   //执行的动作,通过zone名字对应，此时如果有20个并发请求，那么只能处理总共10个，完成5个，有5个被放到缓存队列中等待处理
                }
            }
    }
```
http段内定义触发条件，可以有多个条件，在location内定义达到触发条件时nginx所要执行的动作

配置可以限制特定UA（比如搜索引擎）的访问
```
http {
   limit_req_zone  $anti_spider  zone=one:10m   rate=10r/s;
    ...
    server {
        limit_req zone=one burst=10 nodelay;
        if ($http_user_agent ~* "googlebot|bingbot|Feedfetcher-Google") {
            set $anti_spider $http_user_agent;
        }
        ...
    }
```

对应的日志参数格式：
```
limit_req_log_level info | notice | warn | error;
```

当服务器由于limit被限速或缓存时，配置写入日志，默认值是error


对应的拒绝请求的返回值：
```
limit_req_status code;
```
默认503,只能设置 400 到 599 之间


> limit_req_conn 按连接数限流

用来限制单个IP的请求数，在服务器处理了请求并且已经读取了整个请求头后统计计数

格式：
```
limit_conn zone number;

limit_conn_zone key zone=name:size;
```

实例：
```
limit_conn_zone $binary_remote_addr zone=perip:10m;
limit_conn_zone $server_name zone=perserver:10m;

server {
    ...
    limit_conn perip 10;
    limit_conn perserver 100;
}

```

对应的日志：
```
limit_conn_log_level info | notice | warn | error;
```

默认值：error

对应的日志级别：
```
limit_conn_status code;
```
默认503

> 白名单设置

http_limit_conn和http_limit_req模块限制了单ip单位时间内的并发和请求数，但是如果Nginx前面有lvs或者haproxy之类的负载均衡或者反向代理，nginx获取的都是来自负载均衡的连接或请求，这时不应该限制负载均衡的连接和请求，就需要geo和map模块设置白名单

```
geo $whiteiplist  {
        default 1;
        10.11.15.161 0;
    }
map $whiteiplist  $limit {
        1 $binary_remote_addr;
        0 "";
    }
limit_req_zone $limit zone=one:10m rate=10r/s;
limit_conn_zone $limit zone=addr:10m;
```
geo模块定义了一个默认值是1的变量whiteiplist，当在ip在白名单中，变量whiteiplist的值为0，反之为1

如果在白名单中whiteiplist=0则limit=””，不会存储到10m的会话状态中，反之在白名单中whiteiplist=1则limit=$binary_remote_addr，会存储到10m的会话状态中，也就被限制


