同步redis数据到mysql 方案

使用canal解析binlog的方案，实现redis和mysql的数据一致性;

Canal主要是基于数据库的日志解析，获取增量变更进行同步，由此衍生出了增量订阅&消费的业务，核心基本就是模拟MySql中Slave节点请求。

