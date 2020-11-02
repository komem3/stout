package main

import (
	"bufio"
	"flag"
	"fmt"
	"os"
	"time"
)

var srcPath = flag.String("path", "", "")
var noformat = flag.Bool("no-format", false, "")

var now = time.Now()

const usage = `usage: stout -path $struct_path [-no-format] $struct_name
options:
  -path string
        File path of defined struct. (required)
  -no-format bool
        Not format the output json.
`

func main() {
	flag.Usage = func() { fmt.Fprint(os.Stderr, usage) }
	flag.Parse()
	fromSt := flag.Arg(0)

	if fromSt == "" {
		Fatalf("required target struct name.\n%s", usage)
	}
	if *srcPath == "" {
		Fatalf("required target struct file path.\n%s", usage)
	}

	op := newJsonOption(*srcPath, fromSt, *noformat)
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
