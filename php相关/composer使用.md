#### 安装

```
全局方式安装
# 下载composer.phar 
curl -sS https://getcomposer.org/installer | php

# 把composer.phar移动到环境下让其变成可执行 
mv composer.phar /usr/local/bin/composer

# 测试
composer -V 

使用国内镜像
composer config -g repo.packagist composer https://packagist.phpcomposer.com

查看配置信息：
composer config -g -l

```

#### 使用

```
配置文件
composer.json

安装扩展命令：
composer install

```

以thinkphp5框架为例：

```
composer create-project topthink/think tp
```

参考：

https://packagist.org/packages/topthink/framework#v5.0.9

https://pkg.xyz/#how-to-use-packagist-mirror

#### 其他

插件搜索：

https://packagist.org/

```
版本符：
~和^的意思很接近，在x.y的情况下是一样的都是代表x.y <= 版本号 < (x+1).0，但是在版本号是x.y.z的情况下有区别，举个例子吧：

~1.2.3 代表 1.2.3 <= 版本号 < 1.3.0
^1.2.3 代表 1.2.3 <= 版本号 < 2.0.0

```

参考：

https://getcomposer.org/doc/articles/versions.md