package vm

import (
	. "github.com/LuaProject/api"
)

func move(i Instruction,vm LuaVM){
	a , b ,_ := i.ABC()
	//Lua api从1开始索引 寄存器索引从0开始
	a += 1; b += 1
	vm.Copy(b,a)
}
func jmp(i Instruction,vm LuaVM){
	a, sBx := i.AsBx()
	vm.AddPc(sBx)
	if a != 0{
		vm.CloseUpvalues(a)
	}
}

func length(i Instruction,vm LuaVM){
	a , b ,_ := i.ABC()
	a +=1;b +=1
	vm.Len(b)
	vm.Replace(a)
}