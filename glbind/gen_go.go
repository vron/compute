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
)

func generateGo(inp Input) {
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

type kernel struct {
	k unsafe.Pointer
	dead bool
}

`)

	// Write the struct definitions
	for _, st := range types.ExportedStructTypes() {
		fmt.Fprintf(buf, "type %v struct {\n", st.GoName())
		for _, a := range st.CType().Fields {
			fmt.Fprintf(buf, "\t%v %v\n", a.GoName(), a.Ty.GoName())
		}
		fmt.Fprintf(buf, "}\n\n")
	}

	buf.WriteString(`type Data struct {
`)
	for _, arg := range inp.Arguments {
		cf := CField{Name: arg.Name, Ty: maybeCreateArrayType(arg.Ty, arg.Arrno)}
		buf.WriteString("\t" + cf.GoName() + " " + cf.Ty.GoName() + "\n")
	}
	buf.WriteString(`}

// New creates a kernel using at most numCPU+1 threads. If numCPU <= 0 the
// number of threads to use will be calculated automatically. All kernels
// must be explicitly freed using kernel.Free to avoid memory leaks.
func New(numCPU int) (k *kernel, err error) {
	k = &kernel{}
	if numCPU <= 0 {
		numCPU = runtime.NumCPU()+2
	}
	k.k = C.cpt_new_kernel(C.int(numCPU));
	if k.k == nil {
		return nil, errors.New("failed to create kernel structure")
	}
	runtime.SetFinalizer(k, freeKernel)
	return k, nil
}

// Dispatch a kernel calculation of the specified size. The caller must ensure
// that the data provided in bind matches the kernel's assumptions and that any
// []byte field represents properly aligned data. Not data in bind must
// be accessed (read or write) until Dispatch returns.
func (k *kernel) Dispatch(bind Data, numGroupsX, numGroupsY, numGroupsZ int) error {
	if k.dead {
		panic("cannot use a kernel where Free() has been called")
	}
	cbind := C.cpt_data{
`)

	// Create a c-struct from the provided that that is uploaded to the ernel
	chc := bytes.NewBuffer(nil)
	for _, arg := range inp.Arguments {
		cf := CField{Name: arg.Name, Ty: maybeCreateArrayType(arg.Ty, arg.Arrno)}
		fmt.Fprintf(buf, "\t%v: %v", cf.Name, go2c(chc, inp, cf, "\t\t", "bind"))
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

// Free dealocates any data allocated by the underlying kernel. Note that
// a kernel on which Free has been called can no longer be used.
func (k *kernel) Free() {
	freeKernel(k)
}


func freeKernel(k *kernel) {
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

	for _, st := range types.ExportedStructTypes() {
		if st.userType {
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
	/*

	   // Encode encodes e as a std430 glsl struct to buf. If buf if not long
	   // enough to hold e it will panic. Encode is an expensive call and should
	   // normally not be used for a large number of Elements.
	   type (e *Element) Encode(buf []byte) {

	   }


	   // Decodes data from buf and fills e. If buf is not long enough the call
	   // will panic.
	   type (e *Element) Decode(buf []byte) {

	   }

	*/

}

func printStructDecodes(buf io.Writer, parentPos int, t *Type, head string) int {
	for _, f := range t.CType().Fields {
		if f.Ty.IsSlice {
			panic("cannot encode decode struct with slice, we do not have the size")
		}
		if !f.Ty.ty.userType && len(f.Ty.Fields) > 0 {
			panic("cannot have struct with build in complex type to Encode / Decode")
		}
		no := f.Ty.ArrayLen
		if no == 0 {
			no = 1
		}
		for i := 0; i < no; i++ {
			// special case matrices since they are not contigiuous but have 0 elements in them that we need handle...

			arrelem := ""
			if f.Ty.ArrayLen != 0 {
				arrelem = fmt.Sprintf("[%v]", i)
			}
			tt := f.Ty.BasicGoType()
			if f.Ty.ty.Name == "mat2" || f.Ty.ty.Name == "mat3" || f.Ty.ty.Name == "mat4" {
				size, _ := strconv.Atoi(f.Ty.ty.Name[len(f.Ty.ty.Name)-1:])
				nonopass := 0
				if size == 3 {
					nonopass = i / size
				}
				ii := i + nonopass
				fmt.Fprintf(buf, "\te%v.%v%v = math.Float32frombits(bo.Uint32(d[%v:]))\n", head, f.GoName(), arrelem, parentPos+f.ByteOffset+ii*4)
			} else {
				underlysize := f.Ty.Size.ByteSize
				if f.Ty.ArrayLen != 0 {
					underlysize /= f.Ty.ArrayLen
				}
				switch tt {
				case "bool":
					fmt.Fprintf(buf, "\tif bo.Uint32(d[%v:]) == 0 {\n", parentPos+f.ByteOffset+i*4)
					fmt.Fprintf(buf, "\t\te%v.%v%v = false\n", head, f.GoName(), arrelem)
					fmt.Fprintf(buf, "\t} else {\n")
					fmt.Fprintf(buf, "\t\te%v.%v%v = true\n", head, f.GoName(), arrelem)
					fmt.Fprintf(buf, "\t}\n")
				case "float32":
					fmt.Fprintf(buf, "\te%v.%v%v = math.Float32frombits(bo.Uint32(d[%v:]))\n", head, f.GoName(), arrelem, parentPos+f.ByteOffset+i*4)
				case "int32":
					fmt.Fprintf(buf, "\te%v.%v%v = int32(bo.Uint32(d[%v:]))\n", head, f.GoName(), arrelem, parentPos+f.ByteOffset+i*4)
				case "uint32":
					fmt.Fprintf(buf, "\te%v.%v%v = bo.Uint32(d[%v:])\n", head, f.GoName(), arrelem, parentPos+f.ByteOffset+i*4)
				default:
					// So this is a struct type we have defined, we thus need to use it's Decode method in turn to get it correctly.
					fmt.Fprintf(buf, "\t(&e.%v%v).Decode(d[%v:]) \n", f.GoName(), arrelem, parentPos+f.ByteOffset+i*underlysize)
				}
			}
		}
	}
	return 0
}

func printStructEncodes(buf io.Writer, parentPos int, t *Type, head string) int {
	for _, f := range t.CType().Fields {
		if f.Ty.IsSlice {
			panic("cannot encode decode struct with slice, we do not have the size")
		}
		if !f.Ty.ty.userType && len(f.Ty.Fields) > 0 {
			panic("cannot have struct with build in complex type to Encode / Decode")
		}
		no := f.Ty.ArrayLen
		if no == 0 {
			no = 1
		}
		for i := 0; i < no; i++ {
			arrelem := ""
			if f.Ty.ArrayLen != 0 {
				arrelem = fmt.Sprintf("[%v]", i)
			}
			tt := f.Ty.BasicGoType()
			if f.Ty.ty.Name == "mat2" || f.Ty.ty.Name == "mat3" || f.Ty.ty.Name == "mat4" {
				size, _ := strconv.Atoi(f.Ty.ty.Name[len(f.Ty.ty.Name)-1:])
				nonopass := 0
				if size == 3 {
					nonopass = i / size
				}
				ii := i + nonopass
				fmt.Fprintf(buf, "\tbo.PutUint32(d[%v:], math.Float32bits(e%v.%v%v))\n", parentPos+f.ByteOffset+ii*4, head, f.GoName(), arrelem)
			} else {
				underlysize := f.Ty.Size.ByteSize
				if f.Ty.ArrayLen != 0 {
					underlysize /= f.Ty.ArrayLen
				}
				switch tt {
				case "bool":
					fmt.Fprintf(buf, "\tbo.PutUint32(d[%v:], uint32(cBool(e%v.%v%v)))\n", parentPos+f.ByteOffset+i*4, head, f.GoName(), arrelem)
				case "float32":
					fmt.Fprintf(buf, "\tbo.PutUint32(d[%v:], math.Float32bits(e%v.%v%v))\n", parentPos+f.ByteOffset+i*4, head, f.GoName(), arrelem)
				case "int32":
					fmt.Fprintf(buf, "\tbo.PutUint32(d[%v:], uint32(e%v.%v%v))\n", parentPos+f.ByteOffset+i*4, head, f.GoName(), arrelem)
				case "uint32":
					fmt.Fprintf(buf, "\tbo.PutUint32(d[%v:], e%v.%v%v)\n", parentPos+f.ByteOffset+i*4, head, f.GoName(), arrelem)
				default:
					// So this is a struct type we have defined, we thus need to use it's Decode method in turn to get it correctly.
					fmt.Fprintf(buf, "\t(&e.%v%v).Encode(d[%v:]) \n", f.GoName(), arrelem, parentPos+f.ByteOffset+i*underlysize)
				}
			}
		}
	}
	return 0
}

func goNameCType(s string) string {
	s = "C." + s
	for strings.HasSuffix(s, "*") {
		s = "*" + string(s[:len(s)-1])
	}
	return s
}

// writes the right hand side of the : only
func go2c(chc io.Writer, inp Input, cf CField, indent, head string) (str string) {
	buf := bytes.NewBuffer(nil)
	if cf.Ty.IsSlice {
		// we can do nothing but cast it to unsafe - else we would have to copy all the data...
		fmt.Fprintf(chc, "\tif err := ensureLength(\"%v.%v\", len(%v.%v), %v, %v); err != nil { return err }\n",
			head, cf.GoName(), head, cf.GoName(), cf.Ty.Size.ByteSize, -cf.Ty.ArrayLen)
		return fmt.Sprintf("unsafe.Pointer(&%v.%v[0]),\n", head, cf.GoName())
	}
	if len(cf.Ty.Fields) > 0 {
		// struct type
		if cf.Ty.ArrayLen > 0 {
			fmt.Fprintf(buf, indent+"%v{\n", cf.Ty.GoCTypeName())
			for i := 0; i < cf.Ty.ArrayLen; i++ {
				fmt.Fprintf(buf, " C.%v{\n", cf.Ty.Name)
				for _, f := range cf.Ty.Fields {
					fmt.Fprintf(buf, indent+"%v: %v", f.Name, go2c(chc, inp, f, indent+"\t", head+"."+cf.GoName()+fmt.Sprintf("[%v]", i)))
				}
				fmt.Fprintf(buf, indent+"},\n")
			}
			fmt.Fprintf(buf, indent+"},\n")
			return buf.String()
		} else {
			fmt.Fprintf(buf, " C.%v{\n", cf.Ty.Name)
			for _, f := range cf.Ty.Fields {
				fmt.Fprintf(buf, indent+"%v: %v", f.Name, go2c(chc, inp, f, indent+"\t", head+"."+cf.GoName()))
			}
			fmt.Fprintf(buf, indent+"},\n")
			return buf.String()
		}
	}

	if cf.Ty.ArrayLen > 0 {
		// this is an array - cannot cast directly so must create
		fmt.Fprintf(buf, indent+"%v{\n", cf.Ty.GoCTypeName())
		for i := 0; i < cf.Ty.ArrayLen; i++ {
			fmt.Fprintf(buf, indent+"\t(%v)(%v.%v[%v]),\n", goNameCType(cf.Ty.Name), head, cf.GoName(), i)
		}
		fmt.Fprintf(buf, indent+"},\n")
		return buf.String()
	}

	if cf.Ty.ty.Name == "Bool" {
		return fmt.Sprintf("(%v)(cBool(%v.%v)),\n", cf.Ty.GoCTypeName(), head, cf.GoName())
	}

	// TODO: Need to handle struct type here also

	return fmt.Sprintf("(%v)(%v.%v),\n", cf.Ty.GoCTypeName(), head, cf.GoName())

	//panic("unhandled type in go2c conversion")
}
