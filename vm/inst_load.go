package vm
import . "github.com/LuaProject/api"

//用于给连续n个寄存器放置nil值 iABC模式 起始操作数 由A指定 B指定寄存器数量
func loadNil(i Instruction,vm LuaVM){
	a,b,_ := i.ABC()
	//栈API 从1开始 所以+1
	a += 1
	//栈顶推入 nil 用于拷贝复制
	vm.PushNil()
	//从a 开始 到 a+b 的连续空间复制 nil
	for i := a ; i <= a + b;i++{
		//拷贝栈顶nil
		vm.Copy(-1,i)
	}
	//把 栈顶的 nil 出栈
	vm.Pop(1)
}

//IABC模式 给单个寄存器设置布尔值  寄存器索引右操作数A指定 布尔值右寄存器B指定(0 false 1 true)
//如果 C非0 则 跳过下一条指令
// R(A) := (bool)B; if(C) pc++
func loadBool(i Instruction,vm LuaVM){
	//读取IABC下制定操作数
	a ,b ,c := i.ABC()
	//栈API 从1开始 所以+1
	a += 1
	//1为 true 0位false
	vm.PushBoolean(b!= 0)
	//弹出栈顶的值 替换掉a指向栈的值
	vm.Replace(a)
	//如果C不等于0 寄存器跳过下一条指令
	if c != 0 {
		vm.AddPc(1)
	}
}

//iABx模式 将常量表里的某个常量加载到指定寄存器
//寄存器索引有A指定 常量表索引由操作数Bx指定
func loadK(i Instruction,vm LuaVM){
	a,bx := i.ABx()
	//栈API 从1开始 所以+1
	a += 1
	//获取制定常量加入栈顶
	vm.GetConst(bx)
	//取出栈顶值 替换掉a索引指向栈的值
	vm.Replace(a)
}
//当常量表大于 262143时 可以使用这个方法来获取常量池 替代上面loadK()
func loadkx(i Instruction,vm LuaVM){
	//这里和上面不同只是取值 只是去取寄存器地址
	a ,_ := i.ABx()
	a += 1
	//上面一个方法 最多能表示262143个数 使用这个方法能够索引 到 最大 索引为 67108864
	ax := Instruction(vm.Fetch()).Ax()
	//获取指定常量加入栈顶
	vm.GetConst(ax)
	//取出栈顶值 替换掉a索引指向栈的值
	vm.Replace(a)
}
