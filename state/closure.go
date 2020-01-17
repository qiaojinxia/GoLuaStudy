package state

import (
	. "github.com/LuaProject/api"
	"github.com/LuaProject/binchunk"
)

type closure struct {
	proto *binchunk.Prototype
	goFunc GoFunction
}

func newLuaClosure(proto *binchunk.Prototype) *closure{
	return &closure{proto:proto}
}

func newGoClosure(f GoFunction) *closure{
	return &closure{goFunc:f}
}
