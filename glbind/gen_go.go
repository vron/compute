package main

import (
	"bufio"
	"bytes"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
	"strconv"
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

	buf.WriteString("// Package kernel is a wrapper to execute a particular GLSL compute shader\n")
	buf.WriteString("package kernel" + "\n\n")
	buf.WriteString(`// #cgo darwin LDFLAGS: -L${SRCDIR} -L. build/shader.so
// #cgo linux LDFLAGS: -L${SRCDIR}/build -L. build/shader.so
// #cgo windows LDFLAGS: -L. -lshader
// #include "shared.h"
import "C"

// Code generated DO NOT EDIT

import (
	"errors"
	"fmt"
	"runtime"
	"strconv"
	"unsafe"
	"math"
	"encoding/binary"
)

type Kernel struct {
	k unsafe.Pointer
	dead bool
}

`)

	// Write the struct definitions
	for _, st := range ts.ExportedStructTypes() {
		fmt.Fprintf(buf, "type %v struct {\n", st.GoName())
		for _, a := range st.CType().Fields {
			fmt.Fprintf(buf, "\t%v %v\n", a.GoName(), a.Ty.GoName())
		}
		fmt.Fprintf(buf, "}\n\n")
	}

	buf.WriteString(`type Data struct {
`)
	for _, arg := range inp.Arguments {
		cf := types.CField{Name: arg.Name, Ty: ts.MaybeCreateArrayType(arg.Ty, arg.Arrno)}
		buf.WriteString("\t" + cf.GoName() + " " + cf.Ty.GoName() + "\n")
	}
	buf.WriteString(`}

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

// Dispatch a Kernel calculation of the specified size. The caller must ensure
// that the data provided in bind matches the kernel's assumptions and that any
// []byte field represents properly aligned data. Not data in bind must
// be accessed (read or write) until Dispatch returns.
func (k *Kernel) Dispatch(bind Data, numGroupsX, numGroupsY, numGroupsZ int) error {
	if k.dead {
		panic("cannot use a Kernel where Free() has been called")
	}
	cbind := C.cpt_data{
`)

	// Create a c-struct from the provided that that is uploaded to the ernel
	chc := bytes.NewBuffer(nil)
	for _, arg := range inp.Arguments {
		cf := types.CField{Name: arg.Name, Ty: ts.MaybeCreateArrayType(arg.Ty, arg.Arrno)}
		fmt.Fprintf(buf, "\t%v: %v", cf.Name, go2c(chc, cf, cf.Ty, fmt.Sprintf("bind.%v", cf.GoName()))+",\n")
	}

	buf.WriteString(`	}
`)

	// Add the checs for the lengths of the provided datas, to decrease the prob.
	// of bad data provided...
	buf.Write(chc.Bytes())
	buf.WriteString(`
	errno := C.cpt_dispatch_kernel(k.k, cbind, C.int(numGroupsX), C.int(numGroupsY), C.int(numGroupsZ))`)

	// decode the error message
	buf.WriteString(`
	if errno.code == 0 {
		return nil
	}
	errstr := C.GoString(errno.msg)
	return errors.New(strconv.Itoa(int(errno.code)) + ": " + errstr)
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

func cBool(b bool) int32 {
	if b {
		return 1
	}
	return 0
}

func ensureLength(f string, l, s, arr int) error {
	if arr > 0 {
		if l != s*arr {
			return fmt.Errorf("bad data for %v, expected length %v*%v=%v but got %v", f, s, arr, s*arr, l)
		}
	}
	if arr < 0 {
		if l % s != 0 {
			return fmt.Errorf("bad data for %v, expected length to be multiple of %v but got %v", f, s, l)
		}
	}
	return nil
}

`)

	// Also create the Encode Decode Methods for types that are referred in arrays
	buf.WriteString("var bo = binary.LittleEndian\n\n")

	for _, st := range ts.ExportedStructTypes() {
		if st.UserDefined {
			fmt.Fprintf(buf, "func (d %v) Stride() int { return %v }\n", st.GoName(), st.CType().Size.ByteSize)
			fmt.Fprintf(buf, "func (d %v) Alignment() int { return %v }\n\n", st.GoName(), st.CType().Size.ByteAlignment)

			// Create a Encode function for the element
			fmt.Fprintf(buf, "func (e *%v) Encode(d []byte) {\n", st.GoName())
			printStructEncodes(buf, 0, st, "")
			buf.WriteString("}\n\n")
			// Create a Decode function for the element
			fmt.Fprintf(buf, "func (e *%v) Decode(d []byte) {\n", st.GoName())
			printStructDecodes(buf, 0, st, "")
			buf.WriteString("}\n\n")
		}
	}

}

