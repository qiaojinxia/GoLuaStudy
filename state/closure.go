package state

import (
	. "github.com/LuaProject/api"
	"github.com/LuaProject/binchunk"
)

type upvalue struct {
	val *luaValue
}


type closure struct {
	proto *binchunk.Prototype
	goFunc GoFunction
	upvals []*upvalue
}

func newLuaClosure(proto *binchunk.Prototype) *closure{
	c := &closure{proto:  proto}
	if nUpvals := len(proto.Upvalues);nUpvals > 0{
		c.upvals = make([]*upvalue,nUpvals)
	}
	return c
}
//初始化Go闭包函数  go函数 闭包
func newGoClosure(f GoFunction, nUpvals int) *closure{
	c := &closure{goFunc:f}
	if nUpvals > 0 {
		c.upvals =make([]*upvalue,nUpvals)
	}
	return c
	return &closure{goFunc:f}
}
