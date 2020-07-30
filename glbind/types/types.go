// TODO: To e.g support multiple levels of arrays we must do proper parsing and not merge array lengths etc...

package types

import (
	"fmt"
	"sort"
	"strings"

	"github.com/vron/compute/glbind/input"
)

type Types struct {
	types map[string]*GlslType
}

func New(inp input.Input) *Types {
	types := &Types{types: map[string]*GlslType{}}
	types.createBasicBuiltinTypes()
	types.createComplexBuiltinTypes()
	for i, str := range inp.Structs {
		types.recurseCreateStructTypes(i+1, inp, str)
	}
	types.calculateAlignments()
	types.findApiExportedTypes(inp)

	return types
}

type GlslType struct {
	Name        string
	Export      bool
	UserDefined bool
	C           *CType
	GlslOrder   int
}

type CType struct {
	GlslType *GlslType
	Name     string

	Fields []CField
	Array  ArrayType

	Size Alignment
}

func (ct *CType) ArrayLen() int {
	return ct.Array.Len
}

func (ct *CType) ArraySize() int {
	if ct.ArrayLen() != 0 {
		return 0
	}
	// recursively get the size
	size := 1
	for el := ct; el.ArrayLen() != 0; el = el.Array.Ty {
		size *= el.Array.Len
	}
	return size
}

func (ct *CType) IsSlice() bool {
	return ct.Array.Len == -1
}

