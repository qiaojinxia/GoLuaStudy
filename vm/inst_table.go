package vm

import . "github.com/LuaProject/api"
const LFIELDS_PER_FLUSH = 50
/**
iABC()模式 创建空表
a 存入寄存器索引
b 表的初始容量
c hash表容量
这里注意一点 Fb2int 是将
9个比特能表示的最大数 2^9 方 最多也只能表示成 512
但在某些特殊情况下 初始容量可能也需要很大 一开始调的过小就需要频繁地扩充大小
这样会加载效率 会变慢 所以lua 用了一种叫做浮点字节的编码方式
可以用一个字节表示 更大的容量 格式:eeeeexxxx 当eeeee表示为0的时候大小 就是实际xxxx这
低四位的大小 否则表示 1xxx * 2 ^(eeeee -1)
 */
func newTable(i Instruction,vm LuaVM){
	a, b, c := i.ABC()
	a += 1
	vm.CreateTable(Fb2int(b),Fb2int(c))
	vm.Replace(a)
}

/**
根据键从 表里取值 ,并放入目标寄存器中
iABC()模式
A 将拿到的值存入的目标寄存器地址索引
B 表在寄存器中的地址索引
C 要查询的key 从寄存器或Const表中去取
 */
func getTable(i Instruction, vm LuaVM) {
	a, b, c := i.ABC()
	a += 1;b += 1
	vm.GetRK(c)
	vm.GetTable(b)
	vm.Replace(a)
}
/**
获取 键 值 push 到 表里
iABC 模式
a 表在寄存器中的地址索引
b 要查询的value 从寄存器或const表中去取
c 要查询的key 从寄存器或const表中去取
 */
func setTable(i Instruction,vm LuaVM){
	a, b, c := i.ABC()
	a += 1
	vm.GetRK(b)
	vm.GetRK(c)
	vm.SetTable(a)
}

/**
与setTable 不同 这个操作是写入数组的
iABC 模式
a 表(数组)在栈中的索引地址
b 数组的其实索引 c + 1....b 为要存放的数组索引
c 数组在寄存器中紧挨着的一系列值的数量
但是C数组索引只能存放9个比特 512 个索引大小 再打的索引无法表示 所以
为了扩大数组 的可能容量
 */
func setList(i Instruction,vm LuaVM){
	a,b,c := i.ABC()
	a += 1
	bIsZero := b == 0
	if bIsZero {
		b = int(vm.Tointeger(-1) )- a - 1
		vm.Pop(1)
	}
	if c > 0 {
		c = c - 1
	} else {
		c = Instruction(vm.Fetch()).Ax()
	}

	vm.CheckStack(1)
	idx := int64(c * LFIELDS_PER_FLUSH)
	for j := 1; j <= b; j++ {
		idx++
		vm.PushValue(a + j)
		vm.SetI(a, idx)
	}
	if bIsZero {
		for j := vm.RegisterCount() + 1; j <= vm.GetTop(); j++ {
			idx++
			vm.PushValue(j)
			vm.SetI(a, idx)
		}

		// clear stack
		vm.SetTop(vm.RegisterCount())
	}
}
