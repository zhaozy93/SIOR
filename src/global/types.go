package global

type Config struct {
	Service struct {
		Port    int
		LogFile string
	}
	Raft struct {
		Heartbeat          int
		HeartbeatTimeout   int
		ElectionTimeoutMin int
		ElectionTimeoutMax int
	}
	Hosts struct {
		Cluster []string
	}
}
