package vm

import . "github.com/LuaProject/api"

func getTabUp(i Instruction,vm LuaVM){
	a, _, c := i.ABC()
	a += 1
	vm.PushGlobalTable()
	vm.GetRK(c)
	vm.GetTable(-2)
	vm.Replace(a)
	vm.Pop(1)

}
