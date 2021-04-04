package socket

//Context ...
type Context interface {
	GetMessage() interface{}
	SetValue(key string, val interface{})
	GetValue(key string) interface{}
	Writer() Writer
	Reader() Reader
	Close() error
	Server() *Server
	Disconnect()
}
