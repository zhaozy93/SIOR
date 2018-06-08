package global

import (
	"net"
)

const Version = 0.1
const CONF_FILE = "./conf/service.conf"

var Cfg Config

var FinalPort string

var IPPort = ""

func SetIPPort(port string) {
	ip := getLocalIP()
	if ip == "" {
		ip = "127.0.0.1"
	}
	IPPort = ip + ":" + port
}

func getLocalIP() string {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return "2"
	}
	for _, address := range addrs {
		// check the address type and if it is not a loopback the display it
		if ipnet, ok := address.(*net.IPNet); ok && !ipnet.IP.IsLoopback() {
			if ipnet.IP.To4() != nil {
				return ipnet.IP.String()
			}
		}
	}
	return "1"
}
