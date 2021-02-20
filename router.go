package cnet

import (
	"errors"
	"fmt"
)

type Router struct {
	Handlers map[uint32] func(context *Context) error
}

func NewRouter() *Router {
	return &Router{Handlers: map[uint32]func(context *Context) error{}}
}

// 添加handler
func (r *Router) AddHandler(hid uint32, handler func(context *Context) error)  {
	r.Handlers[hid] = handler
}

// 根据请求数据中的DataId来匹配相应handler
func (r *Router) match(context *Context) error {
	for k, v := range r.Handlers {
		if k == context.Request.ID {
			if err := v(context);err!=nil{
				return errors.New(fmt.Sprintf("handle error: %s\n",err))
			}
			return nil
		}
	}
	return errors.New("no match handler")
}
