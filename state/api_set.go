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
	self.setTable(t,k,v,false)
}
/**
将 k v 存入表中
 */
func (self *luaState) setTable(t,k,v luaValue,raw bool){
	if tbl,ok :=t.(*luaTable); ok{
		if raw || tbl.get(k) != nil || !tbl.hasMetafield("__newindex"){
			tbl.put(k,v)
			return
		}
		if !raw{
			if mf := getMetafield(t,"__newindex",self);mf != nil{
				switch x :=mf.(type) {
				case *luaTable:
					self.setTable(x,k,v,false)
				case *closure:
					self.stack.push(mf)
					self.stack.push(t)
					self.stack.push(k)
					self.stack.push(v)
					self.Call(3,0)
					return

				}
			}
		}
		panic("index error!")
	}
	panic("not a table")
}

/**
给定 表在栈中 的索引 从栈顶取出value 然后更新到 表中 key为传入的 k
 */
func (self *luaState) SetField(idx int, k string) {
	t := self.stack.get(idx)
	v := self.stack.pop()
	self.setTable(t,k,v,false)
}
/**
给定表在栈中的位置 从栈中取出value 值赋值给 指定索引 i为参数
 */
func (self *luaState) SetI(idx int,i int64){
	t := self.stack.get(idx)
	v := self.stack.pop()
	self.setTable(t, i, v,false)
}
//从栈顶取 出 函数 放入注册表
func (self *luaState) SetGlobal(name string){
	t := self.registry.get(api.LUA_RIDX_GLOBALS)
	v := self.stack.pop()
	self.setTable(t, name, v,false)

}
//将Go函数 的名字和函数存放到 注册表
func(self *luaState) Register(name string,f api.GoFunction){
	self.PushGoFunction(f)
	self.SetGlobal(name)
}
