#SIOR
Simple Implementation Of Raft in Go Language With HTTP 

简单实现了Raft的Election Term 与 简易版数据同步

[Raft 动画演示](http://thesecretlivesofdata.com/raft/)

[Raft 论文](https://web.stanford.edu/~ouster/cgi-bin/papers/raft-atc14)


## Election Term Flow
0. 计算一个随机时间 ElectionTimeoutMin 与 ElectionTimeoutMax之间，定位变量为ElectionTimeout
1. 每个程序都是一个RaftClient,在启动的时候都是Follower状态
2. 如果能在ElectionTimeout时间内收到来自Leader的ttl则保持Follower
3. 之后只要ElectionTimeout时间段内没有收到来自Leader的ttl则变为Candidate
4. Candidate状态如果收到ttl则继续转变为Follower
5. Candidate身份会立即发起vote，得到多数票则变为Leader，否则持续发起Vote，直到变为Leader或者收到ttl
6. 每一个RaftClient都只能在同一个Term内投一次票+大多数结果可以保证投出Leader

## Example
0. 程序在service.conf已经配置了三个port
1. 打开三个控制台，分别执行三次make start即可启动三个端口监听以模仿三个Client
2. 查看三个的日志 即可看明白
3. 中途可以干掉Leader那个终端，查看剩余两个的变化
4. 在将刚刚干掉的终端重新make start，他又恢复到cluster中
