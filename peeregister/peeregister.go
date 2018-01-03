package peeregister

import (
	"errors"
	"github.com/Zumium/fyer/cfg"
	db_center "github.com/Zumium/fyer/db/center"
	"net"
)

var (
	ErrNoIPFound = errors.New("cannot found any IP addresses")
)

//InitServer initializes and starts the tcp registration service
func InitServer() error {
	lis, err := net.ListenTCP("tcp", &net.TCPAddr{IP: net.ParseIP("::"), Port: cfg.PeerRegisterPort()}) //Listens for all address of IPv4 and IPv6
	if err != nil {
		return err
	}

	go func() {
		for {
			conn, err := lis.AcceptTCP()
			if err != nil {
				continue
			}
			go handlePeerRegistering(conn)
		}
	}()

	return nil
}

func handlePeerRegistering(conn *net.TCPConn) {
	defer conn.Close()

	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	if err != nil {
		conn.Write([]byte(err.Error()))
		return
	}
	peerId := string(buf[:n])
	dbHandler, err := db_center.ToPeerID(peerId)
	if err != nil {
		conn.Write([]byte(err.Error()))
		return
	}
	if err := dbHandler.Edit().SetAddress(conn.RemoteAddr().(*net.TCPAddr).IP.String()).Done(); err != nil {
		conn.Write([]byte(err.Error()))
		return
	}

	conn.Write([]byte("OK"))
}

//RegisterLocal registers the local peer's id at the center peer
func RegisterLocal() error {
	ips, err := net.LookupIP(cfg.CenterAddress())
	if err != nil {
		return err
	}
	if len(ips) == 0 {
		return ErrNoIPFound
	}
	conn, err := net.DialTCP("tcp", nil, &net.TCPAddr{IP: ips[0], Port: cfg.PeerRegisterPort()})
	if err != nil {
		return err
	}
	defer conn.Close()

	_, err = conn.Write([]byte(cfg.PeerID()))
	if err != nil {
		return err
	}

	buf := make([]byte, 256)
	n, err := conn.Read(buf)
	retStr := string(buf[:n])
	if retStr != "OK" {
		return errors.New(retStr)
	}

	return nil
}
