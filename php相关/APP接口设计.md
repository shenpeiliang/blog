#### 需求

- 需要知道哪个版本的客户端使用了哪个版本的接口访问了哪里，为了方便线下调试通常需要保存到系统日志

- 访问的参数只能在规定的时间内有效，不符合条件则拒绝访问，通常需要携带当前的本地时间戳（如果本地时间和服务器的时间不一致，或手动修改，我们可以每次返回服务器的当前时间作为依据）

- 最低接口版本控制


#### 原理

> 必须字段
```
client_type 客户端类型 例如"android", "ios"
client_version 客户端版本号 例如 10.2.3
uuid 移动设备唯一标识符  例如安卓手机的IMEI，苹果手机的UUID
api_version 三位数字的版本号 1.2.10
timestamp 请求时间戳
signature 签名
```

> 检查客户端接口版本号是否是有效
```
//最低版本
if (version_compare($this->api_version, $this->allow_client_version[$this->client_type]['minimum_api_version'], '<'))
    $this->_fail('ERROR_VERSION', '当前版本过低，为确保正常使用，请您升级~');
//最高版本
if (version_compare($this->api_version, $this->allow_client_version[$this->client_type]['latest_api_version'], '>'))
    $this->_fail('ERROR_VERSION', '当前版本暂未发布，请耐心等待~');
```

> 签名规则
```
使用client_type、client_version、uuid、api_version、timestamp参数进行签名

$params = [
    "client_type" => $this->client_version,
    "client_version" => $this->client_version,
    "uuid" => $this->uuid,
    "api_version" => $this->api_version,
    "timestamp" => $timestamp
];
ksort($params);  //按数组的键排序
$sign = ''; //需要签名加密组合的字符串
foreach ($params as $key => $val) {
    $sign .= $key . '=' . $val;
}
//sha1加密
$sign = sha1($sign);
```

认证成功之后就可以进行正常的逻辑了，如果我们在控制器中定义了属性$_filters，例如login必须登录，则需要判断是否登录

用户登录成功后应该返回给客户端，需要验证登录的时候需要提交
```
uid 用户编号
token 用户令牌
```

#### 应用

模拟请求
```
curl -d 'client_type=ios&client_version=1.0.1&uuid=13jdds333ljj3hhhjl323&api_version=1.1.2&timestamp=1642472443&signature=233232jjjj' http://shphp.cn/article/lists
```

返回参考

```
{"code":"SIGN_MISMATCH","msg":"签名不匹配","data":[]}
```

- [代码](/harry_shen/SHPhp/blob/master/app/Controller/Appserver.class.php)
