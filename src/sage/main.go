package main

import (
	"os"
	"sage/src/sage/controller"
)

func main() {
	os.Exit(runMain())
}

func runMain() int {
	return controller.RunCLIController()
}
