package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
)

func generateAlignH(inp Input) {
	f, err := os.Create(filepath.Join(fOut, "align.hpp"))

	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	buf := bufio.NewWriter(f)
	defer buf.Flush()

	buf.WriteString("#include <cerrno>\n\nbool kernel::ensure_alignments(cpt_data d) {\n\t(void)d;\n")

	cSize := func(ty string, al int) {
		msg := fmt.Sprintf("static check failed: sizeof(%v) != %v", ty, al)
		fmt.Fprintf(buf, "\tif(sizeof(%v) != %v) { return this->set_error(EINVAL, \"%v\"); };\n", ty, al, msg)
	}
	cAlign := func(ty string, al int) {
		msg := fmt.Sprintf("static check failed: alignof(%v) != %v", ty, al)
		fmt.Fprintf(buf, "\tif(alignof(%v) != %v) { return this->set_error(EINVAL, \"%v\"); };\n", ty, al, msg)
	}
	for _, v := range types.AllTypes() {
		cSize(v.Name, v.CType().Size.ByteSize)
		cAlign(v.Name, v.CType().Size.ByteAlignment)
	}

	// ensure alignments of incoming pointers since provided by user
	for _, a := range inp.Arguments {
		cf := CField{
			Name: a.Name,
			Ty:   maybeCreateArrayType(a.Ty, a.Arrno),
		}
		recChecAlignment(buf, inp, cf, "d.")
	}

	buf.WriteString("\n\treturn true;\n")
	buf.WriteString("}\n")

}

func recChecAlignment(buf *bufio.Writer, inp Input, cf CField, head string) {
	if cf.Ty.IsSlice {
		// this is a slice, will be provided as pointer so chec it
		ai := cf.Ty.Size
		msg := fmt.Sprintf("the argument %v provided was not aligned to a %v byte boundary as required", head+cf.Name, ai.ByteAlignment)
		fmt.Fprintf(buf, "\tif((((uintptr_t)(const void *)(%v)) %% (%v)) != 0) { return this->set_error(EINVAL, \"%v\"); };\n", head+cf.Name, ai.ByteSize, msg)

		return
	}
	// ensure all fields for structs
	for _, f := range cf.Ty.Fields {
		recChecAlignment(buf, inp, f, head+cf.Name+".")
	}
}
