package state

import . "github.com/LuaProject/api"

type luaState struct {
	registry *luaTable //注册表
	stack *luaStack
}

//初始化栈
func New() *luaState {
	registry := newLuaTable(0,0)
	registry.put(LUA_RIDX_GLOBALS, newLuaTable( 0,0))
	ls := &luaState{registry:registry}
	ls.pushLuaStack(newLuaStack(LUA_MINSTACK,ls))
	return ls
}
//获取栈顶索引
func (self *luaState) GetTop() int{
	return self.stack.top
}
//索引转换为绝对值索引
func (self *luaState) AbsIndex(idx int) int{
	return self.stack.absIndex(idx)
}
//检查栈是否 足够n个空间不足则添加
func (self *luaState)CheckStack(n int) bool{
	self.stack.check(n)
	return true
}

//复制栈的值 先get 再 set
func (self *luaState) Copy(fromIdx ,toIdx int) {
	val:= self.stack.get(fromIdx)
	self.stack.set(toIdx,val)
}
//pushValue 将制定栈位置值推入栈顶
func (self *luaState) PushValue(idx int){
	val := self.stack.get(idx)
	self.stack.push(val)
}
//弹出栈顶值 后 再插入 覆盖到栈制定位置
func (self *luaState) Replace(idx int){
	val := self.stack.pop()
	self.stack.set(idx,val)
}
//插入
func (self *luaState) Insert(idx int){
	self.Rotate(idx,1)
}
//删除制定索引 索引上面的值下移一格
func (self *luaState) Remove(idx int){
	self.Rotate(idx,-1)
	self.Pop(1)
}
//
func (self *luaState) Rotate(idx,n int) {
	//获取栈顶元素
	t := self.stack.top -1
	//获取指定索引元素
	p := self.stack.absIndex(idx) -1
	var m int
	if n >= 0 {
		m = t - n
	}else {
		m = p - n - 1
	}
	//进行三次旋转
	//说明 https://blog.csdn.net/weixin_41315492/article/details/103495924
	self.stack.reverse(p,m)
	self.stack.reverse(m+1,t)
	self.stack.reverse(p,t)

}

func (self *luaStack) reverse(from,to int){
	slots := self.slots
	for from < to {
		slots[from],slots[to] = slots[to],slots[from]
		from ++
		to --
	}
}
//如果栈顶 top 大于 setTop(n) name 就相当于弹出 n - top 指定0 时清空栈
//如果 top < n 相当于 压入 n - top 个nill 入栈
func (self *luaState) SetTop(idx int){
	newTop := self.stack.absIndex(idx)
	if newTop < 0 {
		panic("statck underflow!")
	}
	n := self.stack.top - newTop
	if n > 0 {
		for i:=0 ;i< n;i++{
			self.stack.pop()
		}
	}else {
		for i:=0 ;i< -n;i++{
			self.stack.push(nil)
		}
	}
}

func (self *luaState) Pop(n int){
	// 传入负索引
	self.SetTop(-n-1)
}
/**
将新的实例化栈的指针推入
并将当前luaState 的 stack 设置为新推入的stack的前一个栈
形成链表
 */
func (self *luaState) pushLuaStack(stack *luaStack) {
	stack.prev = self.stack
	self.stack = stack
}
/**
从链表中推出一个栈
 */
func (self *luaState) popLuaStack(){
	stack := self.stack
	self.stack = stack.prev
	stack.prev = nil
}

