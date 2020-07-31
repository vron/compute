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

func (ct *CType) CString(prefix, name string) (s string) {
	if ct.IsBasic() {
		return ct.Basic.Name + "\t" + name
	}
	if ct.IsVector() {
		return prefix + ct.GlslType.Name
	}
	if ct.IsStruct() {
		return prefix + ct.GlslType.Name
	}
	if ct.IsArray() {
		tt := ct
		for ; tt.IsArray(); tt = tt.Array.CType {
		}
		s = tt.CString(prefix, "")
		s += "\t"
		if ct.Array.Len == -1 {
			s += "*"
		}
		s += name
		for ; ct.IsArray(); ct = ct.Array.CType {
			if ct.Array.Len > 0 {
				s += fmt.Sprintf("[%v]", ct.Array.Len)
			}
		}
		return s
	}
	panic("what type is this?")
}

/*
func (ct CType) Shared() string {
	s := ct.Name + " "
	if ct.IsSlice() {
		s += "*"
	}
	s += ct.Name
	if ct.ArrayLen() != 0 {
		for el := ct; el.ArrayLen() != 0; el = *el.Array.CType {
			s += fmt.Sprintf("[%v]", el.Array.Len)
		}
	}
	return s
}

func (ct CType) GoCTypeName() string {
	s := ""
	if ct.IsSlice() {
		s += "*"
	}
	if ct.ArrayLen() != 0 {
		for el := ct; el.ArrayLen() != 0; el = *el.Array.CType {
			s += fmt.Sprintf("[%v]", el.Array.Len)
		}
	}
	s += "C." + ct.Name
	return s
}



func (cf CField) CxxFieldString() string {
	if cf.CType.IsSlice() {
		// there are no built-in types with arbitrary length so this is fine
		return fmt.Sprintf("\t%v* %v;\n", cf.CType.GlslType.Name, cf.Name)
	}
	s := "\t" + cf.CType.GlslType.Name + " " + cf.Name

	// NoElem here is in C world, e.g. a vec has a length, so what we want
	// to chec is if the length is equal to the one for the defined type or
	// not
	if cf.CType.ArrayLen() != 0 {
		for el := cf.CType; el.ArrayLen() != 0; el = el.Array.CType {
			s += fmt.Sprintf("[%v]", el.Array.Len)
		}
	}
	return s + ";"
}

func (cf CField) CxxFieldStringRef() string {
	if cf.CType.IsSlice() {
		// there are no built-in types with arbitrary length so this is fine
		return fmt.Sprintf("\t%v* (&%v);\n", cf.CType.GlslType.Name, cf.Name)
	}
	s := "\t" + cf.CType.GlslType.Name + "(& " + cf.Name + ")"

	// NoElem here is in C world, e.g. a vec has a length, so what we want
	// to chec is if the length is equal to the one for the defined type or
	// not
	if cf.CType.ArrayLen() != 0 {
		for el := cf.CType; el.ArrayLen() != 0; el = el.Array.CType {
			s += fmt.Sprintf("[%v]", el.Array.Len)
		}
	}
	return s + ";"
}

func (cf CField) GoName() string {
	return strings.Title(cf.Name)
}

func (ct CType) BasicGoType() string {
	if ct.IsSlice() {
		// this is a slice, so will be passed as pointer, but since
		// we cannot ensure that the caller has the same struct alignment
		// that we require we pass it as a void* and require the caller
		// to ensure that the data has the same layout as what we request.
		return "byte"
	}
	var goTypeName string
	if len(ct.Fields) > 0 {
		goTypeName = strings.Title(strings.TrimPrefix(ct.Name, "cpt_"))
	} else if ct.GlslType.Name == "Bool" {
		goTypeName = "bool"
	} else {
		goTypeName = strings.TrimSuffix(ct.Name, "_t")
		if !strings.HasSuffix(goTypeName, "32") {
			goTypeName += "32"
		}
	}
	return goTypeName

}

func (ct CType) GoName() string {
	if ct.IsSlice() {
		return "[]" + ct.BasicGoType()
	}
	var goTypeName = ct.BasicGoType()
	if ct.ArrayLen() != 0 {
		goTypeName = ""
		for el := &ct; el.ArrayLen() != 0; el = el.Array.CType {
			goTypeName += fmt.Sprintf("[%v]", el.Array.Len)
		}
		goTypeName += fmt.Sprintf("%v", ct.BasicGoType())
	}
	return goTypeName
}

*/
