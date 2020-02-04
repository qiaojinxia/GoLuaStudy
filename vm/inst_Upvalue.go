package vm

import . "github.com/LuaProject/api"
//upvalue访问方法 伪索引
func LuaUpvalueIndex(i int) int{
	return LUA_REGISTRYINDEX - i
}
//iabc 把当前闭包的某个Upvalue值拷贝到目标寄存器中
//A指定目标寄存器 Upvalue索引B指定 C不用
func getUpval(i Instruction,vm LuaVM){
	a,b,_ := i.ABC()
	a += 1 ; b += 1
	//将 指定索引处table
	vm.Copy(LuaUpvalueIndex(b),a)
}
func setUpval(i Instruction,vm LuaVM){
	a,b,_ := i.ABC()
	a += 1; b += 1
	vm.Copy(a,LuaUpvalueIndex(b))
}
//
func getTabUp(i Instruction,vm LuaVM){
	a,b,c := i.ABC()
	a += 1; b += 1
	vm.GetRK(c)
	vm.GetTable(LuaUpvalueIndex(b))
	vm.Replace(a)

}
func setTabUp(i Instruction,vm LuaVM){
	a, b, c :=i.ABC()
	a += 1
	vm.GetRK(b)
	vm.GetRK(c)
	vm.SetTable(LuaUpvalueIndex(a))
}