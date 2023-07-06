# 三范式和BCNF
## 三范式

1. 字段具有原子性，不可再分。
2. 表一定要有一个主键，非主键属性一定要完全依赖于主键，而不能只依赖主键的一部分。

    可以理解为field = f(primary_key1, primary_key2...,primary_key3)，而不能仅通过一部分primary_key就能够唯一确定一个field。

    典型的就是t_student(stu_id, name, teacher_id, teacher_name) 如果想要知道一个老师叫什么名字，你只要去表里找到与他teacher_id相同的记录，而不需要管这条记录的stu_id是多少。

    本质上就是这些部分函数依赖于主属性的字段，在这个表里冗余存储了，应该摘出去。
3. 消除非主属性对主键的传递依赖。


```go
// 注：候选码 可以理解为 用于确定表中具体某一行元素，所需要用到的所有字段的集合（通常来说一张表的候选码可以有很多组）。
// 主码可以简单理解为表的主键（强调复合主键）
// t_score(student_id, course_id, score, teacher_id, leader_id) 这表里 需要用到学生id和课程id才能确定某行数据，那么student_id和course_id就是一组候选码。
```

要知道一行记录的score值是多少，就一定要同时知道student_id和course_id，所以socre完全依赖与候选码。

显然，score就是完全依赖于候选码的。

而teacher_id其实是可以有course_id来唯一确定的，即在这张表里随便找到course_id为这个值的记录，teacher_id一定是一样的。这说明老师的id在分数表里是冗余的，teacher_id部分依赖于候选码。

再看leader_id老师的上级领导id，我们一旦知道student_id和course_id，就能知道teacher_id是多少，而我们一旦知道teacher_id是多少，就可以去表里找任意teacher_id与其相同的，这些记录的leader_id一定是一样的。说明老师的领导id也是被冗余在score表里的，它传递依赖于候选码。

三范式的目的就是让一张表只记录一种实体，尽可能消除数据冗余，以减少插入、更新和删除异常的发生。

## BCNF
任意一组候选码中的任意一个字段，都不能部分或传递依赖于其他任意一组候选码。

假定 一个教师只会教一门课程，一门课可以有很多老师教且当学生选定某门课，就对应一个固定的教师
举例：学生id 课程id 教师id
那我们可以 确定两个候选键（学生id,课程id）,(学生id，教师id)
但我们发现 一个教师只会教一门课程，那么 教师id就决定了课程id 可以认为课程id依赖于教师id 那么这个表就不符合 bcnf。

# SQL优化思路
从优化成本由低到高来考虑，大概会有以下几个优化步骤：
1. 首先看sql本身的问题，这里分成三部分来考虑
    - 选择的字段，按需取数，仅查询那些业务真正需要用到的字段，避免手写sql时select *，因为db的执行器还需要把*替换成表里的全部字段，这里开销不小。
    - 看联表，优先使用带有方向的join，这样能明确把控驱动方向，由小表驱动大表，即把left join左边的想象为外层循环(驱动)，右边的内层循环(被驱动)，20w条数据里根据索引关联1条数据循环20次，与20条数据里关联1条数据循环20w次，正常情况下前者效率高，注意关联条件中的字段类型要一致，保证能走索引进行匹配。另一方面，业务上需要对关联表的数量做出严格限制，不允许超过三张表的关联查询，可以做字段冗余、内存缓存或者sql拆分在业务层组装。
    - 看查询条件，避免写出复杂的子查询、嵌套查询、CTE等骚操作，维护成本高、性能隐患大。使用explain查看执行计划，关注type、key和extra，保证所有查询子句尽可能达到range及以上的访问类型，针对使用到的各个key，关注extra中有没有利用到覆盖索引(using index)、涉及到排序的是利用了索引顺序还是内存排序(using file sort)、看能否开启并利用索引下推来优化并不完全适配索引结构的业务查询;另一方面，所创建的索引需要将区分度高的字段放前面，尽可能地复用和拓展索引，而不是新建索引，一张表的辅助索引数量不超过5个，每个复合索引的字段值不超过5个，字符串类型（尤其指基数巨大的字符串）的字段添加索引时要指定长度，控制整个索引的长度。
