# 为消息队列添加监控

修改文件：
basic-go/pkg/saramax/consumer_handler.go，新增72行，及相关builder和初始化工作
basic-go/pkg/saramax/batch_consumer_handler.go，新增84行，及相关初始化工作


监控思路：
需要监控消息消费的错误，正常而言，消息都应该被正确消费掉，不应该出现报错，如果报错，就需要被观测到

告警设置思路：
计算错误增长的速率，如果速率超过阈值，则要报警

