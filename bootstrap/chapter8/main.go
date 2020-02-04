package main

import (
	"fmt"
	. "github.com/LuaProject/api"
	"github.com/LuaProject/state"
	"io/ioutil"
	"os"
	"time"
)

func main() {
	if len(os.Args) > -1 {
		data, err := ioutil.ReadFile("/Users/qiao/go/src/github.com/LuaProject/bootstrap/chapter8/luac.out")
		if err != nil {
			panic(err)
		}

		ts := time.Now()
		ls := state.New()
		ls.Register("print",print)
		ls.Register("getmetatable",getMetatable)
		ls.Register("setmetatable",setMetatable)
		ls.Load(data,"/Users/qiao/go/src/github.com/LuaProject/bootstrap/chapter8/luac.out","b")
		ls.Call(0,0)
		fmt.Println(time.Now().Sub(ts))
	}
}

func getMetatable(ls LuaState) int{
		if !ls.GetMetatable(1){
			ls.PushNil()
		}
		return 1
	}

func setMetatable(ls LuaState) int{
	ls.SetMetatable(1)
	return 1
}

func print(ls LuaState) int {
	nArgs := ls.GetTop()
	for i := 1; i <= nArgs; i++ {
		if ls.IsBoolean(i) {
			fmt.Printf("%t", ls.ToBoolean(i))
		} else if ls.IsString(i) {
			fmt.Print(ls.ToString(i))
		} else if ls.IsNumber(i) {
			fmt.Print(ls.ToString(i))
		}else if ls.IsTable(i) {
			if v,ok :=ls.ToStringX(i);ok{
				fmt.Print(v)
			}

		}else {
			fmt.Print(ls.TypeName(ls.Type(i)))
		}
		if i < nArgs {
			fmt.Print("\t")
		}
	}
	return 0
}
