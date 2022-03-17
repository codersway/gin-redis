# seckill

这里只实现秒杀券业务，不实现秒杀商品，因为秒杀商品就引入了商品、订单、支付、库存的锁定和释放、支付的时效性等问题。这里只聚焦秒杀系统最关键的高并发场景。

另外，秒杀券跟其他业务是隔离的，如果以后需要，也可以直接复用。 当然，这样的话，也无法直面很多秒杀商品的真正需求，对该业务缺少理解

## 参考文档

[niantianlei/second: 秒杀抢购系统，通过Rocketmq、redis、布隆过滤器、主从服务器、验证码等技术进行优化](https://github.com/niantianlei/second)

[qiurunze123/miaosha: ⭐⭐⭐⭐秒杀系统设计与实现.互联网工程师进阶与分析🙋🐓](https://github.com/qiurunze123/miaosha)

[zhongzhh8/SecKill-System: 秒杀系统，高并发，Redis，Lua脚本，Golang，Gin](https://github.com/zhongzhh8/SecKill-System)
