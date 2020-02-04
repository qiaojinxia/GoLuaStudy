package state

import (
	"fmt"
	. "github.com/LuaProject/api"
)
import "github.com/LuaProject/number"


//表示lua值的luaValue类型
type luaValue interface {}

func typeOf(val luaValue) LuaType {
	switch val.(type){
	case nil: return LUA_TNIL
	case bool: return LUA_TBOOLEAN
	case int64: return LUA_TNUMBER
	case float64: return LUA_TNUMBER
	case string: return LUA_TSTRING
	case *luaTable:return LUA_TTABLE
	case *closure: return LUA_TFUNCTION
	default: panic("todo!")
	}
}

func convertToBoolean(val luaValue) bool{
	switch x:= val.(type) {
		case nil : return false
		case bool: return x
		default: return true
	}
}
func convertToFloat(val luaValue) (float64, bool) {
	switch x:= val.(type){
	case float64: return x,true
	case int64: return float64(x),true
	case string: return number.ParseFloat(x)
	default: return 0,false
	}
}
func convertToInteger(val luaValue) (int64,bool) {
	switch x :=val.(type) {
	case int64:return x,true
	case float64: return number.FloatToInteger(x)
	case string: return _stringToInteger(x)
	default: return 0,false
	}
}
func _stringToInteger(s string) (int64 ,bool){
	if i,ok := number.ParseInteger(s); ok{
		return i,true
	}
	if f,ok := number.ParseFloat(s); ok{
		return number.FloatToInteger(f)
	}
	return 0,false
}

func setMetatable(val luaValue, mt *luaTable, ls *luaState){
	//判断值是否是表 如果是 直接修改原表字段
	if t,ok := val.(*luaTable);ok{
		t.metatable = mt
		return
	}
	//否则更具变量类型把原表存储在注册表里
	key := fmt.Sprintf("_MT%d",typeOf(val))
	ls.registry.put(key,mt)
}
//获取元表
func getMetatable(val luaValue, ls *luaState) *luaTable{
	//如果是表类型直接返回其原表
	if t,ok := val.(*luaTable);ok{
		return t.metatable
	}
	//如果不是原表类型 就查找 注册表找到其 原表
	key := fmt.Sprintf("_MT%d",typeOf(val))
	if mt := ls.registry.get(key);mt != nil{
		return mt.(*luaTable)
	}
	//找不到返回nil
	return nil
}
//调用元表方法
func callMetamethod(a,b luaValue,mmName string,ls *luaState) (luaValue,bool){
	var mm luaValue
	//下面的逻辑 按原方法名查找原表 2者之一 找到了 对应的方法
	// 就执行下面的语句 否则 返回 false代表找不到原表
	if mm = getMetafield(a,mmName,ls); mm == nil {
		if mm = getMetafield(b, mmName, ls); mm == nil {
			return nil, false
		}
	}
		//检查栈是否有多于的空间
		ls.stack.check(4)
		//将原方法推入栈中
		ls.stack.push(mm)
		//将参数1推入栈中
		ls.stack.push(a)
		//将参数2推入栈中
		ls.stack.push(b)
		//调用栈顶的 函数和2个参数 返回一个值
		ls.Call(2,1)
		//将栈顶返回值 传出
		return ls.stack.pop(),true
	}


func getMetafield(val luaValue,fieldName string,ls *luaState) luaValue{
	if mt := getMetatable(val,ls);mt != nil{
		return mt.get(fieldName)
	}
	return nil
}