#### PHP7.4的Preloading特性

优点：在PHP 7.4以前，如果你使用了框架来开发，每次请求文件就必须加载和重新编译。预加载在框架启动时在内存中加载文件，而且在后续请求中永久有效。

缺点：性能的提升会在其他方面花费很大的代价，每次预加载的文件发生改变时，框架需要重新启动。

实现：
```
直接通过php.ini的opcache.preload=/www/sk/preload.php

[opcache]
zend_extension=opcache.so
opcache.enable=1
opcache.enable_cli=1
opcache.preload=preload.php

然后在preload.php中包含需要加载的文件（例如框架）：
$files = /* An array of files you want to preload */;

foreach ($files as $file) {

    opcache_compile_file($file);

}
```

另外基于composer的自动预加载解决方案期待有很好的支持！

参考：

https://stitcher.io/blog/preloading-in-php-74

#### PHP8

PHP8的JIT只对密集计算型的应用提升较大，而这类PHP应用很少，除了PHP-Parser或者静态分析这类类库会感觉到有明显性能提升以外，其它应用感知并不会很明显。JIT在国内热度和关注度很高，但好像PHP官方并不这么认为，PHP8官方介绍页里JIT并不是C位（https://www.php.net/releases/8.0/zh.php）。

PHP8影响最大的特性是命名参数，讨论最热烈的特性是注解。PHP8的主旋律是语法改进，而性能提升是长期工作，JIT带来的性能提升不一定有其它各种零碎的小优化加起来带来的高（长期以来，PHP每个版本都会有个位数百分比的性能提升）。对于很多PHP应用来说，PHP7.4的Preloading特性对性能提升巨大，远超JIT，这才是大部分人需要多多关注的特性

参考：

https://www.zhihu.com/question/441607128