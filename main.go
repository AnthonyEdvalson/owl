package main

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/AnthonyEdvalson/owl/exec"
	"github.com/AnthonyEdvalson/owl/repl"
)

func main() {
	argc := len(os.Args)

	if argc == 1 {
		repl.Start(os.Stdin, os.Stdout)
	}

	if argc == 2 {
		path := filepath.Join(os.Args[1], "main.hoot")

		ok, params, parseErr := exec.LoadProgramFromPath(path)
		if !ok {
			if parseErr == nil {
				fmt.Println("Failed to locate program")
				return
			}
			for _, e := range parseErr {
				fmt.Printf("%d:%d: %s\r\n", e.Token.Line, e.Token.Column, e.Message)
			}
			return
		}

		_, _ = exec.ExecuteProgram(params)
	}
}
