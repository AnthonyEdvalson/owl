package exec

import (
	"fmt"
	"testing"
)

func TestEngine(t *testing.T) {
	// Test that the engine can load a program, execute it, and return the
	// correct result.

	contents := `return "Hello, " + globalVar + "!"`
	fileName := "test.hoot"

	params, errs := LoadProgram(contents, fileName)
	if len(errs) > 0 {
		for _, err := range errs {
			fmt.Println(err)
		}
		panic("Failed to load program")
	}

	params.globals["globalVar"] = NewString("world")

	result, _ := ExecuteProgram(params)
	if result.TrueStr() != "Hello, world!" {
		t.Errorf("Expected result to be \"Hello, world!\", got %v", result.TrueStr())
	}
}
