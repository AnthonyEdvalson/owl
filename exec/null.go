package exec

var nullCache *OwlObj

func NewNull() *OwlObj {
	if nullCache != nil {
		return nullCache
	}

	n := NewOwlObj()
	n.SetDeepAttr("eq", NewCallBridge(nullEq))
	n.SetDeepAttr("ne", NewCallBridge(nullNe))
	n.SetDeepAttr("str", NewCallBridge(nullStr))
	n.SetDeepAttr("null", NewBool(true))
	n.DeleteDeepAttr("iter")
	n.DeleteDeepAttr("index")
	n.DeleteDeepAttr("setIndex")
	n.DeleteDeepAttr("has")

	nullCache = n

	return n
}

func nullEq(args []*OwlObj) (*OwlObj, bool) {
	return NewBool(args[1] == args[2]), true
}

func nullNe(args []*OwlObj) (*OwlObj, bool) {
	return NewBool(args[1] != args[2]), true
}

func nullStr(args []*OwlObj) (*OwlObj, bool) {
	return NewString("null"), true
}
