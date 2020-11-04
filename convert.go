package main

import (
	"encoding/json"
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

func StType2Map(path, stType string) (orderDefine, error) {
	fset := token.NewFileSet()
	return genOrderDefine(fset, path, stType, 0)
}

func genOrderDefine(fset *token.FileSet, path, stType string, ktype keyType) (orderDefine, error) {
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
		return convertMap(fset, imports, v, ktype)
	case *ast.InterfaceType:
		return nil, fmt.Errorf("%w: %s", ErrInterfaceType, stType)
	case *ast.Ident:
		return orderDefine{getStructDefine("", "", stType, basicMap[v.Name], nil, ktype)}, nil
	case *ast.ArrayType:
		var ident *ast.Ident
		if star, ok := v.Elt.(*ast.StarExpr); ok {
			ident = star.X.(*ast.Ident)
		} else {
			ident = v.Elt.(*ast.Ident)
		}
		ds, err := identToDefines(fset, imports, ident, "", "", nil, ktype)
		if err != nil {
			return nil, err
		}
		return orderDefine{getStructDefine("", "", stType, ds, allocArray(), ktype)}, nil
	default:
		return nil, fmt.Errorf("%#v is not sport type", fromSpec.Type)
	}
}

func convertMap(fset *token.FileSet, imap importMap, st *ast.StructType, ktype keyType) (orderDefine, error) {
	var defines orderDefine
	for _, f := range st.Fields.List {
		var tag string
		if f.Tag != nil {
			tag = f.Tag.Value
		}
		if len(f.Names) == 0 { // Embedded type
			var ident *ast.Ident
			ftype := f.Type
			if ptr, ok := ftype.(*ast.StarExpr); ok { // pointer
				ftype = ptr.X
			}
			switch v := ftype.(type) {
			case *ast.Ident:
				ident = v
			case *ast.SelectorExpr: // other package
				ident := v.X.(*ast.Ident)
				ds, err := findStructDefine(fset, imap[ident.Name], v.Sel.Name, ktype+exactType)
				if err != nil {
					return nil, err
				}
				if len(ds) == 1 && ds[0].field == "" {
					ds[0].field = v.Sel.Name
				}
				defines = append(defines, ds...)
				continue
			default:
				panic(fmt.Sprintf("unexpected type: %#v", f.Type))
			}

			if ident.Obj == nil ||
				strings.ToUpper(ident.Name)[0] != ident.Name[0] {
				// built-in type or private type
				continue
			}
			ds, err := identToDefines(fset, imap, ident, tag, "", nil, ktype+exactTagType)
			if err != nil {
				return nil, err
			}
			defines = append(defines, ds...)
			continue
		}

		if strings.ToUpper(f.Names[0].Name)[0] != f.Names[0].Name[0] {
			// private field
			continue
		}

		var (
			fname = f.Names[0].Name
			ident *ast.Ident
			array arrayDefines
		)

		ftype := f.Type
		if ptr, ok := ftype.(*ast.StarExpr); ok { // pointer
			ftype = ptr.X
		}
		switch v := ftype.(type) {
		case *ast.Ident:
			ident = v
		case *ast.ArrayType: // array
			if ptr, ok := v.Elt.(*ast.StarExpr); ok {
				ident = ptr.X.(*ast.Ident)
			} else {
				ident = v.Elt.(*ast.Ident)
			}
			array = make(arrayDefines, 1)
		case *ast.SelectorExpr: // other package type
			ident = v.X.(*ast.Ident)
			ds, err := findStructDefine(fset, imap[ident.Name], v.Sel.Name, ktype+normalType)
			if err != nil {
				return nil, err
			}
			typ := ident.Name + "." + v.Sel.Name
			defines = append(defines, getStructDefine(fname, tag, typ, ds, array, ktype+normalType))
			continue
		default:
			panic(fmt.Sprintf("unexpected type: %#v", f.Type))
		}

		if basicValue, ok := basicMap[ident.Name]; ok { // built-in type
			defines = append(defines,
				getStructDefine(fname, tag, ident.Name, basicValue, array, ktype+normalType))
			continue
		}

		ds, err := identToDefines(fset, imap, ident, tag, fname, array, ktype+normalType)
		if err != nil {
			return nil, err
		}
		defines = append(defines, ds...)
	}
	return defines, nil
}

