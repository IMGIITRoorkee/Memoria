// This file contains code for connection to the client
package protocol

// repsresenting the connection to the client
type Conn interface {
	RemoteAddr() string // get the remote addr of the client
	Close() string      // close the connection the client

	WriteError(msg string) //
}
