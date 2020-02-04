package state

import (
	"fmt"
	"github.com/LuaProject/api"
	"github.com/LuaProject/binchunk"
	"github.com/LuaProject/vm"
)

func(self *luaState) Load(chunk []byte,chunkName,mode string) int{
	//解析chunk函数
	proto := binchunk.Undump(chunk)
	//新建函数原型
	c := newLuaClosure(proto)
	//将加载完成的函数原型推入栈顶
	self.stack.push(c)
	if len(proto.Upvalues) >0{
		//获取全局注册表
		env := self.registry.get(api.LUA_RIDX_GLOBALS)
		//upvals设置成全局环境表
		c.upvals[0] = &upvalue{&env}
	}
	return 0
}
//第一次初始化时 默认将函数原型存在栈底
//nArgs 被调用函数参数数量 nResults返回参数数量(多退少补 nil)
func (self *luaState) Call(nArgs, nResults int){
	//从栈中取出被调用lua函数原型
	val := self.stack.get(-(nArgs +1))
	//判断是否是 lua函数类型
	c, ok := val.(*closure);
	if !ok {
		if mf := getMetafield(val,"__call",self);mf != nil{
			if c,ok = mf.(*closure) ;ok{
				self.stack.push(val)
				self.Insert(-(nArgs + 2))
				nArgs += 1
			}
		}
	}
	if ok{
		if c.proto !=nil {
			//调用lua函数
			self.callLuaClosure(nArgs,nResults, c)
		}else{
			//调用go函数
			self.callGoClosure(nArgs,nResults, c)
		}
		if c.proto == nil{
			return
		}
		if api.DEBUG == 0{
			//打印出函数的一些参数
			fmt.Printf("call %s<%d,%d>\n",c.proto.Source,
				c.proto.LineDefined, c.proto.LastLineDefined)
		}

	}else{
		//如果取出的不是 lua函数类型
		panic("not function!")
	}
}
//function f(a,... )
//	local b, c =
//	local t = {a,... }
//		print( a, b, c, #t, ...)
//end 可变参函数 内部才可以定义可变长参数
//function f ()
//return 1, 2, 3
//end
//如果期望返回值 是三四个 但是上面函数实际返回3个 那么剩余一个用nil补充
//如果 期望返回 2个 实际三个 那么 剩余一个丢弃
// function f () return 3, 2, 1 end
// function g() return 4, f() end
// function h(a, b, c, d) print(a, b, c, d) end
// h(4,f()) --> 4 3 2 1
// h(g()) --> 4 3 2 1
// print(table.unpack({4, f()})) --> 4 3 2 1
//函数如果接收到 其他函数的多个返回值那么会全部接受
func (self *luaState) callLuaClosure(nArgs, nResults int, c *closure) {
	//取到执行函数所需寄存器的数量大小 编译时确认
	nRegs := int(c.proto.MaxStackSize)
	//定义函数时的固定参数数量
	nParams := int(c.proto.NumParams)
	//是否vararg函数 可变参函数
	isVararg := c.proto.IsVararg == 1
	//根据所需寄存器数量创建栈空间 + 同时预留一定的寄存器  创建调用帧
	newStack := newLuaStack(nRegs+ api.LUA_MINSTACK, self)
	//闭包 和 调用帧联系起来
	newStack.closure = c
	// 调用子函数 要把要调用的参数和 函数原型取出
	//而这个子函数原型 和上面的c是一样的
	funcAndArgs := self.stack.popN(nArgs + 1)
	//将参数传入 新调用栈 0是函数原型 所以从1开始
	newStack.pushN(funcAndArgs[1:], nParams)
	//修改新栈的栈顶指针指向 最后一个寄存器
	//指向栈顶 上面是存放参数的
	newStack.top = nRegs

	//如果被调用的是 vararg函数 且传入参数的数量多于固定参数数量，需要把vararg参数记下来
	//存在调用帧里 多于的参数记录下来
	if nArgs > nParams && isVararg {
			newStack.varargs = funcAndArgs[nParams+1:]
	}
	//把新创建的帧推入调用栈顶,让它成为当前帧
	// run closure
	self.pushLuaStack(newStack)
	//执行背调函数指令 返回值值将会在栈顶部
	self.runLuaClosure()
	//将被调用函数的弹出 并且将返回值 返回 多退少补
	//此时 调用栈又变成了当前栈
	self.popLuaStack()
	//如果返回参数 大于0个 推出参数
	if nResults != 0 {
		//从子函数栈推出 nRegs 个参数
		results := newStack.popN(newStack.top - nRegs)
		//检查父函数栈是否有足够的空间存放 子函数的返回值
		self.stack.check(len(results))
		//推出的参数再推回栈中 后面再通过pop 推入到a指定寄存器
		self.stack.pushN(results, nResults)
	}
}

//循环读取指令集 直到遇到return
func (self *luaState) runLuaClosure() {
	for {
		//取出指令
		inst := vm.Instruction(self.Fetch())
		//执行指令
		inst.Execute(self)
		//如果遇到 返回指令 执行完毕 返回
		if inst.Opcode() == vm.OP_RETURN {
			break
		}
		if api.DEBUG ==0{
			fmt.Printf("[%02d] %s ", self.stack.pc+1, inst.OpName())
			printStack(*self)
		}


	}
}

func printStack(ls luaState) {
	top := ls.GetTop()
	for i := 1; i <= top; i++ {
		t := ls.Type(i)
		switch t {
		case api.LUA_TBOOLEAN:
			fmt.Printf("[%t]", ls.ToBoolean(i))
		case api.LUA_TNUMBER:
			fmt.Printf("[%g]", ls.ToNumber(i))
		case api.LUA_TSTRING:
			fmt.Printf("[%q]", ls.ToString(i))
		default: // other values
			fmt.Printf("[%s]", ls.TypeName(t))
		}
	}
	fmt.Println()
}

func (self *luaState) callGoClosure(nArgs, nResults int,c *closure){
	//创建一个新栈
	newStack := newLuaStack(nArgs + api.LUA_MINSTACK, self)
	//设置栈的原型函数
	newStack.closure = c
	//将参数从调用栈弹出
	args := self.stack.popN(nArgs)
	//将调用参数传入创建的新栈
	newStack.pushN(args,nArgs)
	//从栈中弹出 Go函数
	self.stack.pop()
	//将栈指针 指向当前栈
	self.pushLuaStack(newStack)
	//将栈传入调用函数
	r := c.goFunc(self)
	//当运算完后 将栈弹出
	self.popLuaStack()
	//如果返回参数 不等于0 将参数返回 调用栈
	if nResults != 0 {
		//将返回值 从栈顶取出
		resultes := newStack.popN(r)
		//检查调用栈 是否有足够的空间 放入返回值
		self.stack.check(len(resultes))
		//将返回值 压入调用栈
		self.stack.pushN(resultes,nResults)
	}
}