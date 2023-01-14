#### 简介

5中数据结构
```
String 字符串
Hash 字典
List 列表
Set 集合
Sorted Set 有序集合
```
#### 使用场景

- String 字符串

String 数据结构是简单的 key-value 类型，可以完全实现目前 Memcached 的功能

另外可以做并发锁的功能，以及计数器等功能

- Hash 字典

在 Memcached 中，我们经常将一些结构化的信息打包成 hashmap，在客户端序列化后存储为一个字符串的值（一般是 JSON 格式）。这时候在需要修改其中某一项时，通常需要将字符串（JSON）取出来，然后进行反序列化，修改某一项的值，再序列化成字符串（JSON）存储回去。Redis 的 Hash 结构可以像在数据库中 Update 一个属性一样只修改某一项属性值

- List 列表

最新消息排行榜或者消息队列，可以利用 List 的 *PUSH 操作，将任务存在 List 中，然后工作线程再用 POP 操作将任务取出进行执行

- Set 集合

Set 就是一个集合，集合的概念就是一堆不重复值的组合

比如在微博应用中，可以将一个用户所有的关注人存在一个集合中，将其所有粉丝存在一个集合

集合提供了求交集、并集、差集等操作，那么就可以非常方便的实现如共同关注、共同喜好、二度好友等功能，对上面的所有集合操作，你还可以使用不同的命令选择将结果返回给客户端还是存集到一个新的集合中

- Sorted Set 有序集合

和Sets相比，Sorted Sets是将 Set 中的元素增加了一个权重参数 score，使得集合中的元素能够按 score 进行有序排列

比如一个存储全班同学成绩的 Sorted Sets，其集合 value 可以是同学的学号，而 score 就可以是其考试得分

也可用于游戏中用户得分排行榜等

```
成员的位置按 score 值递增(从小到大)来排序，具有相同 score 值的成员按字典序(lexicographical order )来排列
ZRANGE key start stop [WITHSCORES]

成员的位置按 score 值递减(从大到小)来排列， 具有相同 score 值的成员按字典序的逆序(reverse lexicographical order)排列
ZREVRANGE key start stop [WITHSCORES]
```

- 订阅-发布

Pub/Sub 从字面上理解就是发布（Publish）与订阅（Subscribe），在 Redis 中，你可以设定对某一个 key 值进行消息发布及消息订阅，当一个 key 值上进行了消息发布后，所有订阅它的客户端都会收到相应的消息。这一功能最明显的用法就是用作实时消息系统，比如普通的即时聊天，群聊等功能


- Setbit 位图

Setbit 命令用于对 key 所储存的字符串值，设置或清除指定偏移量上的位(bit)

参考：

http://redisdoc.com/bitmap/setbit.html

key 所储存的字符串值，设置或清除指定偏移量上的位(bit)。位的设置或清除取决于 value 参数，可以是 0 也可以是 1 。当 key 不存在时，自动生成一个新的字符串值。字符串会进行伸展(grown)以确保它可以将 value 保存在指定的偏移量上，当字符串值进行伸展时，空白位置以 0 填充

可用于用户签到、用户上线次数统计等功能

例如这个月上线了签到功能，第一天用户做了签到的动作，则执行命令 SETBIT sign_1 1 1，如果明天也继续签到，那么执行命令 SETBIT sign_1 2 1 ，第16天签到则执行命令 SETBIT sign_1 16 1，以此类推

使用命令BITCOUNT key [start] [end]计算给定字符串中，被设置为 1 的比特位的数量