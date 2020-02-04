package state

import (
	"fmt"
	"github.com/LuaProject/number"
	"math"
)

type luaTable struct {
	metatable *luaTable //原表
	//存放数组内容
	arr []luaValue
	//存放hash表
	_map map[luaValue]luaValue
}
//创建一个表 如果 nArr大于0 初始化数组大小 如果nRec大于0  初始化map 可以存放lua变量
func newLuaTable(nArr,nRec int) *luaTable {
	t := &luaTable{}
	if nArr > 0{
		t.arr = make([]luaValue, 0, nArr)
	}
	if nRec > 0 {
		t._map = make(map[luaValue]luaValue, nRec)
	}
	return t
}

func (self *luaTable) printTbale() string{
	if len(self.arr) !=0 {
		return fmt.Sprintf("%v",self.arr)
	}else{
		return fmt.Sprintf("%v",self._map)
	}
}

func (self *luaTable) get(key luaValue) luaValue{
	key = _floatToInteger(key)
	//如果能转化成 int类型 那么就是数组查询
	if idx, ok := key.(int64); ok{
		//判断要查询的索引是否在 数组内
		if idx >= 1 && idx <= int64(len(self.arr)){
			//lua 语言 索引从1开始 所以需要减去1
			return self.arr[idx -1]
		}
	}
	//如果 索引是其他类型 就查询hash表查找
	return self._map[key]
}
func _floatToInteger(key luaValue) luaValue{
	if f,ok := key.(float64);ok{
		if i,ok := number.FloatToInteger(f); ok{
			return i
		}
	}
	return key
}

func (self *luaTable) put(key,val luaValue) {
	//判断 key是否为  nil类型
	if key == nil {
		panic("table index is nil!")
	}
	//判断 值的实际类型
	if f, ok := key.(float64); ok && math.IsNaN(f) {
		panic("table index is Nan")
	}
	key = _floatToInteger(key)
	//判断值是否是int类型 并且 索引大于等于1
	if idx, ok := key.(int64); ok && idx >= 1 {
		//取数组长度
		arrLen := int64(len(self.arr))
		//如果索引小于数组的长度
		if idx <= arrLen {
			//替换掉数组里面原来的值
			self.arr[idx-1] = val
			//如果索引 和 数组长度相同 且 val 内容为 nil
			if idx == arrLen && val == nil {
				//删除数组 尾部的洞 考虑到
				// 假设 1 2 3 4 5 → 1 2 3 nil 5 → 1 2 3 nil nil
				//上面的数组第一次替换为nil时  调用_shrinkArray 不会删除末尾的洞
				//但当第二次 末尾 放入nil时 就会把 后两个nil一起删除了
				self._shrinkArray()
			}
			return
		//如果长度超出索引
		}else if (idx  == arrLen + 1 ){
			//先扩容
			self.arr = append(self.arr,val)
			self._expandArray()
			return
		}
	}
	//判断值是否为nil
	if val != nil {
		if self._map == nil {
			self._map = make(map[luaValue]luaValue, 7)
		}
		self._map[key] = val
	}else{
		//如果这个值val存在于map 且更新成了 nil 那么就没必要在map里去记录
		//删除可以节约空间
		delete(self._map,key)
	}
}
//循环遍历删除数组末尾的nil
func (self *luaTable) _shrinkArray() {
	for i:= len(self.arr) -1;i>0;i--{
		if self.arr[i] == nil {
			self.arr =self.arr[0:i]
		}
	}
}
//a=[1,3,4] 如果此时插入 索引为 a[5]='x' 此时由于只有3个元素 插入不进数组
//所以按照插入逻辑会保存在 map中 如果后面数组有了足够的长度足够插入这个元素了 那么就把map里的取出插入list
func (self *luaTable) _expandArray() {
	for idx := int64(len(self.arr)) + 1; true; idx++ {
		//若果在map里面找到了 索引为 arr长度
		if val, found := self._map[idx]; found {
			//删除map里的值
			delete(self._map, idx)
			//把map里的值放到数据里
			self.arr = append(self.arr, val)
		} else {
			break
		}
	}
}

func (self *luaTable) len() int{
	return len(self.arr)
}
func (self *luaTable) hasMetafield(fieldName string) bool{
	return self.metatable != nil &&
		self.metatable.get(fieldName)!=nil
}