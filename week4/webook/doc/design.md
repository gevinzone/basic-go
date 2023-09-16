# 同步转异步


# 方案

## 1. 设计如何判断服务商崩溃

计算请求成功次数和请求失败次数的比值，如果超过预期，则判定服务商崩溃，切换下一个服务商。主要思路为：

1. idx 表示当前的供应商
2. 提供`failureCnt`和`successCnt`两个变量，分别记录失败和成功的请求个数
3. if failureCnt/(failureCnt + successCnt) >= threshold，idx++
4. if svc[idx].send(...): successCnt++ else failureCnt++

```go
func (f *FailureRateFailOverService) Send(ctx context.Context, tpl string, args []string, numbers ...string) error {
	idx := atomic.LoadInt64(&f.idx)
	failureCnt := atomic.LoadInt64(&f.failureCnt)
	successCnt := atomic.LoadInt64(&f.successCnt)
	if float64(failureCnt)/float64(failureCnt+successCnt) >= f.rate {
		newIdx := (idx + 1) % int64(len(f.svcs))
		if atomic.CompareAndSwapInt64(&f.idx, idx, newIdx) {
			atomic.StoreInt64(&f.failureCnt, 0)
			atomic.StoreInt64(&f.successCnt, 0)
		} else {
			// cas 不成功的并发请求，等上面重置完成再继续
			// 如果没有这个分支，虽然有并发不安全隐患，但几乎不可能发生，且即便发生，对本业务影响比较小，可以忽略本分支
			time.Sleep(time.Millisecond)
		}
		idx = atomic.LoadInt64(&f.idx)
	}
	err := f.svcs[idx].Send(ctx, tpl, args, numbers...)
	switch {
	case err == nil:
		atomic.AddInt64(&f.successCnt, 1)
		return nil
	case errors.Is(err, context.DeadlineExceeded):
		atomic.AddInt64(&f.failureCnt, 1)
		return err
	default:
		return err
	}
}
```

## 2. 如何既能触发限流判断，又能判断服务商崩溃

依然用装饰器模式，里面供应商为RateLimitService，外面为FailOverService。

```go
type FailureRateFailOverService struct {
	svcs       []sms.RateLimitService
	idx        int64
	failureCnt int64
	successCnt int64
	rate       float64
}
```

## 3. 失败记录的数据库转存设计

数据库表结构为：

```sql
create table if not exists webook.sms
(
    id         bigint auto_increment primary key,
    tpl        varchar(1000),
    args       varchar(1000),
    numbers    text,
    processing tinyint,
    retry      smallint,
    ctime      bigint null,
    utime      bigint null

)
```

## 4. 异步设计

1. 提供一个守护态对象，里面包含N个goroutine作为消费者，每隔固定时间（如1s）检查数据库表中有无满足要求、有待处理的记录，
   1. 若无，则继续轮训，若有，则计算记录跟新时间与当前时间的duration，并与延迟处理时间（waitDuration）做比较
   2. 若duration < waitDuration，则让goroutine sleep(waitDuration-duration)后再唤醒
   3. 若duration > waitDuration，则用数据库乐观锁并发改写该记录状态，改写不成功的，按上面逻辑开始处理下一个记录
   4. 改写成功的消费者，用atomic找到当前使用的svc，执行svc.send()操作，并删除该记录
   5. 若svc.send()失败，重试次数+1，数据库写会该请求
2. send()方法里，调用err:= svc.send(), 当err != nil 时，写该请求到数据库
3. 判断数据库是否有满足条件记录的判断逻辑为，processing=1, retry<threshold


# 优点

1. 设计如何判断服务商崩溃的算法实现，用了atomic，考虑了无锁并发安全，性能比较高，虽然有一定的并发隐患，但对本业务影响不大，且本业务本身的性质也决定了并发问题的危害较小，故这里用atomic而非锁是平衡折中后的结果，可以接受 
2. 利用生产者消费者模式实现同步转异步，解耦错误处理的同时，实现了性能的可配
3. 使用类似延迟队列的逻辑，进行错误处理


# 缺点

1. 设计如何判断服务商崩溃的算法实现，用了atomic，但业务逻辑不是原子操作，还允许并行执行，虽然代码中尽量去避免并发问题了，但依然不是很准确，如果想精确保证业务互斥，需要用锁，若考虑多实例，这里的锁需要用分布式锁。
2. 未处理重试超过最大次数的请求
3. 守护态的消费者协程，会一直存在
4. 数据库设计时，numbers 用text字段进行存储
5. 用乐观锁抢记录不太合适，乐观锁适用于并发不太激烈的场景，这里每个记录，都会让全部消费者去抢，如果消费者较多，性能浪费更多
6. workshop 和 RateLimitFailOverService 的创建循环依赖，耦合实现

# 改进方向

1. 将守护态的消费者协程，改造为类似Java线程池的逻辑
2. 引入类似“死信队列”的逻辑，处理超过最大重试次数的记录
3. 可以限制numbers的最大数量
4. 使用分布式锁解决消费者的并发抢记录问题

# 方案试用场景

该方案适用于用户的峰值并发较高，且对可用性要求较高的场景