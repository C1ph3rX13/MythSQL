package mysql

import (
	"encoding/hex"
	"github.com/charmbracelet/log"
	"io"
	"net"
)

// send 发送字符串数据
func send(conn net.Conn, s string) error {
	buff, err := hex.DecodeString(s)
	if err != nil {
		log.Warnf("Error decoding string: %v", err)
		return err
	}
	return bSend(conn, buff)
}

// bSend 发送字节数据
func bSend(conn net.Conn, b []byte) error {
	_, err := conn.Write(b)
	return err
}

// recv 接收字符串数据
func recv(conn net.Conn) (string, error) {
	data, err := bRecv(conn)
	if err != nil {
		return "", err
	}
	return string(data), nil
}

// bRecv 接收字节数据
func bRecv(conn net.Conn) ([]byte, error) {
	buf := make([]byte, 10240000)
	n, err := conn.Read(buf)
	if err != nil {
		if err == io.EOF {
			log.Warnf("Connection closed by client.")
		} else {
			log.Warnf("Error reading data: %v", err)
		}
		return nil, err
	}
	log.Warnf("Received data: %s", string(buf[:n]))
	return buf[:n], nil
}
