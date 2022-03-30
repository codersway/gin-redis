# lua

1. EVAL
2. EVALSHA
3. SCRIPT_LOAD
4. SCRIPT_EXISTS
5. SCRIPT_FLUSH
6. SCRIPT_KILL

---


[Lua在Redis的应用 - 后端 - 掘金](https://juejin.im/entry/5bee2e116fb9a049a5707a2a)

*普通的redis操作不是原子性的，想要实现redis的原子性操作，只有使用`mset`, `msetnx`之类的原子性api，或者使用cas锁，或者lua脚本来解决。*




[redis中使用lua脚本代替redis锁，提高操作原子性 - checkboxMan的个人空间 - OSCHINA](https://my.oschina.net/forever9999/blog/2218939)
[如何使用Redis执行Lua脚本 - 云+社区 - 腾讯云](https://cloud.tencent.com/developer/article/1414871)

[Redis 使用 Lua 脚本替代 SETNX / DECR 保证原子性 | Laravel China 社区](https://learnku.com/articles/39265)

[新姿势！Redis中调用Lua脚本以实现原子性操作 - 掘金](https://juejin.im/post/5cb97846f265da036023ac82)

---

[Redis Lua 脚本调试器用法说明 — blog.huangz.me](https://blog.huangz.me/2017/redis-lua-debuger-introduction.html)



