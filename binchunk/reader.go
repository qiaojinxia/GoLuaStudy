package binchunk

import (
	"encoding/binary"
	"math"
)

type reader struct {
	data []byte
}
//从data里读取一个字节 并向后移一个字节
func(self *reader) readByte() byte{
	b := self.data[0]//读取一个字节
	self.data =self.data[1:]//读取byte后位移一位
	return b
}
//从data里读取四个字节 并向后移四个字节
func (self *reader) readUnit32() uint32{
	i := binary.LittleEndian.Uint32(self.data)//小端方式读取cint存储类型 4个字节
	self.data = self.data[4:]//索引向后4个字节
	return i
}

//读取8个字节 并向后移动8个字节
func (self *reader) readUnit64() uint64{
	i := binary.LittleEndian.Uint64(self.data)//小端方式存储 8个字节 低位→高位
	self.data =self.data[8:] //索引向后移动8个字节
	return i
}
//通过readUnit64 读取一个Lua整数 8个字节 转为 GO的int64类型
func (self *reader) readLuaInteger() int64{
	return int64(self.readUnit64())
}
//从字节流读取一个Lua浮点数(占8个字节，映射为Go语言float64类型)
func (self *reader) readeLuaNumber() float64{
	return math.Float64frombits(self.readUnit64())
}

/**
从方法字节流读取字符串 →Go string类型
 */
func (self *reader) readString() string{
	// 读取方式(大小+内容)  先取 表示大小的1个字节
	size := uint(self.readByte())
	//如果大小为0返回 空字符串
	if size == 0 {
		return ""
	}
	//如果大小为 255 表示长字符串
	if size == 0xFF{
		size = uint(self.readUnit64())
	}
	//
	bytes := self.readBytes(size -1)
	return string(bytes)
}
/**
读取n个字节
 */
func (self *reader) readBytes(n uint) []byte{
	bytes := self.data[:n]
	self.data =self.data[n:]
	return bytes
}

func (self *reader) checkHeader(){
	//读取前4个字节 判断魔数
	if string(self.readBytes(4)) != LUA_SIGNATURE {
		panic("not a precompiled chunk!")
	}else if self.readByte() != LUAC_VERSION {
		panic("version mismatch!")
	}else if self.readByte() != LUAC_FORMAT {
		panic("format mismatch!")
	}else if string(self.readBytes(6)) != LUAC_DATA {
		panic("corrupted!")
	}else if self.readByte() != CINT_SIZE {
		panic("int size mismatch!")
	}else if self.readByte() != CSZIET_SIZE {
		panic("size_t size mismatch!")
	}else if self.readByte() != INSTRUCTION_SIZE {
		panic("instruction size mismatch!")
	} else if self.readByte() != LUA_INTEGER_SIZE {
		panic("lua_Integer size mismatch")
	}else if self.readByte() != LUA_NUMBER_SIZE {
		panic("lua_Number size mismatch!")
	}else if self.readLuaInteger() != LUAC_INT {
		panic("endianness  mismatch!")
	}else if self.readeLuaNumber() != LUAC_NUM {
		panic("float format mismatch!")
	}
}
//初始化 函数原型
func (self *reader) readProto(parentSource string) *Prototype{
	source := self.readString()
	if source =="" {source =parentSource}
	return &Prototype{
		Source: source,
		LineDefined: self.readUnit32(),
		LastLineDefined:self.readUnit32(),
		NumParams:self.readByte(),
		IsVararg:self.readByte(),
		MaxStackSize:self.readByte(),
		Code:self.readCode(),
		Constants:self.readConstants(),
		Upvalues:self.readUpvalues(),
		Protos:self.readProtos(source),
		LineInfo:self.readLineInfo(),
		LocVars:self.readLocVars(),
		UpvalueNames:self.readUpvalueNames(),
	}

}

//从字节流读取指令表
func (self *reader) readCode() []uint32 {
	code := make([]uint32,self.readUnit32())
	for i := range code{
		code[i] = self.readUnit32()
	}
	return code
}
//从字节流里读取常量表
func(self *reader) readConstants() []interface{}{
	constants := make([]interface{},self.readUnit32())
	for i := range constants{
		//读取常量 按照不同的tag 初始化成不同的 变量类型 存储在 []interface{}中
		constants[i] = self.readConstant()
	}
	return constants
}

//读取不同类型的常量
func (self *reader) readConstant() interface{}  {
	switch self.readByte() { //tag
	case TAG_NIL: return nil
	case TAG_BOOLEAN: return self.readByte() != 0
	case TAG_INTEGER: return self.readLuaInteger()
	case TAG_NUMBER: return self.readeLuaNumber()
	case TAG_SHORT_STR: return self.readString()
	case TAG_LONG_STR: return self.readString()
	default: panic("corrupted!")
	}
}

//从字节流里读取Upvalues表
func (self *reader) readUpvalues() []Upvalue{
	//获取字节表里的upvalue数量 初始化 Upvalue数组
	upvalues :=make([]Upvalue,self.readUnit32())
	for i := range upvalues{
		upvalues[i] = Upvalue{
			Instack: self.readByte(),
			Idx:     self.readByte(),
		}
	}
	return upvalues
}

//递归调用 读取函数原型
func (self *reader) readProtos(parentSource string) []*Prototype{
	protos := make([]*Prototype,self.readUnit32())
	for i:= range protos{
		protos[i] =self.readProto(parentSource)
	}
	return protos
}
//方法从字节流里读取行号表
func (self *reader) readLineInfo() []uint32 {
	lineInfo := make([]uint32,self.readUnit32())
	for i := range lineInfo {
		lineInfo[i] =self.readUnit32()
	}
	return lineInfo
}
//从发从字节流里读取局部变量表
func (self *reader) readLocVars() []LocVar{
	locVars := make([]LocVar,self.readUnit32())
	for i := range locVars {
		locVars[i] = LocVar{
			VarName: self.readString(),
			StartPC: self.readUnit32(),
			EndPC:   self.readUnit32(),
		}
	}
	return locVars
}
//从方法字节流里读取Upvalue名列表
func (self *reader) readUpvalueNames() []string {
	names := make([]string, self.readUnit32())
	for i := range names{
		names[i] =self.readString()
		}
	return names
}
