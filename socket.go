package socket

import (
	"net"
)

//Reader ...
type Reader interface {
	ReadBytes(length int) (b []byte, err error)
}

//Writer ...
type Writer interface {
	WriteBytes(data []byte) (length int, err error)
}

type socketCtx struct {
	Context
	msg        interface{}
	conn       net.Conn
	buff       *ConnBuffer
	value      map[string]interface{}
	server     *Server
	disconnect bool
}

func (s *socketCtx) Server() *Server {
	return s.server
}

//GetMessage ...
func (s *socketCtx) GetMessage() interface{} {
	return s.msg
}

func (s *socketCtx) buffer() *ConnBuffer {
	if s.buff == nil {
		s.buff = &ConnBuffer{}
	}
	s.buff.conn = s.conn
	return s.buff
}

//Writer ...
func (s *socketCtx) Writer() Writer {
	return s.buffer()
}

//Reader ...
func (s *socketCtx) Reader() Reader {
	return s.buffer()
}

//Close ...
func (s *socketCtx) Close() error {
	return s.conn.Close()
}

func (s *socketCtx) SetValue(key string, val interface{}) {
	if s.value == nil {
		s.value = make(map[string]interface{})
	}
	s.value[key] = val
}
func (s *socketCtx) GetValue(key string) interface{} {
	if s.value == nil {
		s.value = make(map[string]interface{})
	}
	return s.value[key]
}

func (s *socketCtx) Disconnect() {
	s.disconnect = true
}
