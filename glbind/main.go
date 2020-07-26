// tool to write out part of cpp files to build the shader
package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
)

var (
	fIn  string
	fOut string
)

type Input struct {
	Arguments []InputArgument
	Shared    []InputArgument
	Structs   []InputStruct
	Wg_size   [3]int
	Body      string
}

type InputArgument struct {
	Name  string
	Ty    string
	Arrno []int
}

type InputStruct struct {
	Name   string
	Fields []InputArgument
}

func init() {
	flag.StringVar(&fIn, "in", "./build/kernel.json", "input file")
	flag.StringVar(&fOut, "out", "./build/", "output folder")
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

var inp Input

func main() {
	flag.Parse()

	// read the input json we should use
	b, err := ioutil.ReadFile(fIn)
	expect(err)
	expect(json.Unmarshal(b, &inp))

	parseTypeInfo(inp)

	for _, s := range inp.Shared {
		if len(s.Arrno) == 0 {
			// This is since we need to access it shared - maybe we can replace it by using
			// c++ references instead of pointers and thusly achieve what we want?
			//panic("thus far we only support shared arrays")
		}
	}

	generateSharedH(inp)
	generateTypes(inp)
	generateComp(inp)
	generateAlignH(inp)
	generateGo(inp)
}

func expect(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
