# 基础八股

## make 和 new 的区别

- 同： 都会给变量分配内存
- 异：
  - 适用的类型不同：
    - make 仅仅可用于初始化 map、channel、slice(此时必须指定大小)
    - new 用于为任意 Type 分配内存，返回一个该 Type 的零值对象的指针。注意到 map、channel、slice 属于引用类型，如果使用 new，返回的会是 nil(引用类型的零值)，是不能直接赋值的。
  - make 返回所构建数据结构对象本身，而 new 会返回所构建类型零值对象的指针。
  - new 分配内存后会设置为零值，而 make 分配空间后会进行初始化。
    ![零值](zero-value.png)

## 数组和切片的区别

- 同：
  - 只能存储一组相同数据类型的数据结构
  - 都是通过下标来访问，有容量长度，长度通过 len 获取，容量通过 cap 获取
- 异：
  - 数组定长，访问和复制不能超过数组定义的长度，否则就会下标越界，切片长度和容量可以自动扩容。
  - 数组是值类型，而切片是引用类型，每一个切片都引用了一个底层数组。修改切片的时候，改的是底层数组中的数据，切片一旦扩容，就会指向一个新的底层数组，内存地址也会随之改变。（slice = append(slice, newObj)）
- 切片的截取操作

  - s[n] slice 中 index 为 n 的元素
  - s[low:] 从 index=low 开始截取到尾部
  - s[:high] 从头部开始截取到 index 为 high-1 的元素（左闭右开）
  - s[low:high] 从 low 开始截取到 index 为 high-1 的元素（左闭右开）
  - s[low:high:max] 同上，同时会将返回的新 slice 的 capacity 设置为 max-low

  在 golang 中一个 slice 类型的元素，其本质是一个结构体，其中包含一个元素数组、长度 length、容量 capacity。
  当使用截取操作创建新的 slice 时，会复用原 slice 底层的那个元素数组，同时维护新的 length 和 capacity，当新 slice 发生扩容之后，底层的元素数组才会改变。

## for range

在 for a,b := range c 遍历中， a 和 b 在内存中只会存在一份，即之后每次循环时遍历到的数据都是以值覆盖的方式赋给 a 和 b，a，b 的内存地址始终不变。由于有这个特性，for 循环里面如果开协程，不要直接把 a 或者 b 的地址传给协程。解决办法：① 在每次循环时，创建一个临时变量。② 把 a、b 的值当做参数传给协程。

- for range 用在数组上时，每一轮的 val 都是对循环开始之前数组的 index 位置值的拷贝，而不是在本次循环时原数组 index 位置的值。
- 用在 slice 上时，每一轮的 val 都是切片 idx 位置当前元素的值。
- 在 for range 的过程中向 slice 中 append 元素，并不会增加循环的次数。

## defer 与 return 机制

defer 是 Go 语言中的一个关键字（延迟调用），一般用于释放资源和连接、关闭文件、释放锁等。和 defer 类似的有 java 的 finally 和 C++的析构函数，这些语句一般是一定会执行的（某些特殊情况后文会提到），不过析构函数析构的是对象，而 defer 后面一般跟函数或方法。

- 在一个函数里，可以多次 defer some_func()，当函数 return 或者出现 panic 后，以后入先出的方式执行。
- 在声明 defer 的时候，给函数传入的参数是值传递，即在 defer 声明的那一刻就确定了入参的值；若是通过闭包的方式在 defer 后的函数里引用了外部的变量，那么在 defer 执行时，使用到的就是此时此刻该变量的值。
- 那些声明在 return 之后的，或者在发生 panic 位置之后的 defer 不会被加入栈中被执行。
- 在不考虑 panic 的情况下，执行顺序为：return -> defer_late -> ... defer_early.
  - 当函数返回值匿名时，return 后跟着的值就是函数最终将返回的值，此时 defer 里不管如何操作都不会影响返回值
  - 当函数返回值具名时，真正返回的值就是 return 后跟着的变量的值，此时如果通过闭包的方式在 defer 里将这个变量的值改变了，那么会影响最终返回值
  - 如果返回值是一个指针，那么在 defer 里对这个指针指向的对象做出了修改，也会影响最终返回指针指向的对象的值。
  - 调用 os.Exit 时 defer 不会被执行（直接把进程干掉了）
  - defer 的底层是 gorutine 维护的一个\_defer 链表，每声明一个 defer，便会创建一个\_defer 并放在链表的头部，在 goroutine 返回或者 panic 的时候，从头部开始执行。

