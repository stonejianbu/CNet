package CNet

import (
	"fmt"
	"log"
	"net"
	"testing"
)

// 测试服务端
func TestNewCNet(t *testing.T) {
	// 启动服务
	cn := NewCNet(":5678")
	cn.AddHandler(100, func(context *Context) error {
		req := context.Request.Content
		context.Response = []byte("hello world, "+string(req))
		return nil
	})
	// 通过钩子去注册方法
	cn.Hook.RegisterOnRequest(func(ctx *Context) {
		fmt.Println("已经开始发起请求了...")
	})
	// 通过钩子函数去注册属性
	go cn.Serve()

	// 启动客户端
	for i:=0;i<100;i++ {
		go Client("127.0.0.1:5678", func(conn net.Conn) error {
			t := NewTransfer(conn)
			err := t.Write(DataPackage{
				ID:      100,
				Len:     0,
				Content: []byte("stonejianbu"),
			})
			if err != nil {
				log.Println("[client]",err)
				return err
			}
			data, err := t.Read()
			if err != nil {
				log.Println(err)
				return err
			}
			log.Printf("[client] recv Len = %d, ID = %d, data =%v \n", data.Len,data.ID,string(data.Content))
			return nil
		})
	}
	// 阻塞等待
	select{}
}

func Client(address string, handle func(conn net.Conn) error)  {
	conn, err := net.Dial("tcp", address)
	if err != nil {
		log.Fatalln("net dial error:", err)
	}
	err = handle(conn)
	if err != nil {
		log.Fatalln("client handle error:",err)
	}
}