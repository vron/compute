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

func generateGo(inp input.Input, ts *types.Types) {
	f, err := os.Create(filepath.Join(fOut, "kernel.go"))

	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	buf := bufio.NewWriter(f)
	defer buf.Flush()

	fmt.Fprintf(buf, "// Package kernel is a wrapper to execute a particular GLSL compute shader\n")
	fmt.Fprintf(buf, "package kernel"+"\n\n")

	writePreamble(buf, inp, ts)

	fmt.Fprintf(buf, `
// Code generated DO NOT EDIT

type Kernel struct {
	k unsafe.Pointer
	dead bool
}

// New creates a Kernel using at most numCPU+1 threads. If numCPU <= 0 the
// number of threads to use will be calculated automatically. All kernels
// must be explicitly freed using Kernel.Free to avoid memory leaks.
  func New(numCPU int, stackSize int) (k *Kernel, err error) {
	k = &Kernel{}
	if numCPU <= 0 {
		numCPU = runtime.NumCPU()+2
	}
	k.k = C.cpt_new_kernel(C.int32_t(numCPU), C.int32_t(stackSize));
	if k.k == nil {
		return nil, errors.New("failed to create kernel structure")
	}
	runtime.SetFinalizer(k, freeKernel)
	return k, nil
}



// Free dealocates any data allocated by the underlying Kernel. Note that
// a kernel on which Free has been called can no longer be used.
func (k *Kernel) Free() {
	freeKernel(k)
}


func freeKernel(k *Kernel) {
	if k.dead {
		return
	}
	k.dead = true
	C.cpt_free_kernel(k.k);
}
`)

	writeTypeDefinitions(buf, inp, ts)
	writeDataStruct(buf, inp, ts)
	writeDataStructRaw(buf, inp, ts)
	writeDispatch(buf, inp, ts)
	writeDispatchRaw(buf, inp, ts)
	writeSizeAlign(buf, inp, ts)
	writeEncode(buf, inp, ts)
	writeDecode(buf, inp, ts)
	writeEnsureAlign(buf, inp, ts)
	writeSupportFuncs(buf, inp, ts)
	/*
		// Also create the Encode Decode Methods for types that are referred in arrays
		fmt.Fprintf(buf,"var bo = binary.LittleEndian\n\n")

		for _, st := range ts.ExportedStructTypes() {
			if st.UserDefined {
				fmt.Fprintf(buf, "func (d %v) Stride() int { return %v }\n", st.GoName(), st.CType().Size.ByteSize)
				fmt.Fprintf(buf, "func (d %v) Alignment() int { return %v }\n\n", st.GoName(), st.CType().Size.ByteAlignment)

				// Create a Encode function for the element
				fmt.Fprintf(buf, "func (e *%v) Encode(d []byte) {\n", st.GoName())
				printStructEncodes(buf, 0, st, "")
				fmt.Fprintf(buf,"}\n\n")
				// Create a Decode function for the element
				fmt.Fprintf(buf, "func (e *%v) Decode(d []byte) {\n", st.GoName())
				printStructDecodes(buf, 0, st, "")
				fmt.Fprintf(buf,"}\n\n")
			}
		}
	*/
}

func writePreamble(buf io.Writer, inp input.Input, ts *types.Types) {
	fmt.Fprintf(buf, `/*
#cgo darwin LDFLAGS: -L${SRCDIR} -L. build/shader.so
#cgo linux LDFLAGS: -L${SRCDIR}/build -L. build/shader.so
#cgo windows LDFLAGS: -L. -lshader

#include "shared.h"

`)
	fmt.Fprintf(buf, `struct cpt_error_t wrap_dispatch(void *k, `+"\n")
	tab := "                                 "
	for _, arg := range inp.Arguments {
		fmt.Fprintf(buf, tab+"void* "+arg.Name+", ")
		fmt.Fprintf(buf, "int64_t "+arg.Name+"_len,\n")
	}
	fmt.Fprintf(buf, tab+`int32_t x, int32_t y, int32_t z) {
	cpt_data d;
`)
	for _, arg := range inp.Arguments {
		fmt.Fprintf(buf, "\td."+arg.Name+" = "+arg.Name+";\n")
		fmt.Fprintf(buf, "\td."+arg.Name+"_len = "+arg.Name+"_len;\n")
	}
	fmt.Fprintf(buf, `	return cpt_dispatch_kernel(k, d, x, y, z);
}
*/
import "C"

`)
}

