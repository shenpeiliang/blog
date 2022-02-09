#### 简介

- php存储session的方式默认使用的是文件，在很多情况下我们需要改变session的存储方式为redis

- 可以通过修改php.ini文件来改变session的存储方式，但一般情况下没有在服务器修改的权限

```
session.save_handler = redis
session.save_path = "tcp://127.0.0.1:6379"
或者使用密码
session.save_path = "tcp://127.0.0.1:6379?auth=123456"
```

- 通过PHP提供的session_set_save_handler()函数来重写session

#### 原理

PHP 5 > 5.4和php 7 支持SessionHandlerInterface 接口，我们可以实现这个接口中的所有方法，然后通过session_set_save_handler()函数来使方法生效

```
interface SessionHandlerInterface {
    /* Methods */
    public close(): bool
    public destroy(string $id): bool
    public gc(int $max_lifetime): int|false
    public open(string $path, string $name): bool
    public read(string $id): string|false
    public write(string $id, string $data): bool
}
```

#### 实现

> open函数

```
public open(string $path, string $name): bool
```

当session_start()函数被调用的时候该函数被触发

$path对应php.ini中的session.save_path选项值

默认情况下，php.ini中session.save_path这个选项是被注释的，所以$path的值为空

举个例子：session.save_path设置为/tmp ，则$save_path的值为/tmp

$name对应php.ini中的 session.name选项值

默认情况下session.name设置为PHPSESSID，所以说$name参数值为PHPSESSID

```
public function open($save_path, $name): bool
{
    //创建连接
    $handle = new \Core\Cache\Driver\RedisHandler();
    $this->driver = $handle->get_handler();

    //是否指定数据库
    if ($this->db)
        $this->driver->select($this->db);

    return true;
}
```

> close函数

```
public close(): bool
```

关闭当前的session，当session关闭的时候该函数自动被触发，或者在程序中调用session_write_close()函数是触发close()函数

```
public function close(): bool
{
    return true;
}
```

> read函数

```
public read(string $id): string|false
```

当session_start()函数被调用的时候先触发open函数，再触发该函数

$id由客户端传过来的sessionId

```
public function read($session_id)
{
    $key = $this->prefix . $session_id;

    //读取当前sessionid下的data数据
    $res = $this->driver->get($key . '.data');

    //读取完成以后 更新时间，说明已经操作过session
    $this->driver->set($key, 'last_time', time());

    return $res;

}
```

> write函数

```
public write(string $id, string $data): bool
```
将session的数据写入到session的存储空间内，当session准备好存储和关闭的时候调用该函数

```
public function write($session_id, $session_data): bool
{
    $key = $this->prefix . $session_id;

    return $this->driver->save($key, ['last_time' => time(), 'data' => $session_data]);

}
```

> destroy函数

```
public destroy(string $id): bool
```

当函数session_destroy()调用的时候触发该函数，这时我们可以在该函数中将$session_id对应的数据销毁掉

```
public function destroy($session_id): bool
{
    $key = $this->prefix . $session_id;

    return $this->driver->delete($key);
}
```

> gc函数

```
public gc(int $max_lifetime): int|false
```
清除垃圾session，也就是清除过期的session

$max_lifetime参数的值就是session.gc_lifetime选项所设置的值

这个函数是否被触发要取决于session.gc_divisor和session.gc_probability这两个选项。该函数被触发的概率为 session.gc_probability/session.gc_divisor。如果probability设置为1，divisor设置为100。那么gc函数被触发的概率就是1%。也就是说在100个请求中可能会在某一个请求过程中触发这个函数。从这里我们可以知道，如果客户端一直没有请求，那这个函数就永远不会被触发。即使有些session信息没被操作的时间已经超过了session.gc_lifetime所设置的时间

```
public function gc($maxlifetime)
{
    /*
        * 取出所有的 带有指定前缀的键
        */
    $keys = $this->driver->keys($this->prefix . '*');

    $now = time(); //取得现在的时间
    foreach ($keys as $key) {
        //取得当前key的最后更新时间
        $last_time = $this->driver->get($key, 'last_time');
        /*
            * 查看当前时间和最后的更新时间的时间差是否超过最大生命周期
            */
        if (($now - $last_time) > $maxlifetime) {
            //超过了最大生命周期时间 则删除该key
            $this->driver->delete($key);
        }

    }

}
```

参考代码：

[代码1](../../../../SHPhp/tree/master/system/Core/Session/Driver/RedisHandler.class.php)