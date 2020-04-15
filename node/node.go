package node

import (
	"fmt"
	"net"
	"regexp"
	"time"
)

func ByIP() int64 {
	nodeNumber := (time.Now().UnixNano() & 3) << 8
	ifaces, err := net.Interfaces()
	if nil != err {
		panic(fmt.Errorf("net interface error %w", err))
	}
	// handle err
	for _, i := range ifaces {
		if (i.Flags & (net.FlagBroadcast | net.FlagUp)) != (net.FlagBroadcast | net.FlagUp) {
			continue
		}
		addrs, err := i.Addrs()
		if nil != err {
			panic(fmt.Errorf("ip address error %w", err))
		}
		// handle err
		for _, addr := range addrs {
			var ip net.IP
			switch v := addr.(type) {
			case *net.IPNet:
				ip = v.IP
			case *net.IPAddr:
				ip = v.IP
			}
			if nil == ip {
				continue
			}
			to4 := ip.To4()
			if nil == to4 {
				continue
			}
			reg := regexp.MustCompilePOSIX(`^172\.`)
			if reg.MatchString(to4.String()) {
				continue
			}
			nodeNumber = nodeNumber | int64(to4[3])
		}
	}
	return nodeNumber
}
