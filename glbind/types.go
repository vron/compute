// TODO: To e.g support multiple levels of arrays we must do proper parsing and not merge array lengths etc...

package main

import (
	"fmt"
	"sort"
	"strings"
)

var types *Types

func parseTypeInfo(inp Input) {
	types = &Types{types: map[string]*Type{}}
	createBasicBuiltinTypes()
	createComplexBuiltinTypes()
	// create the user defined struct types
	for _, str := range inp.Structs {
		recurseCreateStructTypes(inp, str)
	}
	types.calculateAlignments()
	findApiExportedTypes(inp)
}

// Types contains infomration about all types in use.
type Types struct {
	types map[string]*Type
}

// A Type represents the GLSL type as parsed from the shader code
type Type struct {
	Name     string // the name as refered to it in GLSL
	apiType  bool   // true if the type is used in the api and should be exported
	userType bool   // type created in glsl
	cType    *CType
}

// A CType represents the type defined in shared.h that will translate
// to a go type.
type CType struct {
	ty       *Type
	Name     string
	Fields   []CField // len = 0 not an struct type
	ArrayLen int
	IsSlice  bool
	Size     Alignment
}

type CField struct {
	Name string
	Ty   *CType
	// offset in parent
	ByteOffset int
}

type Alignment struct {
	ByteSize      int
	ByteAlignment int
}

func (ct CType) GetSize() int {
	size := ct.Size.ByteSize
	if ct.ArrayLen > 0 {
		size *= ct.ArrayLen
	}
	if ct.IsSlice {
		panic("getsize")
		return 8 // TODO: This will brea on non-64 bit platforms!
	}
	return size
}

func (tt *Types) Get(ty string) *Type {
	v, o := tt.types[ty]
	if !o {
		panic("tried to access undefined type: " + ty)
	}
	return v
}

func (tt *Types) ExportedStructTypes() (ts []*Type) {
	for _, v := range tt.types {
		if v.apiType && len(v.CType().Fields) > 0 {
			ts = append(ts, v)
		}
	}
	sort.Slice(ts, func(i, j int) bool { return ts[i].Name < ts[j].Name })
	return
}

func (tt *Types) AllTypes() (ts []*Type) {
	for _, v := range tt.types {
		ts = append(ts, v)
	}
	sort.Slice(ts, func(i, j int) bool { return ts[i].Name < ts[j].Name })
	return
}

func (tt *Types) UserStructs() (ts []*Type) {
	for _, v := range tt.types {
		if v.userType {
			ts = append(ts, v)
		}
	}
	sort.Slice(ts, func(i, j int) bool { return ts[i].Name < ts[j].Name })
	return
}

func (tt *Types) put(ty string, t *Type) {
	_, o := tt.types[ty]
	if o {
		panic("trying to create type with same name: " + ty)
	}
	tt.types[ty] = t
}

func (ct CType) String() string {
	s := ct.Name + " "
	if ct.IsSlice {
		s += "*"
	}
	s += ct.Name
	if ct.ArrayLen > 0 {
		s += fmt.Sprintf("[%v]", ct.ArrayLen)
	}
	return s
}

func (ct CType) GoCTypeName() string {
	s := ""
	if ct.IsSlice {
		s += "*"
	} else if ct.ArrayLen > 0 {
		s += fmt.Sprintf("[%v]", ct.ArrayLen)
	}
	s += "C." + ct.Name
	return s
}

func (t Type) CType() *CType {
	return t.cType
}

func (t Type) GoName() string {
	return strings.Title(t.Name)
}

func (cf CField) String() string {
	if cf.Ty.IsSlice {
		// this is a slice, so will be passed as pointer, but since
		// we cannot ensure that the caller has the same struct alignment
		// that we require we pass it as a void* and require the caller
		// to ensure that the data has the same layout as what we request.
		return "void* " + cf.Name
	}
	s := cf.Ty.Name + " " + cf.Name
	if cf.Ty.ArrayLen > 0 {
		s += fmt.Sprintf("[%v]", cf.Ty.ArrayLen)
	}
	return s
}

