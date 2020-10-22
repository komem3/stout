package main

import (
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"
)

var basicMap = map[string]interface{}{
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
	"complex64":  "1i",
	"complex128": "1i",
	"bool":       true,
	"byte":       byte('a'),
	"rune":       rune('a'),
}

func stType2Map(path, stType string) (orderDefine, error) {
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
		return nil, fmt.Errorf("not found %s\n", *fromSt)
	}

	st, ok := fromSpec.Type.(*ast.StructType)
	if !ok {
		return nil, fmt.Errorf("%s is not struct\n", *fromSt)
	}

	return convertMap(st), nil
}

func convertMap(st *ast.StructType) orderDefine {
	var defines orderDefine
	for _, f := range st.Fields.List {
		if len(f.Names) == 0 { // Embedded type
			var ident *ast.Ident
			if ptr, ok := f.Type.(*ast.StarExpr); ok {
				ident = ptr.X.(*ast.Ident)
			} else {
				ident = f.Type.(*ast.Ident)
			}

			if ident.Obj == nil { // built-in type
				continue
			}

			spec := ident.Obj.Decl.(*ast.TypeSpec)
			switch v := spec.Type.(type) {
			case *ast.Ident: // type built-in type
				defines = append(defines,
					getStructDefine(ident.Name, "", v.Name, basicMap[v.Name], nil))
			case *ast.StructType: // type struct
				ds := convertMap(v)
				for _, d := range ds {
					defines = append(defines,
						getStructDefine(d.field, "", d.typ, d.define, nil))
				}
			default:
				panic(fmt.Sprintf("unexpected type: %#v", spec.Type))
			}
			continue
		}
		if strings.ToUpper(f.Names[0].Name)[0] != f.Names[0].Name[0] {
			// private field
			continue
		}

		var (
			fname = f.Names[0].Name
			tag   string
			ident *ast.Ident
			array []interface{}
		)
		if f.Tag != nil {
			tag = f.Tag.Value
		}

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
		case *ast.SelectorExpr:
			panic(fmt.Sprintf("other package type: %#v", f.Type))
		default:
			panic(fmt.Sprintf("unexpected type: %#v", f.Type))
		}

		if basicValue, ok := basicMap[ident.Name]; ok { // built-in type
			defines = append(defines,
				getStructDefine(fname, tag, ident.Name, basicValue, array))
			continue
		}

		spec := ident.Obj.Decl.(*ast.TypeSpec)
		if ident, ok := spec.Type.(*ast.Ident); ok { // type built-in type
			defines = append(defines,
				getStructDefine(fname, tag, ident.Name, basicMap[ident.Name], array))
			continue
		}

		if strct, ok := spec.Type.(*ast.StructType); ok { // type struct
			defines = append(defines,
				getStructDefine(fname, tag, ident.Name, convertMap(strct), array))
			continue
		}
		panic(fmt.Sprintf("not convert %v", f))
	}
	return defines
}

func getStructDefine(fname, tag, typ string, value interface{}, array []interface{}) structDefine {
	if len(array) != 0 {
		array[0] = value
		return structDefine{
			field:  fname,
			tag:    tag,
			typ:    typ,
			define: array,
		}
	}
	return structDefine{
		field:  fname,
		tag:    tag,
		typ:    typ,
		define: value,
	}
}
