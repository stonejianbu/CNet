package cnet

import (
	log "github.com/sirupsen/logrus"
	"net"
)

func init()  {
	// init log
	log.SetFormatter(&log.JSONFormatter{})
	log.SetLevel(log.DebugLevel)
}

var Network = "tcp"
var Name 	= "CNet server"

type CNet struct {
	Network string // 网络协议，默认tcp
	Name    string // 服务名称
	Address string // 服务监听地址
	MaxConn int    // 最大连接数量
	*Router        // 路由管理
	*Hook          // 钩子管理：允许用户注册钩子函数
}

func NewCNet(address string) *CNet {
	return &CNet{
		Network: Network,
		Name:    Name,
		Address: address,
		Router:  NewRouter(),
		Hook:    NewHook(),
		MaxConn: 100,
	}
}

func (cn *CNet) Serve() {
	log.Debugf("[server] %s start to listening", cn.Name)
	// 开启服务监听
	listener, err := net.Listen(cn.Network, cn.Address)
	if err != nil {
		log.Fatal("[server] net listen error", err)
		return
	}
	defer listener.Close()

	// 初始化worker，避免重复创建和销毁goroutine，这里通过提前创建一批goroutine用于处理到来的请求
	worker := NewWorker(cn.MaxConn)
	go worker.Start()
	defer worker.Stop()

	// 初始化连接ID
	var connectId uint32 = 0
	log.Debug("[server] waiting for connections in a loop")
	for {
		// 阻塞等待客户端请求
		conn, err := listener.Accept()
		if err != nil {
			// 出现异常则跳过该连接处理，继续等待新的客户端连接
			log.Error("[server] listener accept error", err)
			_ = conn.Close()
			continue
		}
		// 将处理新连接的业务方法和conn进行绑定 得到我们的连接模块
		// 为了让每个连接conn你能够访问到服务资源（如router,connManage,Hook）,于是将cn传入到connection中
		worker.add(NewConnection(connectId, conn, cn))
		connectId++
	}
}

func (cn *CNet) Stop()  {
	log.Debugf("[server] %s had stop", cn.Name)
}