func (cf CField) CxxArrayLen() int {
	rt := types.Get(cf.Ty.ty.Name)
	if rt.CType().ArrayLen != 0 && rt.CType().ArrayLen != cf.Ty.ArrayLen {
		size := cf.Ty.ArrayLen / rt.CType().ArrayLen
		if size*rt.CType().ArrayLen != cf.Ty.ArrayLen {
			panic("this cannot be an array, what is happening?")
		}
		return size
	} else {
		return 0
	}
}

func (cf CField) CxxFieldString() string {
	if cf.Ty.IsSlice {
		// there are no built-in types with arbitrary length so this is fine
		return fmt.Sprintf("\t%v* %v;\n", cf.Ty.ty.Name, cf.Name)
	}
	s := "\t" + cf.Ty.ty.Name + " " + cf.Name

	// NoElem here is in C world, e.g. a vec has a length, so what we want
	// to chec is if the length is equal to the one for the defined type or
	// not
	if ll := cf.CxxArrayLen(); ll > 0 {
		s += fmt.Sprintf("[%v]", ll)
	}
	return s + ";"
}

func (cf CField) GoName() string {
	return strings.Title(cf.Name)
}

func (ct CType) BasicGoType() string {
	if ct.IsSlice {
		// this is a slice, so will be passed as pointer, but since
		// we cannot ensure that the caller has the same struct alignment
		// that we require we pass it as a void* and require the caller
		// to ensure that the data has the same layout as what we request.
		return "byte"
	}
	var goTypeName string
	if len(ct.Fields) > 0 {
		goTypeName = strings.Title(strings.TrimPrefix(ct.Name, "cpt_"))
	} else if ct.ty.Name == "Bool" {
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
	if ct.IsSlice {
		return "[]" + ct.BasicGoType()
	}
	var goTypeName = ct.BasicGoType()
	if ct.ArrayLen > 0 {
		goTypeName = fmt.Sprintf("[%v]%v", ct.ArrayLen, goTypeName)
	}
	return goTypeName
}

func createBasicBuiltinTypes() {
	cbt := func(name string, cname string, noelc int, size int) {
		tt := Type{
			Name: name,
		}
		tt.cType = &CType{
			ty:       &tt,
			Name:     cname,
			Fields:   []CField{},
			ArrayLen: noelc,
			Size:     Alignment{ByteSize: size, ByteAlignment: size},
		}
		types.types[name] = &tt
	}
	cbt("Bool", "int32_t", 0, 4) // here - how do we create go name?
	cbt("int32_t", "int32_t", 0, 4)
	cbt("uint32_t", "uint32_t", 0, 4)
	cbt("float", "float", 0, 4)
	cbt("vec2", "float", 2, 8)
	cbt("vec3", "float", 3, 16)
	cbt("vec4", "float", 4, 16)
	cbt("ivec2", "int32_t", 2, 8)
	cbt("ivec3", "int32_t", 3, 16)
	cbt("ivec4", "int32_t", 4, 16)
	cbt("uvec2", "uint32_t", 2, 8)
	cbt("uvec3", "uint32_t", 3, 16)
	cbt("uvec4", "uint32_t", 4, 16)
}

func createComplexBuiltinTypes() {
	co := func(name string, args ...interface{}) {
		tt := Type{
			Name: name,
		}
		ct := CType{
			ty:     &tt,
			Name:   "cpt_" + name,
			Fields: []CField{},
		}
		tt.cType = &ct
		for i := 0; i+2 < len(args); i += 3 {
			noel := args[i+2].(int)
			name := args[i].(string)
			ty := args[i+1].(string)
			f := CField{
				Name: name,
				Ty:   types.Get(ty).CType(),
			}
			if noel != 0 {
				f.Ty = maybeCreateArrayType(ty, []int{noel})
			}
			ct.Fields = append(ct.Fields, f)
		}
		types.put(name, &tt)
	}
	co("image2Drgba32f", "data", "float", -1, "width", "int32_t", 0)
}

func maybeCreateArrayType(ty string, noels []int) *CType {
	// create a new ctype representing this type with an array one
	if len(noels) == 0 {
		return types.Get(ty).CType()
	}

	// accumulate the number of elelemnts
	noel := 1
	for i, v := range noels {
		if v < 0 && i < len(noels)-1 {
			panic("only support slice as last type of array yet")
		}
		noel *= v
	}
	iss := false
	if noel < 0 {
		noel *= -1
		iss = true
	}

	bt := *types.Get(ty).CType()
	par := *bt.ty
	if bt.ArrayLen == 0 {
		bt.ArrayLen = noel
	} else if bt.IsSlice {
		panic("we cannot slice a slice")
	} else {
		bt.ArrayLen *= noel
	}

	bt.Size.ByteSize *= noel

	bt.ty = &par
	bt.IsSlice = iss

	if bt.IsSlice {
		// TODO: 64 bit only
		bt.Size.ByteAlignment = 8
		bt.Size.ByteSize = 8
	}

	return &bt
}

func recurseCreateStructTypes(inp Input, str InputStruct) {
	if _, o := types.types[str.Name]; o {
		return
	}
	st := Type{
		Name:     str.Name,
		userType: true,
	}
	ct := CType{
		ty:     &st,
		Name:   "cpt_" + str.Name,
		Fields: []CField{},
	}
	for _, f := range str.Fields {
		if _, o := types.types[f.Ty]; !o {
			for _, s := range inp.Structs {
				if s.Name == f.Ty {
					recurseCreateStructTypes(inp, s)
					break
				}
			}
			panic("user struct refered to unspecified type: " + f.Ty)
		}
		cf := CField{
			Name: f.Name,
			Ty:   types.Get(f.Ty).CType(),
		}
		if len(f.Arrno) != 0 {
			cf.Ty = maybeCreateArrayType(f.Ty, f.Arrno)
		}
		ct.Fields = append(ct.Fields, cf)
	}
	st.cType = &ct
	types.put(str.Name, &st)
}

func (tt *Types) calculateAlignments() {
	for _, t := range tt.types {
		// all basic types have alignment, so this must be a struct type:
		alignVisit(t.CType())
	}
}

func alignVisit(ct *CType) Alignment {
	// from https://www.oreilly.com/library/view/opengl-programming-guide/9780132748445/app09lev1sec3.html
	// but we treat vec3 as vec4 to allow for vector operations
	if ct.Size != (Alignment{}) {
		return ct.Size
	}

	offset := 0
	maxFieldAlignment := 0
	for i, f := range ct.Fields {
		fa := alignVisit(f.Ty) // => we cannot handle recursive types, thats fine...

		// TOOD: We must handle arrays here!

		for ; offset%fa.ByteAlignment != 0; offset++ {
		}

		ct.Fields[i].ByteOffset = offset
		if fa.ByteAlignment > maxFieldAlignment {
			maxFieldAlignment = fa.ByteAlignment
		}
		offset += fa.ByteSize
	}

	for ; offset%maxFieldAlignment != 0; offset++ {
	}

	ct.Size = Alignment{
		ByteSize:      offset,
		ByteAlignment: maxFieldAlignment,
	}

	rsize := ct.Size
	if ct.ArrayLen > 0 {
		rsize.ByteSize *= ct.ArrayLen
	}
	if ct.IsSlice {
		// TODO: this is 64 bit only...
		rsize.ByteSize = 8
		rsize.ByteAlignment = 8
	}

	ct.Size = rsize
	return rsize
}

func findApiExportedTypes(inp Input) {
	for _, arg := range inp.Arguments {
		exportType(arg.Ty)
	}
}

func exportType(s string) {
	t := types.Get(s)
	if t.apiType {
		return
	}
	t.apiType = true
	for _, f := range t.cType.Fields {
		exportType(f.Ty.ty.Name)
	}
}
