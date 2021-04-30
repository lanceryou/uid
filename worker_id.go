package uid

import (
	"errors"
	"math/big"
	"net"
)

type Worker interface {
	WorkerID() int64
}

type WorkerFunc func() int64

func (w WorkerFunc) WorkerID() int64 {
	return w()
}

func defaultWorker() int64 {
	ip, err := getLocalIP()
	if err != nil {
		return 0
	}

	workerID := big.NewInt(0)
	workerID.SetBytes(ip)
	return workerID.Int64()
}

func getLocalIP() (ip net.IP, err error) {
	addrs, err := net.InterfaceAddrs()
	if err != nil {
		return
	}

	for _, rawAddr := range addrs {
		switch addr := rawAddr.(type) {
		case *net.IPAddr:
			ip = addr.IP.To4()
		case *net.IPNet:
			ip = addr.IP.To4()
		default:
			continue
		}

		if ip == nil || ip.IsLoopback() {
			continue
		}

		return
	}

	return nil, errors.New("failed get ip address")
}
