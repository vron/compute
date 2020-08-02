package types

import "strings"

type GlslType struct {
	Name   string
	Export bool
	C      *CType
}

func (t GlslType) CName(prefix string) string {
	return prefix + t.Name
}

func (t GlslType) GoName() string {
	if t.Name == "Bool" {
		return "bool"
	}
	if t.Name == "uint32_t" {
		return "uint32"
	}
	if t.Name == "int32_t" {
		return "int32"
	}
	if t.Name == "float" {
		return "float32"
	}
	return strings.Title(t.Name)
}
