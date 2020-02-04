package state

import (
	"fmt"
	. "github.com/LuaProject/api"
)

func (self *luaState) TypeName(tp LuaType) string {
	switch tp{
	case LUA_TNONE: return "no value"
	case LUA_TNIL: return "nil"
	case LUA_TBOOLEAN: return "boolean"
	case LUA_TNUMBER: return "number"
	case LUA_TSTRING: return "string"
	case LUA_TTABLE: return "table"
	case LUA_TFUNCTION: return "function"
	case LUA_TTHREAD: return "thread"
	default: return "userdata"
	}
}
//根据栈索引 获取指 将类型转换为 基本类型 如果 解析失败返回LUA_TNONE
func (self *luaState) Type(idx int) LuaType {
	if self.stack.isValid(idx) {
		val := self.stack.get(idx)
		return typeOf(val)
	}
	return LUA_TNONE
}
func (self *luaState) IsNone(idx int) bool{
	return self.Type(idx) == LUA_TNONE
}
func (self *luaState) IsNil(idx int) bool{
	return self.Type(idx) ==LUA_TNIL
}
func (self *luaState) IsNoneOrNil(idx int) bool{
	return self.Type(idx) <= LUA_TNIL
}
func (self *luaState) IsBoolean(idx int) bool{
	return  self.Type(idx) == LUA_TBOOLEAN
}
func (self *luaState) IsString(idx int) bool {
	t := self.Type(idx)
	return t == LUA_TSTRING
}
func (self *luaState) IsNumber(idx int) bool{
	_,ok := self.ToNumberX(idx)
	return ok
}
func (self *luaState) IsInteger(idx int) bool {
	val:= self.stack.get(idx)
	_,ok := val.(int64)
	return ok
}

func (self *luaState) IsTable(idx int) bool {
	t := self.Type(idx)
	return t == LUA_TTABLE

}


//从指定索引获取 一个值转换成布尔值 转换成布尔类型
func (self *luaState) ToBoolean(idx int)bool{
	val := self.stack.get(idx)
	return convertToBoolean(val)
}

// 从指定索引取出一个数字 如果值不是数字类型，则需要进行类型转换
func (self *luaState) ToNumber(idx int)float64{
	n,_ := self.ToNumberX(idx)
	return n
}

func (self *luaState) ToNumberX(idx int) (float64,bool){
	val := self.stack.get(idx)
	return convertToFloat(val)
}
//从索引取出一个整数型 如果值不是整数类型 则需要进行类型转换
func (self *luaState) Tointeger(idx int) int64{
	i,_ := self.ToIntegerX(idx)
	return i
}
func (self *luaState) ToIntegerX(idx int) (int64,bool){
	val:= self.stack.get(idx)
	return convertToInteger(val)
}
func (self *luaState) ToTable(idx int) (*luaTable,bool) {
	val:= self.stack.get(idx)
	if t,ok := val.(*luaTable);ok{
		return t,ok
	}
	return nil,false
}
//从指定索引处取出一个值 如果是字符串直接返回 如果是数字转换成字符串
func (self *luaState) ToStringX(idx int) (string,bool) {
	val:= self.stack.get(idx)
	switch x:= val.(type) {
	case string: return x,true
	case int64,float64:
		s := fmt.Sprintf("%v", x)
		self.stack.set(idx, s)
		return s,true
	case *luaTable:
		v,_ := val.(*luaTable)
		ss := v.printTbale()
		return ss,true
	default:
		return "",false
	}
}
func (self *luaState) ToString(idx int) string{
	s,_ := self.ToStringX(idx)
	return s
}
func (self *luaState) IsGoFunction(idx int) bool{
	val := self.stack.get(idx)
	if c,ok := val.(*closure);ok{
		return c.goFunc !=  nil
	}
	return false
}
func (self *luaState) ToGoFunction(idx int) GoFunction{
	val := self.stack.get(idx)
	if c,ok := val.(*closure);ok{
		return c.goFunc
	}
	return nil
}
