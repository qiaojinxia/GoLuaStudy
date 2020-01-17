package vm

import "github.com/LuaProject/api"
//四个字节的指令
type Instruction uint32

func (self Instruction) Opcode() int{
	//取6位操作码 32位  3f=111111 &操作
	return int(self & 0x3F)
}

//iABC分别 代表 A 8个比特 B 9个比特 C 9个比特 的操作数
func (self Instruction) ABC() (a, b, c int){
	//去出操作指令 低6位 取低8位
	a = int(self >> 6 & 0xFF)
	//去除 8位 取左9位
	c = int(self >> 14 & 0x1FF)
	//去除23位 取 9位
	b = int(self >> 23 & 0x1FF)
	return
}
//a bx 2个操作数 分别 A占用 8个比特 Bx占用18个比特
func (self Instruction) ABx() (a ,bx int) {
	//取出操作指令 低6位 取低8位
	a = int(self >> 6 & 0xFF)
	//取出前 高 18位
	bx = int(self >> 14)
	return
}
//S 和 sBx 两个操作数 分别占用 8个和 18个比特
func(self Instruction) AsBx() (a,sbx int){
	//复用上面ABx()代码
	a,bx := self.ABx()
	//读取的是无符号的整数 将其转化为 有符号 需要将其减去 2^18 -1/2
	return a,bx - MAXARG_sBx
}
//iAx 只有一个操作数 占用全部26个比特
func (self Instruction) Ax() int{
	//去出操作指令 低6位
	return int(self >> 6)
}
const MAXARG_Bx = 1 << 18 -1 //2 ^18 -1 =262143

const MAXARG_sBx = MAXARG_Bx >> 1//262143 I 2 = 131071

//返回指令的操作码名字
func (self Instruction) OpName() string{
	return opcodes[self.Opcode()].name
}
//编码模式
func(self Instruction) OpMode() byte{
	return opcodes[self.Opcode()].opMode
}
//操作数B使用模式
func (self Instruction) BMode() byte {
	return opcodes[self.Opcode()].argCMode
}
//操作数C的使用模式
func (self  Instruction) CMode() byte {
	return opcodes[self.Opcode()].opMode
}

func (self Instruction) Execute(vm api.LuaVM){
	action := opcodes[self.Opcode()].action
	if action != nil{
		action(self,vm)
	}else{
		panic(self.OpName())
	}
}