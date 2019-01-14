package main

import (
	"github.com/porkbeans/hashi/internal/cmd"
	"os"
)

func main() {
	os.Exit(cmd.Execute(os.Args[1:]))
}
