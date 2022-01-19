#### 需求

- 需要知道哪个版本的客户端使用了哪个版本的接口访问了哪里，为了方便线下调试通常需要保存到系统日志

- 访问的参数只能在规定的时间内有效，不符合条件则拒绝访问

- 来源及IP限制


#### 原理

> 必须字段
```
data 传递的json数据
timestamp 请求时间戳
source 来源编号
```

> IP拦截
```
//ip地址
$ip = get_client_ip();

//ip白名单
foreach ($this->_config['allowed_ips'] as $pattern) {
    if (!fnmatch($pattern, $ip))
        $this->_fail('IP_NOT_ALLOWED', '访问不允许');
}
```

> 签名规则
```
使用data、timestamp、source、authkey（加密秘钥）参数进行签名

// 检查签名
$authkey = $this->_config['authkey'];

//需要签名加密的字段值
$args = compact('data', 'timestamp', 'source', 'authkey');

ksort($args);  //按数组的键排序
$signature = ''; //需要签名加密组合的字符串
foreach ($args as $key => $val) {
    $signature .= $key . '=' . $val;
}
//sha1加密
$signature = sha1($signature);
```

认证成功之后解析数据
```
// 解码请求数据
$this->_data = json_decode($data, true);

//json_decode出错时返回null或false
if (!$this->_data)
    $this->_fail('DATA_INVALID', '数据不合法');
```

#### 应用

模拟请求
```
curl -d 'data={uid:200}&timestamp=1642472443&source=22&signature=233232jjjj' http://shphp.cn/article/lists
```

返回参考

```
{"code":"SIGN_MISMATCH","msg":"签名不匹配","data":[]}
```

[代码](../../../../SHPhp/tree/master/app/Controller/Api.class.php)
