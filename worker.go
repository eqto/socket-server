package socket

import (
	"net"
	"sync"
	"time"

	log "github.com/eqto/go-logger"
)

const (
	modeWrite = iota
	modeRead
)

//Worker ...
type Worker struct {
	svr  *Server
	conn net.Conn
	done chan uint8
	wg   sync.WaitGroup
}

func (w *Worker) stop() {
	defer recover()
	w.done <- 1
	w.wg.Wait()
}

func (w *Worker) run() {
	defer func() {
		if r := recover(); r != nil {
			log.E(r)
		}
	}()
	w.svr.loggerDebug(`Worker start`)
	ctx := &socketCtx{conn: w.conn, server: w.svr}
	status := statusConnected
	readerCh := make(chan interface{}, 10)

	//reading routine
	go func() {
		for status != statusDisconnected {
			msg, e := w.svr.dec(ctx.conn)
			if e != nil {
				w.svr.loggerDebug(e)
				readerCh <- e
				time.Sleep(1 * time.Second)
				return
			} else {
				readerCh <- msg
			}
		}
	}()
	//end

	for status != statusDisconnected {
		timeout := w.svr.timeout
		if status == statusConnected {
			timeout = w.svr.connectTimeout
		} else if status == statusPing {
			timeout = 3
		}
		ctx.msg = nil
		select {
		case <-time.After(time.Duration(timeout) * time.Second):
			if status == statusConnected || status == statusPing {
				w.svr.loggerDebug(`Timeout`, ctx.conn.RemoteAddr())
				status = statusDisconnected
			} else {
				msg := w.svr.ping(ctx)
				w.svr.enc(msg, w.conn)
				// ctx.Writer().WriteBytes(w.svr.enc(msg, w.conn))
				status = statusPing
				continue
			}
		case req := <-readerCh:
			if _, ok := req.(error); ok {
				status = statusDisconnected
				continue
			} else {
				status = statusReceiving
				ctx.msg = req
			}
		case <-w.done:
			status = statusDisconnected
		}
		if status != statusDisconnected {
			w.wg.Add(1)
			defer w.wg.Done()
			msg, e := w.svr.handle(ctx, ctx.msg)

			if e != nil {
				w.svr.loggerErr(e)
			}
			if msg != nil {
				w.svr.enc(msg, w.conn)
			}
			if ctx.disconnect {
				status = statusDisconnected
			}
		}
	}
	w.svr.loggerDebug(`Connection closed:`, w.conn.RemoteAddr())
	w.conn.Close()
}

func newWorker(svr *Server, conn net.Conn) *Worker {
	return &Worker{svr: svr, conn: conn, done: make(chan uint8, 1)}
}
