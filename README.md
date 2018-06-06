#SIOR
Simple Implementation Of Raft in Go Language With HTTP 



## Example Flow
1. 计算一个随机时间 ElectionTimeoutMin 与 ElectionTimeoutMax之间，定位变量为ElectionTimeout
1. 每个程序都是一个RaftClient,在启动的时候都是Follower状态
2. 如果能在ElectionTimeout时间内收到来自Leader的ttl则保持Follower
3. 之后只要ElectionTimeout时间段内没有收到来自Leader的ttl则变为Candidate
4. Candidate状态如果收到ttl则继续转变为Follower
5. Candidate身份会立即发起vote，得到多数票则变为Leader，否则持续发起Vote，直到变为Leader或者收到ttl
6. 每一个RaftClient都只能在500毫秒内投一次票，因此基本可以保证投出Leader
