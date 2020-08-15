package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strings"

	"github.com/vron/compute/glbind/input"
	"github.com/vron/compute/glbind/types"
)

func generateGoTypes(inp input.Input, ts *types.Types) {
	f, err := os.Create(filepath.Join(fOut, "types.go"))

	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	buf := bufio.NewWriter(f)
	defer buf.Flush()

	fmt.Fprintf(buf, "package kernel"+"\n\n")

	fmt.Fprintf(buf, `
// Code generated DO NOT EDIT

`)

	writeTypeDefinitions(buf, inp, ts)
	writeSizeAlign(buf, inp, ts)
	writeEncode(buf, inp, ts)
	writeDecode(buf, inp, ts)
	writeEnsureAlign(buf, inp, ts)
	writeSupportFuncs(buf, inp, ts)
}

func writeTypeDefinitions(buf io.Writer, inp input.Input, ts *types.Types) {
	for _, st := range ts.ListExportedTypes() {
		if st.C.IsBasic() {
			if st.Name == "Bool" {
				fmt.Fprintf(buf, "type Bool struct{ B bool\n_ [3]bool}\n\n")
				fmt.Fprintf(buf, "var True = Bool{B: true}\n\n")
				fmt.Fprintf(buf, "var False = Bool{}\n\n")
			}
			// do nothing, we use built in go types here
		} else if st.C.IsVector() {
			fmt.Fprintf(buf, "type %v [%v]%v\n\n", st.GoName(), st.C.Size.ByteSize/st.C.Vector.Basic.Size.ByteSize, st.C.Vector.Basic.GlslType.GoName())
		} else if st.C.IsStruct() {
			offset := 0
			fmt.Fprintf(buf, "type %v struct {\n", st.GoName())
			for _, f := range st.C.Struct.Fields {
				if offset != f.ByteOffset {
					fmt.Fprintf(buf, "\t_\t[%v]byte\n", f.ByteOffset-offset)
				}
				offset = f.ByteOffset
				fmt.Fprintf(buf, "  "+f.CType.GoString(strings.Title(f.Name))+"\n")
				if f.CType.IsArray() && f.CType.Array.Len == -1 {
					offset += 8 * 3 // a slice has pointer, cap and lenb: TODO: 32bit
				} else {
					offset += f.CType.Size.ByteSize
				}
			}
			if offset != st.C.Size.ByteSize {
				fmt.Fprintf(buf, "\t_\t[%v]byte\n", st.C.Size.ByteSize-offset)
			}
			fmt.Fprintf(buf, "}\n\n")
		} else {
			panic("cannot have an exported array type? what happened?")
		}
	}
}

func writeEnsureAlign(buf io.Writer, inp input.Input, ts *types.Types) {
	fmt.Fprintf(buf, `

// ensure that the go structs mem layout match those in the C shader code
func init() {
`)
	for _, st := range ts.ListExportedTypes() {
		if st.C.IsBasic() {
			// do nothing, we use built in go types here
		} else if st.C.IsVector() {
			fmt.Fprintf(buf, `if unsafe.Sizeof(%v{}) != %v { panic("sizeof(%v) != %v") }`+"\n", st.GoName(), st.C.Size.ByteSize, st.GoName(), st.C.Size.ByteSize)
		} else if st.C.IsStruct() {
			fmt.Fprintf(buf, `if unsafe.Sizeof(%v{}) != %v { panic("sizeof(%v) != %v") }`+"\n", st.GoName(), st.C.Size.ByteSize, st.GoName(), st.C.Size.ByteSize)
			for _, f := range st.C.Struct.Fields {
				fmt.Fprintf(buf, `if unsafe.Offsetof(%v{}.%v) != %v { panic("offsetof(%v.%v) != %v") }`+"\n", st.GoName(), strings.Title(f.Name), f.ByteOffset, st.GoName(), strings.Title(f.Name), f.ByteOffset)
			}
		} else {
			panic("cannot have an exported array type? what happened?")
		}
	}
	fmt.Fprintf(buf, `}
`)
}

func writeSupportFuncs(buf io.Writer, inp input.Input, ts *types.Types) {
	fmt.Fprint(buf, `
var bo = binary.LittleEndian

func cBool(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

func writeByte(b []byte, by byte) {
	b[0] = by
}

func readByte(b []byte) byte {
	return b[0]
}


func iBool(v uint32) bool {
	return !(v==0)
}

	
`)
}