func writeTypeDefinitions(buf io.Writer, inp input.Input, ts *types.Types) {
	for _, st := range ts.ListExportedTypes() {
		if st.C.IsBasic() {
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
				// we need to special case bool since that is the only basic type we cannot
				// map  1-1 since it is one byte in go and 4 in glsl / c imple
				if f.CType.IsBasic() && f.CType.GlslType.Name == "Bool" {
					offset += 1
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

func writeDataStruct(buf io.Writer, inp input.Input, ts *types.Types) {
	fmt.Fprintf(buf, "type Data struct {\n")
	for _, arg := range inp.Arguments {
		ty := types.CreateArray(ts.Get(arg.Ty).C, arg.Arrno)
		if len(arg.Arrno) > 0 && arg.Arrno[0] == -1 {
			fmt.Fprintf(buf, "  "+ty.GoString(strings.Title(arg.Name))+"\n")
		} else {
			fmt.Fprintf(buf, "  "+pointify(ty.GoString(strings.Title(arg.Name)))+"\n")
		}
	}
	fmt.Fprintf(buf, "}\n\n")
}

func writeDataStructRaw(buf io.Writer, inp input.Input, ts *types.Types) {
	fmt.Fprintf(buf, "type DataRaw struct {\n")
	for _, arg := range inp.Arguments {
		fmt.Fprintf(buf, strings.Title(arg.Name)+" []byte\n")
	}
	fmt.Fprintf(buf, "}\n\n")
}

func writeDispatch(buf io.Writer, inp input.Input, ts *types.Types) {
	fmt.Fprintf(buf, `
// Dispatch a Kernel calculation of the specified size. The caller must ensure
// that the data provided in bind matches the kernel's assumptions and that any
// []byte field represents properly aligned data. Not data in bind must
// be accessed (read or write) until Dispatch returns.
func (k *Kernel) Dispatch(bind Data, numGroupsX, numGroupsY, numGroupsZ int) error {
	if k.dead {
		panic("cannot use a Kernel where Free() has been called")
	}
`)
	fmt.Fprintf(buf, ` errno := C.wrap_dispatch(k.k,
`)
	for _, arg := range inp.Arguments {
		if len(arg.Arrno) > 0 && arg.Arrno[0] == -1 {
			// slice data, we need the size of the entire thing...
			fmt.Fprintf(buf, "\tunsafe.Pointer(&bind."+strings.Title(arg.Name)+"[0]), ")
			fmt.Fprintf(buf, "C.int64_t(int64(len(bind.%v))*int64(unsafe.Sizeof(bind.%v[0]))),\n", strings.Title(arg.Name), strings.Title(arg.Name))
		} else {
			fmt.Fprintf(buf, "\tunsafe.Pointer(bind."+strings.Title(arg.Name)+"), ")
			fmt.Fprintf(buf, "C.int64_t(unsafe.Sizeof(*bind."+strings.Title(arg.Name)+")),\n")
		}
	}
	fmt.Fprintf(buf, `C.int(numGroupsX), C.int(numGroupsY), C.int(numGroupsZ))`)

	// decode the error message
	fmt.Fprintf(buf, `
	if errno.code == 0 {
		return nil
	}
	errstr := C.GoString(errno.msg)
	return errors.New(strconv.Itoa(int(errno.code)) + ": " + errstr)
}
`)
}

func writeDispatchRaw(buf io.Writer, inp input.Input, ts *types.Types) {
	fmt.Fprintf(buf, `
func (k *Kernel) DispatchRaw(bind DataRaw, numGroupsX, numGroupsY, numGroupsZ int) error {
	if k.dead {
		panic("cannot use a Kernel where Free() has been called")
	}
`)
	fmt.Fprintf(buf, ` errno := C.wrap_dispatch(k.k,
`)
	for _, arg := range inp.Arguments {
		fmt.Fprintf(buf, "\tunsafe.Pointer(&bind."+strings.Title(arg.Name)+"[0]), ")
		fmt.Fprintf(buf, "C.int64_t(int64(len(bind.%v))),\n", strings.Title(arg.Name))
	}
	fmt.Fprintf(buf, `C.int(numGroupsX), C.int(numGroupsY), C.int(numGroupsZ))`)

	// decode the error message
	fmt.Fprintf(buf, `
	if errno.code == 0 {
		return nil
	}
	errstr := C.GoString(errno.msg)
	return errors.New(strconv.Itoa(int(errno.code)) + ": " + errstr)
}

`)
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

func iBool(v uint32) bool {
	return !(v==0)
}

// AlignedSlice returns a byte slice where the first element has a minimum
// alignment of align and a length if size.
func AlignedSlice(size, align int) (b []byte) {
	if align < 1 {
		panic("align must be > 0")
	}
	b = make([]byte, size+align-1)
	adr := uintptr(unsafe.Pointer(&b[0]))
	diff := 0
	if int(adr) % align != 0 {
		diff = align - int(adr) % align
	}
	return b[diff:diff+size]
}
	
`)
}

func writeSizeAlign(buf io.Writer, inp input.Input, ts *types.Types) {
	for _, st := range ts.ListExportedTypes() {
		if st.C.IsBasic() {
		} else if st.C.IsVector() || st.C.IsStruct() {
			fmt.Fprintf(buf, "func (v *%v) Alignof() int { return %v }\n\n", st.GoName(), st.C.Size.ByteAlignment)
			fmt.Fprintf(buf, "func (v *%v) Sizeof() int { return %v }\n\n", st.GoName(), st.C.Size.ByteSize)
		} else {
			panic("cannot have an exported array type? what happened?")
		}
	}
}

func writeSingle(buf io.Writer, parentPos string, ty *types.GlslType, head string) {
	switch ty.Name {
	case "Bool":
		fmt.Fprintf(buf, "\tbo.PutUint32(%v, uint32(cBool(%v)))\n", parentPos, head)
	case "float":
		fmt.Fprintf(buf, "\tbo.PutUint32(%v, math.Float32bits(%v))\n", parentPos, head)
	case "int32_t":
		fmt.Fprintf(buf, "\tbo.PutUint32(%v, uint32(%v))\n", parentPos, head)
	case "uint32_t":
		fmt.Fprintf(buf, "\tbo.PutUint32(%v, %v)\n", parentPos, head)
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
					// this one can be multiple levels, but cannot be -1 here
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
		if ty.IsBasic() {
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
		fmt.Fprintf(buf, "\t%v = iBool(bo.Uint32(%v))\n", head, parentPos)
	case "float":
		fmt.Fprintf(buf, "\t%v = math.Float32frombits(bo.Uint32(%v))\n", head, parentPos)
	case "int32_t":
		fmt.Fprintf(buf, "\t%v = int32(bo.Uint32(%v))\n", head, parentPos)
	case "uint32_t":
		fmt.Fprintf(buf, "\t%v = bo.Uint32(%v)\n", head, parentPos)
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
		if ty.IsBasic() {
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

func pointify(s string) string {
	// hac*y but wor*s
	return strings.Replace(s, "\t", "\t*", 1)
}
