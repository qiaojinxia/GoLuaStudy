package api
type LuaVM interface {
	LuaState
	PC() int //返回当前PC
	AddPc(n int) //修改PC(用于实现跳转指令)
	Fetch() uint32 //取出当前指令；将PC指向吓一条指令
	GetConst(idx int)//将制定常量推入栈顶
	GetRK(rk int)//将指定常量或栈值推入栈顶
	RegisterCount() int
	LoadVararg(n int)
	LoadProto(idx int)
}
