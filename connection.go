/*
连接模块：
属性：
	1. 连接ID、连接本身Conn、连接超时
方法：
	1. 读取用户请求数据
	2. 解析数据
	3. 选择处理函数handler
	4. 返回响应数据

与server模块集成思路
1. 每来一个客户端请求则new一个连接对象，并添加到连接管理中
2. 连接对象开始启动服务
3. 连接对象异常或这推出，则删除连接管理中的保存的key-value

创建一个连接管理模块
在连接之前广播上线HOOK函数
连接之后广播下线的HOOK函数
*/
package CNet

import (
	"errors"
	log "github.com/sirupsen/logrus"
	"net"
)

type Context struct {
	Host string
	Request *DataPackage
	Response []byte
}

type Connection struct {
	server *CNet
	ID uint32
	Conn net.Conn
	Timeout int
	cChan chan []byte		// reader把数据放入管道cChan,writer从读取并返回客户端
}

func NewConnection(connectId uint32, conn net.Conn, cn *CNet) *Connection {
	connection := &Connection{
		server: cn,
		ID:       connectId,
		Conn:     conn,
		Timeout:  10,
		cChan: make(chan []byte, 0),
	}
	// 将初始化后的连接添加到连接管理模块（ConnManage）中
	cn.Hook.CallOnConn(connection)
	log.Infof("[Connection:%d] init", connectId)
	return connection
}

func (c *Connection) Start() {
	go c.reader()
	go c.writer()
}

func (c *Connection) reader() {
	defer func() {
		if err:=recover();err!=nil{
			c.Close(errors.New("reader error"))
		}
	}()
	defer log.Infof("[Connection:%d] read data", c.ID)
	t := NewTransfer(c.Conn)
	dataPackage, err := t.Read()
	if err != nil {
		c.Close(err)
		return
	}
	context := Context{
		Host:     c.Conn.RemoteAddr().String(),
		Request:  dataPackage,
		Response: nil,
	}
	// 调用钩子函数
	c.server.CallOnRequest(&context)

	// 交给路由处理
	err = c.server.Router.match(&context)
	if err != nil {
		c.Close(err)
		return
	}
	// 把响应数据写入管道中交给writer返回给客户端
	if len(context.Response) == 0 {
		c.Close(errors.New("no data to write"))
		return
	}

	// 调用钩子函数
	c.server.CallOnResponse(&context)

	// 将响应数据 发送到 保存 处理返回数据 的channel
	c.cChan<-context.Response
}

func (c *Connection) writer() {
	defer func() {
		if err:=recover();err!=nil{
			c.Close(errors.New("writer error"))
		}
	}()
	defer log.Infof("[Connection:%d] response data", c.ID)
	t := NewTransfer(c.Conn)
	data := <- c.cChan
	err := t.Write(DataPackage{
		ID:      c.ID,
		Len:     0,
		Content: data,
	})
	if err != nil {
		c.Close(err)
	}
}

func (c *Connection) Close(err error)  {
	defer log.Infof("[Connection:%d] closed error: %s", c.ID, err)
	c.server.CallOffConn(c)
	// 关闭连接conn
	_ = c.Conn.Close()
	// 关闭 TODO 通道是否需要关闭
	//close(c.cChan)
}
