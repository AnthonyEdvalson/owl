package exec

import (
	"fmt"
	"os/exec"
	"runtime"
	"strings"
)

func OsLibExport() *OwlObj {
	o := NewOwlObj()

	o.SetAttr("Exec", NewCallBridge(execCommand))
	o.SetAttr("Platform", NewCallBridge(execPlatform))

	return o
}

func execCommand(args []*OwlObj) (*OwlObj, bool) {
	if len(args) < 2 {
		return NewString("Not enough arguments, need at least 1"), false
	}

	cmd := args[1].TrueStr()

	cmdArgs := make([]string, len(args)-2)
	for i := 2; i < len(args); i++ {
		cmdArgs[i-2] = args[i].TrueStr()
	}

	if strings.Contains(cmd, " ") {
		return NewString("Command contains spaces, separate into multiple arguments"), false
	}

	cmdOut, err := exec.Command(cmd, cmdArgs...).CombinedOutput()

	if err != nil {
		return NewString("Command failed to run: " + err.Error() + "\r\nOutput: " + fmt.Sprintf("%s", cmdOut)), false
	}

	return NewString(string(cmdOut)), true
}

func execPlatform(args []*OwlObj) (*OwlObj, bool) {
	return NewString(runtime.GOOS), true
}
