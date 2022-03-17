# Dockerfile
FROM redis:latest
# 注意redis.conf路径
COPY redis.conf /usr/local/etc/redis/redis.conf
EXPOSE 6379
ENTRYPOINT [ "redis-server", "/usr/local/etc/redis/redis.conf"]

