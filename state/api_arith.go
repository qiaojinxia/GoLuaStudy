package state

import (
	"github.com/LuaProject/api"
	"github.com/LuaProject/number"
	"math"
)

var (
	iadd = func(a,b int64) int64 {return a + b}
	fadd = func(a,b float64) float64{return a + b}
	isub = func(a,b int64) int64 {return a -b}
	fsub = func(a,b float64) float64 {return a -b}
	imul = func(a,b int64) int64 {return a * b}
	fmul = func(a,b float64) float64{  return a * b }
	imod = number.IMod
	fmod = number.FMod
	pow =math.Pow
	div = func(a ,b float64) float64 {return a/b}
	iidiv = number.IFloorDiv
	fidiv =number.FFloorDiv
	band = func(a,b int64) int64{ return a&b}
	bor = func(a,b int64) int64{return a|b}
	bxor = func(a,b int64) int64{return a^b}
	shl = number.ShiftLeft
	shr = number.ShiftRight
	iunm = func(a,_ int64) int64 { return -a}
	funm = func(a,_ float64) float64 {return -a}
	bnot = func(a,_ int64)  int64 {return ^a}
)
type operator struct{
	metamethod  string
	integerFunc func(int64,int64) int64
	floatFunc func(float64,float64) float64
}
var operators =[]operator{
	operator{"__add", iadd, fadd},
	operator{"__sub", isub, fsub},
	operator{"__mul", imul, fmul},
	operator{"__mod", number.IMod, number.FMod},
	operator{"__pow", nil, math.Pow},
	operator{"__div", nil, div},
	operator{"__idiv", number.IFloorDiv, number.FFloorDiv},
	operator{"__band", band, nil},
	operator{"__bor", bor, nil},
	operator{"__bxor", bxor, nil},
	operator{"__shl", number.ShiftLeft, nil},
	operator{"__shr", number.ShiftRight, nil},
	operator{"__unm", iunm, funm},
	operator{"__bnot", bnot, nil},
}

func (self *luaState) Arith(op api.ArithOp){
	var a,b luaValue
	//去除栈顶第一个元素
	a = self.stack.pop()
	if op != api.LUA_OPUNM && op!=api.LUA_OPBNOT{
		//取出栈顶第二个元素
		b = self.stack.pop()
	}else{
		a = b
	}
	operator := operators[op]
	if result := _arith(a,b,operator);result != nil{
		self.stack.push(result)
	}else{
		panic("arithmetic error!")
	}
}
func _arith(a,b luaValue,op operator) luaValue{
	if op.floatFunc == nil{
		if x,ok := convertToInteger(a);ok{
			if y,ok := convertToInteger(b);ok{
				return op.integerFunc(x,y)
			}
		}
	}else{
		if op.integerFunc != nil{
			if x,ok := a.(int64) ;ok{
				if y,ok := b.(int64);ok{
					return op.integerFunc(x,y)
				}
			}
		}
		if x,ok := convertToFloat(a); ok{
			if y,ok :=convertToFloat(b);ok{
				return op.floatFunc(x,y)
			}
		}
	}
	return nil
}
