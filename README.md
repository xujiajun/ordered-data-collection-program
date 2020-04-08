# ordered data collection program

### 系统描述
一个系统包含多个数据生产者，每个数据生产者都会不停的生产数据，并且生产者会为每个 data 发送两条 messages
* P message(`type == prepare`), 包含了 prepare token
* C message(`type == commit`), 包含了 prepare 和 commit token

### 实现要求
* 按照 commit token 的升序顺序对 data 进行排序，并且尽快的输出 data
* 为你的排序过程实现 `流控` 机制，用来控制内存和 CPU 资源的使用
* 请按照一个完善工程的标准来完成这个小作业

### 实现过程

根据第一条的关键要求是"进行排序，并且尽快的输出 data"，主要考察排序的算法时间复杂度。

* 方案1：用标准库自带的sort包（sort.Slice）

* 方案2：快排之类的O(nlogn) 排序算法

* 方案3：其他算法

首先方案1和方案2可以合并，标准库已经封装了归并排序、堆排序和快速排序等，sort包会根据实际数据自动选择高效的排序算法。

方案3用一些其他算法，需要更快，需要用到多线程，利用多核，但是并发的程序对于排序这种计算密集型，还要考虑到同步等问题，不一定会快，反而会变慢，
实际测试中，也发现还不如用单线程的排序算法。最后我选择多路归并排序算法。算法思路：
* 1、先收集各个数据生产者的数据
* 2、根据CPU核数做数据切分分成N个shard，开对应数量goroutine去并发并行的执行sort
* 3、把各个shard收集起来，push到最小堆（按照每个shard的首部元素作为比较对象），这样保证堆顶的元素首位元素是最小的，把这个值加到结果集，更新这个slice的Idx。
再调整最小堆。
* 4、重复第三步，最后排序完成。

看下benchmark的效果：

```
goos: darwin
goarch: amd64
pkg: github.com/xujiajun/ordered-data-collection-program
BenchmarkMyMergeSort
BenchmarkMyMergeSort-8   	       1	2680632242 ns/op
BenchmarkNormalSort
BenchmarkNormalSort-8    	       1	7350065377 ns/op
PASS
```

MyMergeSort算法比标准库的sort.Slice快了近2倍，就速度而言效果还是很好的。


再看下实际接受方收到1w条数据（5个client，发送2*1000）之后的排序时间：

标准库的sort.Slice，执行几次：
```
3.472551ms
3.081911ms
3.400909ms
3.576235ms
3.541364ms
2.501259ms
3.036942ms
```

结论标准库的sort.Slice 大概在3ms上下波动。

用MyMergeSort算法，执行几次：

```
1.941394ms
1.950004ms
1.62079ms
1.615415ms
2.144688ms
1.955708ms
1.869308ms
```
结论标准库的MyMergeSort大概在2ms上下波动。

以上说明就这个作业的发送1w条数据的场景中`MyMergeSort`还是比`sort.Slice`快的。

关于第二条要求，我们来对接收方做下限流就可以满足。流量限制的手段有很多，最常见的：漏桶、令牌桶两种。

这边使用按照go wiki 上提到的方法，使用time.Ticker简单实现就可满足需求。

先看下没做流控的时候的内存分配和CPU使用情况：



再看下做了流控的时候的内存分配和CPU使用情况：



## 使用指南

### 下载

```
go get github.com/xujiajun/ordered-data-collection-program
```

### 运行

```
go run main.go sort.go
```

### 单元测试

```
go test -v
```

### 性能测试

```
go test -bench=.
```

### 参考资料

* https://xargin.com/talent-plan-week1-solution/
* https://github.com/Deardrops/pingcapAssignment
* https://github.com/golang/go/wiki/RateLimiting