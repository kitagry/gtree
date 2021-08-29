package main

import (
	"context"
	"os"
)

var (
	statusOK  = 0
	statusErr = 1
)

func main() {
	exitCode := run(context.Background())
	os.Exit(exitCode)
}
