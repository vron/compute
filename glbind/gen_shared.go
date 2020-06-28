package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func generateSharedH(inp Input) {
	f, err := os.Create(filepath.Join(fOut, "shared.h"))

	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	buf := bufio.NewWriter(f)
	defer buf.Flush()
	buf.WriteString("#include \"stdint.h\"\n\n")
	//buf.WriteString("#include \"stdio.h\"\n\n")

	// Write all complex types
	for _, st := range types.ExportedStructTypes() {
		fmt.Fprintf(buf, "typedef struct {\n")
		for _, f := range st.CType().Fields {
			buf.WriteString("\t" + f.String() + ";\n")
		}
		fmt.Fprintf(buf, "} %v;\n\n", st.CType().Name)
	}

	// write the actuall data we should export
	buf.WriteString("typedef struct {\n")
	for _, arg := range inp.Arguments {
		ty := maybeCreateArrayType(arg.Ty, arg.Arrno)
		cf := CField{Name: arg.Name, Ty: ty}
		buf.WriteString("\t" + cf.String() + ";\n")
	}
	buf.WriteString(`} cpt_data; 

void *cpt_new_kernel(int32_t);
int cpt_dispatch_kernel(void *, cpt_data, int32_t, int32_t, int32_t);
void cpt_free_kernel(void *);
`)
}
