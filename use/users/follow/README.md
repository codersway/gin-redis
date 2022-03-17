用户关注

# 方案


[关于用户关注粉丝表设计方案的思考 - 四魂之域 - SegmentFault 思否](https://segmentfault.com/a/1190000012639665)


```

关注关系表；（一个关注一条记录；）
用json记录每个用户的的关注者和粉丝；（一个用户一条记录；）
redis的hash；（每个用户一个key来记录关注者列表，一个key来记录粉丝列表；）

```

成熟方案是，用做一个json+redis的同步双写；



# db实现





[50 万用户的应用，怎么设计关注粉丝数据表？ | Laravel China 社区](https://learnku.com/laravel/t/31872)

[实现微博关注关系的方案分析 | 技术鹏飞](https://www.itlipeng.cn/2017/04/06/%E5%AE%9E%E7%8E%B0%E5%BE%AE%E5%8D%9A%E5%85%B3%E6%B3%A8%E5%85%B3%E7%B3%BB%E7%9A%84%E6%96%B9%E6%A1%88%E5%88%86%E6%9E%90/)





## redis实现关注

[Redis实现用户关注功能 - PHP大菜鸡 - 博客园](https://www.cnblogs.com/caiji/p/8395185.html)



[使用redis实现互粉功能 - php开发经历 - SegmentFault 思否](https://segmentfault.com/a/1190000008145959)


[php - 使用Redis存储用户粉丝，该使用哪种数据类型？Hash or Sorted Set ? - SegmentFault 思否](https://segmentfault.com/q/1010000008547001)

[使用Redis实现关注关系（干货；） - 简书](https://www.jianshu.com/p/7b70860d33bf)



```
使用有序集合实现好友关系（对应的score就是关注时间戳）；
使用hash去记录好友的具体信息；
```


