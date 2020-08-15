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

	fmt.Fprint(buf, `
// Code generated DO NOT EDIT

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
	writeTypeDefinitionsToC(buf, inp, ts)
	writeDataStruct(buf, inp, ts)
	writeDataStructRaw(buf, inp, ts)
	writeDispatch(buf, inp, ts)
	writeDispatchRaw(buf, inp, ts)
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
		ty := types.CreateArray(ts.Get(arg.Ty).C, arg.Arrno)
		if ty.IsComplexStruct() {
			fmt.Fprintf(buf, tab+ty.CString("cpt_", "", false)+" "+arg.Name+",\n")
		} else {
			fmt.Fprintf(buf, tab+"void* "+arg.Name+", ")
			fmt.Fprintf(buf, "int64_t "+arg.Name+"_len,\n")
		}
	}
	fmt.Fprintf(buf, tab+`int32_t x, int32_t y, int32_t z) {
	cpt_data d;
`)
	for _, arg := range inp.Arguments {
		ty := types.CreateArray(ts.Get(arg.Ty).C, arg.Arrno)
		fmt.Fprintf(buf, "\td."+arg.Name+" = "+arg.Name+";\n")
		if !ty.IsComplexStruct() {
			fmt.Fprintf(buf, "\td."+arg.Name+"_len = "+arg.Name+"_len;\n")
		}
	}
	fmt.Fprintf(buf, `	return cpt_dispatch_kernel(k, d, x, y, z);
}
*/
import "C"

`)
}

func writeTypeDefinitionsToC(buf io.Writer, inp input.Input, ts *types.Types) {
	for _, st := range ts.ListExportedTypes() {
		if !st.C.IsComplexStruct() {
			continue
		}
		fmt.Fprintf(buf, "func (v %v) toC() C.cpt_%v {\nreturn C.cpt_%v{\n", st.GoName(), st.Name, st.Name)
		for _, f := range st.C.Struct.Fields {
			if f.CType.IsArray() && f.CType.Array.Len == -1 {
				// C expects only a pointer TODO: Should we also send a length here?
				fmt.Fprintf(buf, "%v:(*C.%v)(&v.%v[0]),\n", f.Name, f.CType.Array.CType.Basic.Name, strings.Title(f.Name))
			} else {
				fmt.Fprintf(buf, "%v:(C.%v)(v.%v),\n", f.Name, f.CType.Basic.Name, strings.Title(f.Name))
			}
		}
		fmt.Fprintf(buf, "}\n}\n\n")

	}
}

func writeDataStruct(buf io.Writer, inp input.Input, ts *types.Types) {
	fmt.Fprintf(buf, "type Data struct {\n")
	for _, arg := range inp.Arguments {
		ty := types.CreateArray(ts.Get(arg.Ty).C, arg.Arrno)
		if len(arg.Arrno) > 0 && arg.Arrno[0] == -1 {
			fmt.Fprintf(buf, "  "+ty.GoString(strings.Title(arg.Name))+"\n")
		} else {
			if ty.IsComplexStruct() {
				fmt.Fprintf(buf, "  "+(ty.GoString(strings.Title(arg.Name)))+"\n")
			} else {
				fmt.Fprintf(buf, "  "+pointify(ty.GoString(strings.Title(arg.Name)))+"\n")
			}
		}
	}
	fmt.Fprintf(buf, "}\n\n")
}

func writeDataStructRaw(buf io.Writer, inp input.Input, ts *types.Types) {
	fmt.Fprintf(buf, "type DataRaw struct {\n")
	for _, arg := range inp.Arguments {
		ty := types.CreateArray(ts.Get(arg.Ty).C, arg.Arrno)
		if ty.IsComplexStruct() {
			// this is the sort of complex struct that is handled completely differently in gl anyway
			fmt.Fprintf(buf, "  "+(ty.GoString(strings.Title(arg.Name)))+"\n")
		} else {
			fmt.Fprintf(buf, strings.Title(arg.Name)+" []byte\n")
		}
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
		ty := types.CreateArray(ts.Get(arg.Ty).C, arg.Arrno)
		if len(arg.Arrno) > 0 && arg.Arrno[0] == -1 {
			// slice data, we need the size of the entire thing...
			fmt.Fprintf(buf, "\tunsafe.Pointer(&bind."+strings.Title(arg.Name)+"[0]), ")
			fmt.Fprintf(buf, "C.int64_t(int64(len(bind.%v))*int64(unsafe.Sizeof(bind.%v[0]))),\n", strings.Title(arg.Name), strings.Title(arg.Name))
		} else {
			if ty.IsComplexStruct() {
				fmt.Fprintf(buf, "\tbind."+strings.Title(arg.Name)+".toC(), ")
			} else {
				fmt.Fprintf(buf, "\tunsafe.Pointer(bind."+strings.Title(arg.Name)+"), ")
				fmt.Fprintf(buf, "C.int64_t(unsafe.Sizeof(*bind."+strings.Title(arg.Name)+")),\n")
			}
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
		ty := types.CreateArray(ts.Get(arg.Ty).C, arg.Arrno)
		if ty.IsComplexStruct() {
			fmt.Fprintf(buf, "\tbind."+strings.Title(arg.Name)+".toC(), ")
		} else {
			fmt.Fprintf(buf, "\tunsafe.Pointer(&bind."+strings.Title(arg.Name)+"[0]), ")
			fmt.Fprintf(buf, "C.int64_t(int64(len(bind.%v))),\n", strings.Title(arg.Name))
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

func pointify(s string) string {
	// hac*y but wor*s
	return strings.Replace(s, "\t", "\t*", 1)
}