## recover

当函数主动调用 panic 或者调用的方法抛出了 panic，会立即停止函数的执行，转而后进先出地执行在 panic 发生之前所有声明的 defer，其中有 recover 就能抓到，如果都没有，则会向本函数的调用者抛出 panic。

## 函数接收对象

- struct 接受
  某个对象在调用该方法的时候，本质会拷贝一个当前对象再执行，即如果方法里改变了 struct 内字段的值，并不会影响到该调用对象的值。
- struct 指针接受
  某个对象在调用该方法的时候，如果修改了 struct 的字段值，会直接影响被调用对象。

  ```golang
  func testFuncReceiver() {
    p := Person{ "georgayang" }
    p.SetName("123")			//  这个SetName只会修改拷贝后p的Name  而不会真正修改p的Name
    p.SayHello()	// georgayang
  }

  type Person struct {
    Name string
  }

  func (p Person) SayHello() {
    fmt.Println("我是", p.Name)
  }

  func (p Person) SetName(newName string) {
    p.Name = newName
  }
  ```

## uint 类型溢出问题

```golang
func testUintOverflow() {
	 var a uint32 = 0			// 32位 无符号数   0000 0000 0000 0000 0000 0000 0000 0000
	 var b uint32 = 1			// 				 0000 0000 0000 0000 0000 0000 0000 0001
	 // a - b = a + 补码(b)    补码 = 正数 ? 本身 : 反码 + 1    反码 = 符号位不变 + 其余位置取反
	 // 补码(b) = 反码(b) + 1 = 1111 1111 1111 1111 1111 1111 1111 1110 + 1 = 1111 1111 1111 1111 1111 1111 1111 1111 = 无符号的2^32 - 1 = 有符号的 - 2^31
	 fmt.Println(a - b) // 2^32 - 1
}
```

## rune 类型

相当 int32

```golang
// rune is an alias for int32 and is equivalent to int32 in all ways. It is
// used, by convention, to distinguish character values from integer values.
type rune = int32
```

golang 中的字符串底层实现是通过 byte 数组的，中文字符在 unicode 下占 2 个字节，在 utf-8 编码下占 3 个字节，而 golang 默认编码正好是 utf-8

byte 等同于 int8，常用来处理 ascii 字符

rune 等同于 int32,常用来处理 unicode 或 utf-8 字符

## 单引号、双引号、反引号的区别

- 单引号表示 byte 类型或者 rune 类型，默认是后者。
- 双引号指的是字符串，底层是 byte 数组。
- 反引号指的是字符串字面量，不支持任何转义序列（可以比较舒服地换行）。

# 内存模型
# 协程模型
GMP模型
# 并发模型
CSP模型   Communicating Sequential Process
# 垃圾回收
https://blog.csdn.net/Dong_chongwu/article/details/128710443
# 常用包

# 数据结构原理

## slice

答：Go 的 slice 底层数据结构是由一个 array 指针指向底层数组，len 表示切片长度，cap 表示切片容量。

slice 的主要实现是扩容。对于 append 向 slice 添加元素时，假如 slice 容量够用，则追加新元素进去，slice.len++，返回原来的 slice。

当原容量不够，则 slice 先扩容，扩容之后 slice 得到新的 slice，将元素追加进新的 slice，slice.len++，返回新的 slice。

### 为什么 slice 作为参数传入函数内，有时候外部的数组内容会被修改，有时候又不会？

slice 是引用类型，传入的指针会指向函数外 slice 底层的数组，如果在函数内部没有发生扩容，那么对 slice 做出修改就是在原数组上修改，反之会在扩容之后对新数组做修改，所以不会影响到外部的 slice。

### 对于切片的扩容规则：

