package main

import (
	"errors"
	"fmt"
	"go/ast"
	"go/parser"
	"go/token"
	"strings"

	"golang.org/x/tools/go/packages"
)

type importMap map[string]string

var (
	ErrNotFoundStruct = fmt.Errorf("not found struct")
	ErrInterfaceType  = fmt.Errorf("interface type can not convert json")
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

func StType2Map(path, stType string) (orderDefine, error) {
	fset := token.NewFileSet()
	return genOrderDefine(fset, path, stType)
}

func genOrderDefine(fset *token.FileSet, path, stType string) (orderDefine, error) {
	f, err := parser.ParseFile(fset, path, nil, 0)
	if err != nil {
		return nil, err
	}

	imports := make(importMap)
	for _, imp := range f.Imports {
		path := strings.Trim(imp.Path.Value, "\"")
		index := strings.LastIndex(path, "/")
		pkgName := path
		if index > 0 {
			pkgName = path[index+1:]
		}
		imports[pkgName] = path
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
		return nil, fmt.Errorf("%w：%s", ErrNotFoundStruct, stType)
	}

	switch v := fromSpec.Type.(type) {
	case *ast.StructType:
		return convertMap(fset, imports, v)
	case *ast.InterfaceType:
		return nil, fmt.Errorf("%w: %s", ErrInterfaceType, stType)
	case *ast.Ident:
		return append(orderDefine{}, getStructDefine("", "", stType, basicMap[v.Name], nil)), nil
	default:
		return nil, fmt.Errorf("%#v is not sport type", fromSpec.Type)
	}
}

func convertMap(fset *token.FileSet, imap importMap, st *ast.StructType) (orderDefine, error) {
	var defines orderDefine
	for _, f := range st.Fields.List {
		var tag string
		if f.Tag != nil {
			tag = f.Tag.Value
		}
		if len(f.Names) == 0 { // Embedded type
			var ident *ast.Ident
			switch v := f.Type.(type) {
			case *ast.Ident:
				ident = v
			case *ast.StarExpr: // pointer
				if selExpr, ok := v.X.(*ast.SelectorExpr); ok {
					ident := selExpr.X.(*ast.Ident)
					ds, err := findStructDefine(fset, imap[ident.Name], selExpr.Sel.Name)
					if err != nil {
						return nil, err
					}
					defines = append(defines, ds...)
					continue
				}
				ident = v.X.(*ast.Ident)
			case *ast.SelectorExpr: // other package
				ident := v.X.(*ast.Ident)
				ds, err := findStructDefine(fset, imap[ident.Name], v.Sel.Name)
				if err != nil {
					return nil, err
				}
				defines = append(defines, ds...)
				continue
			default:
				panic(fmt.Sprintf("unexpected type: %#v", f.Type))
			}

			if ident.Obj == nil { // built-in type
				continue
			}

			spec := ident.Obj.Decl.(*ast.TypeSpec)
			switch v := spec.Type.(type) {
			case *ast.Ident: // type built-in type
				defines = append(defines,
					getStructDefine(ident.Name, tag, v.Name, basicMap[v.Name], nil))
			case *ast.StructType: // type struct
				ds, err := convertMap(fset, imap, v)
				if err != nil {
					return nil, err
				}
				for _, d := range ds {
					defines = append(defines,
						getStructDefine(d.field, tag, d.typ, d.define, nil))
				}
			case *ast.ArrayType: // type array
				if ptr, ok := v.Elt.(*ast.StarExpr); ok { // pointer
					sel := ptr.X.(*ast.SelectorExpr)
					arrayIdent := sel.X.(*ast.Ident)
					ds, err := findStructDefine(fset, imap[arrayIdent.Name], sel.Sel.Name)
					if err != nil {
						return nil, err
					}
					defines = append(defines,
						getStructDefine(ident.Name, tag, ident.Name, ds, make([]interface{}, 1)))
					continue
				}
				arrayident := v.Elt.(*ast.Ident)
				ds, err := convertMap(fset, imap, arrayident.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType))
				if err != nil {
					return nil, err
				}
				defines = append(defines,
					getStructDefine(ident.Name, tag, ident.Name, ds, make([]interface{}, 1)))
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
			ident *ast.Ident
			array []interface{}
		)

		switch v := f.Type.(type) {
		case *ast.StarExpr: // pointer
			ident = v.X.(*ast.Ident)
		case *ast.Ident:
			ident = v
		case *ast.ArrayType: // array
			if ptr, ok := v.Elt.(*ast.StarExpr); ok {
				ident = ptr.X.(*ast.Ident)
			} else {
				ident = v.Elt.(*ast.Ident)
			}
			array = make([]interface{}, 1)
		case *ast.SelectorExpr: // other package type
			ident := v.X.(*ast.Ident)
			ds, err := findStructDefine(fset, imap[ident.Name], v.Sel.Name)
			if err != nil {
				return nil, err
			}
			typ := ident.Name + "." + v.Sel.Name
			defines = append(defines, getStructDefine(fname, tag, typ, ds, array))
			continue
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

		switch v := spec.Type.(type) {
		case *ast.StructType:
			ds, err := convertMap(fset, imap, v)
			if err != nil {
				return nil, err
			}
			defines = append(defines,
				getStructDefine(fname, tag, ident.Name, ds, array))
		case *ast.InterfaceType:
			defines = append(defines,
				getStructDefine(fname, tag, ident.Name, nil, array))
		default:
			panic(fmt.Sprintf("not convert %#v", f))
		}
	}
	return defines, nil
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

func otherPkgDefine(
	fset *token.FileSet,
	imap importMap,
	v *ast.SelectorExpr,
	fname string,
	tag string,
	array []interface{},
) (structDefine, error) {
	ident := v.X.(*ast.Ident)
	ds, err := findStructDefine(fset, imap[ident.Name], v.Sel.Name)
	if err != nil {
		return structDefine{}, err
	}
	typ := ident.Name + "." + v.Sel.Name
	return getStructDefine(fname, tag, typ, ds, array), nil
}

func findStructDefine(fset *token.FileSet, dir, typ string) (orderDefine, error) {
	cfg := &packages.Config{
		Mode: packages.LoadFiles | packages.LoadImports,
		Fset: fset,
	}
	pkgs, err := packages.Load(cfg, dir)
	if err != nil {
		return nil, err
	}

	if packages.PrintErrors(pkgs) > 0 {
		return nil, fmt.Errorf("package load error")
	}

	if len(pkgs) == 0 {
		return nil, fmt.Errorf("%s package is not found", dir)
	}
	for _, f := range pkgs[0].GoFiles {
		define, err := genOrderDefine(fset, f, typ)
		if errors.Is(err, ErrNotFoundStruct) {
			continue
		}
		if err != nil {
			return nil, err
		}
		return define, nil
	}

	return nil, fmt.Errorf("%w: %s", ErrNotFoundStruct, typ)
}
