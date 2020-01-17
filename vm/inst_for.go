package vm

import . "github.com/LuaProject/api"


// R(A)-=R(A+2); pc+=sBx
// int i:=0;i<n;i++  这一步事先 将 i 减去了 步骤数 1 = -1
func forPrep(i Instruction, vm LuaVM) {
	a, sBx := i.AsBx()
	a += 1

	//下面三个 if 如果栈 内存储的是 字符串类型的话 转化为数字类型
	//  for 循环三个变量 int i:=0;i<n;i++  a数值 a+1限制 a+2步长
	if vm.Type(a) == LUA_TSTRING {
		vm.PushNumber(vm.ToNumber(a))
		vm.Replace(a)
	}

	if vm.Type(a+1) == LUA_TSTRING {
		vm.PushNumber(vm.ToNumber(a + 1))
		vm.Replace(a + 1)
	}

	if vm.Type(a+2) == LUA_TSTRING {
		vm.PushNumber(vm.ToNumber(a + 2))
		vm.Replace(a + 2)
	}


	isPositiveStep := vm.ToNumber(a+2) >= 0
	if(!isPositiveStep){
		//把要压入的 数值 和 步长
		vm.PushValue(a + 2)
		vm.PushValue(a)
	}else{
		//把要压入的 数值 和 步长
		vm.PushValue(a)
		vm.PushValue(a + 2)
	}
	//对栈顶2元素执行计算
	vm.Arith(LUA_OPSUB)
	//跟新 数值
	vm.Replace(a)
	//对pc计数器加中间代码条 数据
	vm.AddPc(sBx)
}

// R(A)+=R(A+2);
// if R(A) <?= R(A+1) then {
//   pc+=sBx; R(A+3)=R(A)
// }

//这个逻辑类似于  do while() 但 不同的是第一次 forPrep()做好初始化后直接将pc指令集增加一定数值 指向这个指令 来做判断
//这条指令 判断如果没有 符合条件 就将 pc 从新指向 循环体里面第一条 指令
func forLoop(i Instruction, vm LuaVM) {
	a, sBx := i.AsBx()
	a += 1
	// R(A)+=R(A+2);
	vm.PushValue(a + 2)
	vm.PushValue(a)
	vm.Arith(LUA_OPADD)
	vm.Replace(a)

	isPositiveStep := vm.ToNumber(a+2) >= 0

	//如果 a 数值 大于 a+1 限制 将pc计数器
	if isPositiveStep && vm.Compare(a, a+1, LUA_OPLE) ||
		//步长是负数的情况下 a >=  a+1 的话 不满足条件 i=5000;i<0;i--
		!isPositiveStep && vm.Compare(a+1, a, LUA_OPLE) {
		// pc+=sBx; R(A+3)=R(A)
		vm.AddPc(sBx)
		vm.Copy(a, a+3)
	}
}