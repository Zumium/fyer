package connectionmngr

import (
	"google.golang.org/grpc"
)

type Connection struct {
	*grpc.ClientConn
	refCount uint32

	addr string
}

var connectionPool = make(map[string]*Connection)

func ConnectTo(destAddr string) (*Connection, error) {
	conn, exist := connectionPool[destAddr]
	if exist {
		conn.refCount++
		return conn, nil
	}

	clientConn, err := grpc.Dial(destAddr)
	if err != nil {
		return nil, err
	}

	newConn := &Connection{clientConn, 1, destAddr}
	connectionPool[destAddr] = newConn
	return newConn, nil
}

//Close discount the reference counter by 1 and closes the underlying connection once the reference counter equals 0
func (conn *Connection) Close() bool {
	conn.refCount--
	if conn.refCount == 0 {
		conn.ClientConn.Close()
		delete(connectionPool, conn.addr)
		return false
	}
	return true
}
