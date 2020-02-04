package state

import (
	. "github.com/LuaProject/api"
)

func (self *luaState) Compare(idx1,idx2 int,op CompareOp) bool{
	a := self.stack.get(idx1)
	b := self.stack.get(idx2)
	switch op {
	case LUA_OPEQ:return _eq(a,b,self)
	case LUA_OPLT: return _lt(a,b,self)
	case LUA_OPLE: return _le(a, b,self)
	default: panic("invalid compare op!")
	}
}

func _eq(a,b luaValue,ls *luaState) bool{
	switch x:= a.(type){
	case nil: return b==nil
	case *luaTable:
		if y,ok := b.(*luaTable);ok && x!= y&& ls!= nil{
			if result,ok := callMetamethod(x,y,"__eq",ls);ok{
				return convertToBoolean(result)
			}
		}
	return a == b
	case bool: y,ok := b.(bool)
	return ok&&x == y
	case string:
		y,ok := b.(string)
		return ok && x== y
	case int64:
		switch y:= b.(type) {
		case int64: return x==y
		case float64: return float64(x)==y
		default: return false
		}
	default:
		return a== b
	}
}
func _lt(a,b luaValue,ls *luaState) bool{
	switch x := a.(type){
	case string:
		if y,ok := b.(string); ok{
			return x < y
		}
	case int64:
		switch y := b.(type){
		case int64: return x < y
		case float64: return float64(x) < y
		}
	case float64:
		switch y := b.(type){
		case float64: return x<y
		case int64:return x<float64(y)
		}
	}
	if result,ok := callMetamethod(a,b,"__lt",ls);ok{
		return convertToBoolean(result)

	}else {
		panic("comparsion error!")
	}
	
}

func _le(a, b luaValue,ls *luaState) bool {
	switch x := a.(type) {
	case string:
		if y, ok := b.(string); ok {
			return x <= y
		}
	case int64:
		switch y := b.(type) {
		case int64:
			return x <= y
		case float64:
			return float64(x) <= y
		}
	case float64:
		switch y := b.(type) {
		case float64:
			return x <= y
		case int64:
			return x <= float64(y)
		}
	}
	//这下面的判断  a<=b 和  !a>b相等 所以 当 <= 找不到这个元方法时 可以找等价 > 取反
	if result,ok := callMetamethod(a,b,"__le",ls);ok {
		return convertToBoolean(result)
	}else if result,ok := callMetamethod(a,b,"__lt",ls);ok{
		return !convertToBoolean(result)
	}
	panic("comparison error!")
}
