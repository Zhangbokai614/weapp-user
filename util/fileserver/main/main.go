package main

import (
	"os"

	"github.com/dovics/wx-demo/util/fileserver"
)

func main() {
	wdir, _ := os.Getwd()
	fileserver.StartFileServer(":9573", wdir)
}
