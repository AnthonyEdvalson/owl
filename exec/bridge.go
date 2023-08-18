package exec

type BridgeData struct {
	This       *OwlObj
	BridgeCall func(args []*OwlObj) (*OwlObj, bool)
}

func NewCallBridge(f func(args []*OwlObj) (*OwlObj, bool)) *OwlObj {
	b := &OwlObj{}
	data := &BridgeData{nil, f}

	b.BridgeCall = func(a *OwlObj) (*OwlObj, bool) { return bridgeCall(b, a) }
	b.Bind = func(this *OwlObj) { data.This = this }
	b.Raw = data

	return b
}

func flattenArg(arg *OwlObj) []*OwlObj {
	var args []*OwlObj

	if arg == nil {
		args = []*OwlObj{}
	} else {
		list, ok := arg.Raw.([]*OwlObj)

		if ok {
			args = list
		} else {
			args = []*OwlObj{arg}
		}
	}

	return args
}

func bridgeCall(callBridge *OwlObj, arg *OwlObj) (*OwlObj, bool) {
	data := callBridge.Raw.(*BridgeData)
	this := data.This

	args := append([]*OwlObj{this}, flattenArg(arg)...)

	return data.BridgeCall(args)
}
