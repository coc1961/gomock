package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/coc1961/gomock/internal/mockmaker"
)

func main() {
	var src = flag.String("s", "", "source go file")
	var name = flag.String("n", "", "interface name")
	var addPackage = flag.Bool("p", false, "(Optional) add package name in local objects")

	flag.Parse()

	if *src == "" || *name == "" {
		flag.CommandLine.Usage()
		return
	}
	fileExists := func(filename string) bool {
		info, err := os.Stat(filename)
		if os.IsNotExist(err) {
			return false
		}
		return !info.IsDir()
	}

	if !fileExists(*src) {
		fmt.Fprintf(os.Stderr, "file not found %v", *src)
		return
	}

	mm := mockmaker.MockMaker{}
	x := mm.CreateMock(*src, *name, *addPackage)
	fmt.Fprint(os.Stdout, x.String())
}
