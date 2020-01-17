package state

import . "github.com/LuaProject/api"

func (self *luaState) PushNil(){ self.stack.push(nil)}//推入空指针
func (self *luaState) PushBoolean(b bool) {self.stack.push(b)}//推入布尔值
func (self *luaState) PushInteger(n int64){self.stack.push(n)}//推入整数类型
func (self *luaState) PushNumber(n float64){self.stack.push(n)}//推入浮点数
func (self *luaState) PushString(s string){self.stack.push(s)}//推入字符串类型
func (self *luaState) PushGoFunction(f GoFunction){self.stack.push(newGoClosure(f))}//推入Go闭包函数

func (self *luaState) PushGlobalTable(){
	global := self.registry.get(LUA_RIDX_GLOBALS)
	self.stack.push(global)
}

func (self *luaState) GetGlobal(name string) LuaType{
	t := self.registry.get(LUA_RIDX_GLOBALS)
	return self.getTable(t,name)
}

