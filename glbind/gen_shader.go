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
	f, err := os.Create(filepath.Join(fOut, "generated/shader.hpp"))
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	buf := bufio.NewWriter(f)
	defer buf.Flush()

	fmt.Fprintf(buf, `#pragma once

#define _cpt_WG_SIZE_X %v
#define _cpt_WG_SIZE_Y %v
#define _cpt_WG_SIZE_Z %v
#define _cpt_WG_SIZE %v

#include <cmath>
#include "../types/types.hpp"
#include "usertypes.hpp"
#include "../co/routines.hpp"

`, inp.Wg_size[0], inp.Wg_size[1], inp.Wg_size[2], inp.Wg_size[0]*inp.Wg_size[1]*inp.Wg_size[2])

	// ugly hac to manage array initialization of the shaders
	buf.WriteString("#define _cpt_REPEAT_WG_SIZE(x) x")
	for i := 1; i < inp.Wg_size[0]*inp.Wg_size[1]*inp.Wg_size[2]; i++ {
		buf.WriteString(", x")
	}
	buf.WriteString("\n")

	writeSharedStruct(buf, inp)

	buf.WriteString("struct shader {\n")
	buf.WriteString("\tuvec3 gl_NumWorkGroups;\n")
	buf.WriteString("\tuvec3 gl_WorkGroupSize;\n")
	buf.WriteString("\tuvec3 gl_WorkGroupID;\n")
	buf.WriteString("\tuvec3 gl_LocalInvocationID;\n")
	buf.WriteString("\tuvec3 gl_GlobalInvocationID;\n")
	buf.WriteString("\tuint32_t gl_LocalInvocationIndex;\n")
	buf.WriteString("\tco::Routine<struct shader*>  *invocation;\n\n")

	// write all the globals we should be able to access
	for _, arg := range inp.Arguments {
		cf := CField{Name: arg.Name, Ty: maybeCreateArrayType(arg.Ty, arg.Arrno)}
		fmt.Fprintf(buf, cf.CxxFieldString()+"\n")
	}

	// also write all the shared variabels we should be able to access
	for _, arg := range inp.Shared {
		cf := CField{Name: arg.Name, Ty: maybeCreateArrayType(arg.Ty, arg.Arrno)}
		fmt.Fprintf(buf, cf.CxxFieldStringRef()+"\n")
	}

	buf.WriteString(`
#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wunused-parameter"
`)
	buf.WriteString("\tshader(shared_data_t *sd) ")
	if len(inp.Shared) > 0 {
		buf.WriteString(": ")
		for i, arg := range inp.Shared {
			fmt.Fprintf(buf, "%v(sd->%v)", arg.Name, arg.Name)
			if i != len(inp.Shared)-1 {
				buf.WriteString(", ")
			}
		}
	}
	buf.WriteString("{};\n")

	// copy constructor is used for array initialization
	buf.WriteString("\n")
	buf.WriteString("\tshader(const shader& org) ")
	if len(inp.Shared) > 0 {
		buf.WriteString(": ")
		for i, arg := range inp.Shared {
			fmt.Fprintf(buf, "%v(org.%v)", arg.Name, arg.Name)
			if i != len(inp.Shared)-1 {
				buf.WriteString(", ")
			}
		}
	}
	buf.WriteString(`{};
#pragma clang diagnostic pop
`)

	buf.WriteString("\n")

	buf.WriteString(inp.Body)
	buf.WriteString("\n")
	buf.WriteString("\n")
	buf.WriteString(`
	void set_data(cpt_data d) {
	auto me = this;
`)
	generateSetData(buf, inp)
	buf.WriteString(`
	}
`)
	buf.WriteString(`
	void barrier();
`)

	buf.WriteString(`
};
`)
}

func writeSharedStruct(buf io.Writer, inp Input) {
	fmt.Fprintf(buf, "class  shared_data_t {\npublic:\n")
	for _, arg := range inp.Shared {
		cf := CField{Name: arg.Name, Ty: maybeCreateArrayType(arg.Ty, arg.Arrno)}
		fmt.Fprintf(buf, cf.CxxFieldString()+"\n")
	}
	fmt.Fprintf(buf, "} ;\n\n")
}

func generateSetData(buf *bufio.Writer, inp Input) {
	// we have an variable d of the data struct type that we need
	// to translate to the member variables.

	for _, a := range inp.Arguments {
		cf := CField{
			Name: a.Name,
			Ty:   maybeCreateArrayType(a.Ty, a.Arrno),
		}
		cf.CxxBinding(buf)
	}

	buf.WriteString("\treturn;")
}
