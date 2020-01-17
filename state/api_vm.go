package state
func (self *luaState) PC() int{
	return self.stack.pc
}
func (self *luaState) AddPc(n int) {
	self.stack.pc += n
}
//从函数原型指令表取出当前指令 把PC加1
func (self *luaState) Fetch() uint32{
	i := self.stack.closure.proto.Code[self.stack.pc]
	self.stack.pc ++
	return i
}
//从原型常量表取出一个常量值 推入栈顶
func (self *luaState) GetConst(idx int){
	c := self.stack.closure.proto.Constants[idx]
	self.stack.push(c)
}

func (self *luaState) GetRK(rk int){
	//如果最高位是1 那么从常量池去取
	if rk > 0xFF{
		//取不含1的低8位
		self.GetConst(rk & 0xFF)
	}else{
		//Lua API 栈索引 从1开始 要+1
		self.PushValue(rk +1)
	}
}

func (self *luaState) RegisterCount() int{
	return int(self.stack.closure.proto.MaxStackSize)
}

func(self *luaState) LoadVararg(n int){
	if n < 0 {
		n = len(self.stack.varargs)
	}
	self.stack.check(n)
	self.stack.pushN(self.stack.varargs,n)
}

func (self *luaState) LoadProto(idx int){
	proto := self.stack.closure.proto.Protos[idx]
	closure := newLuaClosure(proto)
	self.stack.push(closure)
}