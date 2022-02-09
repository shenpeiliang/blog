#### Cookie与Session的关系

cookie和session的分别属于客户端和服务端，服务端的session的实现需要客户端的cookie信息

服务端执行session机制时会生成sessionID，然后把该ID发送给客户端并让客户端保存下来，客户端每次请求都会把这个ID放到http请求头发送给服务端，因此当我们完全禁掉浏览器的cookie的时候，服务端的session也会不能正常使用


#### 如何解决Session共享问题

在集群服务中，客户端的请求会被随机分配到各个服务器中，这就回导致Session丢失，客户端的表现将是一直跳转到登录页面，通常有如下几种方案解决：

> 1. 使用Cookie实现

原理是将系统用户的Session信息加密、序列化后，以Cookie的方式统一存在客户端，并以根域名保存

优点：

1. 不需要额外的服务器资源
2. 简单方便

缺点：

1. 受http协议头信息长度的限制，仅能够存储小部分的用户信息
2. 需要对内容进行加密解密，存在安全隐患
3. 占用一定的带宽资源


> 2. 使用Nginx中的负载均衡算法

1. ip hash，根据客户端的IP，将请求分配到不同的服务器上

2. sticky，根据服务器给客户端的cookie，客户端再次请求时会带上此cookie，nginx会把有此cookie(该cookie记录了具体服务器信息,比如服务器IP的Hash值等)的请求转发到颁发cookie的服务器上

缺点：

1. 如果指定的那台服务器宕机了，session既然也就不存在了

2. 变成了单个节点，不能做负载均衡

参考：

http://tengine.taobao.org/book/chapter_05.html#id6

https://zhuanlan.zhihu.com/p/194088100

> 3. 使用数据库同步session

每次将session数据存到数据库中，但它的缺点也非常明显，对于Session的并发读写能力取决于MySQL数据库的性能，对数据库的压力非常大，同时需要自己实现Session淘汰逻辑，以便定时从数据表中更新、删除 Session记录，当并发过高时即使使用了行级锁也是很大问题

> 4. Session replication方式

所有服务器上通过应用之间的session同步都存储一份所有用户的session，但明显占用服务器资源，因为同步所以将会消耗大量的带宽


> 5. Session数据集中存储

再构建一个专门的服务器(为了高可用该服务器也是集群部署)用来存储所有的会话信息,所有服务器获取会话信息都从该服务器上获取，比如可以利用Redis集群来做，这也是我们目前常用的方案


> 6. 使用token代替Session

使用类似JWT的方式来替换掉Session实现数据共享

大概的流程是这样的： 
```
1、客户端通过用户名和密码登录服务器；
2、服务端对客户端身份进行验证；
3、服务端对该用户生成Token，返回给客户端；
4、客户端将Token保存到本地浏览器，一般保存到cookie中；
5、客户端发起请求，需要携带该Token；
6、服务端收到请求后，首先验证Token，之后返回数据
```

优点：

1. 无状态、可扩展 ：在客户端存储的Token是无状态的，并且能够被扩展。基于这种无状态和不存储Session信息，负载均衡器能够将用户信息从一个服务传到其他服务器上

2. 安全：请求中发送token而不再是发送cookie能够防止CSRF(跨站请求伪造)

3. 可提供接口给第三方服务：使用token时，可以提供可选的权限给第三方应用程序

4. 多平台跨域

对应用程序和服务进行扩展的时候，需要介入各种各种的设备和应用程序。 假如我们的后端api服务器a.com只提供数据，而静态资源则存放在cdn 服务器b.com上。当我们从a.com请求b.com下面的资源时，由于触发浏览器的同源策略限制而被阻止

我们通过CORS（跨域资源共享）标准和token来解决资源共享和安全问题

举个例子，我们可以设置b.com的响应首部字段为：

```
Access-Control-Allow-Origin: http://a.com

Access-Control-Allow-Headers: Authorization, X-Requested-With, Content-Type, Accept

Access-Control-Allow-Methods: GET, POST, PUT,DELETE
```

第一行指定了允许访问该资源的外域 URI。

第二行指明了实际请求中允许携带的首部字段，这里加入了Authorization，用来存放token

第三行用于预检请求的响应。其指明了实际请求所允许使用的 HTTP 方法

然后用户从a.com携带有一个通过了验证的token访问B域名，数据和资源就能够在任何域上被请求到
