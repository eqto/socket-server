package socket

import (
	"errors"
	"fmt"
	"io"
	"log"
	"net"
	"strconv"
	"sync"
)

const (
	//DefConnectTimeout ...
	DefConnectTimeout = 5
	//DefTimeout ...
	DefTimeout = 60
)
const (
	statusConnected = iota
	statusReceiving
	statusDisconnected
	statusPing
	statusSvrListening
	statusSvrStopping
)

//Server ...
type Server struct {
	port         uint16
	listener     net.Listener
	dec          func(io.Reader) (interface{}, error)
	enc          func(interface{}, io.Writer)
	handle       func(Context, interface{}) (interface{}, error)
	ping         func(Context) interface{}
	serverStatus int
	workers      []*Worker
	asciiMode    bool

	loggerDebug, loggerInfo, loggerErr func(...interface{})

	connectTimeout, timeout uint16
}

//SetDecoder ...
func (s *Server) SetDecoder(dec func(r io.Reader) (interface{}, error)) {
	s.dec = dec
}

//SetConnectTimeout waiting before sending ping, in seconds unit
func (s *Server) SetConnectTimeout(sec uint16) {
	s.connectTimeout = sec
}

//SetTimeout waiting before sending ping, in seconds unit
func (s *Server) SetTimeout(sec uint16) {
	s.timeout = sec
}

//SetEncoder ...
func (s *Server) SetEncoder(enc func(msg interface{}, w io.Writer)) {
	s.enc = enc
}

//SetHandle ...
func (s *Server) SetHandle(handle func(ctx Context, msg interface{}) (interface{}, error)) {
	s.handle = handle
}

//SetPing ...
func (s *Server) SetPing(ping func(ctx Context) interface{}) {
	s.ping = ping
}

//Start ...
func (s *Server) Start() error {
	if e := s.startListening(); e != nil {
		return e
	}
	s.loggerInfo(`Started at`, s.port)
	go s.listenerRoutine()
	return nil
}

//Log ...
func (s *Server) Log(d, i, e func(...interface{})) {
	s.loggerDebug = d
	s.loggerInfo = i
	s.loggerErr = e
}

func (s *Server) listenerRoutine() {
	s.serverStatus = statusSvrListening

	for s.serverStatus == statusSvrListening {
		cn, e := s.listener.Accept()
		if e == nil {
			s.loggerDebug(fmt.Sprintf(`Accept connection: %s`, cn.RemoteAddr()))
			w := newWorker(s, cn)
			s.workers = append(s.workers, w)
			go w.run()
		} else {
			s.loggerInfo(`Stop receiving connection`)
			s.serverStatus = statusSvrStopping
		}
	}
}

func (s *Server) startListening() error {
	if s.dec == nil && s.enc == nil {
		return errors.New(`cannot start without decoder and encoder`)
	}

	listener, e := net.Listen(`tcp`, `:`+strconv.Itoa(int(s.port)))
	if e != nil {
		return e
	}
	s.listener = listener
	return nil
}

//Stop ...
func (s *Server) Stop() error {
	s.listener.Close()
	wg := sync.WaitGroup{}
	wg.Add(len(s.workers))
	for _, w := range s.workers {
		go func(w *Worker) {
			defer wg.Done()
			w.stop()
		}(w)
	}
	wg.Wait()
	return nil
}

//NewServer ...
func NewServer(port int) *Server {
	return &Server{port: uint16(port),
		loggerErr:      log.Println,
		loggerInfo:     log.Println,
		loggerDebug:    log.Println,
		connectTimeout: DefConnectTimeout,
		timeout:        DefTimeout,
	}
}
