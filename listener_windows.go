package getlistener

import (
	"fmt"
	"net"
)

func GetListener() (net.Listener, error) {
	host, port, err := loadConfig()
	if err != nil {
		return nil, err
	}
	return net.Listen("tcp", fmt.Sprintf("%s:%d", host, port))
}
