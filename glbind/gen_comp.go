package main

import (
	"bufio"
	"fmt"
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

`, inp.Wg_size[0], inp.Wg_size[1], inp.Wg_size[2])

	buf.WriteString("struct kernel_comp {\n")
	buf.WriteString("\tuvec3 gl_GlobalInvocationID;\n")
	buf.WriteString("\tuvec3 gl_LocalInvocationID;\n\n")

	// write all the globals we should be able to access
	for _, arg := range inp.Arguments {
		cf := CField{Name: arg.Name, Ty: maybeCreateArrayType(arg.Ty, arg.Arrno)}
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
};
`)
}
