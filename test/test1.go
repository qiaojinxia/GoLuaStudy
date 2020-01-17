package main

import "fmt"

type lists struct {
	top int
	slots []interface{}
}
func main(){
	ls := lists{14,nil}
	for i := 0;i<ls.top;i++ {
		ls.slots = append(ls.slots, i);
	}
	ls.Rotate(2,-2)
	for i := 0;i<ls.top;i++ {
		fmt.Println(ls.slots[i].(int))
		//switch ls.slots[i].(type){
		//case float64: fmt.Printf("%f\n" ,ls.slots[i])
		//case int: fmt.Printf("%d\n" ,ls.slots[i])
		//default: fmt.Printf("null")
		//}
	}
}

func (self *lists) reverse(from,to int){
	slots := self.slots
	for from < to {
		slots[from],slots[to] = slots[to],slots[from]
		from ++
		to --
	}
}
func (self *lists) Rotate(idx,n int) {
	//获取栈顶元素
	t := self.top -1
	//获取指定索引元素
	p := idx
	var m int
	if n >= 0 {
		m = t - n
	}else {
		m = p - n - 1
	}
	self.reverse(p,m)
	self.reverse(m+1,t)
	self.reverse(p,t)
}
