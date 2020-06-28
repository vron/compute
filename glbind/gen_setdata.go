package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
)

func generateSetData(inp Input) {
	f, err := os.Create(filepath.Join(fOut, "setdata.hpp"))
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	buf := bufio.NewWriter(f)
	defer buf.Flush()

	// we have an variable d of the data struct type that we need
	// to translate to the member variables.

	for _, a := range inp.Arguments {
		cf := CField{
			Name: a.Name,
			Ty:   maybeCreateArrayType(a.Ty, a.Arrno),
		}
		recChecAlignment(buf, inp, cf, "d.")
		cf.CxxBinding(buf)
	}

	buf.WriteString("\treturn 0;")
}