type ArrayType struct {
	Ty  *CType
	Len int
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

func (tt *Types) Get(ty string) *GlslType {
	v, o := tt.types[ty]
	if !o {
		panic("tried to access undefined type: " + ty)
	}
	return v
}

func (tt *Types) ExportedStructTypes() (ts []*GlslType) {
	for _, v := range tt.types {
		if v.Export && len(v.CType().Fields) > 0 {
			ts = append(ts, v)
		}
	}
	sort.Slice(ts, less(ts))
	return
}

func (tt *Types) AllTypes() (ts []*GlslType) {
	for _, v := range tt.types {
		ts = append(ts, v)
	}
	sort.Slice(ts, less(ts))
	return
}

func (tt *Types) UserStructs() (ts []*GlslType) {
	for _, v := range tt.types {
		if v.UserDefined {
			ts = append(ts, v)
		}
	}
	sort.Slice(ts, less(ts))
	return
}

func (tt *Types) put(ty string, t *GlslType) {
	_, o := tt.types[ty]
	if o {
		panic("trying to create type with same name: " + ty)
	}
	tt.types[ty] = t
}

func (ct CType) String() string {
	s := ct.Name + " "
	if ct.IsSlice() {
		s += "*"
	}
	s += ct.Name
	if ct.ArrayLen() != 0 {
		for el := ct; el.ArrayLen() != 0; el = *el.Array.Ty {
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
		for el := ct; el.ArrayLen() != 0; el = *el.Array.Ty {
			s += fmt.Sprintf("[%v]", el.Array.Len)
		}
	}
	s += "C." + ct.Name
	return s
}

func (t GlslType) CType() *CType {
	return t.C
}

func (t GlslType) GoName() string {
	return strings.Title(t.Name)
}

func (cf CField) String() string {
	if cf.Ty.IsSlice() {
		// this is a slice, so will be passed as pointer, but since
		// we cannot ensure that the caller has the same struct alignment
		// that we require we pass it as a void* and require the caller
		// to ensure that the data has the same layout as what we request.
		return "void* " + cf.Name
	}
	s := cf.Ty.Name + " " + cf.Name
	if cf.Ty.ArrayLen() != 0 {
		for el := cf.Ty; el.ArrayLen() != 0; el = el.Array.Ty {
			s += fmt.Sprintf("[%v]", el.Array.Len)
		}
	}
	return s
}

func (cf CField) CxxFieldString() string {
	if cf.Ty.IsSlice() {
		// there are no built-in types with arbitrary length so this is fine
		return fmt.Sprintf("\t%v* %v;\n", cf.Ty.GlslType.Name, cf.Name)
	}
	s := "\t" + cf.Ty.GlslType.Name + " " + cf.Name

	// NoElem here is in C world, e.g. a vec has a length, so what we want
	// to chec is if the length is equal to the one for the defined type or
	// not
	if cf.Ty.ArrayLen() != 0 {
		for el := cf.Ty; el.ArrayLen() != 0; el = el.Array.Ty {
			s += fmt.Sprintf("[%v]", el.Array.Len)
		}
	}
	return s + ";"
}

func (cf CField) CxxFieldStringRef() string {
	if cf.Ty.IsSlice() {
		// there are no built-in types with arbitrary length so this is fine
		return fmt.Sprintf("\t%v* (&%v);\n", cf.Ty.GlslType.Name, cf.Name)
	}
	s := "\t" + cf.Ty.GlslType.Name + "(& " + cf.Name + ")"

	// NoElem here is in C world, e.g. a vec has a length, so what we want
	// to chec is if the length is equal to the one for the defined type or
	// not
	if cf.Ty.ArrayLen() != 0 {
		for el := cf.Ty; el.ArrayLen() != 0; el = el.Array.Ty {
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
		for el := &ct; el.ArrayLen() != 0; el = el.Array.Ty {
			goTypeName += fmt.Sprintf("[%v]", el.Array.Len)
		}
		goTypeName += fmt.Sprintf("%v", ct.BasicGoType())
	}
	return goTypeName
}

func (types *Types) createBasicBuiltinTypes() {
	cbt := func(name string, cname string, noelc int, size int, align int) {
		tt := GlslType{
			Name: name,
		}
		ut := CType{
			GlslType: &tt,
			Name:     cname,
			Size:     Alignment{ByteSize: 4, ByteAlignment: 4},
		}
		ttt := CType{
			GlslType: &tt,
			Name:     cname,
			Array:    ArrayType{Len: noelc, Ty: &ut},
			Size:     Alignment{ByteSize: size, ByteAlignment: align},
		}

		tt.C = &ttt
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

func (types *Types) createComplexBuiltinTypes() {
	co := func(name string, args ...interface{}) {
		tt := GlslType{
			Name: name,
		}
		ct := CType{
			GlslType: &tt,
			Name:     "cpt_" + name,
			Fields:   []CField{},
		}
		tt.C = &ct
		for i := 0; i+2 < len(args); i += 3 {
			noel := args[i+2].(int)
			name := args[i].(string)
			ty := args[i+1].(string)
			f := CField{
				Name: name,
				Ty:   types.Get(ty).CType(),
			}
			if noel != 0 {
				f.Ty = types.MaybeCreateArrayType(ty, []int{noel})
			}
			ct.Fields = append(ct.Fields, f)
		}
		types.put(name, &tt)
	}
	co("image2Drgba32f", "data", "float", -1, "width", "int32_t", 0)
}

func (types *Types) recCreateArrayType(ty *CType, noels []int) {
	if len(noels) == 0 {
		return
	}
	nt := *(ty) // a copy of it

	noel := noels[0]
	noels = noels[1:]

	ty.Array.Len = noel
	ty.Array.Ty = &nt
	types.recCreateArrayType(&nt, noels)

	// propagate the size all the way bac up
	ty.Size.ByteSize = nt.Size.ByteSize * noel

}

func (types *Types) MaybeCreateArrayType(ty string, noels []int) *CType {
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
		types.recCreateArrayType(&bt, noels)
		bt.Size.ByteAlignment = 8
		bt.Size.ByteSize = 8
	} else {
		types.recCreateArrayType(&bt, noels)
	}
	return &bt
}

func (types *Types) recurseCreateStructTypes(i int, inp input.Input, str input.InputStruct) {
	if _, o := types.types[str.Name]; o {
		return
	}
	st := GlslType{
		Name:        str.Name,
		UserDefined: true,
		GlslOrder:   i,
	}
	ct := CType{
		GlslType: &st,
		Name:     "cpt_" + str.Name,
		Fields:   []CField{},
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
			cf.Ty = types.MaybeCreateArrayType(f.Ty, f.Arrno)
		}
		ct.Fields = append(ct.Fields, cf)
	}
	st.C = &ct
	types.put(str.Name, &st)
}

func (types *Types) calculateAlignments() {
	for _, t := range types.types {
		// all basic types have alignment, so this must be a struct type:
		types.alignVisit(t.CType())
	}
}

func (types *Types) alignVisit(ct *CType) Alignment {
	// from https://www.oreilly.com/library/view/opengl-programming-guide/9780132748445/app09lev1sec3.html
	// but we treat vec3 as vec4 to allow for vector operations
	if ct.Size != (Alignment{}) {
		return ct.Size
	}

	offset := 0
	maxFieldAlignment := 0
	for i, f := range ct.Fields {
		fa := types.alignVisit(f.Ty) // => we cannot handle recursive types, thats fine...

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
	if ct.IsSlice() {
		// TODO: this is 64 bit only...
		rsize.ByteSize = 8
		rsize.ByteAlignment = 8
	}

	ct.Size = rsize
	return rsize
}

func (types *Types) findApiExportedTypes(inp input.Input) {
	for _, arg := range inp.Arguments {
		types.exportType(arg.Ty)
	}
}

func (types *Types) exportType(s string) {
	t := types.Get(s)
	if t.Export {
		return
	}
	t.Export = true
	for _, f := range t.C.Fields {
		types.exportType(f.Ty.GlslType.Name)
	}
}
