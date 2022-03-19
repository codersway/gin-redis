# 查看redis.conf
function config_path() {
    source_host="127.0.0.1"
    source_port=6379
    source_pwd=123456
    server=$(redis-cli -h $source_host -p $source_port -a $source_pwd INFO)
#    echo $server
    result=$(echo $server | grep "config_file")
    echo $result
}

function source_backup() {
    source_host="127.0.0.1"
    source_port=6379
    source_pwd=123456
    # 后台备份
    redis-cli -h $source_host -p $source_port -a $source_pwd bgsave
    sleep 5s

    saveTime=$(redis-cli -h $source_host -p $source_port -a $source_pwd lastsave)
    echo $saveTime
    # osx下的
    dateStr=`date -r${saveTime} +"%Y-%m-%d %H:%M:%S"`
    echo $dateStr

    # 获取RDB文件夹
    redis-cli -h $source_host -p $source_port -a $source_pwd CONFIG GET dir
}

function restore() {
    source_host="127.0.0.1"
    source_port=6379
    source_pwd=123456

    redis-cli -h $source_host -p $source_port -a $source_pwd SHUTDOWN
    # todo 把dump.rdb复制到redis备份目录

    # 并重启redis
    redis-server
}

case $1 in
  config_path)
    config_path;;
  source_backup)
    source_backup;;
  restore)
    restore;;
esac