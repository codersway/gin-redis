function get_aof_config() {

  lk | redis-cli -h 127.0.0.1 -p 6379 -a 123456
}

#1. appendonly no：是否开启 AOF 持久化功能
#2. appendfilename "appendonly.aof"：AOF 文件的名称
#3. dir ./：RDB 文件和 AOF 文件所在目录
#4. appendfsync everysec：fsync 持久化策略
#5. no-appendfsync-on-rewrite no：重写 AOF 文件期间是否禁止 fsync 操作。如果开启该选项，可以减轻文件重写时 CPU 和磁盘的负载（尤其是磁盘），但是可能会丢失 AOF 重写期间的数据，需要在负载和安全性之间进行平衡
#6. auto-aof-rewrite-percentage 100：AOF 文件重写触发条件之一
#7. auto-aof-rewrite-min-size 64mb：AOF 文件重写触发条件之一
#8. aof-load-truncated yes：如果 AOF 文件结尾损坏，Redis 服务器在启动时是否仍载入 AOF 文件
function lk() {
    printf "%s\n%s\n%s\n%s\n%s\n%s\n%s\n%s\n" "config get appendonly" "config get appendfilename" "config get dir" "config get appendfsync" "config get no-appendfsync-on-rewrite" "config get auto-aof-rewrite-percentage" "config get auto-aof-rewrite-min-size" "config get aof-load-truncated"
}

function test() {
    redis-cli -h 127.0.0.1 -p 6379 -a 123456 -n 1 keys  '*' |
while read key
do
    key_val=`redis-cli -h 127.0.0.1 -p 6379 -a 123456 -n 1 get ${key}`
    echo ${key}  ${key_val}
done
}

get_aof_config

#test