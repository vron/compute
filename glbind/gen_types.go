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

func generateTypes(inp input.Input, ts *types.Types) {
	f, err := os.Create(filepath.Join(fOut, "generated/usertypes.hpp"))
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	buf := bufio.NewWriter(f)
	defer buf.Flush()

	buf.WriteString("#pragma once\n")
	buf.WriteString("// Code generated DO NOT EDIT\n")
	buf.WriteString("#include \"../types/types.hpp\"\n")
	buf.WriteString("#include \"./shared.h\"\n\n")

	// generate all the srructs needed (vectors allread defined, so ar matrices) and chec all the sizes
	// and alignments, similiarly chec the shared ones here since we do not want to polute the shared.h
	// file with all that info.
	for _, st := range ts.ListExportedTypes() {
		if st.C.IsStruct() && !st.Builtin {
			w := tabwriter.NewWriter(buf, 0, 1, 1, ' ', 0)
			fmt.Fprintf(buf, "typedef struct {  // size = %v, align = %v\n", st.C.Size.ByteSize, st.C.Size.ByteAlignment)
			for _, f := range st.C.Struct.Fields {
				fmt.Fprintf(w, "  "+f.CType.CString(CPTC, f.Name, false)+";\t// offset =\t%v\t\n", f.ByteOffset)
			}
			w.Flush()
			fmt.Fprintf(buf, "} %v;\n", st.CName(CPTC))
		}

		// assert the sizes as we have the numbers
		if st.C.IsBasic() {
			fmt.Fprintf(buf, `static_assert (sizeof(%v) == %v, "Size of %v is not correct");%v`, st.CName(""), st.C.Size.ByteSize, st.CName(""), "\n")
			fmt.Fprintf(buf, `static_assert (alignof(%v) == %v, "Align of %v is not correct");%v`, st.CName(""), st.C.Size.ByteAlignment, st.CName(""), "\n\n")
		} else if st.C.IsVector() {
			fmt.Fprintf(buf, `static_assert (sizeof(%v) == %v, "Size of %v is not correct");%v`, st.CName(""), st.C.Size.ByteSize, st.CName(""), "\n")
			fmt.Fprintf(buf, `static_assert (alignof(%v) == %v, "Align of %v is not correct");%v`, st.CName(""), st.C.Size.ByteAlignment, st.CName(""), "\n")
			fmt.Fprintf(buf, `static_assert (sizeof(%v) == %v, "Size of %v is not correct");%v`, st.CName("cpt_"), st.C.Size.ByteSize, st.CName("cpt_"), "\n")
			fmt.Fprintf(buf, `static_assert (alignof(%v) == %v, "Align of %v is not correct");%v`, st.CName("cpt_"), st.C.Size.ByteAlignment, st.CName("cpt_"), "\n\n")
		} else if st.C.IsStruct() {
			fmt.Fprintf(buf, `static_assert (sizeof(%v) == %v, "Size of %v is not correct");%v`, st.CName(CPTC), st.C.Size.ByteSize, st.CName(CPTC), "\n")
			fmt.Fprintf(buf, `static_assert (alignof(%v) == %v, "Align of %v is not correct");%v`, st.CName(CPTC), st.C.Size.ByteAlignment, st.CName(CPTC), "\n")
			fmt.Fprintf(buf, `static_assert (sizeof(%v) == %v, "Size of %v is not correct");%v`, st.CName("cpt_"), st.C.Size.ByteSize, st.CName("cpt_"), "\n")
			fmt.Fprintf(buf, `static_assert (alignof(%v) == %v, "Align of %v is not correct");%v`, st.CName("cpt_"), st.C.Size.ByteAlignment, st.CName("cpt_"), "\n\n")
		} else {
			fmt.Println(st)
			panic("cannot have an exported array type? what happened?")
		}
	}

	buf.WriteString("typedef struct {\n")
	w := tabwriter.NewWriter(buf, 0, 1, 1, ' ', 0)
	for _, arg := range inp.Arguments {
		ty := types.CreateArray(ts.Get(arg.Ty).C, arg.Arrno)
		if ty.IsComplexStruct() {
			fmt.Fprintf(w, "  "+ty.CString("", arg.Name, true)+";\n")
		} else {
			name := arg.Name
			if !(ty.IsArray() && ty.Array.Len == -1) {
				name = "(*" + name + ")"
			}
			fmt.Fprintf(w, "  "+ty.CString(CPTC, name, false)+";\n")
			fmt.Fprintf(w, "  int64_t "+arg.Name+"_len;\n\n")
		}
	}
	w.Flush()
	buf.WriteString("} cptc_data;\n")
	fmt.Fprintf(buf, `static_assert (sizeof(cptc_data) == sizeof(cpt_data), "Size of cptc_data != cpt_data");%v`, "\n\n")

}