func decodeType(buf io.Writer, parentPos int, ty *types.CType, head string) (offset int) {
	if ty.ArrayLen() != 0 {
		// TODO: write this into an array in go instead so the generated code does not become so long?
		for i := 0; i < ty.ArrayLen(); i++ {
			offset += decodeType(buf, parentPos+offset, ty.Array.Ty, head+fmt.Sprintf("[%v]", i))
		}
		return
	}

	if len(ty.Fields) > 0 {
		for _, f := range ty.Fields {
			offset += decodeType(buf, parentPos+offset, f.Ty, head+"."+f.GoName())
		}
		return
	}

	if ty.Name == "mat2" || ty.Name == "mat3" || ty.Name == "mat4" {
		size, _ := strconv.Atoi(ty.Name[len(ty.Name)-1:])
		for i := 0; i < size*size; i++ {
			offset += 4 // float is always 4...
			if size == 3 && i > 0 && i%3 == 0 {
				fmt.Fprintf(buf, "// 4 byte spill for mat3 alignment\n")
				offset += 4
			}
			fmt.Fprintf(buf, "bo.PutUint32(d[%v:], math.Float32bits(e%v[%v]))\n", parentPos+offset, head, i)
		}
		return
	}

	switch ty.BasicGoType() {
	case "bool":
		fmt.Fprintf(buf, "\tif bo.Uint32(d[%v:]) == 0 {\n", parentPos)
		fmt.Fprintf(buf, "\t\te%v = false\n", head)
		fmt.Fprintf(buf, "\t} else {\n")
		fmt.Fprintf(buf, "\t\te%v = true\n", head)
		fmt.Fprintf(buf, "\t}\n")
		offset += 4
	case "float32":
		fmt.Fprintf(buf, "\te%v = math.Float32frombits(bo.Uint32(d[%v:]))\n", head, parentPos)
		offset += 4
	case "int32":
		fmt.Fprintf(buf, "\te%v = int32(bo.Uint32(d[%v:]))\n", head, parentPos)
		offset += 4
	case "uint32":
		fmt.Fprintf(buf, "\te%v = bo.Uint32(d[%v:])\n", head, parentPos)
		offset += 4
	default:
		// So this is a struct type we have defined, we thus need to use it's Decode method in turn to get it correctly.
		fmt.Fprintf(buf, "\t(&e%v).Decode(d[%v:]) \n", head, parentPos)
		offset += ty.Size.ByteSize
	}
	return
}

func printStructDecodes(buf io.Writer, parentPos int, t *types.GlslType, head string) int {
	return decodeType(buf, parentPos, t.C, head)
}

