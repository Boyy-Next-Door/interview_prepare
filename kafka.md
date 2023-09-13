# kafka & rocketmq
- 前者像是水渠，流量大，流速慢，强调更高的吞吐量
- 后者像是高压水枪，吞吐量小但是流速快，强调瞬时的并发处理能力

# kafka为什么快
1. 磁盘顺序读写， 采用append only的方式向partition里追加，只写不删，顺序读写磁盘
2. 页缓存 kafka越过JVM，直接通过操作系统的PageCache，当需要读写磁盘上的文件时，以pageCache作为磁盘的缓存。越过JVM是为了避免它的GC机制带来的开销；另一方面所有In-Process Cache其实都会在OS中有一份一样的PageCache，越过JVM直接使用PageCache至少可以翻倍利用缓存空间。
3. 零拷贝 kafka运行在linux系统上，一般来说当一个进程需要从磁盘读取文件并通过socket发送出去时，需要先将数据读到内核态的pageCache，然后再拷贝到用户态的用户进程的内存空间中，之后再由用户线程将数据写入到内核态的socket缓冲区进行发送；kafka使用了sendfile技术，直接在内核态读取文件到PageCache，然后直接发送到socket缓冲区，避免了拷贝到用户内存这一步，大大提升了I/O效率。
4. 批量处理 提供大量的批处理API，能够基于对数据的压缩合并，以更小的数据包、更快的传输速度对数据进行处理。

# 基础
## kafka简介
- 分布式消息系统 -> 分布式流处理平台
- 吞吐量高、性能好
- 伸缩性好，支持在线水平扩展
- 容错性、可靠性
- 与大数据生态紧密结合，如hadoop、spark等（可以无压力处理TB级别的数据）

## 消息模型
- JMS规范 --- Java Message Service API （Java消息服务）
- AMQP  Advanced Message Queueing Protocol 高级消息队列协议
    - 模型
        - 队列 queues
        - 信箱 exchanges
        - 绑定 bindings
    - 支持事务，数据一致性高，不强调性能，多用于银行、金融行业
    - RabbitMQ、Spring AMQP、Spring JMS
- MQTT Message Queueing Telementry Transport 是一种轻量级的发布/订阅消息传输协议
    - 模型 
        - Client
        - Broker
        - Topic
    - 非常轻量级，消息头只有固定2个字节，适用于带宽和处理能力有限的嵌入式设备上，广泛应用于物联网中。

## 基本概念
  ![kafka集群架构](images/kafka_cluster.png)
- Brocker 消息代理
- Topic主题
- Partition分区
- Replication副本
- Segment段
- Producer生产者
- Consumer消费者
- Offset偏移量
- Consumer Group消费者组


    可以理解为关系型数据库中的表，从整体上来看，kafka集群内可以有很多个topic，这是个半结构化的容器，可以将相同/不同类型的消息塞到一个topic里去。

    具体地，在kafka中的每一个topic都可以进行分区，即分散成多个partition存在于集群里的不同brocker上。

    对于每一个partition（可以理解为一个topic的一部分），kafka只会往其尾部写数据，而不会修改或删除之前的数据，每一条消息都会有一个自增的、partition内不重复、partition间可能重复的编号，称之为offset，可用于快速定位某条消息在这个partition中的位置。

    采用partition进行数据分片能够一定程度上提高整体数据的可用性，并让系统能够横向拓展，但是也许需要引入replications来做数据冗余；在kafka里每个partition的备份会分散到其他brocker上，并且会带有leader或follower的标签，集群内对这个partition的读写请求都只会打到leader上，follower只是单纯地同步数据；kafka还会为每这个partition维护一个In-Sync Replica List(ISR)来记录当前同步状况良好的分区集合。

    在每个broker里会有一个特殊的topic，用于存储各个consumer对当前broker上topic的读取offset。（在老版本的kafka中这个信息是存在zookeeper里的）
    ```bash
    [root@VM-141-14-centos ~/kafka/kafka_2.12-2.8.1/bin]# ./kafka-topics.sh --zookeeper localhost:2181 --describe --topic test
    Topic: test     TopicId: 8i9nK73iS6Coa40mO1fF2A PartitionCount: 3       ReplicationFactor: 2    Configs: 
            Topic: test     Partition: 0    Leader: 2       Replicas: 2,1   Isr: 2,1
            Topic: test     Partition: 1    Leader: 0       Replicas: 0,2   Isr: 0,2
            Topic: test     Partition: 2    Leader: 1       Replicas: 1,0   Isr: 1,0
    ```