func writeSizeAlign(buf io.Writer, inp input.Input, ts *types.Types) {
	for _, st := range ts.ListExportedTypes() {
		if st.C.IsBasic() {
		} else if st.C.IsVector() || st.C.IsStruct() {
			fmt.Fprintf(buf, "func (v %v) Alignof() int { return %v }\n\n", st.GoName(), st.C.Size.ByteAlignment)
			fmt.Fprintf(buf, "func (v %v) Sizeof() int { return %v }\n\n", st.GoName(), st.C.Size.ByteSize)
		} else {
			panic("cannot have an exported array type? what happened?")
		}
	}
}

func writeSingle(buf io.Writer, parentPos string, ty *types.GlslType, head string) {
	switch ty.Name {
	case "Bool":
		fmt.Fprintf(buf, "\tbo.PutUint32(%v, uint32(cBool((%v).B)))\n", parentPos, head)
	case "float":
		fmt.Fprintf(buf, "\tbo.PutUint32(%v, math.Float32bits(%v))\n", parentPos, head)
	case "int32_t":
		fmt.Fprintf(buf, "\tbo.PutUint32(%v, uint32(%v))\n", parentPos, head)
	case "uint32_t":
		fmt.Fprintf(buf, "\tbo.PutUint32(%v, %v)\n", parentPos, head)
	case "uint8_t":
		fmt.Fprintf(buf, "\twriteByte(%v, %v)\n", parentPos, head)
	default:
		fmt.Fprintf(buf, "\t(%v).Encode(%v) \n", head, parentPos)
	}
}

func recEncodeArray(buf io.Writer, lvl int, ty *types.CType, offset string, heado string, name string) {
	head := heado + fmt.Sprintf("[i%v]", lvl)
	st := ty.Array.CType
	if st.IsBasic() || st.IsStruct() || st.IsVector() {
		if ty.Array.Len == -1 {
			fmt.Fprintf(buf, "for i%v := 0; i%v < len(%v); i%v++ {\n", lvl, lvl, heado, lvl)
			writeSingle(buf, fmt.Sprintf("%v[%v+i%v*%v:]", name, offset, lvl, st.Size.ByteSize), st.GlslType, head)

		} else {
			fmt.Fprintf(buf, "for i%v := 0; i%v < %v; i%v++ {\n", lvl, lvl, ty.Array.Len, lvl)
			writeSingle(buf, fmt.Sprintf("%v[%v+i%v*%v:]", name, offset, lvl, st.Size.ByteSize), st.GlslType, head)
		}
	} else if st.IsArray() {
		if ty.Array.Len == -1 {
			fmt.Fprintf(buf, "for i%v := 0; i%v < len(%v); i%v++ {\n", lvl, lvl, heado, lvl)
			recEncodeArray(buf, lvl+1, st, fmt.Sprintf("(%v)+i%v*%v", offset, lvl, st.Size.ByteSize), head, name)
		} else {
			fmt.Fprintf(buf, "for i%v := 0; i%v < %v; i%v++ {\n", lvl, lvl, ty.Array.Len, lvl)
			recEncodeArray(buf, lvl+1, st, fmt.Sprintf("(%v)+i%v*%v", offset, lvl, st.Size.ByteSize), head, name)
		}
	}
	fmt.Fprintf(buf, "}\n")
}

