package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func generateSharedH(inp Input) {
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

#include "errno.h"
#include "stdint.h"

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

	// Write all complex types
	for _, st := range types.ExportedStructTypes() {
		fmt.Fprintf(buf, "typedef struct {\n")
		for _, f := range st.CType().Fields {
			buf.WriteString("  " + f.String() + ";\n")
		}
		fmt.Fprintf(buf, "} %v;\n\n", st.CType().Name)
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
	for _, arg := range inp.Arguments {
		ty := maybeCreateArrayType(arg.Ty, arg.Arrno)
		cf := CField{Name: arg.Name, Ty: ty}
		buf.WriteString("  " + cf.String() + ";\n")
	}
	buf.WriteString(`} cpt_data;

/*
  cpt_new_kernel creates a new computational kernel using a maximum of num_t
  threads for kernel calculation, returning a reference to the kernel created.
  If there is insufficient memory available to create a new kernel 0 is
  returned. For all other possible errors a kernel reference is returned and
  the next call to cpt_dispatch_kernel will return the error information.
  cpt_new_kernel is safe for concurrent use from multiple threads.
*/
exported_func void *cpt_new_kernel(int32_t num_t);

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
