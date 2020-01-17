package main

import (
	"fmt"
	"github.com/LuaProject/api"
	"github.com/LuaProject/state"
	"io/ioutil"
	"time"
)

func main04(){
	data, err := ioutil.ReadFile("/Users/qiao/go/src/github.com/LuaProject/bootstrap/luac.out")
	if err != nil {
		panic(err)
	}
	ls := state.New()
	ls.Load(data,"123","b")
	ls.Call(0,0)
}


func main05() {
	data, err := ioutil.ReadFile("/Users/qiao/go/src/github.com/LuaProject/bootstrap/luac.out")
	if err != nil {
		panic(err)
	}
		ls := state.New()
		ls.Register("print", print)
		ls.Register("now", now)
		ls.Load(data,"123","b")
		ls.Call(0, 0)
	}


func print(ls api.LuaState) int {
	nArgs := ls.GetTop()
	for i := 1; i <= nArgs; i++ {
		if ls.IsBoolean(i) {
			fmt.Printf("%t", ls.ToBoolean(i))
		} else if ls.IsString(i) {
			fmt.Print(ls.ToString(i))
		} else if ls.IsNumber(i) {
			fmt.Print(ls.ToString(i))
		} else {
			fmt.Print(ls.TypeName(ls.Type(i)))
		}
		if i < nArgs {
			fmt.Print("\t")
		}
	}
	fmt.Println()
	return 0
}



func now(ls api.LuaState) int {
	fmt.Println(time.Now())
	return 0
}