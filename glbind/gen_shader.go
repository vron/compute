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

const CPTC = ""

func generateShader(inp input.Input, ts *types.Types) {
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
	fmt.Fprintf(buf, "#define _cpt_REPEAT_WG_SIZE(x) x")
	for i := 1; i < inp.Wg_size[0]*inp.Wg_size[1]*inp.Wg_size[2]; i++ {
		fmt.Fprintf(buf, ", x")
	}
	fmt.Fprintf(buf, "\n")

	writeSharedStruct(buf, inp, ts)

	fmt.Fprintf(buf, "struct shader {\n")
	fmt.Fprintf(buf, "  uvec3 gl_NumWorkGroups;\n")
	fmt.Fprintf(buf, "  uvec3 gl_WorkGroupSize;\n")
	fmt.Fprintf(buf, "  uvec3 gl_WorkGroupID;\n")
	fmt.Fprintf(buf, "  uvec3 gl_LocalInvocationID;\n")
	fmt.Fprintf(buf, "  uvec3 gl_GlobalInvocationID;\n")
	fmt.Fprintf(buf, "  uint32_t gl_LocalInvocationIndex;\n")
	fmt.Fprintf(buf, "  co::Routine<struct shader*>  *invocation;\n\n")

	// write all the globals we should be able to access
	for _, arg := range inp.Arguments {
		t := ts.Get(arg.Ty)
		cf := types.CField{Name: arg.Name, CType: types.CreateArray(t.C, arg.Arrno)}
		fmt.Fprintf(buf, "  "+refify(cf.CType.CString(CPTC, cf.Name, false))+";\n")
	}
	fmt.Fprintf(buf, "\n")

	// also write all the shared variabels we should be able to access
	for _, arg := range inp.Shared {
		t := ts.Get(arg.Ty)
		cf := types.CField{Name: arg.Name, CType: types.CreateArray(t.C, arg.Arrno)}
		fmt.Fprintf(buf, "  "+refify(cf.CType.CString(CPTC, cf.Name, false))+";\n")
	}

	writeConstructors(buf, inp, ts)

	fmt.Fprint(buf, inp.Body)
	fmt.Fprintf(buf, "\n")
	fmt.Fprintf(buf, "\n")
	fmt.Fprintf(buf, `
	void barrier();
`)

	fmt.Fprintf(buf, `
};
`)
}

func writeSharedStruct(buf io.Writer, inp input.Input, ts *types.Types) {
	fmt.Fprintf(buf, "class  shared_data_t {\npublic:\n")
	for _, arg := range inp.Shared {
		t := ts.Get(arg.Ty)
		cf := types.CField{Name: arg.Name, CType: types.CreateArray(t.C, arg.Arrno)}
		fmt.Fprintf(buf, "  "+cf.CType.CString(CPTC, cf.Name, false)+";\n")
	}
	fmt.Fprintf(buf, "} ;\n\n")
}

func writeConstructors(buf io.Writer, inp input.Input, ts *types.Types) {

	fmt.Fprintf(buf, `
#pragma clang diagnostic push
#pragma clang diagnostic ignored "-Wunused-parameter"
`)
	fmt.Fprintf(buf, "  shader(cptc_data *d, shared_data_t *sd) ")
	if len(inp.Shared)+len(inp.Arguments) > 0 {
		fmt.Fprintf(buf, ":\n")
		for i, arg := range inp.Arguments {
			ty := types.CreateArray(ts.Get(arg.Ty).C, arg.Arrno)
			deref := "*"
			if (len(arg.Arrno) > 0 && arg.Arrno[0] == -1) || ty.IsComplexStruct() {
				deref = ""
			}
			fmt.Fprintf(buf, "        %v(%v(d->%v))", arg.Name, deref, arg.Name)
			if i != len(inp.Arguments)-1 || len(inp.Shared) != 0 {
				fmt.Fprintf(buf, ",\n")
			}
		}
		for i, arg := range inp.Shared {
			fmt.Fprintf(buf, "        %v(sd->%v)", arg.Name, arg.Name)
			if i != len(inp.Shared)-1 {
				fmt.Fprintf(buf, ",\n")
			}
		}
	}
	fmt.Fprintf(buf, "{};\n")

	// copy constructor is used for array initialization
	fmt.Fprintf(buf, "\n")
	fmt.Fprintf(buf, "\tshader(const shader& org) ")
	if len(inp.Shared)+len(inp.Arguments) > 0 {
		fmt.Fprintf(buf, ":\n")
		for i, arg := range inp.Arguments {
			fmt.Fprintf(buf, "       %v(org.%v)", arg.Name, arg.Name)
			if i != len(inp.Arguments)-1 || len(inp.Shared) != 0 {
				fmt.Fprintf(buf, ",\n")
			}
		}
		for i, arg := range inp.Shared {
			fmt.Fprintf(buf, "        %v(org.%v)", arg.Name, arg.Name)
			if i != len(inp.Shared)-1 {
				fmt.Fprintf(buf, ",\n")
			}
		}
	}
	fmt.Fprintf(buf, `{};
#pragma clang diagnostic pop
`)

	fmt.Fprintf(buf, "\n")
}

func refify(s string) string {
	// hac*y but wor*s
	sp := strings.SplitN(s, "\t", 2)
	tp, nm := sp[0], sp[1]
	if strings.HasPrefix(strings.TrimSpace(nm), "(*") {
		return s // pointer, no need to add reference
	}
	if !strings.Contains(nm, "[") {
		return tp + "&" + "\t" + nm
	}
	// so this is an array type, need to add it innermost (ref spirlal rule)
	sp = strings.SplitN(nm, "[", 2)
	return tp + "\t" + "(&" + strings.TrimSpace(sp[0]) + ")[" + sp[1]
}
