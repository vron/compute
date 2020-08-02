package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"text/tabwriter"

	"github.com/vron/compute/glbind/input"
	"github.com/vron/compute/glbind/types"
)

func generateShared(inp input.Input, ts *types.Types) {
	f, err := os.Create(filepath.Join(fOut, "generated/shared.h"))

	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	buf := bufio.NewWriter(f)
	defer buf.Flush()
	buf.WriteString(`#pragma once
/*
  This header and associated library was generated from a GLSL compute shader
  to be executed on a CPU as static code. The library is safe for threaded use
  as further specified below.
*/

/*
 In order to make the library useful on multiple platforms we define some support
 macros that will optionally be used.
*/
#ifdef _WIN64
#define exported_func __declspec(dllexport)
#else
#define exported_func
#endif

#include <errno.h>
#include <stdint.h>
#include <stdalign.h>

/*
  cpt_error_t represents an error as reported from cpt_dispatch_kernel. The 
  possible errors can mostly be classified as either user errors or underlying
  system errors. In case of underlying errors, such as insufficient resources,
  the .code field will be set to an error code from errno.h. In case of user
  errors, such as providing data with bad alignment, .code will be set to
  EINVAL with further description given in .msg. The data pointed to by .msg is
  only accessible until the next call to cpt_dispatch_kernel or cpt_free_kernel
  for the same kernel reference.
*/
struct cpt_error_t {
  int code;
  char* msg;
};

`)

	// Write all exported types complete with alignment info for user reference.
	// TODO: Add padding to all these such that they actually get what we expect
	for _, st := range ts.ListExportedTypes() {
		if st.C.IsStruct() {
			w := tabwriter.NewWriter(buf, 0, 1, 1, ' ', 0)
			fmt.Fprintf(buf, "typedef struct {  // size = %v, align = %v\n", st.C.Size.ByteSize, st.C.Size.ByteAlignment)
			offset := 0
			for i, f := range st.C.Struct.Fields {
				if offset != f.ByteOffset {
					fmt.Fprintf(w, "  char\t _pad%v[%v];\t\t\t\n", i, f.ByteOffset-offset)
				}
				offset = f.ByteOffset
				fmt.Fprintf(w, "  "+alignas(i, st.C.Size.ByteAlignment)+f.CType.CString("cpt_", f.Name, true)+";\t// offset =\t%v\t\n", f.ByteOffset)
				offset += f.CType.Size.ByteSize
			}
			if offset != st.C.Size.ByteSize {
				fmt.Fprintf(w, "  char\t _pad[%v];\t\t\t\n", st.C.Size.ByteSize-offset)
			}
			w.Flush()
			fmt.Fprintf(buf, "} %v;\n\n", st.CName("cpt_"))
		} else if st.C.IsVector() {
			w := tabwriter.NewWriter(buf, 0, 1, 1, ' ', 0)
			fmt.Fprintf(w, "typedef struct {  // size = %v, align = %v\n", st.C.Size.ByteSize, st.C.Size.ByteAlignment)
			offset := 0
			for i := 0; i < st.C.Vector.Len; i++ {
				fmt.Fprintf(w, "  %v%v\t%v;\t// offset =\t%v\t\n", alignas(i, st.C.Size.ByteAlignment), st.C.Vector.Basic.CString("", "", true), string('x'+i), st.C.Vector.Basic.Size.ByteSize*i)
				offset += st.C.Vector.Basic.Size.ByteSize
			}
			if offset != st.C.Size.ByteSize {
				fmt.Fprintf(w, "  char\t _pad[%v];\t\t\t\n", st.C.Size.ByteSize-offset)
			}
			w.Flush()
			fmt.Fprintf(buf, "} %v;\n\n", st.CName("cpt_"))
		}
	}

	// write the actually data we should export
	buf.WriteString(`/*
  cpt_data consists of all the input/output required by the compute kernel. All
  fixed sized fields (including arrays) will be copied internally to ensure
  correct alignment. For all variable sizes fields (type void*) the user must
  ensure sufficient length and data alignment for the relevant use.
*/
`)
	buf.WriteString("typedef struct {\n")
	w := tabwriter.NewWriter(buf, 0, 1, 1, ' ', 0)
	for _, arg := range inp.Arguments {
		ty := types.CreateArray(ts.Get(arg.Ty).C, arg.Arrno)
		name := arg.Name
		if !(ty.IsArray() && ty.Array.Len == -1) {
			name = "(*" + name + ")"
		}
		fmt.Fprintf(w, "  "+ty.CString("cpt_", name, true)+";\n")
		fmt.Fprintf(w, "  int64_t "+arg.Name+"_len;\n\n")
	}
	w.Flush()
	buf.WriteString(`} cpt_data;

/*
  cpt_new_kernel creates a new computational kernel using a maximum of num_t
  threads for kernel calculation, returning a reference to the kernel created.
  If there is insufficient memory available to create a new kernel 0 is
  returned. For all other possible errors a kernel reference is returned and
  the next call to cpt_dispatch_kernel will return the error information.
  cpt_new_kernel is safe for concurrent use from multiple threads. The stack size
  that each shader invocation should have access to can be specified in the last
  argument. If negative a default value of 16kB will be used.
*/
exported_func void *cpt_new_kernel(int32_t num_t, int32_t stack_size);

/*
  cpt_dispatch_kernel issues a calculation of the compute shader using x, y, z
  work groups in x, y, z directions respectively. The kernel reference k passed
  must have been created using cpt_new_kernel and not subsequently deallocated
  using cpt_free_kernel. It is the callers responsibility to ensure that any
  data of non-fixed size in d is properly aligned as required by the kernel and
  of sufficient length for the number of work groups issued. Any error message
  description returned in cpt_error_t.msg is only accessible until the next call to
  cpt_dispatch_kernel or cpt_free_kernel for the same kernel reference k.
  cpt_dispatch_kernel is safe for concurrent use by multiple threads for
  different kernel references (k) but must not be called concurrently for the
  same k.
*/
exported_func struct cpt_error_t cpt_dispatch_kernel(void *k, cpt_data d, int32_t x, int32_t y, int32_t z);

/*
  cpt_free_kernel must be called for any non-null kernel k created to avoid
  leaks. Note that any k for which cpt_free_kernel has been called is unsafe for
  any further use.
*/
exported_func void cpt_free_kernel(void *k);
`)
}

func alignas(i, v int) string {
	if i == 0 {
		return fmt.Sprintf("alignas(%v) ", v)
	}
	return ""
}
