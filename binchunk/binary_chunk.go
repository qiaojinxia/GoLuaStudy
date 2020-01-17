package binchunk

type binaryChunk struct {
	header //header
	sizeUPvalues byte //主函数upvalue 数量
	mainFunc *Prototype //主函数原型
}

type header struct {
	signature       [4]byte //魔数 4个字节
	version         byte    //版本号 一个字节
	format          byte    //格式号 一个字节
	luacData        [6]byte //前2个字节代表lua发布年份0x1993 后几个字节 依次是 0x0d 0a 1a 0a 6个字节
	cintSize        byte    //一个字节 表示cint占用的字节数 4个字节
	sizetSize       byte    //用一个字节表示 占用八个字节 0x08
	instructionSize byte    //lua虚拟机指令长度 4个字节
	luIntegerSize   byte //lua整形大小 八个字节
	luaNumberSize   byte//lua number大小
	luacInt         int64//判断 大小端 加载方式 intel是小端存储
	luacNum         float64//检查二进制加载浮点数的格式 是否一致 否则拒绝加载
}
type Prototype struct {
	Source string//源文件名 记录chunk是由哪个源文件编译出来的
	LineDefined uint32//起止行号 1个cint记录函数开始行号
	LastLineDefined uint32//结束行号 1个字节记录函数结束行号
	NumParams byte//一个字节  记录函数的固定参数个数
	IsVararg byte//是否是Vararg函数 即是否有变长参数 0代表否 1代表是
	MaxStackSize byte//寄存器数量 1个字节 编译器对lua函数生成指令表 lua虚拟机基于栈的结构实现弹入弹出功能
	Code []uint32//指令表 每条指令4个字节
	Constants []interface{}//常量表 用于表示lua常见数据类型  nil 布尔值 整数 浮点数 字符串 表示方法: tag(1个字节表示类型) + 内容
	Upvalues []Upvalue //Upvalue表 占用2个字节
	Protos []*Prototype//子函数原型表
	LineInfo []uint32//行号表 cint类型存储  行号 和指令表中的指令一一对应
	LocVars []LocVar//局部变量表 用于记录局部变量名 包含：变量名(字符串存储) +起止指令索引(cint存储)
	UpvalueNames []string//Upvalue名列表 字符串类型存储 对应Upvalues 在源代码中的名字
	//行号表 局部变量表 和Upvalue名列表存储的都是调试信息，对程序执行不必要 使用-s 编译lua时会在chunk中清空这三个表内容
}

const(
	LUA_SIGNATURE = "\x1bLua" //魔数
	LUAC_VERSION = 0x53 //主版本号
	LUAC_FORMAT = 0 //格式号 虚拟机识别匹配
	LUAC_DATA = "\x19\x93\r\n\x1a\n" //固定的  用做虚拟机校验
	CINT_SIZE = 4 //cint 4个字节
	CSZIET_SIZE = 8
	INSTRUCTION_SIZE = 4 //虚拟机指令大小
	LUA_INTEGER_SIZE = 8 //LU整数大小
	LUA_NUMBER_SIZE = 8 //浮点数大小
	LUAC_INT = 0x5678 //校验大小端存储方式
	LUAC_NUM =370.5//校验浮点数存储格式
)

const (
	TAG_NIL = 0x00 //tag 表示nil
	TAG_BOOLEAN = 0x01//tag 表示布尔类型
	TAG_NUMBER = 0x03 //tag 表示浮点数
	TAG_INTEGER = 0x13//tag 表示整数
	TAG_SHORT_STR =0x04 //tag 字符串
	TAG_LONG_STR =0x14
)
type Upvalue struct {
	Instack  byte
	Idx byte
}

//局部变量表
type LocVar struct {
	//变量名
	VarName string
	//其实索引
	StartPC uint32
	//结束索引
	EndPC uint32
}
//使用reader 结构体来处理 chunk文件
func Undump(data []byte) *Prototype {
	reader := &reader{data}
	reader.checkHeader()        //校验头部
	reader.readByte()           //跳过Upvalue数量
	return reader.readProto("") //读函数原型
}

