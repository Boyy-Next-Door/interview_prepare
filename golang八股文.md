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

## for range
在 for a,b := range c 遍历中， a 和 b 在内存中只会存在一份，即之后每次循环时遍历到的数据都是以值覆盖的方式赋给 a 和 b，a，b 的内存地址始终不变。由于有这个特性，for 循环里面如果开协程，不要直接把 a 或者 b 的地址传给协程。解决办法：①在每次循环时，创建一个临时变量。②把a、b的值当做参数传给协程。

- for range 用在数组上时，每一轮的val都是对循环开始之前数组的index位置值的拷贝，而不是在本次循环时原数组index位置的值。
- 用在slice上时，每一轮的val都是切片idx位置当前元素的值。
- 在for range的过程中向slice中append元素，并不会增加循环的次数。

## defer与return机制
defer是Go语言中的一个关键字（延迟调用），一般用于释放资源和连接、关闭文件、释放锁等。和defer类似的有java的finally和C++的析构函数，这些语句一般是一定会执行的（某些特殊情况后文会提到），不过析构函数析构的是对象，而defer后面一般跟函数或方法。
- 在一个函数里，可以多次 defer some_func()，当函数return或者出现panic后，以后入先出的方式执行。
- 在声明defer的时候，给函数传入的参数是值传递，即在defer声明的那一刻就确定了入参的值；若是通过闭包的方式在defer后的函数里引用了外部的变量，那么在defer执行时，使用到的就是此时此刻该变量的值。
- 那些声明在return之后的，或者在发生panic位置之后的defer不会被加入栈中被执行。
- 在不考虑panic的情况下，执行顺序为：return -> defer_late -> ... defer_early.
  - 当函数返回值匿名时，return后跟着的值就是函数最终将返回的值，此时defer里不管如何操作都不会影响返回值
  - 当函数返回值具名时，真正返回的值就是return后跟着的变量的值，此时如果通过闭包的方式在defer里将这个变量的值改变了，那么会影响最终返回值
  - 如果返回值是一个指针，那么在defer里对这个指针指向的对象做出了修改，也会影响最终返回指针指向的对象的值。
  - 调用os.Exit时defer不会被执行（直接把进程干掉了）
  - defer的底层是gorutine维护的一个_defer链表，每声明一个defer，便会创建一个_defer并放在链表的头部，在goroutine返回或者panic的时候，从头部开始执行。
## recover
当函数主动调用panic或者调用的方法抛出了panic，会立即停止函数的执行，转而后进先出地执行在panic发生之前所有声明的defer，其中有recover就能抓到，如果都没有，则会向本函数的调用者抛出panic。
## 函数接收对象
- struct接受
  某个对象在调用该方法的时候，本质会拷贝一个当前对象再执行，即如果方法里改变了struct内字段的值，并不会影响到该调用对象的值。
- struct指针接受
  某个对象在调用该方法的时候，如果修改了struct的字段值，会直接影响被调用对象。
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

## uint类型溢出问题
```golang
func testUintOverflow() {
	 var a uint32 = 0			// 32位 无符号数   0000 0000 0000 0000 0000 0000 0000 0000
	 var b uint32 = 1			// 				 0000 0000 0000 0000 0000 0000 0000 0001
	 // a - b = a + 补码(b)    补码 = 正数 ? 本身 : 反码 + 1    反码 = 符号位不变 + 其余位置取反
	 // 补码(b) = 反码(b) + 1 = 1111 1111 1111 1111 1111 1111 1111 1110 + 1 = 1111 1111 1111 1111 1111 1111 1111 1111 = 无符号的2^32 - 1 = 有符号的 - 2^31
	 fmt.Println(a - b) // 2^32 - 1
}
```

## rune类型
相当int32
```golang
// rune is an alias for int32 and is equivalent to int32 in all ways. It is
// used, by convention, to distinguish character values from integer values.
type rune = int32
```

golang中的字符串底层实现是通过byte数组的，中文字符在unicode下占2个字节，在utf-8编码下占3个字节，而golang默认编码正好是utf-8

byte 等同于int8，常用来处理ascii字符

rune 等同于int32,常用来处理unicode或utf-8字符
## 单引号、双引号、反引号的区别
- 单引号表示byte类型或者rune类型，默认是后者。
- 双引号指的是字符串，底层是byte数组。
- 反引号指的是字符串字面量，不支持任何转义序列（可以比较舒服地换行）。
# 内存模型

# 并发模型

# 垃圾回收

# 常用包

# 数据结构原理
## slice
答：Go 的 slice 底层数据结构是由一个 array 指针指向底层数组，len 表示切片长度，cap 表示切片容量。

slice 的主要实现是扩容。对于 append 向 slice 添加元素时，假如 slice 容量够用，则追加新元素进去，slice.len++，返回原来的 slice。

当原容量不够，则 slice 先扩容，扩容之后 slice 得到新的 slice，将元素追加进新的 slice，slice.len++，返回新的 slice。

### 为什么slice作为参数传入函数内，有时候外部的数组内容会被修改，有时候又不会？
slice是引用类型，传入的指针会指向函数外slice底层的数组，如果在函数内部没有发生扩容，那么对slice做出修改就是在原数组上修改，反之会在扩容之后对新数组做修改，所以不会影响到外部的slice。

### 对于切片的扩容规则：
- 当切片比较小时（容量小于 1024），则采用较大的扩容倍速进行扩容（新的扩容会是原来的 2 倍），避免频繁扩容，从而减少内存分配的次数和数据拷贝的代价。
- 当切片较大的时（原来的 slice 的容量大于或者等于 1024），采用较小的扩容倍速（新的扩容将扩大大于或者等于原来 1.25 倍），主要避免空间浪费，网上其实很多总结的是 1.25 倍，那是在不考虑内存对齐的情况下，实际上还要考虑内存对齐，扩容是大于或者等于 1.25 倍。


## map
1. 并发安全吗？
  - 使用时一定要先用make做初始化，不然报空指针
  - 是并发不安全的，并发读写时会出现panic
2. 循环是有序的吗？
  - 是无序的，for range map在开始处理循环逻辑的时候就会做随机播种，避免顺序遍历
3. map中删除一个key，它的内存会释放吗？
  - 如果删除的元素是值类型，如int，float，bool，string以及数组和struct，map的内存不会自动释放

  - 如果删除的元素是引用类型，如指针，slice，map，chan等，map的内存会自动释放，但释放的内存是子元素应用类型的内存占用

  - 将map设置为nil后，内存被回收。
4. 如何并发访问map
  - sync.Map
  - 加读写锁
  - 乐观锁 原子操作 ChangeAndSwap 
5. 底层数据结构是什么？ 怎么扩容的？
  - 底层是hash table，用链表来解决冲突，这里说的table是一个bucket数组，每一个bucket都是一个bmap链表，一个bmap可以放8个kv
# 一些核心概念
- select
golang中的IO多路复用机制，主要针对多个chan同时读取的场景：
  - 每个case里智能处理一个channel，要么读要么写

  - 多个case的执行顺序是随机的
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
