package main

import (
	"fmt"
	"os"
	"owl/exec"
	"owl/repl"
	"path/filepath"
)

func main() {
	argc := len(os.Args)

	if argc == 1 {
		repl.Start(os.Stdin, os.Stdout)
	}

	if argc == 2 {
		path := filepath.Join(os.Args[1], "main.hoot")

		_, _, errors := exec.ExecuteFile(path)

		for _, error := range errors {
			fmt.Printf("%d:%d: %s\r\n", error.Token.Line, error.Token.Column, error.Message)
		}
	}
}
