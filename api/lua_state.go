package api

type LuaType = int
type ArithOp = int
type CompareOp = int

type GoFunction func(LuaState) int
type LuaState interface {
	GetTop() int
	AbsIndex(idx int) int
	CheckStack(n int) bool
	Pop(n int)
	Copy(fromidx, toIdx int)
	PushValue(idx int)
	Replace(idx int)
	Insert(idx int)
	Remove(idx int)
	Rotate(idx,n int)
	SetTop(idx int)
	/* access functions(stack -> Go) */
	TypeName(tp LuaType) string
	Type(idx int) LuaType
	IsNone(idx int)bool
	IsNil(idx int)bool
	IsNoneOrNil(idx int) bool
	IsBoolean(Idx int) bool
	IsInteger(idx int)bool
	IsNumber(idx int) bool
	IsString(idx int) bool
	IsTable(idx int) bool
	ToBoolean(idx int) bool
	Tointeger(idx int) int64
	ToIntegerX (idx int) (int64, bool)
	ToNumber(idx int) float64
	ToNumberX(idx int) (float64 , bool)
	ToString(idx int) string
	ToStringX(idx int) (string, bool)

	/* push f unctions (Go -> stack)*/
	PushNil()
	PushBoolean(b bool)
	PushInteger(n int64)
	PushNumber(n float64)
	PushString(s string)
	Arith(op ArithOp)
	Compare(idx1,idx2 int,op CompareOp) bool
	Len(idx int)
	Concat(n int)
	/*table func*/
	NewTable()
	CreateTable(nArr,nRec int)
	GetTable(idx int)LuaType
	GetField(idx int,k string) LuaType
	GetI(idx int,i int64) LuaType
	SetTable(idx int)
	SetField(idx int,k string)
	SetI(idx int, n int64)
	//加载二进制chunk,把主函数原型实例化为闭包推入栈顶
	//三个参数 chunck数据  chunck名字 加载模式 b二进制编译好的文件 bt源文件或者二进制
	Load(chunk []byte, chunkName, mode string) int
	Call(nArgs, nResults int)

	PushGoFunction(f GoFunction)
	IsGoFunction(idx int)bool
	ToGoFunction(idx int) GoFunction

	PushGlobalTable()
	GetGlobal(narne string) LuaType
	SetGlobal(name string)
	Register(name string, f GoFunction)

	PushGoClosure(f GoFunction,n int)

	GetMetatable(idx int) bool  //获取原方法
	SetMetatable(idx int) //设定原方法
	RawLen(idx int) uint //
	RawEqual(idxl , idx2 int)bool
	RawGet(idx int)LuaType
	RawSet(idx int )
	RawGetI( idx int, i int64)LuaType
	RawSetI(idx int, i int64)
}
