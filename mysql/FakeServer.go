package mysql

import (
	"github.com/charmbracelet/log"
	"net"
	"strings"
)

// honeypotPayload 处理连接并发送蜜罐负载
func honeypotPayload(conn net.Conn) {
	defer conn.Close()

	if err := send(conn, ServerGreeting); err != nil {
		log.Warnf("Error sending greeting: %v", err)
		return
	}

	if err := send(conn, LoginOK); err != nil {
		log.Warnf("Error sending login: %v", err)
		return
	}

	for {
		buf, err := bRecv(conn)
		if err != nil {
			log.Warnf("Error reading data: %v", err)
			return
		}

		if err = processCommands(conn, buf); err != nil {
			log.Warnf("Error processing commands: %v", err)
			return
		}
	}
}

// processCommands 处理请求中的命令
func processCommands(conn net.Conn, buf []byte) error {
	cmd := string(buf)
	var resp string

	switch {
	case strings.Contains(cmd, `utf8mb4`):
		resp = SetName
	case strings.Contains(cmd, `lower_case_%`):
		resp = Ndbcluster
	case strings.Contains(cmd, `information_schema.SCHEMATA`):
		resp = Schemata
	case strings.Contains(cmd, `SHOW DATABASES`):
		resp = ShowDatabases
	default:
		resp = Aggressor
	}

	return send(conn, resp)
}

func SqlHoneypot() {
	listener, err := net.Listen("tcp", ":3306")
	if err != nil {
		log.Fatalf("Error listening on port 3306: %v", err)
	}
	defer listener.Close()

	log.Info("Server is listening on port 3306")

	for {
		conn, err := listener.Accept()
		if err != nil {
			log.Warnf("Error accepting connection: %v", err)
			continue
		}

		go func(conn net.Conn) {
			defer conn.Close()
			log.Warnf("New connection from: %s", conn.RemoteAddr())
			honeypotPayload(conn)
		}(conn)
	}
}
