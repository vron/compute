package main

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/vron/compute/glbind/input"
	"github.com/vron/compute/glbind/types"
)

// can we do this instead? static_assert( sizeof( double ) == sizeof( int64_t ), "" ) ;

func generateAlign(inp input.Input, ts *types.Types) {
	f, err := os.Create(filepath.Join(fOut, "generated/align.hpp"))

	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	buf := bufio.NewWriter(f)
	defer buf.Flush()
	buf.WriteString("#pragma once\n")
	buf.WriteString("#include <cerrno>\n\nbool Kernel::ensure_alignments(cptc_data *d) {\n\t(void)d;\n\n")

	ensureNoNullptr(buf, inp, ts)
	ensureAlignment(buf, inp, ts)
	ensureSizes(buf, inp, ts)

	/*
		cSize := func(ty string, al int) {
			msg := fmt.Sprintf("static check failed: sizeof(%v) != %v", ty, al)
			fmt.Fprintf(buf, "\tif(sizeof(%v) != %v) { return this->set_error(EINVAL, \"%v\"); };\n", ty, al, msg)
		}
		cAlign := func(ty string, al int) {
			msg := fmt.Sprintf("static check failed: alignof(%v) != %v", ty, al)
			fmt.Fprintf(buf, "\tif(alignof(%v) != %v) { return this->set_error(EINVAL, \"%v\"); };\n", ty, al, msg)
		}
		for _, v := range ts.AllTypes() {
			cSize(v.Name, v.CType().Size.ByteSize)
			cAlign(v.Name, v.CType().Size.ByteAlignment)
		}

		// ensure alignments of incoming pointers since provided by user
		for _, a := range inp.Arguments {
			cf := types.CField{
				Name:  a.Name,
				CType: ts.MaybeCreateArrayType(a.Ty, a.Arrno),
			}
			recChecAlignment(buf, inp, cf, "d.")
		}
	*/
	buf.WriteString("\n\treturn true;\n")
	buf.WriteString("}\n")

}

func ensureNoNullptr(buf *bufio.Writer, inp input.Input, ts *types.Types) {
	buf.WriteString("\t// first ensure the user has provided all data\n")
	for _, arg := range inp.Arguments {
		fmt.Fprintf(buf, `	if(d->%v == nullptr) return this->set_error(EINVAL, "no data was provided for %v");%v`, arg.Name, arg.Name, "\n")
	}
	buf.WriteString("\n")
}

func ensureAlignment(buf *bufio.Writer, inp input.Input, ts *types.Types) {
	buf.WriteString("\t// ensure that the provided pointers have the expected alignment\n")
	for _, arg := range inp.Arguments {
		ty := types.CreateArray(ts.Get(arg.Ty).C, arg.Arrno)
		fmt.Fprintf(buf, `	if((((uintptr_t)(const void *)(d->%v)) %% (%v)) != 0) return this->set_error(EINVAL, "%v was not aligned to %v byte address");%v`, arg.Name, ty.Size.ByteAlignment, arg.Name, ty.Size.ByteAlignment, "\n")
	}
	buf.WriteString("\n")
}

func ensureSizes(buf *bufio.Writer, inp input.Input, ts *types.Types) {
	buf.WriteString("\t// ensure that the provided data have lengths that match the expected\n")
	for _, arg := range inp.Arguments {
		ty := types.CreateArray(ts.Get(arg.Ty).C, arg.Arrno)
		op := "-"
		if ty.IsArray() && ty.Array.Len == -1 {
			op = "%"
		}
		fmt.Fprintf(buf, `	if(d->%v_len %v %v != 0) return this->set_error(EINVAL, "the data provided for %v must have a length compatible with %v");%v`, arg.Name, op, ty.Size.ByteSize, arg.Name, ty.Size.ByteSize, "\n")
	}
	buf.WriteString("\n")
}