func getStructDefine(fname, tag, typ string, value json.Marshaler, array arrayDefines, ktype keyType) fieldDefine {
	if len(array) != 0 {
		array[0] = value
		return fieldDefine{
			field:   fname,
			tag:     tag,
			typ:     typ,
			define:  array,
			enable:  true,
			keyType: ktype,
		}
	}
	return fieldDefine{
		field:   fname,
		tag:     tag,
		typ:     typ,
		define:  value,
		enable:  true,
		keyType: ktype,
	}
}

func identToDefines(fset *token.FileSet, imap importMap, ident *ast.Ident, tag string, fname string, array arrayDefines, ktype keyType) (orderDefine, error) {
	spec := ident.Obj.Decl.(*ast.TypeSpec)
	switch v := spec.Type.(type) {
	case *ast.Ident: // type built-in type
		if fname != "" {
			return orderDefine{getStructDefine(fname, tag, v.Name, basicMap[v.Name], array, ktype)}, nil
		}
		return orderDefine{getStructDefine(ident.Name, tag, v.Name, basicMap[v.Name], nil, ktype)}, nil

	case *ast.StructType: // type struct
		ds, err := convertMap(fset, imap, v, ktype)
		if err != nil {
			return nil, err
		}
		if fname != "" {
			return orderDefine{getStructDefine(fname, tag, ident.Name, ds, array, ktype)}, nil
		}
		return ds, nil
	case *ast.SelectorExpr: // other package
		ident := v.X.(*ast.Ident)
		ds, err := findStructDefine(fset, imap[ident.Name], v.Sel.Name, ktype)
		if err != nil {
			return nil, err
		}
		if fname != "" {
			return orderDefine{getStructDefine(fname, tag, ident.Name, ds, array, ktype)}, nil
		}
		return ds, nil
	case *ast.InterfaceType: // type interface
		return orderDefine{getStructDefine(fname, tag, ident.Name, nil, array, ktype)}, nil
	case *ast.ArrayType: // type array
		var arrayident *ast.Ident
		if ptr, ok := v.Elt.(*ast.StarExpr); ok { // pointer
			switch x := ptr.X.(type) {
			case *ast.SelectorExpr:
				arrayIdent := x.X.(*ast.Ident)
				ds, err := findStructDefine(fset, imap[arrayIdent.Name], x.Sel.Name, ktype)
				if err != nil {
					return nil, err
				}
				if fname != "" {
					return orderDefine{getStructDefine(fname, tag, ident.Name, ds, allocArray(), ktype)}, nil
				}
				return orderDefine{getStructDefine(ident.Name, tag, ident.Name, ds, allocArray(), ktype)}, nil
			case *ast.Ident:
				arrayident = x
			default:
				panic(fmt.Sprintf("not implemented type: %+v", x))
			}
		} else {
			arrayident = v.Elt.(*ast.Ident)
		}
		ds, err := convertMap(fset, imap, arrayident.Obj.Decl.(*ast.TypeSpec).Type.(*ast.StructType), ktype)
		if err != nil {
			return nil, err
		}
		if fname != "" {
			return orderDefine{getStructDefine(fname, tag, ident.Name, ds, allocArray(), ktype)}, nil
		}
		return orderDefine{getStructDefine(ident.Name, tag, ident.Name, ds, allocArray(), ktype)}, nil
	default:
		panic(fmt.Sprintf("unexpected type: %#v", spec.Type))
	}
}

func findStructDefine(fset *token.FileSet, dir, typ string, ktype keyType) (orderDefine, error) {
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
		define, err := genOrderDefine(fset, f, typ, ktype)
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
