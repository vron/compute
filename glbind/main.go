// tool to write out part of cpp files to build the shader
package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"

	"github.com/vron/compute/glbind/input"
	"github.com/vron/compute/glbind/types"
)

var (
	fIn  string
	fOut string
)

func init() {
	flag.StringVar(&fIn, "in", "./build/kernel.json", "input file")
	flag.StringVar(&fOut, "out", "./build/", "output folder")
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

var inp input.Input

func main() {
	flag.Parse()

	// read the input json we should use
	b, err := ioutil.ReadFile(fIn)
	expect(err)
	expect(json.Unmarshal(b, &inp))

	ts := types.New(inp)

	generateSharedH(inp, ts)
	//generateTypes(inp, ts)
	//generateComp(inp, ts)
	//generateAlignH(inp, ts)
	//generateGo(inp, ts)
}

func expect(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}
