# redis与传统关系型数据库的区别
# 常用的数据结构
- string 

    最常用的单值存储，底层的RedisObject有三种encoding方式：
    - int  
    
    当value为整数且能够被当前操作系统中C语言long类型所存储的话，那么就会按照int的格式存储（long可能是32位也可能是64位，取决于redis编译时的编译器）。当对int编码的值进行修改致使它超过long表示的范围后，encoding就会变成embstr或者raw。
    - embstr
    
    redisObject的结构分为【type、encoding、ptr】，当encoding为embstr时，ptr存储的是一个SDS结构体指针（Simple Dynamic String）简单动态字符串。
    - raw 

    当SDS的长度超过某个阈值时
    
    > redis 2.+  --- 32 字节

    > redis 3.0-4.0 --- 39 字节

    > redis 5.0 ---  44 字节

    encoding会选择raw，redisObject底层其实还是SDS，与embstr的区别在于，raw格式的object在创建时，redisObject结构体和SDS结构体是分两次申请内存的，而embstr是一次申请。

    embstr优点在于创建和销毁的时候都只需要操作内存一次，且空间上object和SDS是连续的，更有利于利用CPU缓存。

    缺点在于embstr格式的对象是只读的，即如果使用append等指令修改其value时，底层其实将encoding升级为raw，重新分配SDS空间并复制之前的值再做出追加的修改。

    SDS相较于C语言中字符串的优势：

    1. SDS结构体内部会维护当前字符串的长度len，使用strlen命令获取长度时时间复杂度仅为O(1)，而C语言中仅通过字符串末尾的'\0'来标致结束
    2. 基于1中的设计，SDS的末尾也有一个'\0'，目的是能够复用一部分string.h中的库函数
    3. 由于维护了长度len，SDS也能存储二进制数据，如图片、音频、视频、压缩文件等
    4. SDS的API非常安全，不用担心对字符串进行操作会造成内存溢出，因为SDS的API底层会校验剩余可用空间，不够时会进行扩容。

    SDS的结构体：

    在redis3.2版本之前，SDS的实现还比较简单，内有len、free、buf三个元素
    ```c
    struct sdshdr {
        //记录buf数组中已使用字节的数量
        //等于SDS所保存字符串的长度
        unsigned int len;

        //记录buf数组中未使用字节的数量
        unsigned int free;

        //char数组，用于保存字符串
        char buf[];
    };
    ```
    显然，这一版本的实现会有很严重的内存浪费，即非常短的字符串也会在头内用到两个无符号整数来存储（每个2~4个字节，取决于编译器），于是redis在3.2及以后的版本中针对不同长度的字符串，采用了不同的header：

    根据索要存储字符串的长度，分别采用1、2、4、8字节的整型来记录len和buf的总长度alloc；采用一个字节的flag的低3位存储当前header采用的字节长度规格（3bit正好记录0~7 可以对应1到8字节）。

    上面四个规格对应sdshdr8~64，还有个sdshdr5，这个结构里没有len和alloc，而是使用flag的高5位来记录len，不过这个结构只用在了key的存储中，而value的存储最低只使用了sdshdr8

    ```c
    // 注意：sdshdr5从未被使用，Redis中只是访问flags。
    /* Note: sdshdr5 is never used, we just access the flags byte directly.
    * However is here to document the layout of type 5 SDS strings. */
    struct __attribute__ ((__packed__)) sdshdr5 {
        unsigned char flags; /* 低3位存储类型, 高5位存储长度 */
        char buf[];
    };
    struct __attribute__ ((__packed__)) sdshdr8 {
        uint8_t len; /* 已使用 */
        uint8_t alloc; /* 总长度，用1字节存储 */
        unsigned char flags; /* 低3位存储类型, 高5位预留 */
        char buf[];
    };
    struct __attribute__ ((__packed__)) sdshdr16 {
        uint16_t len; /* 已使用 */
        uint16_t alloc; /* 总长度，用2字节存储 */
        unsigned char flags; /* 低3位存储类型, 高5位预留 */
        char buf[];
    };
    struct __attribute__ ((__packed__)) sdshdr32 {
        uint32_t len; /* 已使用 */
        uint32_t alloc; /* 总长度，用4字节存储 */
        unsigned char flags; /* 低3位存储类型, 高5位预留 */
        char buf[];
    };
    struct __attribute__ ((__packed__)) sdshdr64 {
        uint64_t len; /* 已使用 */
        uint64_t alloc; /* 总长度，用8字节存储 */
        unsigned char flags; /* 低3位存储类型, 高5位预留 */
        char buf[];
    };
    ```

    注意到上面struct声明里的 ``__attribute__((__packed__))``，这是在告诉编译器，取消对这个结构体的字节对齐，而是按照实际占用的字节数进行存储，原因是，sds的指针其实指向的不是结构体的起始地址，而是buf的首地址，这是为了让sds指针能直接复用string.h里的函数；

    如果不进行对齐填充，就能保证sds的指针往回退一个字节就能找到flags的8个bit，进而根据header规格找到len和alloc的起始位置；相反如果进行了填充对齐，就破坏了这个巧妙地设计。

    > By default, a single Redis string can be a maximum of 512 MB.
    - set key value [EX seconds|PX milliseconds|EXAT timestamp|PXAT milliseconds-timestamp|KEEPTTL] [NX|XX] 
    - setnx key value   在老版本redis中 setnx不支持指定过期时间，所以实现分布式锁在加锁的环节需要用lua脚本使得setnx和expire原子地执行。
        > As of Redis version 2.6.12, this command is regarded as deprecated.
        It can be replaced by SET with the NX argument when migrating or writing new code.
    - get key
    - mget key1 key2 key3 ... 批量获取多个key
    - mset key1 val1 key2 val2 ... 原子性地批量设置多个kv
    - msetnx key1 val1 key2 val2 ... 原子性地设置多个key（当且仅当key不存在时才写，不会覆盖）
    - incr key  对于存储整数类型的key，可以原子性地加1，注意如果val是个浮点数则不能用incr
    - decr key
    - incrby key intvalue  对于存储整数类型的key，原子性地增加intvalue，同样不适用于浮点数
    - decrby
    - incrbyfloat key floatval 可以将整数型的key升级成浮点数，此后只能用incrbyfloat来修改其值
