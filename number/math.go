package number

import (
	"fmt"
	"math"
)

func IFloorDiv(a, b int64) int64 {
	if a >0 && b > 0 || a < 0 && b < 0 || a%b ==0 {
		return a /b
	}else {
		return a/b -1
	}
}
func FFloorDiv(a , b float64) float64{
	return math.Floor(a/b)
}
//整数取余数
func IMod (a,b int64) int64{
	return a - IFloorDiv(a,b) * b
}
//浮点数计算  math.Floor() 向下取整
func FMod(a,b float64) float64{
	return a - math.Floor(a/b) * b
}
//按位左移
func ShiftLeft(a,n int64) int64{
	if n>= 0 {
		return a << uint64(n)
	}else{
		return ShiftRight(a,-n)
	}
}
func ShiftRight(a,n int64) int64{
	if n>= 0{
		return int64(uint64(a) >> uint64(n))
	}else{
		return ShiftLeft(a,-n)
	}
}

func FloatToInteger(f float64) (int64,bool){
	i:= int64(f)
	return i,float64(i) == f
}

func main() {

	fmt.Println(ParseInteger("0x1231232"))
	fmt.Println(ParseFloat("0x123.23"))

}