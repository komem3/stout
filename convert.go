package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

type orderDefine struct {
	order  uint
	filed  string
	define interface{}
}

type structDefine map[string]interface{}

var basicMap = structDefine{
	"string":     "string",
	"int":        1,
	"int8":       1,
	"int16":      1,
	"int32":      1,
	"int64":      1,
	"uint":       1,
	"uint8":      1,
	"uint16":     1,
	"uint32":     1,
	"uint64":     1,
	"float32":    "1.0",
	"float64":    "1.0",
	"complex64":  "1.0",
	"complex128": "1.0",
	"bool":       true,
	"byte":       byte('a'),
	"rune":       rune('a'),
}

func stType2Map(path string, stType string) (structDefine, error) {
	fset := token.NewFileSet()
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil, err
	}

	var fromSpec *ast.TypeSpec
	for _, d := range f.Decls {
		gencl, ok := d.(*ast.GenDecl)
		if !ok || gencl.Tok != token.TYPE {
			continue
		}
		for _, spec := range gencl.Specs {
			spec := spec.(*ast.TypeSpec)
			if spec.Name.Name == stType {
				fromSpec = spec
			}
		}
	}
	if fromSpec == nil {
		return nil, fmt.Errorf("not found %s\n", fromSt)
	}

	st, ok := fromSpec.Type.(*ast.StructType)
	if !ok {
		return nil, fmt.Errorf("%s is not struct\n", fromSt)
	}

	return convertMap(st), nil
}

// TODO: jsonタグの認識
// TODO: test code
func convertMap(st *ast.StructType) structDefine {
	dstJson := make(map[string]interface{})
	for _, f := range st.Fields.List {
		if len(f.Names) == 0 ||
			strings.ToUpper(f.Names[0].Name)[0] != f.Names[0].Name[0] {
			continue
		}
		fname := f.Names[0].Name

		var (
			ident *ast.Ident
			array []interface{}
		)
		switch v := f.Type.(type) {
		case *ast.StarExpr:
			ident = v.X.(*ast.Ident)
		case *ast.Ident:
			ident = v
		case *ast.ArrayType:
			if ptr, ok := v.Elt.(*ast.StarExpr); ok {
				ident = ptr.X.(*ast.Ident)
			} else {
				ident = v.Elt.(*ast.Ident)
			}
			array = make([]interface{}, 1)
		default:
			panic(fmt.Sprintf("unexpected type: %#v", f.Type))
		}

		if basicValue, ok := basicMap[ident.Name]; ok {
			if len(array) != 0 {
				array[0] = basicValue
				dstJson[fname] = array
				continue
			}
			dstJson[fname] = basicValue
			continue
		}

		spec := ident.Obj.Decl.(*ast.TypeSpec)
		if ident, ok := spec.Type.(*ast.Ident); ok {
			if len(array) != 0 {
				array[0] = basicMap[ident.Name]
				dstJson[fname] = array
				continue
			}
			dstJson[fname] = basicMap[ident.Name]
			continue
		}
		if strct, ok := spec.Type.(*ast.StructType); ok {
			if len(array) != 0 {
				array[0] = convertMap(strct)
				dstJson[fname] = array
				continue
			}
			dstJson[fname] = convertMap(strct)
			continue
		}
		panic(fmt.Sprintf("not convert %v", f))
	}
	return dstJson
}
