package getlistener

import (
	"fmt"
	"net"
)

func GetListener() (net.Listener, error) {
	if initErr != nil {
		return nil, initErr
	}
	return net.Listen("tcp", fmt.Sprintf("%s:%d", HOST, PORT))
}
