# 思路

当前的preempt的接口，没处理抢占后续约失败的情况。在Service层，preempt层，抢占成功时，会修改job的更新时间，然后每间隔refreshInterval 时间，还会修改job的更新时间，以此标记为续约。
故续约失败，意味着 job 更新时间与当前时间间隔大于refreshInterval。

所以抢占job时，先去抢占未调度过的任务，如果没有记录，则去抢占续约失败的job


# 作业

主要改动：

修改webook/internal/repository/dao/job.go 文件，72行的错误处理中，增加抢占续约失败job的逻辑

给Preempt函数增加refreshInterval参数，并按上述进行改动

```go
func (g *GORMJobDAO) Preempt(ctx context.Context, refreshInterval time.Duration) (Job, error) {
    ...
}
```
