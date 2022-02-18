#### 简介

JWT，JSON Web Token（简称 JWT）是目前最流行的跨域认证解决方案，是为了在网络应用环境间传递声明而执行的一种基于 JSON 的开放标准（RFC 7519）。JWT 的声明一般被用来在身份提供者和服务提供者间传递被认证的用户身份信息，以便于从资源服务器获取资源，比如用在用户登录上


注意：

- JWT 不依赖 Cookie，因此可以不需要担心跨域资源共享问题（CORS）

当用户希望访问一个受保护的路由或者资源的时候，可以把它放在 Cookie 里面自动发送，但是这样不能跨域，所以更好的做法是放在 HTTP 请求头信息的 Authorization 字段里，使用 Bearer 模式添加 JWT

参考： https://developer.mozilla.org/en-US/docs/Web/HTTP/Authentication#authentication_schemes

- JWT 最大的优势是服务器不再需要存储 Session，使得服务器认证鉴权业务可以方便扩展。但这也是 JWT 最大的缺点：由于服务器不需要存储 Session 状态，因此使用过程中无法废弃某个 Token 或者更改 Token 的权限。也就是说一旦 JWT 签发了，到期之前就会始终有效，除非服务器部署额外的逻辑

- JWT 本身包含了认证信息，任何人都可以获得该令牌的所有权限。为了减少其被盗用，JWT的有效期应该设置得比较短。对于一些比较重要的权限，使用时应该再次对用户进行认证

- JWT 适合一次性的命令认证，颁发一个有效期极短的 JWT，即使暴露了危险也很小，由于每次操作都会生成新的 JWT，因此也没必要保存 JWT，真正实现无状态

- 为了减少盗用，JWT 不应该使用 HTTP 协议明码传输，要使用 HTTPS 协议传输，返回 jwt 给客户端时设置 httpOnly=true 并且使用 cookie 而不是 LocalStorage 存储 jwt，这样可以防止 XSS 攻击和 CSRF 攻击


#### 认证流程

1. 客户端通过用户名和密码登录服务器

2. 服务端对客户端身份进行验证

3. 服务端对该用户生成jwt，返回给客户端

4. 客户端将token保存到本地浏览器，一般保存到cookie或localstorage

5. 客户端发起请求，需要在请求头的 Authorization 字段中携带该token，或以Cookie的方式存储但服务端需要设置支持CORS(跨来源资源共享)策略

例如：
```
Access-Control-Allow-Origin: *
#或指定具体域名
Access-Control-Allow-Origin: http://www.demo.com

Access-Control-Allow-Headers: Authorization, X-Requested-With, Content-Type, Accept

Access-Control-Allow-Methods: GET, POST, PUT,DELETE
```

6. 服务端收到请求后，首先验证token，之后返回数据


#### 结构

一个 JWT 实际就是一个字符串，它包含三部分，分别是： 头部（header），载荷（payload），签名 （signature），以点号.字符分隔，Base64(header).Base64(payload).H256(Base64(header).Base64(payload))，例如：
```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJpYXQiOjE0MTY5MjkxMDksImp0aSI6ImFhN2Y4ZDBhOTVjIiwic2NvcGVzIjpbInJlcG8iLCJwdWJsaWNfcmVwbyJdfQ.XCEwpBGvOLma4TCoh36FU7XhUbcskygS81HE1uHLf0E
```

> 头部（header）

```
{
    'typ': 'JWT',
    'alg': 'HS256'
}
```

jwt的头部承载两部分信息：

- typ 声明类型

- alg 声明加密的算法


alg属性表示签名的算法（algorithm），默认是 HMAC SHA256（写成 HS256）；typ属性表示这个令牌（token）的类型（type），JWT 令牌统一写为JWT，将json数据进行base64加密构成第一部分

> 载荷（payload）

Payload 部分也是一个 JSON 对象，用来存放实际需要传递的数据

- 标准中注册的声明(建议但不强制使用)

```
iss: jwt签发者

sub: jwt所面向的用户

aud: 接收jwt的一方

exp: jwt的过期时间，这个过期时间必须要大于签发时间

nbf: 定义在什么时间之前，该jwt都是不可用的

iat: jwt的签发时间

jti: jwt的唯一身份标识，主要用来作为一次性token,从而回避重放攻击
```

- 公共的声明

公共的声明可以添加任何的信息，一般添加用户的相关信息或其他业务需要的必要信息

- 私有的声明

私有声明是提供者和消费者所共同定义的声明


如下定义一个简单的payload，将json数据进行base64加密构成第二部分:
```
{
  "sub": "1234567890",
  "name": "John Doe",
  "admin": true
}
```

> 签名 （signature）
jwt的第三部分是一个签证信息，这个签证信息由三部分组成：

- header (base64后的string)

- payload (base64后的string)

- secret (需要指定一个密钥（secret），这个密钥只有服务器才知道，可以设计成和用户相关的属性，而不是一个所有用户公用的统一值，不能泄露给用户)

签名生成算法：
```
// javascript
var encodedString = base64UrlEncode(header) + '.' + base64UrlEncode(payload);

var signature = HMACSHA256(encodedString, 'secret'); // TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ
```

算出签名以后，现在把header和payload用Base64URL 算法对象序列化，然后把这三部分用“.”拼接起来就是生成的jwt:

```
eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJzdWIiOiIxMjM0NTY3ODkwIiwibmFtZSI6IkpvaG4gRG9lIiwiYWRtaW4iOnRydWV9.TJVA95OrM7E2cBab30RMHrHDcEfxjoYZgeFONFh7HgQ
```

