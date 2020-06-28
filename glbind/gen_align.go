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

	// chec alignments and sizes of the computation types we will use
	cSize := func(ty string, al int) {
		erno := newError("static check failed: sizeof(%v) != %v", ty, al)
		fmt.Fprintf(buf, "\tif(sizeof(%v) != %v) { return %v; };\n", ty, al, erno)
		//fmt.Fprintf(buf, "\tif(sizeof(%v) != %v) { printf(\"%v\\n\", sizeof(%v)); return %v; };\n", ty, al, "%d", ty, erno)
	}
	cAlign := func(ty string, al int) {
		erno := newError("static check failed: alignof(%v) != %v", ty, al)
		fmt.Fprintf(buf, "\tif(alignof(%v) != %v) { return %v; };\n", ty, al, erno)
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
	buf.WriteString("\n\treturn 0;\n")

}

func recChecAlignment(buf *bufio.Writer, inp Input, cf CField, head string) {
	if cf.Ty.IsSlice {
		// this is a slice, will be provided as pointer so chec it
		ai := cf.Ty.Size
		erno := newError("the argument %v provided was not aligned to a %v byte boundary as required", head+cf.Name, ai.ByteAlignment)
		fmt.Fprintf(buf, "\tif((((uintptr_t)(const void *)(%v)) %% (%v)) != 0) { return %v; };\n", head+cf.Name, ai.ByteSize, erno)

		return
	}
	// ensure all fields for structs
	for _, f := range cf.Ty.Fields {
		recChecAlignment(buf, inp, f, head+cf.Name+".")
	}
}