## 消息队列的作用

消息队列这种中间件的核心作用有三：
1. 作为异步业务模型中的事件总线
2. 作为各个服务间通信的管道，实现服务间解耦
3. 作为并发业务的缓冲区，达到削峰填谷的效果

> 任何技术架构都不是完美的，消息队列在引入上述优点之后也会带来其他缺点：
1. 系统可用性降低，在整个系统关键路径上每多引入一个中间件，就会带来额外的出现错误的可能性
2. 系统复杂度增加，原本串行的业务逻辑，在经过消息队列解耦之后，需要加入非常多通信逻辑，也会引给业务入很多中间状态，势必增加系统的复杂程度
3. 数据一致性，引入MQ之后，消息的顺序、是否丢失、重复消费等问题，都需要额外考虑，否则就可能出现系统多个服务之间的数据不一致

> 总而言之，需要具体考虑业务模型、系统并发程度、性能瓶颈以及对数据一致性的要求程度，再来考虑是否要引入消息队列。
## 消息模型

最常见的消息模型有两种
1. 点对点模型

    即MQ相当于一个先进先出的队列，生产者不断生产并将消息加入队列中，而多个消费者一次从队列中取出一条消息，一条消息至多被一个消费者所消费，否则就一直停留在队列中或者直到它过期被移除队列。
2. 发布-订阅模型

    该模型引入topic的概念，多个生产者向topic内写入消息，而这个topic可以被不同的消费者所订阅，不同消费者之间的消费进度互不影响，一条消息可以反复被很多个消费者所消费，而不会在第一次被消费后就移出topic

## 几种消息队列对比
- ActiveMQ

    目前已经被淘汰，单机吞吐量在每秒万级，基于主从架构实现高可用。
- RabbitMQ

    基于Erlang开发，语言天生支持高并发，性能极高，延低至微秒级，其单机吞吐量在每秒万级；另一方面其语言门槛高，几乎很少有源码级别的定制二开，如果业务对并发量要求不高，而对时延要求极高，可以考虑RabbitMQ
- RocketMQ

    阿里出品，JAVA系生态首选，吞吐量在单机每秒十万级，基于分布式架构，支持无限动态扩容，k8s友好。

    它有着轻量级、高拓展、高性能的丰富API，有着金融级的稳定性，更是经过阿里系众多产品实践打磨后的通用解决方案。
- Kafka

    kafka目前的定位已经是一个高性能的数据流式处理平台，它用最简单的发布、订阅功能，支持多种灵活多变的消息消费模式。

    它天生分布式的设计加上数据备份机制，在提供单机百万级吞吐量的同时，能够灵活拓展、且保证数据不丢失。

    其丰富的流式处理API，能够非常好的适配大数据处理生态（spark、hadoop等）。

    另外kafka在2.8版本之前，严重依赖于zookeeper，后引入了基于Raft协议的KRaft模式，进而抛去对zk的依赖。大大简化了系统架构。

## kafka的主要应用场景
1. 消息队列
2. 数据处理

## kafka的优势
1. 极致的性能：

    kafka基于scala和java语言开发，设计中大量使用批处理和异步的思想，仅单机就支持每秒百万级别以上的消息吞吐量
2. 生态极其优秀：

    kafka与周边生态系统的兼容性是最好的，尤其是在大数据和流式计算领域

## kafka多副本机制

kafaka中，一个topic会被分成多个partition，并分布在不用的broker上，在此基础之上又为各个partition创建若干个副本并分布在不同broker上，如此提供了极高的容灾能力。

对于一个partition的多个副本来说，只有一个leader，其他的都是follower，其中leader负责对这个partition（逻辑上的）的读写，而其他的partition仅仅作为备份不提供读写功能。

当leader分区挂了之后，kafka会从那些follower中挑选一个作为新的leader（对follower的数据同步程度有要求）。

> topic分区、分区副本的设计有什么好处？
1. 将一个topic进行分片，并分散到集群中的不同broker上，天然地将读写流量在集群上负载均衡，也即提高了集群的并发能力
2. 分片副本极大地保证了kafka中数据的安全性、集群的容灾能力，不过这也会带来额外的存储和网络带宽开销。

## kafka如何保证消息的顺序消费？

