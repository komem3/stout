package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

var srcPath = *flag.String("path", "", "Definetion file. (required)")
var fromSt = *flag.String("type", "", "Target struct type. (required)")

func main() {
	flag.Parse()
	if srcPath == "" || fromSt == "" {
		flag.Usage()
		os.Exit(1)
	}
	writer := bufio.NewWriter(os.Stdout)
	if err := stType2Json(srcPath, fromSt, writer); err != nil {
		Fatalf(err.Error())
	}
	if err := writer.Flush(); err != nil {
		Fatalf(err.Error())
	}
}

func Fatalf(format string, args ...interface{}) {
	fmt.Fprintf(os.Stderr, format, args...)
	os.Exit(1)
}
