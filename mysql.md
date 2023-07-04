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
