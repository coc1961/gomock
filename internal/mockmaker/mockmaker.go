package mockmaker

import (
	"bytes"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type MockMaker struct {
	StructName string
	Funcs      []*Func
}

type Func struct {
	FuncName string
	Params   []*DataType
	Returns  []*DataType
}

type DataType struct {
	Name string
	Type string
}

func (mm *MockMaker) CreateMock(filePath, structName string) *MockMaker {

	m := MockMaker{
		Funcs: make([]*Func, 0),
	}
	m.StructName = structName

	fs := token.NewFileSet()
	f, _ := parser.ParseFile(fs, filePath, nil, 0)

	for _, dec := range f.Decls {
		if gen, ok := dec.(*ast.GenDecl); ok {
			if gen.Tok != token.TYPE {
				continue
			}
			for _, specs := range gen.Specs {
				if ts, ok := specs.(*ast.TypeSpec); ok {
					if ts.Name.String() != structName {
						continue
					}

					if iFace, ok := ts.Type.(*ast.InterfaceType); ok {
						for _, meths := range iFace.Methods.List {
							if len(meths.Names) == 0 {
								break
							}
							for _, name := range meths.Names {
								ff := Func{
									FuncName: name.String(),
									Params:   make([]*DataType, 0),
									Returns:  make([]*DataType, 0),
								}
								m.Funcs = append(m.Funcs, &ff)

								if ft, ok := meths.Type.(*ast.FuncType); ok {
									for _, p := range ft.Params.List {

										dn := ""
										if p.Names != nil && len(p.Names) > 0 {
											dn = p.Names[0].String()
										}
										dt := &DataType{
											Name: dn,
										}
										ff.Params = append(ff.Params, dt)

										switch t := p.Type.(type) {
										case *ast.Ident:
											dt.Type = t.Name
										case *ast.SelectorExpr:
											dt.Type = mm.processSelectorExpr(t)
										case *ast.StarExpr:
											dt.Type = mm.processStarExpr(t)
										}
									}
									for _, r := range ft.Results.List {
										dn := ""
										if r.Names != nil && len(r.Names) > 0 {
											dn = r.Names[0].String()
										}

										dt := &DataType{
											Name: dn,
										}

										ff.Returns = append(ff.Returns, dt)

										switch t := r.Type.(type) {
										case *ast.SelectorExpr:
											dt.Type = mm.processSelectorExpr(t)
										case *ast.StarExpr:
											dt.Type = mm.processStarExpr(t)
										case *ast.Ident:
											dt.Type = t.Name
										}
									}
								}
							}
						}
					}
				}
			}
		}
	}

	return &m
}

func (mm *MockMaker) processSelectorExpr(t *ast.SelectorExpr) string {
	var param bytes.Buffer
	if ident, ok := t.X.(*ast.Ident); ok {
		param.WriteString(ident.Name)
	}
	param.WriteString(".")
	param.WriteString(t.Sel.Name)
	return param.String() // context.Context
}

func (mm *MockMaker) processStarExpr(t *ast.StarExpr) string {
	var param bytes.Buffer
	param.WriteString("*")
	if ident, ok := t.X.(*ast.Ident); ok {
		param.WriteString(ident.Name)
	}
	if expr, ok := t.X.(*ast.SelectorExpr); ok {
		param.WriteString(mm.processSelectorExpr(expr))
	}
	return param.String() // *Message
}

func (mm *MockMaker) String() string {
	if len(mm.Funcs) == 0 {
		return ""
	}

	var str bytes.Buffer

	c := func(s ...string) {
		for _, s1 := range s {
			str.WriteString(s1)
		}
	}

	c("type ", mm.StructName, "Mock struct {\n")
	for _, f := range mm.Funcs {
		c("\tCallback", f.FuncName, " func(")
		coma := ""
		for i, p := range f.Params {
			if p.Name == "" {
				p.Name = fmt.Sprintf("param%v", i)
			}
			c(coma, p.Name, " ", p.Type)
			coma = ", "
		}
		c(") (")
		coma = ""
		for _, p := range f.Returns {
			if p.Name == "" {
				c(coma, p.Type)
			} else {
				c(coma, p.Name, " ", p.Type)
			}
			coma = ", "
		}
		c(")\n")
	}
	c("}\n\n")
	for _, f := range mm.Funcs {
		c("func (m *", mm.StructName, "Mock) ", f.FuncName, "(")
		coma := ""
		for i, p := range f.Params {
			if p.Name == "" {
				p.Name = fmt.Sprintf("param%v", i)
			}
			c(coma, p.Name, " ", p.Type)
			coma = ", "
		}
		c(") (")
		coma = ""
		for _, p := range f.Returns {
			if p.Name == "" {
				c(coma, p.Type)
			} else {
				c(coma, p.Name, " ", p.Type)
			}
			coma = ", "
		}
		c(") {\n")
		c("\tif m.Callback", f.FuncName, " != nil {\n")
		c("\t\treturn m.Callback", f.FuncName, "(")
		coma = ""
		for i, p := range f.Params {
			if p.Name == "" {
				p.Name = fmt.Sprintf("param%v", i)
			}
			c(coma, p.Name)
			coma = ", "
		}
		c(")\n")
		c("\t}\n")

		c("\treturn ")
		coma = ""
		for _, p := range f.Returns {
			c(coma, mm.instance(p.Type))
			coma = ", "
		}
		c("\n")
		c("}\n\n")
	}

	return str.String()
}

var data_type = map[string]string{
	"uint8":      "uint8",
	"uint16":     "uint16",
	"uint32":     "uint32",
	"uint64":     "uint64",
	"int8":       "int8",
	"int16":      "int16",
	"int32":      "int32",
	"int64":      "int64",
	"float32":    "float32",
	"float64":    "float64",
	"complex64":  "complex64",
	"complex128": "complex128",
	"byte":       "byte",
	"rune":       "rune",
	"uint":       "uint",
	"int":        "int",
	"uintptr":    "uintptr",
	"error":      "nil",
}

func (mm *MockMaker) instance(s string) string {
	aster := strings.Contains(s, "*")
	dt := strings.ReplaceAll(s, "*", "")

	if d, ok := data_type[s]; ok {
		dt = d
	} else {
		dt = dt + "{}"
	}
	if aster {
		dt = "&" + dt
	}
	return dt
}
