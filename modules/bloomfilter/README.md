# bloomfilter

`docker run -p 6379:6379 --name redis-redisbloom redislabs/rebloom`

```shell

# 添加一个过滤器和记录
BF.ADD newFilter foo

# 判断记录是否存在
BF.EXISTS newFilter foo

# 新建过滤器时候指定参数，如容错、数据规模
BF.RESERVE {key} {error_rate} {size}

# 指定redis服务器对布隆过滤器的默认参数，在启动命令或配置文件中添加
--loadmodule /usr/local/redis/src/rebloom.so INITIAL_SIZE 10000000 ERROR_RATE 0.0001

```


[//]: # (TODO 怎么手动安装redis布隆过滤器?)


## ref

- [redislabs/rebloom - Docker Image | Docker Hub](https://hub.docker.com/r/redislabs/rebloom#launch-redisbloom-with-docker)
- [docker redis 安装布隆过滤器插件 redBloom filter](https://blog.csdn.net/qq_35425070/article/details/107880501)
- [Redis Module 实现布隆过滤器 - 掘金](https://juejin.im/post/5dcd53346fb9a0202d6ea387)