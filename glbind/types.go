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
	for i, str := range inp.Structs {
		recurseCreateStructTypes(i+1, inp, str)
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
	Name         string // the name as refered to it in GLSL
	apiType      bool   // true if the type is used in the api and should be exported
	userType     bool   // type created in glsl
	cType        *CType
	userStructId int
}

// A CType represents the type defined in shared.h that will translate
// to a go type.
type CType struct {
	ty     *Type
	Name   string
	Fields []CField // len = 0 not an struct type
	Array  ArrayType
	Size   Alignment
}

func (ct *CType) isArray() bool {
	return ct.Array.len != 0
}

func (ct *CType) ArraySize() int {
	if !ct.isArray() {
		return 0
	}
	// recursively get the size
	size := 1
	for el := ct; el.isArray(); el = el.Array.ty {
		size *= el.Array.len
	}
	return size
}

func (ct *CType) isSlice() bool {
	return ct.Array.len == -1
}

type ArrayType struct {
	ty  *CType
	len int
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
	sort.Slice(ts, less(ts))
	return
}

func (tt *Types) AllTypes() (ts []*Type) {
	for _, v := range tt.types {
		ts = append(ts, v)
	}
	sort.Slice(ts, less(ts))
	return
}

func (tt *Types) UserStructs() (ts []*Type) {
	for _, v := range tt.types {
		if v.userType {
			ts = append(ts, v)
		}
	}
	sort.Slice(ts, less(ts))
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
	if ct.isSlice() {
		s += "*"
	}
	s += ct.Name
	if ct.isArray() {
		for el := ct; el.isArray(); el = *el.Array.ty {
			s += fmt.Sprintf("[%v]", el.Array.len)
		}
	}
	return s
}

