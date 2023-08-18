package exec

import (
	"fmt"
	"os"
	"path/filepath"
	"strings"
)

// TODO use plugins to dynamically import go modules
var golib = map[string]*OwlObj{
	"lib_http": HttpLibExport(),
	"fs":       FsLibExport(),
	"os":       OsLibExport(),
}

func NewModule(name string, currentPath string) (*OwlObj, string) {
	pathStr := name

	if pathStr[0] == '.' {
		dir := filepath.Dir(currentPath)
		pathStr = filepath.Clean(filepath.Join(dir, pathStr+".hoot"))
	} else if pathStr[0] == '/' {
		pathStr = filepath.Clean(pathStr + ".hoot")
	} else {
		lib, ok := golib[name]

		if ok {
			return lib, pathStr
		}
		ex, err := os.Executable()
		if err != nil {
			panic(err)
		}
		exPath := filepath.Dir(ex)
		pathStr = filepath.Join(exPath, "lib", pathStr+".hoot")
	}

	alias := strings.TrimSuffix(filepath.Base(pathStr), ".hoot")

	_, e, err := ExecuteFile(pathStr)

	if len(err) > 0 {
		for _, e := range err {
			fmt.Println(e.Token.File + ":" + fmt.Sprint(e.Token.Line) + ":" + fmt.Sprint(e.Token.Column) + ": " + e.Message)
		}

		panic("Failed to load module: " + name)
	}

	o := NewOwlObj()
	o.Attr = e.Frames[0]
	o.SetDeepAttr("name", NewString(alias))
	o.SetDeepAttr("str", NewCallBridge(moduleStr))

	return o, alias
}

func moduleStr(args []*OwlObj) (*OwlObj, bool) {
	name, ok := args[0].GetDeepAttr("name")

	if !ok {
		return nil, false
	}

	return NewString(name.TrueStr()), true
}
