package vm

import (
	. "github.com/LuaProject/api"
)

// R(A+1) := R(B); R(A) := R(B)[RK(C)]
func self(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1
	b += 1
	vm.Copy(b, a+1)
	vm.GetRK(c)
	vm.GetTable(b)
	vm.Replace(a)
}

// R(A) := closure(KPROTO[Bx])
//iABx模式 把Lua函数的子函数原型实例化为闭包
//A 寄存器地址 用来存放 闭包
//Bx 索引 子函数原型来自于当前函数原型的子函数原型表
func closure(i Instruction, vm LuaVM) {
	a, bx := i.ABx()
	a += 1
	vm.LoadProto(bx)
	vm.Replace(a)
}

// R(A), R(A+1), ..., R(A+B-2) = vararg
func vararg(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	if b != 1 {  // b==0 or b>1
		vm.LoadVararg(b - 1)
		_popResults(a, b, vm)
	}
}

// return R(A)(R(A+1), ... ,R(A+B-1))
func tailCall(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	// todo: optimize tail call!
	c := 0
	nArgs := _pushFuncAndArgs(a, b, vm)
	//在当前栈中创建子栈
	vm.Call(nArgs, c-1)
	//将栈顶值
	_popResults(a, c, vm)
}

// R(A), ... ,R(A+C-2) := R(A)(R(A+1), ... ,R(A+B-1))
//iABC模式 调用Lua函数
//A 指定被调用函数所在寄存器索引
//B 调用参数数量
//C 指定返回值的数量
func call(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1
	// println(":::"+ vm.StackToString())
	//将参数取出推到栈顶
	nArgs := _pushFuncAndArgs(a, b, vm)
	//调用call 生成子调用栈
	vm.Call(nArgs, c -1)
	_popResults(a, c, vm)
}

/**
将函数 和参数推入栈顶
但由于传入的参数数量不确定 所以需要checkstack确保栈空间足够
 */
func _pushFuncAndArgs(a, b int, vm LuaVM) (nArgs int) {
	//如果传入的参数大于b - 1个
	if b >= 1 {
		vm.CheckStack(b)
		//包含a 是函数 a+b
		for i := a; i < a + b; i++ {
			vm.PushValue(i)
		}
		return b - 1
	} else {
		_fixStack(a, vm)
		return vm.GetTop() - vm.RegisterCount() - 1
	}
}

func _fixStack(a int, vm LuaVM) {
	//取要 返回寄存器的索引
	x := int(vm.Tointeger(-1))
	// 将索引弹出
	vm.Pop(1)
	//检查有没有空间
	vm.CheckStack(x - a)
	for i := a; i < x; i++ {
		vm.PushValue(i)
	}
	vm.Rotate(vm.RegisterCount() + 1 , x-a)
}

func _popResults(a, c int, vm LuaVM) {
	//c 是 1返回值为空
	if c  ==  1 {
		// no results
		//C大于1返回值 c-1
	} else if c  > 1 {
		for i := a + c - 2; i >= a; i-- {
			vm.Replace(i)
		}
		//如果c -1 < 0 c < 1 返回所有参数
	} else {
		// leave results on stack
		vm.CheckStack(1)
		//记录返回栈的地址
		vm.PushInteger(int64(a))
	}
}

// return R(A), ... ,R(A+B-2)
//iABC模式 把存在连续多个寄存器里的值返回给主调函数
//A 第一个寄存器索引地址
//B 寄存器数量 操作数
func _return(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1

	if b == 1 {
		// no return values
	} else if b > 1 {
		// b-1 return values
		vm.CheckStack(b - 1)
		for i := a; i <= a + b - 2; i++ {
			vm.PushValue(i)
		}
	} else {
		_fixStack(a, vm)
	}
}