2. 看表设计上的问题，业务模型中的一对一、一对多、多对多关系是否使用了恰当的表设计，一对一不需要关系表，一对多、多对多则引入中间关系表。另一方面，又需要适时打破三大范式，对那些不容易变更的、或者一致性要求没那么高的字段做冗余，减少业务中表之间的关联查询。表字段类型和长度的选择上也要充分考虑业务属性和拓展性（曾经因为一个extra字段长度不够，而多次引发系统故障）。
3. 当表结构设计和sql都优化到位了，但sql执行还是很慢，那就说明是表中的数据本身规模太大了，需要采用水平和垂直分表来缩小表的大小。
    - 当单表对应实体过于复杂，属性字段过多时，可以采用垂直分表，对字段进行冷热分离存储，那些不常使用到的字段放入detail表中，常用字段留在主表，这样可以很高程度提高主表的查询效率（可以从主键索引的B+树结构来解释）
    - 当单表数据量巨大时，会导致B+树层级过高（int主键4000W，bigint2000W），影响读写效率，于是可以进行水平分表；这时候则需要结合业务场景（范围查询？排序？精确匹配？）来选出合适的分表键和分表策略，常用的有range和hash两种，前者根据某个字段的数值范围做分表，后者根据某个字段哈希之后的值取模来分表。

4. 如果业务的瓶颈在于db的网络带宽的磁盘I/O效率，那么可能需要考虑进行db分库并引入读写分离的架构。

    通过dao层路由或者中间件的方式实现读写流量分离，让主节点负责全部的写流量，而读流量由从节点分担。
    
    在读写分离模式下，mysql默认使用异步同步策略，即主节点的dump线程将binlog变更发送给从节点的I/O线程之后，并不关心从节点是否写入relaylog并执行成功，可能存在写数操作之后短暂的数据不一致。
    
    在不改变mysql的异步复制策略的情况下，如果业务只要求最终一致性，那其实mysql异步复制也能满足，如果涉及到redis缓存，则只要给缓存设置过期时间，也能保证最终一致性；
    
    如果业务要求强一致性，那只能对所有数据加分布式读写锁，如果要修改db，那就加写锁，让写锁在从节点复制最大时延之后自动过期，就能保证业务读到的是最新的数据，不过这也会大大降低业务的并发度。

    > 以下记录一下使用docker搭建mysql主从架构的过程
    - 主节点
        - .cnf文件加入配置

        ```bash
        # 开启binlog 并设置与从节点不相同的server-id
        [mysqld]
        log-bin=mysql-bin
        server-id=1
        ```
        - 通过docker启动master节点
        ```bash
        docker run --name mysql-master -v /root/mysql/my_master.cnf:/etc/mysql/conf.d/my.cnf -v /root/mysql/master:/var/lib/mysql --network mysql-m-s -e MYSQL_ROOT_PASSWORD=xxx -dp 3306:3306  mysql_5.7_vim:1.0.0
        ```
        - 为从节点创建账号
        ```sql
        CREATE USER 'testslave'@'%' IDENTIFIED BY '123456';
        GRANT REPLICATION SLAVE ON *.* TO 'testslave'@'%' WITH GRANT OPTION;
        ```
        - 查看binlog状态
        ```sql
        # 这里可以看到当前master正在写的binlog文件以及偏移量  slave需要借助这两个参数去做增量同步
        mysql> SHOW MASTER STATUS;
        +------------------+----------+--------------+------------------+-------------------+
        | File             | Position | Binlog_Do_DB | Binlog_Ignore_DB | Executed_Gtid_Set |
        +------------------+----------+--------------+------------------+-------------------+
        | mysql-bin.000004 |     2891 |              |                  |                   |
        +------------------+----------+--------------+------------------+-------------------+
        1 row in set (0.00 sec)
    
        ```
    - 从节点
         - .cnf文件加入配置
        ```bash
        # 开启binlog 并设置与主节点和其他从节点不相同的server-id
        [mysqld]
        server-id=2
        # 忽略master的系统库
        replicate-wild-ignore-table=mysql.*
        replicate-wild-ignore-table=sys.*
        ```
        - 通过docker启动slave节点
        ```bash
        docker run --name mysql-master -v /root/mysql/my_slave_1.cnf:/etc/mysql/conf.d/my.cnf -v /root/mysql/slave1:/var/lib/mysql --network mysql-m-s -e MYSQL_ROOT_PASSWORD=xxx -dp 3316:3306  mysql_5.7_vim:1.0.0
        ```
        - 配置主节点信息、开启I/O线程
        ```sql
        mysql> CHANGE MASTER TO 
            -> MASTER_HOST='mysql-master',
            -> MASTER_PORT=3306,
            -> MASTER_USER='testslave',
            -> MASTER_PASSWORD='123456',
            -> MASTER_LOG_FILE='mysql-bin.000001',   # 如果要从头开始复制 就填000001
            -> MASTER_LOG_POS=0;                     # 如果要从头开始复制 就填0

            mysql> start slave;     # 开始同步主节点
            Query OK, 0 rows affected (0.00 sec)
            mysql> show slave status;  # 查看从节点状态
        ```