package main

import (
	"bufio"
	"fmt"
	"io"
	"log"
	"os"
	"path/filepath"
)

func generateComp(inp Input) {
	f, err := os.Create(filepath.Join(fOut, "kernel.hpp"))
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	buf := bufio.NewWriter(f)
	defer buf.Flush()

	fmt.Fprintf(buf, `
#define _cpt_WG_SIZE_X %v
#define _cpt_WG_SIZE_Y %v
#define _cpt_WG_SIZE_Z %v

#import <math.h>

`, inp.Wg_size[0], inp.Wg_size[1], inp.Wg_size[2])

	writeSharedStruct(buf, inp)

	buf.WriteString("struct kernel_comp {\n")
	buf.WriteString("\tuvec3 gl_NumWorkGroups;\n")
	buf.WriteString("\tuvec3 gl_WorkGroupSize;\n")
	buf.WriteString("\tuvec3 gl_WorkGroupID;\n")
	buf.WriteString("\tuvec3 gl_LocalInvocationID;\n")
	buf.WriteString("\tuvec3 gl_GlobalInvocationID;\n")
	buf.WriteString("\tuint32_t gl_LocalInvocationIndex;\n")
	buf.WriteString("\tthread_data *thread;\n\n")

	// write all the globals we should be able to access
	for _, arg := range inp.Arguments {
		cf := CField{Name: arg.Name, Ty: maybeCreateArrayType(arg.Ty, arg.Arrno)}
		fmt.Fprintf(buf, cf.CxxFieldString()+"\n")
	}

	// also write all the shared variabels we should be able to access
	for _, arg := range inp.Shared {
		cf := CField{Name: arg.Name, Ty: maybeCreateArrayType(arg.Ty, []int{-1})} // TODO: chec this, we always do array here to use pointer
		fmt.Fprintf(buf, cf.CxxFieldString()+"\n")
	}

	buf.WriteString("\n")
	buf.WriteString("\tkernel_comp() {};\n")
	buf.WriteString("\n")

	buf.WriteString(inp.Body)
	buf.WriteString("\n")
	buf.WriteString("\n")
	buf.WriteString(`
	int set_data(cpt_data d) {
		#include "setdata.hpp"
	}
`)
	buf.WriteString(`
	void barrier();
`)

	// Generate a function for allocating and de-allocating shared data, and one for binding it to
	// a class instance.
	// TODO: handle allocation errors...
	buf.WriteString(`
	static shared_data_t* create_shared_data() {
		shared_data_t *sd = new shared_data_t();
		`)
	// allocate each of the shared buffers we need
	for _, arg := range inp.Shared {
		cf := CField{Name: arg.Name, Ty: maybeCreateArrayType(arg.Ty, arg.Arrno)}
		fmt.Fprintf(buf, "sd->%v = (%v*)malloc(%v*sizeof(%v));\n", arg.Name, cf.Ty.ty.Name, cf.Ty.ArrayLen, cf.Ty.ty.Name)
	}
	buf.WriteString(`
		return sd;
	}

	static void free_shared_data(shared_data_t *sd) {
		`)
	for _, arg := range inp.Shared {
		fmt.Fprintf(buf, "free(sd->%v);\n", arg.Name)
	}
	buf.WriteString(`
		delete sd;
	}

	void set_shared_data(shared_data_t *sd) {
		`)

	for _, arg := range inp.Shared {
		fmt.Fprintf(buf, "this->%v = sd->%v;\n", arg.Name, arg.Name)
	}
	buf.WriteString(`
	}
};
`)
}

func writeSharedStruct(buf io.Writer, inp Input) {
	fmt.Fprintf(buf, "class  shared_data_t {\npublic:\n")
	for _, arg := range inp.Shared {
		cf := CField{Name: arg.Name, Ty: maybeCreateArrayType(arg.Ty, []int{-1})}
		fmt.Fprintf(buf, cf.CxxFieldString()+"\n")
	}
	fmt.Fprintf(buf, "} ;\n\n")
}
