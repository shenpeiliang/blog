#### 攻击原理及过程

1. 用户C打开浏览器，访问受信任网站A，输入用户名和密码请求登录网站A；

2. 在用户信息通过验证后，网站A产生Cookie信息并返回给浏览器，此时用户登录网站A成功，可以正常发送请求到网站A；

3. 用户未退出网站A之前，在同一浏览器中，打开一个TAB页访问网站B；

4. 网站B接收到用户请求后，返回一些攻击性代码，并发出一个请求要求访问第三方站点A；

5. 浏览器在接收到这些攻击性代码后，根据网站B的请求，在用户不知情的情况下携带Cookie信息，向网站A发出请求。网站A并不知道该请求其实是由B发起的，所以会根据用户C的Cookie信息以C的权限处理该请求，导致来自网站B的恶意代码被执行

#### 防守策略

> 验证 HTTP Referer 字段

当用户通过黑客的网站发送请求到后台时，该请求的 Referer 是指向黑客自己的网站，检查域名是否合法，但 Referer 的值是由浏览器提供的，因此也可以被修改而变无效

> 在请求地址中或表单中添加 token 并验证

可以在 HTTP 请求中以参数的形式加入一个随机产生的 token 并保存到 Session，并在服务器端建立一个拦截器来验证这个 token，如果请求中没有 token 或者 token 内容不正确，则认为可能是 CSRF 攻击而拒绝该请求

GET 请求将被追加一个参数https://www.baidu.com?csrf=token

POST 则需要在 form 表单中添加隐藏域 <input type="hidden" name="csrf" value="token"/>

这种方式依然存在被黑客拿到token的安全性

> 在 HTTP 头中自定义属性并验证

和上一种方法不同的是，这里并不是把 token 以参数的形式置于 HTTP 请求之中，而是把它放到 HTTP 头中自定义的属性里

通过 XMLHttpRequest 这个类，可以一次性给所有该类请求加上 csrftoken 这个 HTTP 头属性，并把 token 值放入其中

但这种方式只能用于Ajax方法中，普通的页面就行不通了


> 根据业务的重要性做特殊处理

既然都拦截不了，那针对业务来做处理就相对简单得多，例如银行转账等操作可以使用手机验证码的方式进行拦截

#### XSS跨站脚本攻击

主要通过参数过滤，让其变成不可执行的代码，不信任任何用户的数据，严格区分数据和代码

例如前端直接提交过来的值：
```
<script>alert('sss')</script>
```

通常后端要进行：
```
//接收
$name = htmlspecialchars($_POST['name']);

//输出
echo htmlspecialchars_decode($name);
```

1. 直接过滤所有的JavaScript脚本；

2. 转义Html元字符，使用htmlentities、htmlspecialchars等函数；

3. 引用其他类库：HTMLPurifier



