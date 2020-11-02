package testutil

type (
	OtherInteger          int64
	PtrOtherString        string
	OtherArraySamePkg     []OtherPkg
	OtherArrayPtrSampePkg []*PtrOtherPkg
	PtrOtherArraySamePkg  []OtherPkg
)

type OtherPkg struct {
	PkgContent string
}

type PtrOtherPkg struct {
	PkgContentDiff string
}

type CombindPkg struct {
	CombindContent string
}
