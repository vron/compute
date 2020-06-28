// Command test is used to run all the end-to-end tests for project compute.
package main

import (
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
)

func main() {
	// for each file found in the folder, extract the source, dump it, build
	// everything, link it, build go package and test.. Happy days..
	// For now we run the tests in the container such that we are on linux
	// for sure.

	testFiles, err := filepath.Glob("./test/tests/*.go")
	ensure(err)
	for _, testFile := range testFiles {
		runTest(testFile)
	}

	// clean up after outselves
	c := exec.Command("make", "wipe")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	ensure(c.Run())
}

func ensure(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

// TODO: Do not log all the output unless there is an error... (or -v is specified...)

func runTest(p string) {
	fmt.Println("Running test: " + p)
	defer fmt.Println("\n\n ")

	buf, err := ioutil.ReadFile(p)
	ensure(err)

	shader := extractShader(buf)

	// Clean the output directory
	c := exec.Command("make", "wipe")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	ensure(c.Run())

	// Write the shader we want to run
	ensure(ioutil.WriteFile("./data/test.comp", shader, 0777))

	// Copy the test files we want to use
	copy("./data/output/test_test.go", p)

	// issue the docer command to actually build and test it
	c = exec.Command("docker", "build", "-q", ".")
	data, err := c.CombinedOutput()
	ensure(err)
	path, _ := os.Getwd()
	c = exec.Command("docker", "run", "-v", filepath.Join(path, "/data")+":/data", "--rm", string(data[7:17]))
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	ensure(c.Run())
}

func extractShader(b []byte) []byte {
	re := regexp.MustCompile(`(?s)shader\s*\=\s*` + "`" + `(.*)` + "`")
	m := re.FindSubmatch(b)
	if len(m) != 2 {
		log.Fatalln("found no shader source...")
	}
	return m[1]
}

func copy(dst, src string) {
	data, err := ioutil.ReadFile(src)
	ensure(err)
	err = ioutil.WriteFile(dst, data, 0644)
	ensure(err)
}
