package state

import "github.com/LuaProject/api"

type luaStack struct{
	slots []luaValue//栈存储数组
	top int//栈顶索引
	prev *luaStack
	//存放chunk数据
	closure *closure
	varargs []luaValue
	//指向执行到第几条指令
	pc int
	//指向luastate
	state *luaState

	openuvs map[int]*upvalue
}
func newLuaStack(size int,state *luaState) *luaStack{
	return &luaStack{
		slots: make([]luaValue,size),
		top:   0,
		state:state,
	}
}
//检查栈的空闲空间是否还可以 容纳 至少n个值
func (self *luaStack) check(n int) {
	//总容量 减去 占用的容量
	free := len(self.slots) - self.top
	//循环 n - free 次 如果n < free表示空间够用 不需要添加
	for i := free ; i < n ;i++ {
		self.slots =append(self.slots,nil)
	}
}
//push() 将值推入栈顶 ,如果栈溢出,报错处理
func (self *luaStack) push(val luaValue){
	if self.top == len(self.slots) {
		panic("stack overflow!")
	}
	self.slots[self.top] = val
	self.top ++
}
//pop() 将值推出栈
func (self *luaStack) pop() luaValue{
	//如果栈 元素小于1 则无法出栈
	if self.top < 1{
		panic("stack underflow")
	}
	//栈顶高度减一
	self.top --
	//栈顶元素
	val := self.slots[self.top]
	//栈顶元素 置nil
	self.slots[self.top] = nil
	//返回pop值
	return val
}
//absIndex() 方法把索引转换成绝对索引
func (self *luaStack) absIndex(idx int) int{
	if idx >=0 || idx <= api.LUA_REGISTRYINDEX{
		return idx
	}
	//大于零 直接返回
	if idx >= 0 {
		return idx
	}
	//以负数为索引是返回 top = -1
	return idx + self.top + 1
}
//isValid()验证索引是否有效
func (self *luaStack) isValid(idx int) bool{
	if idx < api.LUA_REGISTRYINDEX{
		uvIdx := api.LUA_REGISTRYINDEX - idx -1
		c := self.closure
		return c!= nil && uvIdx < len(c.upvals)
	}
	if idx == api.LUA_REGISTRYINDEX {
		return true
	}
	absIdex :=  self.absIndex(idx)
	return absIdex > 0 && absIdex <= self.top
}
//按索引获取值 不会改变栈内容
func (self *luaStack) get(idx int) luaValue {
	if idx < api.LUA_REGISTRYINDEX {
		uvIdx := api.LUA_REGISTRYINDEX -idx -1
		c := self.closure
		if c == nil || uvIdx >= len(c.upvals){
			return nil
		}
		return *(c.upvals[uvIdx].val)
	}
	if idx == api.LUA_REGISTRYINDEX{
		return self.state.registry
	}
	//对索引进行 绝对值如果为负
	absIndex := self.absIndex(idx)
	//如果索引没有超出范围
	if absIndex > 0 && absIndex <= self.top {
		//返回值
		return self.slots[absIndex -1]
	}
	//返回nil
	return nil
}
//set()方法 根据索引往栈里写入值 如果索引无效 抛出异常 终止程序
func (self *luaStack) set(idx int,val luaValue){
	if idx < api.LUA_REGISTRYINDEX{
		uvIdx := api.LUA_REGISTRYINDEX -idx -1
		c := self.closure
		if c != nil && uvIdx < len(c.upvals){
			*(c.upvals[uvIdx].val) =val
		}
		return
	}
	if idx == api.LUA_REGISTRYINDEX {
		if v,ok:=val.(*luaTable);ok{
			self.state.registry = v
			return
		}else{
			panic("set not register table")
		}
	}
	absIdx := self.absIndex(idx)
	if absIdx > 0 && absIdx <= self.top{
		self.slots[absIdx -1] = val
		return
	}
	panic("invalid index!")
}
//pop n个值
func (self *luaStack) popN(n int) []luaValue{
	vals := make([]luaValue,n)
	for i:=n-1;i>=0;i--{
		vals[i] = self.pop()
	}
	return vals
}
//推入 n 个 值 如果 vals 长度 如果n的长度比 vals长那么剩余的就push nil
func (self *luaStack) pushN(vals []luaValue,n int){
	//假设 实际返回 func a(){ return 1 2 3} 但是 b,c,d,e=a() 那么e 是nil 所谓的多退少补
	nVals := len(vals)
	if n < 0 { n = nVals}
	for i := 0; i < n; i++ {
		if i < nVals {
			self.push(vals[i])
		}else {
			self.push(nil)
		}
	}


}