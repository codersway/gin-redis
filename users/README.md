# users

1. 用户登录
2. redis存储用户数据
3. 附近的人
4. 用户相册
5. 用户签到
6. 用户互相关注


# 用户登录？

redis实现用户登录

[redis实践：用户注册登录功能 - 一束光 - 博客园](https://www.cnblogs.com/Alight/p/4745746.html)


[基于Redis实现的单点登录 - 知乎](https://zhuanlan.zhihu.com/p/76542451)




# redis怎么存储用户数据？






# 附近的人

[使用 Laravel 开发简易的附近动态功能 | Laravel China 社区](https://learnku.com/articles/9273/developing-simple-nearby-dynamic-functions-using-laravel)


[Redis 应用-Geo | Laravel China 社区](https://learnku.com/articles/30821)


[稳了！用Redis实现“附近的人”功能](https://mp.weixin.qq.com/s/VUOGCcLNODp3dUFkP5eA4Q)



# 用户相册的图片



[Mysql数据量在几千万的照片表设计，请大家整理下思路 - SegmentFault 思否](https://segmentfault.com/q/1010000003906588)



## hash实现


[节约内存：Instagram的Redis实践（转） - 奋斗终生 - 博客园](https://www.cnblogs.com/ajianbeyourself/p/4475172.html)


ins没使用数据库，因为用不到数据库的修改操作，事务和关联查询等功能，所以没有使用关系型数据库；直接使用了redis；

ins用hash来存图片，具体做法是`把数据分段，每一段使用一个hash存储，由于hash结构会在单个hash不足一定数量时压缩存储，所以大量节约内存；这个“一定数量”是由配置文件里的hash-zipmap-max-entries控制的，为了解决内存，我们把该值设置为1000，性能比较好，超过1000后，HSET命令就会导致CPU消耗变大；`


```
1155315是图片id，939是用户id

HSET "mediabucket:1155" "1155315" "939"
HGET "mediabucket:1155" "1155315"

```

使用hash，我们可以实现通过图片id反查用户id；






# 用户签到

[基于Redis位图实现用户签到功能（干货，签到相关功能均可参考） - 无心码农 - 博客园](https://www.cnblogs.com/liujiduo/p/10396020.html)

[如何利用 Redis 快速实现签到统计功能 | Laravel China 社区](https://learnku.com/articles/25181)



[利用redis的bitmap实现用户签到功能 - 知乎](https://zhuanlan.zhihu.com/p/83604814)

[redis实现用户签到的功能 · Issue #2 · eqiuno/note](https://github.com/eqiuno/note/issues/2)


[Redis：Bitmaps使用场景-用户签到、统计活跃用户、用户在线状态_Redis,Bitmaps,BITOP _琦彦-CSDN博客](https://blog.csdn.net/fly910905/article/details/82629687)


# 关注功能

[使用redis实现互粉功能 - php开发经历 - SegmentFault 思否](https://segmentfault.com/a/1190000008145959)


[php - 使用Redis存储用户粉丝，该使用哪种数据类型？Hash or Sorted Set ? - SegmentFault 思否](https://segmentfault.com/q/1010000008547001)


1. 使用有序集合实现好友关系（对应的score就是关注时间戳）
2. 使用hash去记录好友的具体信息




