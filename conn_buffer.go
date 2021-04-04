package socket

import (
	"net"
)

//ConnBuffer ...
type ConnBuffer struct {
	Reader
	Writer
	conn net.Conn
}

//ReadBytes ...
func (c *ConnBuffer) ReadBytes(length int) (b []byte, err error) {
	b = make([]byte, length)
	_, e := c.conn.Read(b)
	if e != nil {
		return nil, e
	}
	return b, nil
}

//WriteBytes ...
func (c *ConnBuffer) WriteBytes(data []byte) (length int, err error) {
	return c.conn.Write(data)
}

//NewBuffer ...
func NewBuffer(conn net.Conn) *ConnBuffer {
	return &ConnBuffer{conn: conn}
}
