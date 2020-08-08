package types

import (
	"github.com/vron/compute/glbind/input"
)

type Types struct {
	m map[string]*GlslType
	l []*GlslType
}

func New(inp input.Input) *Types {
	ts := &Types{
		m: map[string]*GlslType{},
		l: []*GlslType{},
	}
	ts.createBasicBuiltinTypes()
	ts.createComplexBuiltinTypes()
	for _, str := range inp.Structs {
		ts.createUserTypes(str)
	}
	ts.calculateAlignments()
	for _, arg := range inp.Arguments {
		ts.exportType(ts.Get(arg.Ty).C)
	}

	return ts
}

func (ts *Types) put(t GlslType) {
	_, o := ts.m[t.Name]
	if o {
		panic("trying to create type with same name: " + t.Name)
	}
	t.C.GlslType = &t
	ts.m[t.Name] = &t
	ts.l = append(ts.l, &t)
}

func (ts *Types) Get(name string) *GlslType {
	v, o := ts.m[name]
	if !o {
		panic("trying to get type: " + name)
	}
	return v
}

func (ts *Types) ListExportedTypes() (l []*GlslType) {
	for _, v := range ts.l {
		if v.Export {
			l = append(l, v)
		}
	}
	return
}

func (tt *Types) ListAllTypes() (ts []*GlslType) {
	return tt.l
}

func (ts *Types) createBasicBuiltinTypes() {
	ts.put(GlslType{Builtin: true, Name: "Bool", C: &CType{Basic: CBasicType{Name: "int32_t"}, Size: align(4, 4)}})
	ts.put(GlslType{Builtin: true, Name: "int32_t", C: &CType{Basic: CBasicType{Name: "int32_t"}, Size: align(4, 4)}})
	ts.put(GlslType{Builtin: true, Name: "uint32_t", C: &CType{Basic: CBasicType{Name: "uint32_t"}, Size: align(4, 4)}})
	ts.put(GlslType{Builtin: true, Name: "float", C: &CType{Basic: CBasicType{Name: "float"}, Size: align(4, 4)}})

	ts.put(GlslType{Builtin: true, Name: "vec2", C: &CType{Vector: CVector{Len: 2, Basic: ts.Get("float").C}, Size: align(8, 8)}})
	ts.put(GlslType{Builtin: true, Name: "vec3", C: &CType{Vector: CVector{Len: 3, Basic: ts.Get("float").C}, Size: align(16, 16)}})
	ts.put(GlslType{Builtin: true, Name: "vec4", C: &CType{Vector: CVector{Len: 4, Basic: ts.Get("float").C}, Size: align(16, 16)}})

	ts.put(GlslType{Builtin: true, Name: "ivec2", C: &CType{Vector: CVector{Len: 2, Basic: ts.Get("int32_t").C}, Size: align(8, 8)}})
	ts.put(GlslType{Builtin: true, Name: "ivec3", C: &CType{Vector: CVector{Len: 3, Basic: ts.Get("int32_t").C}, Size: align(16, 16)}})
	ts.put(GlslType{Builtin: true, Name: "ivec4", C: &CType{Vector: CVector{Len: 4, Basic: ts.Get("int32_t").C}, Size: align(16, 16)}})

	ts.put(GlslType{Builtin: true, Name: "uvec2", C: &CType{Vector: CVector{Len: 2, Basic: ts.Get("uint32_t").C}, Size: align(8, 8)}})
	ts.put(GlslType{Builtin: true, Name: "uvec3", C: &CType{Vector: CVector{Len: 3, Basic: ts.Get("uint32_t").C}, Size: align(16, 16)}})
	ts.put(GlslType{Builtin: true, Name: "uvec4", C: &CType{Vector: CVector{Len: 4, Basic: ts.Get("uint32_t").C}, Size: align(16, 16)}})

	ts.put(GlslType{Builtin: true, Name: "mat2", C: &CType{
		Struct: CStruct{Fields: []CField{
			{Name: "column0", CType: ts.Get("vec2").C, ByteOffset: 0},
			{Name: "column1", CType: ts.Get("vec2").C, ByteOffset: 8},
		}},
		Size: align(16, 8)}})
	ts.put(GlslType{Builtin: true, Name: "mat3", C: &CType{
		Struct: CStruct{Fields: []CField{
			{Name: "column0", CType: ts.Get("vec3").C, ByteOffset: 0},
			{Name: "column1", CType: ts.Get("vec3").C, ByteOffset: 16},
			{Name: "column2", CType: ts.Get("vec3").C, ByteOffset: 16 * 2},
		}},
		Size: align(16*3, 16)}})
	ts.put(GlslType{Builtin: true, Name: "mat4", C: &CType{
		Struct: CStruct{Fields: []CField{
			{Name: "column0", CType: ts.Get("vec4").C, ByteOffset: 0},
			{Name: "column1", CType: ts.Get("vec4").C, ByteOffset: 16},
			{Name: "column2", CType: ts.Get("vec4").C, ByteOffset: 16 * 2},
			{Name: "column3", CType: ts.Get("vec4").C, ByteOffset: 16 * 3},
		}},
		Size: align(16*4, 16)}})
}
func (ts *Types) createComplexBuiltinTypes() {
	ts.put(GlslType{Builtin: true, Name: "image2Drgba32f", C: &CType{Struct: CStruct{
		Fields: []CField{
			{Name: "data", CType: CreateArray(ts.Get("float").C, []int{-1}), ByteOffset: 0},
			{Name: "width", CType: ts.Get("int32_t").C, ByteOffset: 28},
		}}, Size: align(32, 8)}})
}

func (ts *Types) createUserTypes(str input.InputStruct) {
	fields := []CField{}
	for _, f := range str.Fields {
		fields = append(fields, CField{
			Name:  f.Name,
			CType: CreateArray(ts.Get(f.Ty).C, f.Arrno),
		})
	}

	gt := GlslType{
		Name: str.Name,
		C: &CType{
			Struct: CStruct{
				Fields: fields,
			}}}
	ts.put(gt)
	ts.exportType(gt.C)
}

func (ts *Types) exportType(t *CType) {
	if t.GlslType != nil {
		if t.GlslType.Export {
			return
		}
		t.GlslType.Export = true
	}

	if t.IsArray() {
		ts.exportType(t.Array.CType)
	} else if t.IsStruct() {
		for _, f := range t.Struct.Fields {
			ts.exportType(f.CType)
		}
	}
}
