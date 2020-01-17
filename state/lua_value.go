package state

import . "github.com/LuaProject/api"
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
