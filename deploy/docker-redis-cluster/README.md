

### 说明

* 使用Redis版本为5.0.5-alpine,如有需要请自行更改.  
> redis官方在redis3.x和redis4.x提供了redis-trib.rb工具方便我们快速搭建集群,在redis5.x中更是可以直接使用redis-cli命令来直接完成集群的一键搭建,省去了redis-trib.rb依赖ruby环境的问题。
* docker-compose自动组网
> 使用docker-compose up启动容器后，这些容器都会被加入`{app_name}_default`网络中,但是ip不固定.所以该项目在docker-compose内使用了指定ip地址的方式,使每个容器的ip为定值.
* client端节点互通问题  
在redis配置文件中有如下配置:
``` shell
# docker虚拟网卡的ip
cluster-announce-ip 10.1.1.5
# 节点映射端口
cluster-announce-port 6379
# 总线映射端口,通常为节点映射端口前加1
cluster-announce-bus-port 16379
```
> docker虚拟网卡地址为docker与外界互通所使用的虚拟ip,如果没有上述配置,在client端操作时,会有取docker内网ip的问题.
* 关于时区,项目内默认使用了上海时区,需要使用其他时区的同学请自行更改.
```yaml
    environment:
      # 设置时区为上海
      - TZ=Asia/Shanghai
```

redis配置文件中的一些cluster相关配置项

- masterauth: 如果主节点开启了访问认证，从节点访问主节点需要认证；如果主库开启了requirepass这里就必须填相应的密码
- cluster-enabled: 是否开启集群模式，默认 no；
- cluster-config-file: 集群节点信息文件；
- cluster-node-timeout: 集群节点连接超时时间；
- cluster-announce-ip: 集群节点 IP，填写宿主机的 IP；
- cluster-announce-port: 集群节点映射端口；
- cluster-announce-bus-port: 集群节点总线端口


[//]: # (TODO 每次使用cluster都需要执行该命令，怎么办?)

redis-cli --cluster fix 127.0.0.1:6379
只需要修复master机器即可;


redis-cli --cluster check 127.0.0.1:6379
应该修改ip为局域网ip，才能查看所有节点;

cluster meet 172.17.0.1 6379

- cluster nodes
- cluster info

- 怎么验证redis主从复制? 查看在master上写操作，在slave上读操作是否成功;
- 怎么判断cluster搭建是否成功?







## ref

- [MistRay/redis-docker-compose: 使用docker-compose一键搭建Redis集群](https://github.com/MistRay/redis-docker-compose)
- 