#### Token 和 JWT 的区别

> 相同

- 它们都是访问资源的令牌，只有验证成功后，客户端才能访问服务端上受保护的资源

- 都是使服务端无状态化


> 区别

- Token：服务端验证客户端发送过来的 Token 时，还需要查询数据获取用户信息，然后验证 Token 是否有效

- JWT：将 Token 和 Payload 加密后存储于客户端，服务端只需要使用密钥解密进行校验（校验也是 JWT 自己实现的）即可，不需要查询或者减少查询数据库，因为 JWT 自包含了用户信息和加密的数据

- Token只是随机字符串，而JWT则包含可在一个时间范围或域内描述用户身份，授权数据和令牌有效性的信息和元数据，Oauth2兼容，JWT数据可以被检查且有失效控制


#### 应用场景

- 一次性验证以及时效性的特性，可以用于账号注册成功后的激活动作，jwt 的 payload 中固定的参数：iss 签发者（激活账号）和 exp 过期时间可以作为激活链接，先验证链接的合法性，从而避免多余地去查询数据库

- 无状态特性，用于API身份验证，客户端和服务端共享 secret ，过期时间由服务端校验，客户端定时刷新（签名信息不可被修改）

- 单点登录+会话管理（不推荐）

账号注销处理方式：

1. 为了 jwt 的安全性，通常使用 https 协议以及保存 jwt 到 cookie 中，没有会话状态，当用户退出登录时，只能删除 cookie ，但依然可以使用之前的 jwt 继续访问资源

2. secret 设计成和用户相关的值, 清空或修改服务端的用户对应的 secret ， 修改密码的情况应该强制修改模拟重新登录

3. 借助 redis 管理 jwt 的状态， jwt 作为 key 记录每次操作的最新时间，这种方式类似 session 的方式，把无状态变成了有状态

#### 实现

- 编码
```
public function encode($payload, $key, $alg, $head = null): string
{
    $header = [
        'typ' => 'JWT',
        'alg' => $alg
    ];

    if (isset($head) && is_array($head))
        $header = array_merge($head, $header);

    //三部分
    $segments = [];

    $url_safe_base64 = new UrlSafeBase64();

    //头部（header）
    $segments[] = $url_safe_base64->encode(json_encode($header, JSON_UNESCAPED_SLASHES));

    //载荷（payload）
    $segments[] = $url_safe_base64->encode(json_encode($payload, JSON_UNESCAPED_SLASHES));

    //签名 （signature）
    $msg = implode('.', $segments);
    $signature = $this->sign($msg, $key, $alg);
    if(!$signature)
        return false;

    $segments[] = $url_safe_base64->encode($signature, JSON_UNESCAPED_SLASHES);

    //组合
    return implode('.', $segments);
}
```

- 解码
```
public function decode($jwt, $key){
    //当前时间戳
    $timestamp = time();

    //拆分
    list($header64, $payload64, $signature64) = explode('.', $jwt);

    $url_safe_base64 = new UrlSafeBase64();

    //解码
    $header = $url_safe_base64->decode($header64);
    $payload = $url_safe_base64->decode($payload64);
    $signature = $url_safe_base64->decode($signature64);

    if(!isset($header['alg']))
        return $this->_error('参数错误');

    //签名是否合法
    $sign = $this->sign($header64 . '.' . $payload64, $key, $header['alg']);
    if(!$sign) return false;

    if($sign != $signature)
        return $this->_error('签名不匹配');

    //定义在什么时间之前，该jwt都是不可用的
    if (isset($payload['nbf']) && $payload['nbf'] > ($timestamp + static::$leeway))
        return $this->_error('token不可用');

    //jwt的签发时间
    if (isset($payload['iat']) && $payload['iat'] > ($timestamp + static::$leeway))
        return $this->_error('token不合法');

    //jwt的过期时间
    if (isset($payload['exp']) && $payload['exp'] < ($timestamp - static::$leeway))
        return $this->_error('token已过期');

    return $payload;
}
```


[代码](../../../../SHPhp/tree/master/system/Library/Jwt.class.php)


#### php-jwt库的使用

参考：
```
https://github.com/firebase/php-jwt
```

```
<?php

use Firebase\JWT\JWT;
use Firebase\JWT\Key;

require './vendor/autoload.php';


$key = "example_key";
$payload = array(
    "iss" => "http://example.org",
    "aud" => "http://example.com",
    "iat" => 1356999524,
    "nbf" => 1357000000
);

$jwt = JWT::encode($payload, $key, 'HS256');
$decoded = (array)JWT::decode($jwt, new Key($key, 'HS256'));

var_dump($jwt, $decoded);
```

结果打印：
```
eyJ0eXAiOiJKV1QiLCJhbGciOiJIUzI1NiJ9.eyJpc3MiOiJodHRwOi8vZXhhbXBsZS5vcmciLCJhdWQiOiJodHRwOi8vZXhhbXBsZS5jb20iLCJpYXQiOjEzNTY5OTk1MjQsIm5iZiI6MTM1NzAwMDAwMH0.gOEkQc3YCCIIjE-GxU0UTa9Lx6hQwwk5zYfO4pZQZt4" 

array(4) { 
    ["iss"]=> string(18) "http://example.org" 
    ["aud"]=> string(18) "http://example.com" 
    ["iat"]=> int(1356999524) 
    ["nbf"]=> int(1357000000) 
    }
```
