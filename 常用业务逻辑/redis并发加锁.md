#### 作用

- 为防止多个进程同时进行一个操作，而导致出现意想不到的结果，因此对缓存进行操作时自定义加锁

- 适合Redis单例应用场景

#### 原理

涉及到的命令：

```
#以秒为单位返回 key 的剩余过期时间
redis 127.0.0.1:6379> TTL KEY_NAME

#在指定的 key 不存在时，为 key 设置指定的值
redis 127.0.0.1:6379> SETNX KEY_NAME VALUE
```

1、使用setnx更新指定键，如果不可以直接返回false

2、键的过期时间是-1不限制的，而我们最后给它设置了时间，如果进程中途出问题那么它的过期时间
则为-1，如果当前获取到键的过期时间为-1，那么应该删除脏数据返回false以等待下一次请求

3、如果可以获得更新的权限，最后设置一个过期时间即可

#### 实现

```
//键名
$lock_key = $this->prefix . $lock_key;

//锁不住直接false
if(!$this->redis->setnx($lock_key, 1)){
    //处理设置过期时间失败的情况：直接删锁，下一个请求就正常了
    if($this->redis->ttl($lock_key) === -1){
        $this->redis->del($lock_key);
    }
    return FALSE;
}

//锁n秒，注：此时可能进程中断，导致设置过期时间失败，则ttl = -1
$this->redis->expire($lock_key, $lock_time);
return TRUE;
```

[代码](../../../../SHPhp/tree/master/system/Library/RedisLock.class.php)
