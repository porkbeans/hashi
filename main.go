package main

import (
	"os"

	"github.com/porkbeans/hashi/internal/cmd"
)

func main() {
	os.Exit(cmd.Execute(os.Args[1:]))
}
