package mockmaker

import (
	"go/ast"
	"go/parser"
	"go/token"
	"io/fs"
	"os"
	"path/filepath"
	"strings"
)

type Imports struct {
	Package    string
	Path       string
	Interfaces []string
}

func findInterface(filePath string) []string {
	fs1 := token.NewFileSet()
	f, _ := parser.ParseFile(fs1, filePath, nil, 0)

	imports := getImports(f)
	interfaces := getInterfaces(f)
	_ = interfaces

	for i := range imports {
		filepath.Walk(imports[i].Path, func(path string, info fs.FileInfo, err error) error {
			if info.IsDir() {
				return nil
			}
			fs := token.NewFileSet()
			f, err := parser.ParseFile(fs, path, nil, 0)
			if err == nil {
				imports[i].Interfaces = getInterfaces(f)
				for _, x := range imports[i].Interfaces {
					interfaces = append(interfaces, imports[i].Package+"."+x)
				}
			}
			return nil
		})
	}

	return interfaces
}

func getImports(f *ast.File) []Imports {
	ret := []Imports{}

	ast.Inspect(f, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.ImportSpec:
			pt := strings.ReplaceAll(t.Path.Value, "\"", "")
			filePath := os.Getenv("HOME") + "/go/src/" + pt
			if os.Getenv("GOPATH") != "" {
				filePath = os.Getenv("GOPATH") + "/src/" + pt
			}
			pack := strings.Split(filePath, "/")
			ret = append(ret, Imports{
				Package: pack[len(pack)-1],
				Path:    filePath,
			})
		}
		return true
	})
	return ret
}

func getInterfaces(f *ast.File) []string {
	ret := []string{}

	ast.Inspect(f, func(n ast.Node) bool {
		switch t := n.(type) {
		case *ast.TypeSpec:
			if _, ok := t.Type.(*ast.InterfaceType); ok {
				ret = append(ret, t.Name.Name)
			}
		}
		return true
	})
	return ret
}
