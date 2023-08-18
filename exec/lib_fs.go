package exec

import (
	"io/ioutil"
	"os"
)

func read(args []*OwlObj) (*OwlObj, bool) {
	d := args[1].TrueStr()
	bytes, err := os.ReadFile(d)

	if err != nil {
		return nil, false
	}

	return NewString(string(bytes)), true
}

func listDir(args []*OwlObj) (*OwlObj, bool) {
	d := args[1].TrueStr()
	files, err := ioutil.ReadDir(d)

	if err != nil {
		return nil, false
	}

	names := []*OwlObj{}
	for _, file := range files {
		if !file.IsDir() {
			names = append(names, NewString(file.Name()))
		} else {
			names = append(names, NewString(file.Name()+"/"))
		}
	}

	return NewList(names), true
}

func FsLibExport() *OwlObj {
	o := NewOwlObj()

	o.SetAttr("Read", NewCallBridge(read))
	o.SetAttr("ListDir", NewCallBridge(listDir))

	return o
}
