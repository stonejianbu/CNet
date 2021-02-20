/*
为了能够更快的处理用户的请求，又避免goroutine创建和销毁的性能消耗，不必每来一个连接就创建多个线程去处理请求，而是提前创建
*/
package cnet

import (
	log "github.com/sirupsen/logrus"
)

type Worker struct {
	conChan chan *Connection
	endChan chan bool
	MaxConnNum int
}

func NewWorker(num int) *Worker {
	return &Worker{
		conChan: make(chan *Connection, num),
		endChan: make(chan bool),
		MaxConnNum: num,
	}
}

func (w *Worker) add(obj *Connection) {
	w.conChan<-obj
}

func (w *Worker) Start() {
	for i:=0;i< w.MaxConnNum;i++ {
		log.Infof("worker-%d has start and wait for a connection",i)
		go func() {
		loop:
			for {
				select {
				case connection := <-w.conChan:
					connection.Start()
				case <-w.endChan:
					break loop
				}
			}
		}()
	}
}

func (w *Worker) Stop() {
	// 停止工作池goroutine阻塞以结束goroutine
	log.Info("worker pool exit")
	w.endChan <- true
}
