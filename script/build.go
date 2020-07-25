package main

import (
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/termie/go-shutil"
)

func ensure(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func main() {
	SHADER := os.Args[1]
	CLANGPP := "clang++"

	cpfiles, _ := filepath.Glob("runtime/*")
	for _, f := range cpfiles {
		fi, _ := os.Stat(f)
		ft := f[len("runtime/"):]
		if fi.IsDir() {
			ensure(shutil.CopyTree(f, filepath.Join("build", filepath.Base(ft)), nil))
		} else {
			_, e := shutil.Copy(f, filepath.Join("build", ft), false)
			ensure(e)
		}
	}

	ensure(os.MkdirAll("build/generated", 0777))

	// first use lcpp to concatenate all into one file by processing the
	// includes such that we have single file we can build.
	ts := filepath.Base(SHADER)
	ts = filepath.Join("build", ts+".inc.comp")
	ensure(run("lua", "script/lcpp.lua", SHADER, "-I.", "-o", ts))
	SHADER = ts

	ensure(run("glslangValidator", SHADER))

	ensure(runf("gl2c", "cargo", "run", "-q", "--", "../"+SHADER, "../build/kernel.json"))

	files, _ := filepath.Glob("glbind/*.go")
	ensure(run("go", append([]string{"run"}, files...)...))

	ensure(run("goimports", "-w", "build/kernel.go"))

	// TODO: here we want to build for multiple platforms...
	cargs := []string{
		"-std=c++2a",
		//"-fvisibility=hidden",

		"-Wall",
		"-Wextra",
		"-Werror",
		"-Wno-unused-function",

		"-Ofast",
		"-ffast-math",
		"-fno-math-errno",
	}
	asm := ""
	outf := ""
	if runtime.GOOS == "windows" {
		asm = "co/arch/amd64_win.S"
		outf = "shader.dll"
	} else {
		asm = "co/arch/amd64_nix.S"
		outf = "./shader.so"
		cargs = append(cargs, "-fPIC")
	}
	ensure(runf("build", CLANGPP, append([]string{"lib.cpp", asm, "-shared", "-o", outf}, cargs...)...))

	ensure(os.MkdirAll("build/go/build", 0777))

	if ex("build/test_test.go") {
		cp("build/test_test.go", "build/go/test_test.go")
	}
	if ex("build/util_test.go") {
		cp("build/util_test.go", "build/go/util_test.go")
	}
	cp("build/kernel.go", "build/go/kernel.go")

	if runtime.GOOS == "windows" {
		cp("build/shader.dll", "build/go/shader.dll")
	} else {
		cp("build/shader.so", "build/go/shader.so")
		cp("build/shader.so", "build/go/build/shader.so")
	}
	cp("build/generated/shared.h", "build/go/shared.h")
	ensure(runf("build/go", "go", "test", "-v", "."))
}

func run(m string, args ...string) error {
	return runf("", m, args...)
}

func runf(path, m string, args ...string) error {
	c := exec.Command(m, args...)
	c.Dir = path
	c.Stderr = os.Stderr
	c.Stdout = os.Stdout
	return c.Run()
}

func ex(path string) bool {
	_, err := os.Stat(path)
	if err == nil {
		return true
	}
	if os.IsNotExist(err) {
		return false
	}
	ensure(err)
	panic("cannot be reached")
}

func cp(from, to string) {
	_, e := shutil.Copy(from, to, false)
	ensure(e)
}

/*
cp build/shader.so build/go/build
cp build/shader.so build/go
cp build/generated/shared.h build/go/
# build the go pacae and run it to test it wors
(cd build/go && go test -v . )

*/