func writeEncode(buf io.Writer, inp input.Input, ts *types.Types) {
	for _, st := range ts.ListExportedTypes() {
		if st.C.IsBasic() {
		} else if st.C.IsVector() {
			fmt.Fprintf(buf, "func (v *%v) Encode(d []byte) {\n", st.GoName())
			for i := 0; i < st.C.Vector.Len; i++ {
				writeSingle(buf, fmt.Sprintf("d[%v:]", i*st.C.Vector.Basic.Size.ByteSize), st.C.Vector.Basic.GlslType, fmt.Sprintf("v[%v]", i))
			}
			fmt.Fprintf(buf, "}\n\n")
		} else if st.C.IsStruct() {
			fmt.Fprintf(buf, "func (v *%v) Encode(d []byte) {\n", st.GoName())
			for _, f := range st.C.Struct.Fields {
				if f.CType.IsBasic() || f.CType.IsStruct() || f.CType.IsVector() {
					writeSingle(buf, fmt.Sprintf("d[%v:]", f.ByteOffset), f.CType.GlslType, "v."+strings.Title(f.Name))
				} else if f.CType.IsArray() {
					recEncodeArray(buf, 0, f.CType, fmt.Sprint(f.ByteOffset), "v."+strings.Title(f.Name), "d")
				}
			}
			fmt.Fprintf(buf, "}\n\n")
		} else {
			panic("cannot have an exported array type? what happened?")
		}
	}

	// also an unexported encode function for the full data structure that is used in tests
	fmt.Fprintf(buf, "func encodeData(v Data) (r DataRaw) {\n")
	for _, arg := range inp.Arguments {
		ty := types.CreateArray(ts.Get(arg.Ty).C, arg.Arrno)
		if ty.IsComplexStruct() {
			fmt.Fprintf(buf, "  r.%v =v.%v\n", strings.Title(arg.Name), strings.Title(arg.Name))
		} else if ty.IsBasic() {
			fmt.Fprintf(buf, "  r.%v = AlignedSlice(%v, %v)\n", strings.Title(arg.Name), ty.Size.ByteSize, ty.Size.ByteAlignment)
			writeSingle(buf, fmt.Sprintf("r.%v", strings.Title(arg.Name)), ty.GlslType, "*v."+strings.Title(arg.Name))
		} else if ty.IsStruct() || ty.IsVector() {
			fmt.Fprintf(buf, "  r.%v = AlignedSlice(%v, %v)\n", strings.Title(arg.Name), ty.Size.ByteSize, ty.Size.ByteAlignment)
			writeSingle(buf, fmt.Sprintf("r.%v", strings.Title(arg.Name)), ty.GlslType, "v."+strings.Title(arg.Name))
		} else if ty.IsArray() {
			if ty.Array.Len == -1 {
				fmt.Fprintf(buf, "  r.%v = AlignedSlice(len(v.%v)*%v, %v)\n", strings.Title(arg.Name), strings.Title(arg.Name), ty.Array.CType.Size.ByteSize, ty.Array.CType.Size.ByteAlignment)
			} else {
				fmt.Fprintf(buf, "  r.%v = AlignedSlice(%v, %v)\n", strings.Title(arg.Name), ty.Size.ByteSize, ty.Size.ByteAlignment)
			}
			recEncodeArray(buf, 0, ty, fmt.Sprint(0), "v."+strings.Title(arg.Name), "r."+strings.Title(arg.Name))
		}
	}
	fmt.Fprintf(buf, "return\n}\n\n")
}

func readSingle(buf io.Writer, parentPos string, ty *types.GlslType, head string) {
	switch ty.Name {
	case "Bool":
		fmt.Fprintf(buf, "\t%v = Bool{B: iBool(bo.Uint32(%v))}\n", head, parentPos)
	case "float":
		fmt.Fprintf(buf, "\t%v = math.Float32frombits(bo.Uint32(%v))\n", head, parentPos)
	case "int32_t":
		fmt.Fprintf(buf, "\t%v = int32(bo.Uint32(%v))\n", head, parentPos)
	case "uint32_t":
		fmt.Fprintf(buf, "\t%v = bo.Uint32(%v)\n", head, parentPos)
	case "uint8_t":
		fmt.Fprintf(buf, "\t%v = readByte(%v)\n", head, parentPos)
	default:
		fmt.Fprintf(buf, "\t(%v).Decode(%v) \n", head, parentPos)
	}
}

func recDecodeArray(buf io.Writer, lvl int, ty *types.CType, offset string, heado string, name string) {
	head := heado + fmt.Sprintf("[i%v]", lvl)
	st := ty.Array.CType
	if st.IsBasic() || st.IsStruct() || st.IsVector() {
		if ty.Array.Len == -1 {
			fmt.Fprintf(buf, "for i%v := 0; i%v < len(%v); i%v++ {\n", lvl, lvl, heado, lvl)
			readSingle(buf, fmt.Sprintf("%v[%v+i%v*%v:]", name, offset, lvl, st.Size.ByteSize), st.GlslType, head)
		} else {
			fmt.Fprintf(buf, "for i%v := 0; i%v < %v; i%v++ {\n", lvl, lvl, ty.Array.Len, lvl)
			readSingle(buf, fmt.Sprintf("%v[%v+i%v*%v:]", name, offset, lvl, st.Size.ByteSize), st.GlslType, head)
		}
	} else if st.IsArray() {
		if ty.Array.Len == -1 {
			fmt.Fprintf(buf, "for i%v := 0; i%v < len(%v); i%v++ {\n", lvl, lvl, heado, lvl)
			recDecodeArray(buf, lvl+1, st, fmt.Sprintf("(%v)+i%v*%v", offset, lvl, st.Size.ByteSize), head, name)
		} else {
			// this one can be multiple levels, but cannot be -1 here
			fmt.Fprintf(buf, "for i%v := 0; i%v < %v; i%v++ {\n", lvl, lvl, ty.Array.Len, lvl)
			recDecodeArray(buf, lvl+1, st, fmt.Sprintf("(%v)+i%v*%v", offset, lvl, st.Size.ByteSize), head, name)
		}
	}
	fmt.Fprintf(buf, "}\n")
}