- 当切片比较小时（容量小于 1024），则采用较大的扩容倍速进行扩容（新的扩容会是原来的 2 倍），避免频繁扩容，从而减少内存分配的次数和数据拷贝的代价。
- 当切片较大的时（原来的 slice 的容量大于或者等于 1024），采用较小的扩容倍速（新的扩容将扩大大于或者等于原来 1.25 倍），主要避免空间浪费，网上其实很多总结的是 1.25 倍，那是在不考虑内存对齐的情况下，实际上还要考虑内存对齐，扩容是大于或者等于 1.25 倍。

## map

1. 并发安全吗？

- 使用时一定要先用 make 做初始化，不然报空指针
- 是并发不安全的，并发读写时会出现 panic

2. 循环是有序的吗？

- 是无序的，for range map 在开始处理循环逻辑的时候就会做随机播种，避免顺序遍历。
- 从底层结构来看，hmap在发生扩容时，原来bucket中的值会搬迁到新的bucket中去，而for range map的原理就是遍历底层的bucket链表数组，如此看来对它的遍历注定是乱序的。

3. map 中删除一个 key，它的内存会释放吗？

- 如果删除的元素是值类型，如 int，float，bool，string 以及数组和 struct，map 的内存不会自动释放

- 如果删除的元素是引用类型，如指针，slice，map，chan 等，map 的内存会自动释放，但释放的内存是子元素应用类型的内存占用

- 将 map 设置为 nil 后，内存被回收。

4. 如何并发访问 map

- sync.Map
- 加读写锁
- 乐观锁 原子操作 ChangeAndSwap

5. 底层数据结构是什么？ 怎么扩容的？

- 底层是 hash table，用链表来解决冲突，这里说的 table 是一个 bucket 数组，每一个 bucket 底层都是一个 bmap ，一个 bmap 可以放 8 个 kv（key连续存放、val连续存放，更高效利用内存）。
- 当一个bucket中元素存满了之后，bucket会有一个指针指向下一个溢出桶，构成类似于一个bucket链表的结构。
- 整个hmap的结构体中会有一个变量B，表示map底层正维护着2^B个bucket，他们全是正式桶；当B > 4时，hmap会认为很有可能会使用到溢出桶，于是会预先分配2^(B-4)个溢出桶，它们在内存上与正式桶是连续的。
- hmap底层是通过计算负载因子来进行扩容的，即map所存元素数量 / bucket数量，阈值默认为6.5，当执行完一次写操作后该负载因子超过阈值，则会翻倍扩容；如果负载因子没有超过阈值，但是hmap中已使用的溢出桶超过一定数量（当B <= 15, 溢出桶数量大于2^B;当B> 15， 溢出桶数量大于2^15），就会触发等量扩容。
- 翻倍扩容的过程是渐进式的，每一次写操作触发扩容后，至多只会迁移两个bucket中的数据到新桶，这样做可以避免扩容瞬间对整个map读写性能的抖动影响。
- 等量扩容触发是因为hmap在写过程中，一边写一边删，导致每个bucket中的bmap非常稀疏（即8个坑位没占满），等量扩容可以对这些数据进行规整，提高内存利用率和读写的效率。


# 一些核心概念

- select
  golang 中的 IO 多路复用机制，主要针对多个 chan 同时读取的场景：

  - 每个 case 里只能处理一个 channel，要么读要么写

  - 多个 case 的执行顺序是随机的

  ```golang
  select {
    case data1 := <- ch1 :{
      // ...
    }
    case data2 := <- ch2 :{
      // ...
    }
    case ch3 <- :{
      // ...
    }
    default : {
      //
    }
  }
  ```

- reflect
- 闭包
- defer
- recover
- panic
- channel
  - 有缓冲的
  - 无缓冲的
- gin.Context

# 性能排查/调优
https://blog.csdn.net/m0_46251547/article/details/126242661

https://blog.csdn.net/luduoyuan/article/details/128721103?spm=1001.2101.3001.6650.1&utm_medium=distribute.pc_relevant.none-task-blog-2%7Edefault%7ECTRLIST%7ERate-1-128721103-blog-126242661.235%5Ev38%5Epc_relevant_sort_base3&depth_1-utm_source=distribute.pc_relevant.none-task-blog-2%7Edefault%7ECTRLIST%7ERate-1-128721103-blog-126242661.235%5Ev38%5Epc_relevant_sort_base3&utm_relevant_index=2