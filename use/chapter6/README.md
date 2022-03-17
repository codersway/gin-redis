# 前缀匹配

用redis实现前缀匹配，我们一般是用zset来实现，因为我们可以通过zset的score实现对于匹配出数据的排序

1. 前缀补全，输入ab，返回abc、abks
2. 随机补全，输入m、p，返回pskm、aspqbmc



自动补全联系人和通讯录

## 方法1（最佳实践）

[Redis In Action 笔记（六）：使用 Redis 作为应用程序组件 | Redis 技术论坛](https://learnku.com/articles/30573)

[百度面试题-字符串前缀匹配（二分法） - neyer - 博客园](https://www.cnblogs.com/neyer/p/4518624.html)

[Auto Complete with Redis](http://oldblog.antirez.com/post/autocomplete-with-redis.html)

[redis实现自动补全(包含前缀补全和自动补全) - tinysakura的博客 - CSDN博客](https://blog.csdn.net/m0_37556444/article/details/82709764)


## 方法2

直接用zscan也可以实现，但是性能相对会慢

[Redis ZSCAN 命令 - Redis 基础教程 - 简单教程，简单编程](https://www.twle.cn/l/yufei/redis/redis-basic-sorted-sets-zscan.html)



