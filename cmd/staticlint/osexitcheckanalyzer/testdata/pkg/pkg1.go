package pkg1

import (
	"os"
)

func func1() {
	os.Exit(0) // want "expression is os.Exit"
}
