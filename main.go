package main

import (
	"log"

	"github.com/liyu36/kubectl-version-manager/cmd"
)

func init() {
	log.SetFlags(0)
}

func main() {
	cmd.Execute()
}
