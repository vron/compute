// Simple tool to parse a vulkan compute shader and translate it as should be.
package main

import (
	"flag"
	"io/ioutil"
	"log"
	"path/filepath"
	"regexp"
	"strconv"
	"strings"
)

var header = []byte(`
/* BEGIN TRANSLATION HEADER */

// Definitions provided by the runtime
unsigned int __attribute__((overloadable)) __attribute__((const)) get_global_id_(const long _thread_context_, unsigned int a);
unsigned int __attribute__((overloadable)) __attribute__((const)) get_local_id_(const long _thread_context_, unsigned int a);


uint3 gl_GlobalInvocationID(const long _thread_context_) {
	return (uint3)(get_global_id_(_thread_context_, 0), get_global_id_(_thread_context_, 1), get_global_id_(_thread_context_, 2));
}

uint3 gl_LocalInvocationID(const long _thread_context_) {
	return (uint3)(get_local_id_(_thread_context_, 0), get_local_id_(_thread_context_, 1), get_local_id_(_thread_context_, 2));
}


typedef struct {
	__global float* data;
	uint width;
} image2D;

void __attribute__((overloadable)) imageStore(image2D img, int2 pos, float4 pixel) {
	// this one is for a float rgba image - fix the overloadeable ones...

	// TODO: Does this function imply some synchronization in the standard?

	int index = pos.x*4 + pos.y*(int)(img.width*4);
	img.data[index + 0] = pixel.x;
	img.data[index + 1] = pixel.y;
	img.data[index + 2] = pixel.z;
	img.data[index + 3] = pixel.w;
}

/* END TRANSLATION HEADER */
`)

var (
	fOutPath    string
	fBuildPath  string
	fHeaderPath string
)

func init() {
	flag.StringVar(&fOutPath, "outpath", "", "the folder where to write the output")
	flag.StringVar(&fBuildPath, "buildpath", "", "the folder where to write compilation files")
	flag.StringVar(&fHeaderPath, "headerpath", "", "the folder where to write headers files")
}

func main() {
	flag.Parse()
	inf := flag.Arg(0)
	buff, err := ioutil.ReadFile(inf)
	if err != nil {
		log.Fatalln(err)
	}
	info, buff, err := process(buff)
	if err != nil {
		log.Fatalln(err)
	}

	// get the name of the input file
	name := filepath.Base(inf)
	if strings.Contains(name, ".") {
		name = strings.Split(name, ".")[0]
	}
	info.KernelName = name

	err = ioutil.WriteFile(filepath.Join(fBuildPath, "kernel.cl"), buff, 0777)
	if err != nil {
		log.Fatalln(err)
	}

	buildGoPackage(fOutPath, info)
	buildCHeaders(fHeaderPath, info)
}

type Info struct {
	KernelName                         string
	LocalSizeX, LocalSizeY, LocalSizeZ int
	Args                               []Argument
	CallArgs                           []Argument
	Signature                          string
}

type Argument struct {
	Name string
	Type Type
}

type Type struct {
	GL string
	CL string
	C  string
	Go string
}

type Handler func(buff []byte, info *Info) ([]byte, error)

// Process it, extract the infor and dump it
func process(buff []byte) (info Info, b []byte, err error) {
	handlers := []Handler{
		hShaderVersion,
		hDispatch,
		hLayouts,
		hVecTypes,
		hConstructCast,
		hIDs,
		hMain,
	}
	for _, h := range handlers {
		buff, err = h(buff, &info)
		if err != nil {
			return info, nil, err
		}
	}
	return info, append(header, buff...), err
}

// remove all shader version tags, and ensure that they are the version we support
func hShaderVersion(buff []byte, info *Info) ([]byte, error) {
	r := regexp.MustCompile(`\#version\s*(\d+)`)
	finds := r.FindAllSubmatch(buff, -1)
	if len(finds) < 1 {
		log.Fatalln("could not find shader version, expected 430")
	}
	for fi := range finds {
		val := string(finds[fi][1])
		if val != "450" {
			log.Fatalln("expected shader version 450, found " + val)
		}
	}
	return r.ReplaceAll(buff, []byte{}), nil
}

