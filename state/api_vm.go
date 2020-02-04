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

//需要寄存器数量 在newstack() 从chunk中获取
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
//将子函数原型表初始化 推入栈顶
func (self *luaState) LoadProto(idx int){
	stack := self.stack
	subProto := self.stack.closure.proto.Protos[idx]
	closure := newLuaClosure(subProto)
	stack.push(closure)
	//编译器 事先已经 把要引用的外部变量 的索引idx写入upvalues
	//这里 要理解的是每次 闭包加载指令 会得到当前函数的子函数并初始化成闭包
	//同时 stack.openuvs 这个属性表示自己需要的局部变量 实现编译器会在upvalues里记录索引
	//然后通过外部函数得到这个值
	for i,uvInfo := range subProto.Upvalues{
		uvIdx := int(uvInfo.Idx)
		//如果得到的是当前函数的变量
		if uvInfo.Instack == 1{
			//如果当前变量 的局部变量表为空 防止某些情况下 未初始化 下面查找语句报错
			if stack.openuvs == nil{
				stack.openuvs = map[int]*upvalue{}
			}
			//如果找到了闭包 那么就把他放进 子函数闭包里
			// 假设当前父函数 有2 个upvalue  0 = c 1 = a  子函数闭包 如果 要取  0 那么编译器实现生成 idx=0
			if openuv,found := stack.openuvs[uvIdx];found{
				closure.upvals[i] = openuv
			}else{
				//没找到的话 从栈中查找
				closure.upvals[i] = &upvalue{&stack.slots[uvIdx]}
				//由于编译器指明了 这个是当前函数的局部变量所以需要放入局部变量表
				stack.openuvs[uvIdx] = closure.upvals[i]
			}
			//如果 得到的不是当前函数的变量 但是被当前函数所捕获到的变量 那么就从 闭包函数中查找
		}else{
			closure.upvals[i] = stack.closure.upvals[uvIdx]
		}
	}
}
func (self *luaState)CloseUpvalues(a int){
	////这个语句不知道
	//for i, openuv := range self.stack.openuvs {
	//	if i >= a-1 {
	//		//因为要推出函数作用域 所以复制一份局部变量
	//		// 将指针指向新的地址
	//		val := *openuv.val
	//		openuv.val = &val
	//		delete(self.stack.openuvs, i)
	//	}
	//}
}