func encodeType(buf io.Writer, parentPos int, ty *types.CType, head string) (offset int) {
	if ty.ArrayLen() != 0 {
		// TODO: write this into an array in go instead so the generated code does not become so long?
		for i := 0; i < ty.ArrayLen(); i++ {
			offset += encodeType(buf, parentPos+offset, ty.Array.Ty, head+fmt.Sprintf("[%v]", i))
		}
		return
	}

	if len(ty.Fields) > 0 {
		for _, f := range ty.Fields {
			offset += encodeType(buf, parentPos+offset, f.Ty, head+"."+f.GoName())
		}
		return
	}

	if ty.Name == "mat2" || ty.Name == "mat3" || ty.Name == "mat4" {
		size, _ := strconv.Atoi(ty.Name[len(ty.Name)-1:])
		for i := 0; i < size*size; i++ {
			offset += 4 // float is always 4...
			if size == 3 && i > 0 && i%3 == 0 {
				fmt.Fprintf(buf, "// 4 byte spill for mat3 alignment\n")
				offset += 4
			}
			fmt.Fprintf(buf, "bo.PutUint32(d[%v:], math.Float32bits(e%v[%v]))\n", parentPos+offset, head, i)
		}
		return
	}

	switch ty.BasicGoType() {
	case "bool":
		fmt.Fprintf(buf, "\tbo.PutUint32(d[%v:], uint32(cBool(e%v)))\n", parentPos, head)
		offset += 4
	case "float32":
		fmt.Fprintf(buf, "\tbo.PutUint32(d[%v:], math.Float32bits(e%v))\n", parentPos, head)
		offset += 4
	case "int32":
		fmt.Fprintf(buf, "\tbo.PutUint32(d[%v:], uint32(e%v))\n", parentPos, head)
		offset += 4
	case "uint32":
		fmt.Fprintf(buf, "\tbo.PutUint32(d[%v:], e%v)\n", parentPos, head)
		offset += 4
	default:
		// So this is a struct type we have defined, we thus need to use it's Decode method in turn to get it correctly.
		fmt.Fprintf(buf, "\t(&e%v).Encode(d[%v:]) \n", head, parentPos)
		offset += ty.Size.ByteSize
	}
	return
}

// TODO: mae these call each other instead of re-encoding everything...
func printStructEncodes(buf io.Writer, parentPos int, t *types.GlslType, head string) int {
	return encodeType(buf, parentPos, t.C, head)
}

func goNameCType(s string) string {
	s = "C." + s
	for strings.HasSuffix(s, "*") {
		s = "*" + string(s[:len(s)-1])
	}
	return s
}

// how we call it above:
// fmt.Fprintf(buf, "\t%v: %v", cf.Name, go2c(chc, inp, cf, "\t\t", "bind"))

// writes the right hand side of the : only
func go2c(chc io.Writer, cf types.CField, ty *types.CType, head string) (str string) {
	if ty.IsSlice() {
		// we can do nothing but cast it to unsafe - else we would have to copy all the data...
		fmt.Fprintf(chc, "\tif err := ensureLength(\"%v\", len(%v), %v, %v); err != nil { return err }\n",
			head, head, ty.Array.Ty.Size.ByteSize, -1)
		return fmt.Sprintf("unsafe.Pointer(&%v[0])", head)
	}

	// the actual stuff below:
	buf := bytes.NewBuffer(nil)
	if ty.ArrayLen() != 0 {
		fmt.Fprintf(buf, "%v{\n", ty.GoCTypeName())
		for i := 0; i < ty.ArrayLen(); i++ {
			io.WriteString(buf, go2c(chc, cf, ty.Array.Ty, head+fmt.Sprintf("[%v]", i))+",\n")
		}
		fmt.Fprintf(buf, "}")
		return buf.String()
	}

	if len(ty.Fields) > 0 {
		fmt.Fprintf(buf, " C.%v{\n", ty.Name)
		for _, f := range ty.Fields {
			fmt.Fprintf(buf, "%v: %v,\n", f.Name, go2c(chc, f, f.Ty, head+"."+f.GoName()))
		}
		fmt.Fprintf(buf, "}")
		return buf.String()
	}

	if ty.Name == "Bool" {
		return fmt.Sprintf("(%v)(cBool(%v.%v))", cf.Ty.GoCTypeName(), head, cf.GoName())
	}

	return fmt.Sprintf("(%v)(%v)", ty.GoCTypeName(), head)
}
