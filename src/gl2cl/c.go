package main

import (
	"bufio"
	"log"
	"os"
	"path/filepath"
	"strconv"
)

// this file generates the c header files needed for the binding

func buildCHeaders(folder string, info Info) {
	buildShared(folder, info)
	buildkernelHeader(folder, info)
}

func buildShared(folder string, info Info) {
	f, err := os.Create(filepath.Join(folder, "shared.h"))
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	buf := bufio.NewWriter(f)
	defer buf.Flush()

	buf.WriteString(`// Code generated DO NOT EDIT

#define uint unsigned int

typedef struct {
`)

	for _, arg := range info.CallArgs {
		// TOOD: Add comments detailing from where this argument came in the original kernel.
		buf.WriteString("\t" + arg.Type.C + " " + arg.Name + ";\n")
	}
	buf.WriteString(`} kernel_data;

void* cpt_new_kernel(int);
int cpt_dispatch_kernel(void*, kernel_data, int, int, int);
void cpt_free_kernel(void*);
`)
}

func buildkernelHeader(folder string, info Info) {
	f, err := os.Create(filepath.Join(folder, "kernel.h"))
	if err != nil {
		log.Fatalln(err)
	}
	defer f.Close()
	buf := bufio.NewWriter(f)
	defer buf.Flush()

	buf.WriteString(`// Code generated DO NOT EDIT
#include "shared.h"
#define LOCAL_SIZE_X ` + strconv.Itoa(info.LocalSizeX) + `
#define LOCAL_SIZE_Y ` + strconv.Itoa(info.LocalSizeY) + `
#define LOCAL_SIZE_Z ` + strconv.Itoa(info.LocalSizeZ) + `

#define kernel_call 4

void kern(long, `)

	for i, arg := range info.CallArgs {
		// TOOD: Add comments detailing from where this argument came in the original kernel.
		buf.WriteString(arg.Type.C + " " + arg.Name)
		if i != len(info.CallArgs)-1 {
			buf.WriteString(", ")
		}
	}
	buf.WriteString(`);
`)
}
