package types

import "fmt"

type CType struct {
	// TODO: Rename and remove Type in field name
	GlslType *GlslType // nil if not a direct glsl type

	// The different types of types a CType can represent, they
	// are mutually exclusive and only ever is one,
	Struct CStruct
	Array  CArray
	Vector CVector
	Basic  CBasicType

	Size Alignment
}

type CStruct struct {
	Fields []CField
}

type CField struct {
	Name       string
	CType      *CType
	ByteOffset int // offset in parent
}

type CArray struct {
	Len   int
	CType *CType
}

type CVector struct {
	Len   int
	Basic *CType
}

type CBasicType struct {
	Name string
}

func (ct *CType) IsStruct() bool {
	return ct.Struct.Fields != nil
}

func (ct *CType) IsArray() bool {
	return ct.Array.Len != 0
}

func (ct *CType) IsVector() bool {
	return ct.Vector.Len != 0
}

func (ct *CType) IsBasic() bool {
	return ct.Basic.Name != ""
}

func (ct *CType) CString(prefix, name string, prefixVectors bool) (s string) {
	if ct.IsBasic() {
		return ct.Basic.Name + "\t" + name
	}
	if ct.IsVector() {
		if prefixVectors {
			return prefix + ct.GlslType.Name + "\t" + name
		} else {
			return ct.GlslType.Name + "\t" + name
		}
	}
	if ct.IsStruct() {
		return prefix + ct.GlslType.Name + "\t" + name
	}
	if ct.IsArray() {
		tt := ct
		for ; tt.IsArray(); tt = tt.Array.CType {
		}
		s = tt.CString(prefix, "", prefixVectors)
		s += "\t"
		if ct.Array.Len == -1 {
			s += "(*"
		}
		s += name
		if ct.Array.Len == -1 {
			s += ")"
		}
		for ; ct.IsArray(); ct = ct.Array.CType {
			if ct.Array.Len > 0 {
				s += fmt.Sprintf("[%v]", ct.Array.Len)
			}
		}
		return s
	}
	panic("what type is this?")
}

func (ct *CType) GoString(name string) (s string) {
	if name != "" {
		name += "\t"
	}
	if ct.IsBasic() || ct.IsVector() || ct.IsStruct() {
		return name + ct.GlslType.GoName()
	}
	if ct.IsArray() {
		for ; ct.IsArray(); ct = ct.Array.CType {
			if ct.Array.Len == -1 {
				s += "[]"
			} else {
				s += fmt.Sprintf("[%v]", ct.Array.Len)
			}
		}
		return name + s + ct.GoString("")
	}
	panic("what type is this?")
}
