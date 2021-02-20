/*
允许用户注册以下4个钩子函数和属性，由框架内部在合适时机调用，具体调用时机查看相应函数代码注释
1. RegisterOnConn(hook func(conn *Connection))
2. RegisterOnRequest(hook func(ctx *Context))
3. RegisterOnResponse(hook func(ctx *Context))
4. RegisterOffConn(hook func(conn *Connection))
*/
// https://www.bilibili.com/video/BV1wE411d7th?p=47&spm_id_from=pageDriver
package CNet

import log "github.com/sirupsen/logrus"

type Hook struct {
	onConn     func(conn *Connection)
	onRequest  func(ctx *Context)
	onResponse func(ctx *Context)
	offConn    func(conn *Connection)
}

func NewHook() *Hook {
	return &Hook{}
}

// 钩子调用时机：已经连接，未解析请求数据调用的函数，可访问到connection实例
func (h *Hook) RegisterOnConn(hook func(conn *Connection)) {
	h.onConn = hook
	log.Info("[Hook] Registered onConn function")
}

// 钩子调用时机：已经连接，已解析请求数据时调用的函数,可访问到Context实例
func (h *Hook) RegisterOnRequest(hook func(ctx *Context)) {
	h.onRequest = hook
	log.Info("[Hook] Registered onRequest function")
}

// 钩子调用时机：已经连接，已解析请求数据，已备好响应数据，可访问到Context实例
func (h *Hook) RegisterOnResponse(hook func(ctx *Context)) {
	h.onResponse = hook
	log.Info("[Hook] Registered onResponse function")
}

// 钩子调用时机：断开连接，可访问到connection实例，通过RegisterOffConn注册
func (h Hook) RegisterOffConn(hook func(conn *Connection)) {
	h.offConn = hook
	log.Info("[Hook] Registered offConn function")
}

// 封装调用钩子函数
func (h *Hook) CallOnConn(conn *Connection)  {
	if h.onConn != nil {
		log.Debug("[Hook] calling onConn")
		h.onConn(conn)
	}
}

func (h *Hook) CallOnRequest(ctx *Context)  {
	if h.onRequest != nil {
		log.Debug("[Hook] calling onRequest")
		h.onRequest(ctx)
	}
}

func (h *Hook) CallOnResponse(ctx *Context)  {
	if h.onResponse != nil {
		log.Debug("[Hook] calling onResponse")
		h.onResponse(ctx)
	}
}

func (h *Hook) CallOffConn(conn *Connection)  {
	if h.offConn != nil {
		log.Debug("[Hook] calling offConn")
		h.offConn(conn)
	}
}