func (ct CType) GoCTypeName() string {
	s := ""
	if ct.isSlice() {
		s += "*"
	}
	if ct.isArray() {
		for el := ct; el.isArray(); el = *el.Array.ty {
			s += fmt.Sprintf("[%v]", el.Array.len)
		}
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
	if cf.Ty.isSlice() {
		// this is a slice, so will be passed as pointer, but since
		// we cannot ensure that the caller has the same struct alignment
		// that we require we pass it as a void* and require the caller
		// to ensure that the data has the same layout as what we request.
		return "void* " + cf.Name
	}
	s := cf.Ty.Name + " " + cf.Name
	if cf.Ty.isArray() {
		for el := cf.Ty; el.isArray(); el = el.Array.ty {
			s += fmt.Sprintf("[%v]", el.Array.len)
		}
	}
	return s
}

func (cf CField) CxxFieldString() string {
	if cf.Ty.isSlice() {
		// there are no built-in types with arbitrary length so this is fine
		return fmt.Sprintf("\t%v* %v;\n", cf.Ty.ty.Name, cf.Name)
	}
	s := "\t" + cf.Ty.ty.Name + " " + cf.Name

	// NoElem here is in C world, e.g. a vec has a length, so what we want
	// to chec is if the length is equal to the one for the defined type or
	// not
	if cf.Ty.isArray() {
		for el := cf.Ty; el.isArray(); el = el.Array.ty {
			s += fmt.Sprintf("[%v]", el.Array.len)
		}
	}
	return s + ";"
}

func (cf CField) CxxFieldStringRef() string {
	if cf.Ty.isSlice() {
		// there are no built-in types with arbitrary length so this is fine
		return fmt.Sprintf("\t%v* (&%v);\n", cf.Ty.ty.Name, cf.Name)
	}
	s := "\t" + cf.Ty.ty.Name + "(& " + cf.Name + ")"

	// NoElem here is in C world, e.g. a vec has a length, so what we want
	// to chec is if the length is equal to the one for the defined type or
	// not
	if cf.Ty.isArray() {
		for el := cf.Ty; el.isArray(); el = el.Array.ty {
			s += fmt.Sprintf("[%v]", el.Array.len)
		}
	}
	return s + ";"
}

func (cf CField) GoName() string {
	return strings.Title(cf.Name)
}

func (ct CType) BasicGoType() string {
	if ct.isSlice() {
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
	if ct.isSlice() {
		return "[]" + ct.BasicGoType()
	}
	var goTypeName = ct.BasicGoType()
	if ct.isArray() {
		goTypeName = ""
		for el := &ct; el.isArray(); el = el.Array.ty {
			goTypeName += fmt.Sprintf("[%v]", el.Array.len)
		}
		goTypeName += fmt.Sprintf("%v", ct.BasicGoType())
	}
	return goTypeName
}

func createBasicBuiltinTypes() {
	cbt := func(name string, cname string, noelc int, size int, align int) {
		tt := Type{
			Name: name,
		}
		ut := CType{
			ty:   &tt,
			Name: cname,
			Size: Alignment{ByteSize: 4, ByteAlignment: 4},
		}
		ttt := CType{
			ty:    &tt,
			Name:  cname,
			Array: ArrayType{len: noelc, ty: &ut},
			Size:  Alignment{ByteSize: size, ByteAlignment: align},
		}

		tt.cType = &ttt
		types.types[name] = &tt
	}
	cbt("Bool", "int32_t", 0, 4, 4)
	cbt("int32_t", "int32_t", 0, 4, 4)
	cbt("uint32_t", "uint32_t", 0, 4, 4)
	cbt("float", "float", 0, 4, 4)
	cbt("vec2", "float", 2, 8, 8)
	cbt("vec3", "float", 3, 16, 16)
	cbt("vec4", "float", 4, 16, 16)
	cbt("ivec2", "int32_t", 2, 8, 8)
	cbt("ivec3", "int32_t", 3, 16, 16)
	cbt("ivec4", "int32_t", 4, 16, 16)
	cbt("uvec2", "uint32_t", 2, 8, 8)
	cbt("uvec3", "uint32_t", 3, 16, 16)
	cbt("uvec4", "uint32_t", 4, 16, 16)
	cbt("mat2", "float", 4, 16, 8)
	cbt("mat3", "float", 9, 16*3, 16)
	cbt("mat4", "float", 16, 16*4, 16)
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

func recCreateArrayType(ty *CType, noels []int) {
	if len(noels) == 0 {
		return
	}
	nt := *(ty) // a copy of it

	noel := noels[0]
	noels = noels[1:]

	ty.Array.len = noel
	ty.Array.ty = &nt
	recCreateArrayType(&nt, noels)

	// propagate the size all the way bac up
	ty.Size.ByteSize = nt.Size.ByteSize * noel

}

func maybeCreateArrayType(ty string, noels []int) *CType {
	// create a new ctype representing this type with an array one
	if len(noels) == 0 {
		return types.Get(ty).CType()
	}

	// recurse down to create this one as is
	bt := *types.Get(ty).CType()
	noel := noels[0]
	if noel < 0 {
		// a slice can only be the first available one!
		// TODO: 64 bit only - else change alignment?
		recCreateArrayType(&bt, noels)
		bt.Size.ByteAlignment = 8
		bt.Size.ByteSize = 8
	} else {
		recCreateArrayType(&bt, noels)
	}
	return &bt
}

func recurseCreateStructTypes(i int, inp Input, str InputStruct) {
	if _, o := types.types[str.Name]; o {
		return
	}
	st := Type{
		Name:         str.Name,
		userType:     true,
		userStructId: i,
	}
	ct := CType{
		ty:     &st,
		Name:   "cpt_" + str.Name,
		Fields: []CField{},
	}
	for _, f := range str.Fields {
		if _, o := types.types[f.Ty]; !o {
			panic("structs in input must be defined in dependicy order")
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
	if ct.ArraySize() > 0 {
		rsize.ByteSize *= ct.ArraySize()
	}
	if ct.isSlice() {
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

func (cf CField) CxxArrayLen() int {
	rt := types.Get(cf.Ty.ty.Name)
	if rt.CType().ArraySize() != cf.Ty.ArraySize() {
		if rt.CType().ArraySize() != 0 {
			size := cf.Ty.ArraySize() / rt.CType().ArraySize()
			if size*rt.CType().ArraySize() != cf.Ty.ArraySize() {
				panic("this cannot be an array, what is happening?")
			}
			return size
		} else {
			return cf.Ty.ArraySize()
		}
	} else {
		return 0
	}
}
