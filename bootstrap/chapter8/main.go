package main

import (
	"fmt"
	. "github.com/LuaProject/api"
	"github.com/LuaProject/binchunk"
	"github.com/LuaProject/state"
	. "github.com/LuaProject/vm"
	"io/ioutil"
	"os"
	"time"
)

func main3() {
	if len(os.Args) > -1 {
		data, err := ioutil.ReadFile("/Users/qiao/go/src/github.com/LuaProject/bootstrap/luac.out")
		if err != nil {
			panic(err)
		}

		ts := time.Now()
		proto := binchunk.Undump(data)
		luaMain(proto)
		fmt.Println(time.Now().Sub(ts))
	}
}

func luaMain(proto *binchunk.Prototype) {
	nRegs := int(proto.MaxStackSize)
	ls := state.New()
	ls.SetTop(nRegs)
	for {
		pc := ls.PC()
		inst := Instruction(ls.Fetch())
		if inst.Opcode() != OP_RETURN {
			inst.Execute(ls)
			fmt.Printf("[%02d] %s ", pc+1, inst.OpName())
			printStack(ls)
		} else {
			break
		}
	}
}
func printStack(ls LuaState) {
	top := ls.GetTop()
	for i := 1; i <= top; i++ {
		t := ls.Type(i)
		switch t {
		case LUA_TBOOLEAN:
			fmt.Printf("[%t]", ls.ToBoolean(i))
		case LUA_TNUMBER:
			fmt.Printf("[%g]", ls.ToNumber(i))
		case LUA_TSTRING:
			fmt.Printf("[%q]", ls.ToString(i))
		default: // other values
			fmt.Printf("[%s]", ls.TypeName(t))
		}
	}
	fmt.Println()
}