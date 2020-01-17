package state

import . "github.com/LuaProject/api"

/**
创建一个表结构将引用指针推入栈中
 */
func (self *luaState) CreateTable(nArr,nRec int ){
	t := newLuaTable(nArr,nRec)
	self.stack.push(t)
}
/**
初始化一个 0,0大小的表
 */
func (self *luaState) NewTable(){
	self.CreateTable(0,0)
}

/**
获取栈顶key 的类型
 */
func (self *luaState) GetTable(idx int) LuaType{
	//从索引处取 table
	t := self.stack.get(idx)
	//将键从栈顶取出
	k :=self.stack.pop()
	//key 在 表里 对应 value 的返回类型
	return self.getTable(t,k)
}

/**
获取指定表 key 的 value的类型 并将值推入栈顶
 */
func (self *luaState) getTable(t,k luaValue) LuaType{
	//如果 是luavalu类型
	if tbl, ok := t.(*luaTable); ok{
		//从表里取出key对应值
		v := tbl.get(k)
		//将 值推入栈顶
		self.stack.push(v)
		//返回值的类型
		return typeOf(v)
	}
	panic("not a table!")
}
/**
获取指定key 的 value类型
 */
func (self *luaState) GetField(idx int,k string) LuaType{
	//从索引处 获得表
	t := self.stack.get(idx)
	return self.getTable(t,k)
}
/**
传入 表在栈中位置 和 数组索引 返回数组类型
 */
func (self *luaState) GetI(idx int,i int64) LuaType{
	t := self.stack.get(idx)
	return self.getTable(t,i)
}

