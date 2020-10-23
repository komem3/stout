package main

import (
	"go/ast"
	"stout/internal/testutil"
)

type Embbading struct {
	IDString string
}

type PEmb struct {
	Content string
	Next    string
}

type Integer uint64

type PString string

type SampleJson struct {
	ID        int64  `custome:"custome_tag"`
	Content   string `json:"super_content" custom:"custom_tag"`
	Tag       int64  `json:"tag,omitempty,string"`
	Ignore    int64  `json:"-"`
	Hyphen    int8   `json:"-,"`
	PStr      *string
	Pointer   *Embbading
	Struct    Embbading
	Array     []string
	ArraySct  []Embbading
	ArrayPSct []*Embbading
	Other     ast.ArrayType
	Internal  testutil.OtherPkg
	Uinter    Integer
	Embbading
	Integer
	*PString
	*PEmb
	*testutil.OtherPkg
	private string
}

func pStr(s string) *string { return &s }