func writeDecode(buf io.Writer, inp input.Input, ts *types.Types) {
	for _, st := range ts.ListExportedTypes() {
		if st.C.IsBasic() {
		} else if st.C.IsVector() {
			fmt.Fprintf(buf, "func (v *%v) Decode(d []byte) {\n", st.GoName())
			for i := 0; i < st.C.Vector.Len; i++ {
				readSingle(buf, fmt.Sprintf("d[%v:]", i*st.C.Vector.Basic.Size.ByteSize), st.C.Vector.Basic.GlslType, fmt.Sprintf("v[%v]", i))
			}
			fmt.Fprintf(buf, "}\n\n")
		} else if st.C.IsStruct() {
			fmt.Fprintf(buf, "func (v *%v) Decode(d []byte) {\n", st.GoName())
			for _, f := range st.C.Struct.Fields {
				if f.CType.IsBasic() || f.CType.IsStruct() || f.CType.IsVector() {
					readSingle(buf, fmt.Sprintf("d[%v:]", f.ByteOffset), f.CType.GlslType, "v."+strings.Title(f.Name))
				} else if f.CType.IsArray() {
					// this one can be multiple levels, but cannot be -1 here
					recDecodeArray(buf, 0, f.CType, fmt.Sprint(f.ByteOffset), "v."+strings.Title(f.Name), "d")
				}
			}
			fmt.Fprintf(buf, "}\n\n")
		} else {
			panic("cannot have an exported array type? what happened?")
		}
	}

	// also an unexported encode function for the full data structure that is used in tests
	fmt.Fprintf(buf, "func decodeData(r DataRaw) (d Data) {\n")
	for i, arg := range inp.Arguments {
		ty := types.CreateArray(ts.Get(arg.Ty).C, arg.Arrno)
		if ty.IsComplexStruct() {
			fmt.Fprintf(buf, "  d.%v =r.%v\n", strings.Title(arg.Name), strings.Title(arg.Name))
		} else if ty.IsBasic() {
			readSingle(buf, fmt.Sprintf("r.%v", strings.Title(arg.Name)), ty.GlslType, fmt.Sprintf("var t%v %v", i, ty.GlslType.GoName()))
			fmt.Fprintf(buf, "d.%v = &t%v\n", strings.Title(arg.Name), i)
		} else if ty.IsStruct() || ty.IsVector() {
			fmt.Fprintf(buf, "var t%v %v\n", i, ty.GlslType.GoName())
			readSingle(buf, fmt.Sprintf("r.%v", strings.Title(arg.Name)), ty.GlslType, fmt.Sprintf("&t%v", i))
			fmt.Fprintf(buf, "d.%v = &t%v\n", strings.Title(arg.Name), i)
		} else if ty.IsArray() {
			if ty.Array.Len == -1 {
				fmt.Fprintf(buf, "t%v := make(%v, len(r.%v)/%v)\n", i, ty.GoString(""), strings.Title(arg.Name), ty.Array.CType.Size.ByteSize)
				recDecodeArray(buf, 0, ty, fmt.Sprint(0), fmt.Sprintf("t%v", i), "r."+strings.Title(arg.Name))
				fmt.Fprintf(buf, "d.%v = t%v\n", strings.Title(arg.Name), i)
			} else {
				fmt.Fprintf(buf, "var t%v %v\n", i, ty.GoString(""))
				recDecodeArray(buf, 0, ty, fmt.Sprint(0), fmt.Sprintf("t%v", i), "r."+strings.Title(arg.Name))
				fmt.Fprintf(buf, "d.%v = &t%v\n", strings.Title(arg.Name), i)
			}
		}
	}
	fmt.Fprintf(buf, "return\n}\n\n")
}
