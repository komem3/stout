package main

type Embbading struct {
	ID string
}

type PEmb struct {
	Content string
}

type Integer uint64

type SampleJson struct {
	ID        int64
	Content   string `json:"super_content"`
	PStr      *string
	Pointer   *Embbading
	Struct    Embbading
	Array     []string
	ArraySct  []Embbading
	ArrayPSct []*Embbading
	Uinter    Integer
	Embbading
	*PEmb
	private string
}

func pStr(s string) *string { return &s }
