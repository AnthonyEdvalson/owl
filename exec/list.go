package exec

import (
	"sort"
	"strconv"
	"strings"
)

func NewList(values []*OwlObj) *OwlObj {
	l := NewOwlObj()

	l.SetDeepAttr("bool", NewCallBridge(listBool))
	l.SetDeepAttr("index", NewCallBridge(listIndex))
	l.SetDeepAttr("setIndex", NewCallBridge(listSetIndex))
	l.SetDeepAttr("slice", NewCallBridge(listSlice))
	l.SetDeepAttr("str", NewCallBridge(listStr))
	l.SetDeepAttr("has", NewCallBridge(listHas))
	l.SetDeepAttr("iter", NewCallBridge(listIter))
	l.SetAttr("Reverse", NewCallBridge(listReverse))
	l.SetAttr("Add", NewCallBridge(listAppend))
	l.SetAttr("Sort", NewCallBridge(listSort))
	l.SetAttr("Join", NewCallBridge(listJoin))
	l.SetAttr("Len", NewCallBridge(listLen))

	l.SetAttr("Map", NewCallBridge(listMap))
	l.SetAttr("Filter", NewCallBridge(listFilter))
	l.SetAttr("Reduce", NewCallBridge(listReduce))
	l.SetAttr("FlatMap", NewCallBridge(listFlatMap))

	l.Raw = values

	return l
}

func mapIndex(i int64, len int) int {
	if i < 0 {
		return int(i) + len
	}

	return int(i)
}

func listBool(args []*OwlObj) (*OwlObj, bool) {
	raw, ok := args[0].TrueList()

	if !ok {
		return nil, false
	}

	return NewBool(len(raw) > 0), true
}

func listIndex(args []*OwlObj) (*OwlObj, bool) {
	raw, ok := args[0].TrueList()

	if !ok {
		return nil, false
	}

	index, ok := args[1].TrueInt()

	if !ok {
		return nil, false
	}

	index32 := mapIndex(index, len(raw))

	return raw[index32], true
}

func listSetIndex(args []*OwlObj) (*OwlObj, bool) {
	raw, ok := args[0].TrueList()

	if !ok {
		return nil, false
	}

	index, ok := args[1].TrueInt()

	if !ok {
		return nil, false
	}

	index32 := mapIndex(index, len(raw))

	raw[index32] = args[2]

	return NewNull(), true
}

func listSlice(args []*OwlObj) (*OwlObj, bool) {
	list, ok := args[0].TrueList()

	if !ok {
		return nil, false
	}

	start := args[1]
	end := args[2]

	var startIndex64, endIndex64 int64
	var startIndex, endIndex int = 0, len(list)
	var startOk, endOk bool

	if start != nil {
		startIndex64, startOk = start.TrueInt()
		if startOk {
			startIndex = mapIndex(startIndex64, len(list))
		}
	}

	if end != nil {
		endIndex64, endOk = end.TrueInt()
		if endOk {
			endIndex = mapIndex(endIndex64, len(list))
		}
	}

	if startIndex < 0 {
		startIndex = startIndex + len(list)
	}

	if endIndex > len(list) {
		endIndex = len(list)
	}

	l := endIndex - startIndex

	newList := make([]*OwlObj, l)

	copy(newList, list[startIndex:endIndex])

	return NewList(newList), true
}

func listLen(args []*OwlObj) (*OwlObj, bool) {
	raw, ok := args[0].TrueList()

	if !ok {
		return nil, false
	}

	return NewInt(int64(len(raw))), true
}

func listReverse(args []*OwlObj) (*OwlObj, bool) {
	raw, ok := args[0].TrueList()

	if !ok {
		return nil, false
	}

	for i, j := 0, len(raw)-1; i < j; i, j = i+1, j-1 {
		raw[i], raw[j] = raw[j], raw[i]
	}

	return NewList(raw), true
}

func listAppend(args []*OwlObj) (*OwlObj, bool) {
	if len(args) != 2 {
		panic("listAppend expects 2 arguments, got " + strconv.Itoa(len(args)))
	}

	raw, ok := args[0].TrueList()

	if !ok {
		return nil, false
	}

	args[0].Raw = append(raw, args[1])

	return args[0], true
}

func listStr(args []*OwlObj) (*OwlObj, bool) {
	raw, ok := args[0].TrueList()

	if !ok {
		return nil, false
	}

	b := strings.Builder{}

	b.WriteString("[")

	for i, v := range raw {
		b.WriteString(v.TrueStr())

		if i < len(raw)-1 {
			b.WriteString(", ")
		}
	}

	b.WriteString("]")

	return NewString(b.String()), true
}

func listHas(args []*OwlObj) (*OwlObj, bool) {
	list, ok := args[0].TrueList()

	if !ok {
		return nil, false
	}

	for _, v := range list {
		bo, ok := v.Eq(v, args[1])
		if ok && bo.IsTruthy() {
			return NewBool(true), true
		}
	}

	return NewBool(false), true
}

func listIter(args []*OwlObj) (*OwlObj, bool) {
	raw, ok := args[0].TrueList()

	if !ok {
		return nil, false
	}

	return NewList(raw), true
}

func listSort(args []*OwlObj) (*OwlObj, bool) {
	raw, ok := args[0].TrueList()

	if !ok {
		return nil, false
	}

	sort.Slice(raw[:],
		func(i, j int) bool {
			co, ok := raw[i].Lt(raw[i], raw[j])
			return ok && co.IsTruthy()
		},
	)

	return NewList(raw), true
}

func listJoin(args []*OwlObj) (*OwlObj, bool) {
	raw, ok := args[0].TrueList()

	if !ok {
		return nil, false
	}

	delim := args[1].TrueStr()

	b := strings.Builder{}

	for i, v := range raw {
		b.WriteString(v.TrueStr())

		if i < len(raw)-1 {
			b.WriteString(delim)
		}
	}

	return NewString(b.String()), true
}

func listMap(args []*OwlObj) (*OwlObj, bool) {
	raw, ok := args[0].TrueList()

	if !ok {
		return nil, false
	}

	fn := args[1]

	newList := make([]*OwlObj, len(raw))

	for i, v := range raw {
		newList[i], ok = fn.Call(v)
		if !ok {
			return nil, false
		}
	}

	return NewList(newList), true
}

func listFilter(args []*OwlObj) (*OwlObj, bool) {
	raw, ok := args[0].TrueList()

	if !ok {
		return nil, false
	}

	fn := args[1]

	newList := []*OwlObj{}

	for _, v := range raw {
		keep, ok := fn.Call(v)
		if !ok {
			return nil, false
		}
		if keep.IsTruthy() {
			newList = append(newList, v)
		}
	}

	return NewList(newList), true
}

func listReduce(args []*OwlObj) (*OwlObj, bool) {
	raw, ok := args[0].TrueList()

	if !ok {
		return nil, false
	}

	fn := args[1]
	acc := args[2]

	for _, v := range raw {
		arg := NewList([]*OwlObj{acc, v})
		acc, ok = fn.Call(arg)
		if !ok {
			return nil, false
		}
	}

	return acc, true
}

func listFlatMap(args []*OwlObj) (*OwlObj, bool) {
	raw, ok := args[0].TrueList()

	if !ok {
		return nil, false
	}

	fn := args[1]

	newList := []*OwlObj{}

	for _, v := range raw {
		mapped, ok := fn.Call(v)
		if !ok {
			return nil, false
		}
		mappedList, ok := mapped.AsList()
		if !ok {
			return nil, false
		}
		newList = append(newList, mappedList...)
	}

	return NewList(newList), true
}
