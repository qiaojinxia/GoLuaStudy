package state

import "github.com/LuaProject/api"

/**
给定 表索引 推出 栈 前两个值
value 和 key 作为键值对
存入表中
 */
func (self *luaState) SetTable(idx int){
	t := self.stack.get(idx)
	v :=self.stack.pop()
	k := self.stack.pop()
	self.setTable(t,k,v)
}
/**
将 k v 存入表中
 */
func (self *luaState) setTable(t,k,v luaValue){
	if tbl,ok :=t.(*luaTable); ok{
		tbl.put(k,v)
		return
	}
	panic("not a table")
}

/**
给定 表在栈中 的索引 从栈顶取出value 然后更新到 表中 key为传入的 k
 */
func (self *luaState) SetField(idx int, k string) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	self.setTable(t,k,v)
}
/**
给定表在栈中的位置 从栈中取出value 值赋值给 指定索引 i为参数
 */
func (self *luaState) SetI(idx int,i int64){
	t := self.stack.get(idx)
	v := self.stack.pop()
	self.setTable(t, i, v)
}
//从栈顶取 出 函数 放入注册表
func (self *luaState) SetGlobal(name string){
	t := self.registry.get(api.LUA_RIDX_GLOBALS)
	v := self.stack.pop()
	self.setTable(t, name, v)

}
func(self *luaState) Register(name string,f api.GoFunction){
	self.PushGoFunction(f)
	self.SetGlobal(name)
}