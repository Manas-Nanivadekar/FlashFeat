package internal

import (
	"net"

	"github.com/mdlayher/vsock"
)

const VsockPort = 9000

const EnclaveCID = 16_000

func ListenVsock(port uint32) (net.Listener, error) {
	return vsock.Listen(port, nil)
}

func DialEnclave() (net.Conn, error) {
	return vsock.Dial(EnclaveCID, VsockPort, nil)
}

func VsockConnCloseGraceful(c net.Conn) { _ = c.Close() }
