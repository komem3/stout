package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
)

var srcPath = flag.String("path", "", "Definetion file. (required)")
var fromSt = flag.String("type", "", "Target struct type. (required)")
var noformat = flag.Bool("no-format", false, "Not format.")

func main() {
	flag.Parse()
	if *srcPath == "" || *fromSt == "" {
		fmt.Printf("%v, %v\n", srcPath, fromSt)
		flag.Usage()
		os.Exit(1)
	}

	op := newJsonOption(*srcPath, *fromSt, *noformat)
	writer := bufio.NewWriter(os.Stdout)
	if err := stType2Json(writer, op); err != nil {
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