// read the dispatch sizes and use
func hDispatch(buff []byte, info *Info) ([]byte, error) {
	r := regexp.MustCompile(`layout\(local_size_x\s*\=\s*(\d+),\s*local_size_y\s*\=\s*(\d+),\s*local_size_z\s*\=\s*(\d+)\s*\)\s*in;`)
	finds := r.FindAllSubmatch(buff, -1)
	if len(finds) != 1 {
		log.Fatalln("expected one dispatch sioze specifier, found", len(finds))
	}
	x, y, z := string(finds[0][1]), string(finds[0][2]), string(finds[0][3])

	info.LocalSizeX, _ = strconv.Atoi(string(x))
	info.LocalSizeY, _ = strconv.Atoi(string(y))
	info.LocalSizeZ, _ = strconv.Atoi(string(z))

	return r.ReplaceAll(buff, []byte{}), nil
}

// handle layout specifications.
func hLayouts(buff []byte, info *Info) ([]byte, error) {
	r := regexp.MustCompile(`(?s)(layout\(.*?);`)
	finds := r.FindAllSubmatch(buff, -1)
	for _, f := range finds {
		handleLayout(f[1], info)
	}
	// so remove all of them, we have handled them and will specify
	// them as part of the inputs instead.
	return r.ReplaceAll(buff, []byte{}), nil
}

// translate calls to get the incovation id
func hIDs(buff []byte, info *Info) ([]byte, error) {
	// replace them with function calls that we will implement
	// and expect to get inlined...
	r := regexp.MustCompile(`gl_(Global|Local)InvocationID`)
	return r.ReplaceAll(buff, []byte("gl_${1}InvocationID(_thread_context_)")), nil
}

// translate opengl vectors to opencl vectors
func hVecTypes(buff []byte, info *Info) ([]byte, error) {
	for _, prefix := range [][2]string{{"u", "uint"}, {"i", "int"}, {"", "float"}} {
		for _, size := range []string{"2", "3", "4"} {
			// replace literals
			r := regexp.MustCompile(prefix[0] + "vec" + size + "\\(")
			buff = r.ReplaceAll(buff, []byte("("+prefix[1]+size+")("))
			// replace vatiable names
			r = regexp.MustCompile(prefix[0] + "vec" + size)
			buff = r.ReplaceAll(buff, []byte(prefix[1]+size))
		}
	}
	return buff, nil
}

// translate constructors to casts
func hConstructCast(buff []byte, info *Info) ([]byte, error) {
	for _, prefix := range [][2]string{{"uint", "uint"}, {"int", "int"}, {"float", "float"}} {
		// replace constructors
		r := regexp.MustCompile(prefix[0] + "\\(")
		buff = r.ReplaceAll(buff, []byte("("+prefix[1]+")("))
	}
	return buff, nil
}

// translate the main function to a kernel function that we can expose.
func hMain(buff []byte, info *Info) ([]byte, error) {
	r := regexp.MustCompile(`void\s*main\(\)\s*{`)
	repl := "__kernel void kern(const long _thread_context_, "
	intro := ""

	// TODO: Re-work to move this out, and to use them to print instead

	for _, a := range info.Args {
		repl += a.Type.CL + " " + a.Name + "Data, "
		info.CallArgs = append(info.CallArgs, Argument{
			Name: a.Name + "Data",
			Type: convertTypes(a.Type.CL),
		})

		// if this is an image we also need to specify the widht, and first
		// thing in the function convert this to a struct we can later pass around
		if a.Type.GL == "image2D" {
			repl += "const uint " + a.Name + "Width, "
			intro += "\timage2D " + a.Name + ";\n"
			intro += "\t" + a.Name + ".data = " + a.Name + "Data;\n"
			intro += "\t" + a.Name + ".width = " + a.Name + "Width;\n"
			info.CallArgs = append(info.CallArgs, Argument{
				Name: a.Name + "Width",
				Type: convertTypes("uint"),
			})
		}
	}
	repl = repl[:len(repl)-2]
	repl += ")"
	info.Signature = repl
	repl += " {\n" + intro

	return r.ReplaceAll(buff, []byte(repl)), nil
}

func convertTypes(clt string) (t Type) {
	clt = strings.Replace(clt, "const", "", -1)
	clt = strings.Replace(clt, "__global", "", -1)
	clt = strings.TrimSpace(clt)
	t.C = clt

	// convert the c type to go type
	if strings.HasSuffix(clt, "*") {
		clt = clt[:len(clt)-1]
		clt = "*" + clt
	}
	clt = strings.Replace(clt, "float", "float32", -1)
	clt = strings.Replace(clt, "unsigned int", "uint", -1)
	t.Go = clt
	return
}