- list

    线性表结构，按照插入的顺序进行排序，可以用于构建队列/栈的结构，最大存储2^32 - 1个元素，即单列表支持存储超过``40亿个数据``。

    list的数据结构有以下两个阶段：

    - 压缩列表ziplist （元素数少于512、每个元素都小于64字节)    ->    双向链表list
    
    - redis3.2之后 list底层只使用quicklist

        - ziplist
            ![ziplist](ziplist.png)
                        
            entries占用内存连续，能够更好地利用CPU高速缓存。每个entry中存储上个entry的长度以及当前节点的长度，于是可以根据当前entry的指针快速移动到上/下个entry的起点。

            相较于链表而言，它能够更高效地利用内存（链表每个节点都要维护两个指针，指针所占字节是很长的，64位系统一个指针就8字节）；但是当所存储元素数量增多后，ziplist的查询和更新操作效率都很低。

            更致命的是，entry中的prelen在前一个entry长度小于254时，采用一个字节存储，否则会一下扩大到五个字节（第一个字节固定值为0xFE，用后面四个字节存储前一个entry的字节长度），极端情况下，对ziplist中一个元素做出修改使其变长，会导致后面所有节点的prelen都发生扩展，极度影响性能。
        - quicklist
            ![quicklist](quicklist.png)

            quicklist为了减小ziplist产生的连锁更新现象的影响，进行了改造。

            quicklist可以理解为一个链表，每个节点是一个ziplist，它严格限制了每个ziplist的长度，使其即时发生连锁更新，也不会影响太严重，同时也保证了每一个ziplist的查询和更新效率。

            另外，quicklistNode中ziplist还可以采用压缩算法进行压缩，从而生成一个LZF结构，node里的zl会指向这个LZF，进一步压榨内存，提高使用率。
        - listpack
        ![listpack](listpack.png)

        listpack也是优化ziplist后的一种数据结构，它的整体结构与ziplist类似，关键在于调整了entry内的结构，依次维护
        
        encoding --- 不同的编码会占1~5个字节不等，可以理解为前几个bit位用来枚举编码，剩下的比特位用于存储content中内容的长度
        
        data listpack
        
        slen --- 指的是当前字节对于本entry起始地址的偏移量，也是encoding + content的总字节长度

- set
- zset
- hash
- hyperloglog
- bitmap
- bitfield
- stream
- geospatial

# 内存清除策略
redis是一个基于内存的k-v数据库，它的持久化策略更偏向于用于数据的恢复，而不是像innodb那样用于检索，于是当redis进程所分配的内存空间不够用时，它不能将一部分数据持久化到磁盘以腾出内存空间，所以需要引入内存清除策略。

redis中的策略完整有八种： 

1. noeviction --- 不清除策略，当新的写请求到来时，只会去执行那些能让内存占用减少的指令（例如del和其他的删除类操作），其他的则返回内存不足的错误从而拒绝写入
2. volatile-lru --- 针对那些设置了过期时间的key，使用lru算法，移除那些最后一次使用时间最远的key；least recently used 最少最近使用。
3. allkeys-lru  --- 针对所有的key，使用lru算法淘汰
4. volatile-lfu --- 针对那些设置了过期时间的key，使用lfu算法，移除那些一段时间内使用频率最低的key；least frequently used最少频率使用。
5. allkeys-lfu  --- 针对所有的key，使用lfu算法淘汰
6. volatile-ttl --- 针对那些设置了过期时间的key，把那些距离过期时间最近的key移除。
7. volatile-random --- 针对那些设置了过期时间的key，随机删除
8. allkeys-random --- 针对所有key，随机删除

如果业务缓存数据真的数量庞大以至于会占满redis内存，那就应该考虑使用redis集群架构对数据进行分片存储，尽量避免驱逐策略的使用。

实际生产中，不应当把redis当做一个可靠的持久化存储设备，而应该是一个为数据库分担压力的缓存层，以保证业务在某些key被驱逐时的正确运行。

开启方式： 配置文件中设置maxmenmory和maxmemory-policy，另外mapmemory-samples [count]是在lru和ttl的驱逐策略下，会随机去待淘汰队列选取的key数量，在这随机捞出来的数据里选最应该被淘汰的，该数值默认是5，当设置为10时基本就趋近于精确的lru和ttl，但是会消耗更多的cpu资源。 由此可见，redis内存驱逐策略中的lru和ttl都是近似的，而不是完全精确的。 

# key过期策略
redis中的key通过时间戳的方式实现过期，对过期key的删除采用的是定时随机检测和惰性删除。

惰性删除指的是，当读取一个key的时候，redis会检查其过期时间与当前系统的时间戳，如果过期则删除它。

定时随机检测
# 持久化策略

# 主从同步原理

# 哨兵模式

# 集群模式

# 分布式锁
## 与用数据库实现、zookeeper、etcd的异同

# 分布式信号量

# 红锁

# 为什么快？


LSM tree  https://zhuanlan.zhihu.com/p/181498475  
bitcask: 简洁且能快速写入的存储系统模型 https://zhuanlan.zhihu.com/p/551334186