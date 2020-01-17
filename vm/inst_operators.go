package vm

import . "github.com/LuaProject/api"

//lua 语言里 8个算术运算符 和 6个按位运算符 对应 Lua虚拟机指令集里的14条 指令
//二元 运算符指令 a 和 c 运算完的值 存储到 a操作数指定的值
func _binaryArith(i Instruction,vm LuaVM,op ArithOp){
	a, b, c := i.ABC()
	a += 1
	//从常量表 或者栈 里 将 值推入栈顶进行运算
	vm.GetRK(c)
	vm.GetRK(b)
	//将栈顶2个元素 做 op运算
	vm.Arith(op)
	//将栈顶值 推出 赋值给 a指定的寄存器
	vm.Replace(a)
}
//一元算术运算指令 操作数A 用于存放结果到指定寄存器 操作数B 进行运算
func _unaryArith(i Instruction,vm LuaVM,op ArithOp){
	a, b , _ := i.ABC()
	//栈api 从1开始 加一
	a += 1; b += 1
	//将需要操作的数字推入栈顶
	vm.PushValue(a)
	//进行操作
	vm.Arith(op)
	//将栈顶值 推出 赋值给 a指定的寄存器
	vm.Replace(b)

}

/* arith */

func add(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPADD) }  // +
func sub(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPSUB) }  // -
func mul(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPMUL) }  // *
func mod(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPMOD) }  // %
func pow(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPPOW) }  // ^
func div(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPDIV) }  // /
func idiv(i Instruction, vm LuaVM) { _binaryArith(i, vm, LUA_OPIDIV) } // //
func band(i Instruction, vm LuaVM) { _binaryArith(i, vm, LUA_OPBAND) } // &
func bor(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPBOR) }  // |
func bxor(i Instruction, vm LuaVM) { _binaryArith(i, vm, LUA_OPBXOR) } // ~
func shl(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPSHL) }  // <<
func shr(i Instruction, vm LuaVM)  { _binaryArith(i, vm, LUA_OPSHR) }  // >>
func unm(i Instruction, vm LuaVM)  { _unaryArith(i, vm, LUA_OPUNM) }   // -
func bnot(i Instruction, vm LuaVM) { _unaryArith(i, vm, LUA_OPBNOT) }  // ~
//iAbc 模式  取b的长度 存入 a 指定寄存器 *a = len(b)
func _len(i Instruction,vm LuaVM){
	a, b, _ :=i.ABC()
	a += 1;b += 1
	vm.Len(b)
	vm.Replace(a)
}

//iABC模式 将连续n个寄存器(起止索引分别由操作数B和C指定)里的值进行拼接 放入A操作数所制定的寄存器
//* a += [stack[i] for i：=A;i< C ;i++]
func concat(i Instruction, vm LuaVM){
	a,b,c := i.ABC()
	//栈索引从1开始
	a +=1;b+=1;c+=1
	//连续的n个
	n := c - b + 1
	//查看是否超出栈索引
	vm.CheckStack(n)
	//推出栈索引
	for i := b; i <= c; i++ {
		vm.PushValue(i)
	}
	//从栈中取出n个栈顶元素拼接
	vm.Concat(n)
	//将栈顶值 推出 赋值给 a指定的寄存器
	vm.Replace(a)
}
//比较2个数的操作 先调用getRk 方法推入要比较的值入栈顶 然后调用Comapare()方法 执行比较运算 比较结果和操作数A一致则PC加1
//然后pop清理栈顶的2个值
func _compare(i Instruction,vm LuaVM,op CompareOp){
	a,b,c := i.ABC()
	vm.GetRK(b)
	vm.GetRK(c)
	//比较栈顶元素是否 b op c == a
	if vm.Compare(-2,-1,op) != (a != 0){
		vm.AddPc(1)
	}
	vm.Pop(2)
}
/* compare */
func eq(i Instruction, vm LuaVM) { _compare(i, vm, LUA_OPEQ) } // ==
func lt(i Instruction, vm LuaVM) { _compare(i, vm, LUA_OPLT) } // <
func le(i Instruction, vm LuaVM) { _compare(i, vm, LUA_OPLE) } // <=

// R(A) := not R(B)
//对 b进行一元操作 存入a操作数地址寄存器
func not(i Instruction, vm LuaVM) {
	a, b, _ := i.ABC()
	a += 1
	b += 1
	vm.PushBoolean(!vm.ToBoolean(b))
	vm.Replace(a)
}

// if (R(B) <=> C) then R(A) := R(B) else pc++
//判断B操作数寄存器值 是否和C一致 布尔值  一致的话B寄存器 的值 替换掉A寄存器的值 不一致 直接跳过
func testSet(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1
	b += 1

	if vm.ToBoolean(b) == (c != 0) {
		vm.Copy(b, a)
	} else {
		vm.AddPc(1)
	}
}

// if not (R(A) <=> C) then pc++
//iABC 模式  判断寄存器a 布尔值 和 寄存器C表示布尔值是否一致 一致 则跳过下一条指令 不适用操作数B 也不改变寄存器状态
func test(i Instruction, vm LuaVM) {
	a, _, c := i.ABC()
	a += 1
	if vm.ToBoolean(a) != (c != 0) {
		vm.AddPc(1)
	}
}

