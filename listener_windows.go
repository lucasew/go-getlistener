package getlistener

import (
	"fmt"
	"net"
)

func GetListener() (net.Listener, error) {
	return net.Listen("tcp", fmt.Sprintf("%s:%d", HOST, PORT))
}
