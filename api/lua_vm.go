package api
type LuaVM interface {
	LuaState
	PC() int //返回当前PC
	AddPc(n int) //修改PC(用于实现跳转指令)
	Fetch() uint32 //取出当前指令；将PC指向吓一条指令
	GetConst(idx int)//将制定常量推入栈顶
	GetRK(rk int)//将指定常量或栈值推入栈顶
	RegisterCount() int//当前Lua函数所操作的寄存器数量
	LoadVararg(n int)//传递给当前Lua函数的边长参数推入栈顶(多退少补)
	LoadProto(idx int)//把当前Lua函数的子函数的实现原型实例化为闭包推入栈顶
	CloseUpvalues(a int)//闭合闭包函数

}