对于一个topic，kafka会将它划分成多个partition，每当生产者向topic中写入一条消息，这条消息只会被追加到某一个partition的尾部；kafka只保证每一个partition内部的消息是按照时间顺序追加到尾部的，但不能保证它们在时间上是连续到达的，也不能保证不同partition中消息追加写入的顺序。

于是一般有两个方法来保证顺序消费：
1. 修改kafka的配置，每个topic的partition数量设置成1，但是这样违背了其分片的设计。
2. 在通过api写入一条消息时，有四个参数可以指定：topic、partition、key、data；生产者可以在生产环节，将那些需要顺序消费的消息，指定partition；或者不指定partition，但是传入相同的key，kafka会保证key相同的消息总是被写到同一个partition的尾部。

## kafka如何保证消息不丢失

在kafka中一条消息大致会经过： 生产逻辑 -> kafka集群 -> 消费逻辑  三个环节

- 在生产逻辑中，程序调用kafka的client，写入数据
- client通过网络，调用kafka集群的API将数据写入对应的leader分区，再之后将新数据同步到follower分区
- 在消费逻辑中，程序调用api拿到一条数据，然后进行业务处理，当业务处理完毕之后才能确定该条消息是否已经成功消费

在以上三个环节中，都有可能出现消息丢失：
1. 生产者业务调用client新建消息，但是因为网络原因调用失败，消息没能到达kafka集群
2. 生产者client的请求成功到达kafka集群，并写入leader分区，但还没同步到任何一个follower分区，随后leader分区所在的broker宕机，某个follower替补成为新的leader，此时消息丢失。
3. 消费者读取一条消息后，业务处理出现异常并回滚，此时消息并没能被正确地消费，但不会有其他消费者再次对其进行处理。
---
在kafka的通信模型中，支持三种不同的语义，基于生产者client的ACK机制，以及消费者端手动提交offset所实现：
1. 至少一次
    
    一条消息一旦被生产者成功生成，那么它「至少会被」一个消费者所消费。

    具体地，需要将Producer的ACK_CONFIG设置为-1，保证一条消息写入到所有分片副本中。

    考虑ACK失败时的重试，只要成功了，那么该条消息至少会被消费一次；如果写入kafka成功，但是ACK在网络中丢失，客户端进行重试，则可能导致消息重复写入多次，进而被消费多次。
2. 至多一次

    一条消息一旦被生产者成功生成，那么它「至多只被」一个消费者所消费。

    ACK_CONFIG设置为0，即在Producer中不关心是否成功写到kafka中，即使成功了，最多只会被消费一次。
3. 精确一次

    一条消息一旦被生产者成功生成，那么它「只被」一个消费者所消费。

    精确一次是基于「至少一次」语义，加上生产者->broker、brocker->Consumer、Consumer业务逻辑中的消息幂等处理实现的。

对于Producer，支持一个定制参数ProducerConfig.ACKS_CONFIG

- 值为0：不关心broker的响应，即只要生产者程序调用了client的新增消息接口，就返回成功，不管这则消息是否成功写入到kafka中。
- 值为1：写入的消息到达kafka集群后，并成功写入leader分区的本地日志，就会返回成功。
- 值为-1或者all： 写入的消息到达kafka集群后，需要将它完整写到所有分区副本中，才会返回成功。

对于kafka集群本身而言，给每个分区配置备份因子replication-factors，加上leader和follower的选举机制，保证数据本身在集群内部不丢失。

对于Consumer，关闭消息offset自动提交功能enable.auto.commit = false；避免异步处理消息，而要在消费逻辑确保成功之后，再向broker提交该消息的offset。

## kafka如何保证消息幂等性

幂等性指的是一个操作，被执行多次，所造成的影响与只执行一次的形象完全一样。

站在Producer ->  broker这一环节看，「至多一次」的语义在拿不到ack的情况下会进行重试，如果是一个为true的ack因网络而丢失，那么后续的重试可能导致这条消息被两次写入到分区的本地日志中，进而导致后续消费多次。

于是可以在Producer中开启enable.idempotence，Producer会给每条消息添加一些独立的标识字段，于是对于带有相同标识的消息，broker至多只会将它写入日志一次。

在broker->consumer这个环节，kafka也支持具有精确一次语义的流式处理API，可以通过processing.guarantee = exact_one来开启。

另外，对于Consumer业务逻辑，也需要对消息本身进行幂等判断，例如给每一个消息加上唯一id，处理之后将该id写入缓存，每次消费前查询缓存看本消息是否已经被处理过。