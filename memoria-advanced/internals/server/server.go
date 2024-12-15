// / Defines the server for clients
package server

import (
	"bufio"
	"io"
	"net"

	"github.com/utkarshdagoat/memoria/internals/protocol"
)

type Reader struct {
	rd    *bufio.Reader
	buf   []byte
	start int
	end   int
	cmds  []Command
}

type conn struct {
	conn net.Conn
	wr   *io.Writer
	rd   *Reader
}

type Server struct {
	net     string
	laddr   string
	handler func(conn protocol.Conn, cmd Command)
	conns   map[*conn]bool
}
