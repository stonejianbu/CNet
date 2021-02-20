/*
为了能够更快的处理用户的请求，又避免goroutine创建和销毁的性能消耗，不必每来一个连接就创建多个线程去处理请求，
而是提前创建工作池，然后每个连接接收到数据就将其发送到channel,工作池的goroutine去处理用户请求，可给予用户选择开启或关闭工作池

根据CPU核心数来初始化goroutine,并把结果传递给Context结构体的Context.Request.XX

消费者集群

提供给外部的接口是怎么样的？
1.调用方法，添加到指定管道中，内部调用该实例的方法XX.Start()
*/
package CNet

import (
	log "github.com/sirupsen/logrus"
)

var MaxConnNum = 5

type Worker struct {
	conChan chan *Connection
	endChan chan bool
}

func NewWorker() *Worker {
	return &Worker{
		conChan: make(chan *Connection, MaxConnNum),
		endChan: make(chan bool),
	}
}

func (w *Worker) add(obj *Connection) {
	w.conChan<-obj
}

func (w *Worker) Start() {
	for i:=0;i< MaxConnNum;i++ {
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
