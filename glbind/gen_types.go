package main

import (
	"bufio"
	"fmt"
	"io"
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
			fmt.Fprintf(buf, "struct %v {  // size = %v, align = %v\n", st.CName(CPTC), st.C.Size.ByteSize, st.C.Size.ByteAlignment)
			for _, f := range st.C.Struct.Fields {
				fmt.Fprintf(w, "  "+f.CType.CString(CPTC, f.Name, false)+";\t// offset =\t%v\t\n", f.ByteOffset)
			}
			w.Flush()

			writeTypeConstructors(buf, st)

			fmt.Fprintf(buf, "};\n")
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

func hasField(fn string, ct *types.CType) bool {
	if !ct.IsStruct() {
		return false
	}
	for _, f := range ct.Struct.Fields {
		if f.Name == fn {
			return true
		}
	}
	return false
}

func writeTypeConstructors(buf io.Writer, st *types.GlslType) {
	// write function style constructors since they exist in glsl. Especially arrays are realy
	// anoying since we have to provide them value by value?
	fmt.Fprintf(buf, "%v() = default ;\n", st.CName(CPTC))
	fmt.Fprintf(buf, "  %v(", st.CName(CPTC))
	for i, f := range st.C.Struct.Fields {
		// this is a bit subtle - we need to ensure we do not have a field of the
		// same name as the struct as in that case we have to address it by struct begore
		pre := CPTC
		if hasField(f.Name, st.C) {
			pre = "struct " + CPTC
		}
		fmt.Fprintf(buf, "%v", f.CType.CString(pre, f.Name, false))
		if i != len(st.C.Struct.Fields)-1 {
			fmt.Fprint(buf, ", ")
		}
	}
	fmt.Fprintf(buf, ") : ")
	for i, f := range st.C.Struct.Fields {
		if f.CType.IsArray() {
			// we need to write it all out? really? ( do we actually need to do this recursively?)
			fmt.Fprintf(buf, "%v", f.Name)
			recWriteConst(buf, f.CType, fmt.Sprintf("%v", f.Name))
		} else {
			fmt.Fprintf(buf, "%v(%v)", f.Name, f.Name)
		}
		if i != len(st.C.Struct.Fields)-1 {
			fmt.Fprint(buf, ", ")
		}
	}
	fmt.Fprintf(buf, " {};\n")
}

func recWriteConst(buf io.Writer, ct *types.CType, head string) {
	if ct.IsArray() {
		fmt.Fprintf(buf, "{")
		for j := 0; j < ct.Array.Len; j++ {
			recWriteConst(buf, ct.Array.CType, head+fmt.Sprintf("[%v]", j))
			if j != ct.Array.Len-1 {
				fmt.Fprint(buf, ", ")
			}
		}
		fmt.Fprintf(buf, "}")
	} else {
		fmt.Fprintf(buf, head)
	}
}
