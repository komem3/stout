package main

import (
	"go/ast"
	"stout/internal/testutil"
)

type Integer uint64

type PString string

type JsonTag struct {
	ID      int64   `custome:"custome_tag"`
	Content string  `json:"super_content" custom:"custom_tag"`
	Tag     int64   `json:"tag,omitempty,string"`
	Ignore  int64   `json:"-"`
	Hyphen  int8    `json:"-,"`
	Bool    bool    `json:",string"`
	Byte    byte    `json:",string"`
	Error   error   `json:",string"`
	Float32 float32 `json:",string"`
	Float64 float64 `json:",string"`
	Int16   int16   `json:",string"`
	Int32   int32   `json:",string"`
	Int64   int64   `json:",string"`
	Int8    int8    `json:",string"`
	Rune    rune    `json:",string"`
	String  string  `json:",string"`
	Uint    uint    `json:",string"`
	Uint16  uint16  `json:",string"`
	Uint32  uint32  `json:",string"`
	Uint64  uint64  `json:",string"`
	Uint8   uint8   `json:",string"`
	Uintptr uintptr `json:",string"`
}

type BuiltInTypes struct {
	Bool bool
	Byte byte
	// Complex128 complex128 unsupported json
	// Complex64 complex64 unsupported json
	Error   error
	Float32 float32
	Float64 float64
	Int16   int16
	Int32   int32
	Int64   int64
	Int8    int8
	Rune    rune
	String  string
	Uint    uint
	Uint16  uint16
	Uint32  uint32
	Uint64  uint64
	Uint8   uint8
	Uintptr uintptr
}

type private int
type ArraySamePkg []SamePkg
type ArrayPtrSampePkg []*PtrSamePkg
type PtrArraySamePkg []SamePkg

type SamePkg struct {
	IDString string
}

type PtrSamePkg struct {
	Content string
	Next    string
}

type CombineCase testutil.CombindPkg

type SimpleStruct struct {
	pfield private

	// same pkg
	IntF          Integer
	PstrF         *PString
	SameF         SamePkg
	PSameF        *PtrSamePkg
	ArraySameF    ArraySamePkg
	ArrayPtrSameF ArrayPtrSampePkg
	PtrArraySameF *PtrArraySamePkg

	// other pkg
	OtherIntF      testutil.OtherInteger
	PrtOtherStrF   *testutil.PtrOtherString
	OtherF         testutil.OtherPkg
	PtrOtherF      *testutil.PtrOtherPkg
	OtherArryF     testutil.OtherArraySamePkg
	OtherArrayPtrF testutil.OtherArrayPtrSampePkg
	PrtOtherArrayF *testutil.PtrOtherArraySamePkg

	// combine
	CombinedF CombineCase

	AstF ast.ArrayType
}

type Emmbaddings struct {
	int
	private

	// same pkg
	Integer
	*PString
	SamePkg
	*PtrSamePkg
	ArraySamePkg
	ArrayPtrSampePkg
	*PtrArraySamePkg

	// other pkg
	testutil.OtherInteger
	*testutil.PtrOtherString
	testutil.OtherPkg
	*testutil.PtrOtherPkg
	testutil.OtherArraySamePkg
	testutil.OtherArrayPtrSampePkg
	*testutil.PtrOtherArraySamePkg

	// combine
	CombineCase

	ast.ArrayType
}

type DuplicateFields struct {
	Field     int
	Duplicate string `json:"Field"`
}

func pStr(s string) *string { return &s }
