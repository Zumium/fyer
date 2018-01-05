package connectionmngr

import (
	"fmt"
	"github.com/Zumium/fyer/cfg"
	"github.com/xtaci/kcp-go"
	"google.golang.org/grpc"
	_ "google.golang.org/grpc/encoding/gzip"
	"net"
	"sync"
	"time"
)

type Connection struct {
	*grpc.ClientConn
	refCount uint32

	addr string
}

var connectionPool = make(map[string]*Connection)
var serialAccess sync.Mutex

func ConnectTo(destAddr string) (*Connection, error) {
	serialAccess.Lock()
	defer serialAccess.Unlock()

	fmt.Printf("Establishing connection to %s\n", destAddr)

	conn, exist := connectionPool[destAddr]
	if exist {
		conn.refCount++
		return conn, nil
	}

	clientConn, err := grpc.Dial(
		destAddr, grpc.WithBlock(),
		grpc.WithInsecure(),
		grpc.WithDialer(func(addr string, duration time.Duration) (net.Conn, error) {
			return kcp.DialWithOptions(addr, nil, 10, 3)
		}),
		grpc.WithDefaultCallOptions(
			grpc.UseCompressor("gzip"),
			grpc.MaxCallRecvMsgSize(cfg.MaxSendRecvMsgSize()),
			grpc.MaxCallSendMsgSize(cfg.MaxSendRecvMsgSize()),
		),
	)
	if err != nil {
		return nil, err
	}

	newConn := &Connection{clientConn, 1, destAddr}
	connectionPool[destAddr] = newConn
	return newConn, nil
}

//ConnectToCenter is a helper function to conveniently establish a connection towards center
func ConnectToCenter() (*Connection, error) {
	return ConnectTo(fmt.Sprintf("%s:%d", cfg.CenterAddress(), cfg.Port()))
}

func ConnectToWithDefaultPort(addr string) (*Connection, error) {
	return ConnectTo(fmt.Sprintf("%s:%d", addr, cfg.Port()))
}

//Close discount the reference counter by 1 and closes the underlying connection once the reference counter equals 0
func (conn *Connection) Close() bool {
	serialAccess.Lock()
	defer serialAccess.Unlock()

	conn.refCount--
	if conn.refCount == 0 {
		conn.ClientConn.Close()
		delete(connectionPool, conn.addr)
		return false
	}
	return true
}
