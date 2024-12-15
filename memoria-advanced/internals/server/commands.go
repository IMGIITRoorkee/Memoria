// the commands that server can handle
package server

// represents the command
type Command struct {
	Raw  []byte   // encoded message by the RESP Protocol
	Args [][]byte // arguments to the command
}
