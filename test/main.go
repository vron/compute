// Command test is used to run all the end-to-end tests for project compute.
package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"sync"
	"time"
)

var (
	fLocal    bool
	once      sync.Once
	imageName string
)

func init() {
	flag.BoolVar(&fLocal, "local", false, "run on local machine instead of container")
	log.SetFlags(log.LstdFlags | log.Llongfile)
}

func main() {
	// for each file found in the folder, extract the source, dump it, build
	// everything, link it, build go package and test.. Happy days..
	// For now we run the tests in the container such that we are on linux
	// for sure.
	flag.Parse()
	filter := ".*"
	if flag.NArg() > 0 {
		filter = flag.Arg(0)
	}
	re := regexp.MustCompile(filter)

	testFiles, err := filepath.Glob("./test/tests/*.go")
	ensure(err)
	for _, testFile := range testFiles {
		if !re.MatchString(testFile) {
			continue
		}
		runTest(testFile)
	}

	fmt.Println("\nPASS")
}

func ensure(err error) {
	if err != nil {
		log.Fatalln(err)
	}
}

func cleanBuild() {
	c := exec.Command("./script/clean.sh")
	c.Stdout = os.Stdout
	c.Stderr = os.Stderr
	ensure(c.Run())
	ensure(os.MkdirAll("./build", 0777))
}

func buildImage() {
	if fLocal {
		return
	}
	ts := time.Now()
	fmt.Print("~ building image: ")
	c := exec.Command("docker", "build", ".")
	data, err := c.CombinedOutput()
	if err != nil {
		fmt.Println(string(data))
		log.Fatalln("image build failed")
	}
	imageName = string(data[len(data)-13 : len(data)-3])
	fmt.Println("finished in", time.Now().Sub(ts))
}

func getShader(p string) {
	buf, err := ioutil.ReadFile(p)
	ensure(err)

	re := regexp.MustCompile(`(?s)shader\s*\=\s*` + "`" + `(.*)` + "`")
	m := re.FindSubmatch(buf)
	if len(m) != 2 {
		log.Fatalln("found no shader source...")
	}
	shader := m[1]
	ensure(ioutil.WriteFile("./build/test.comp", shader, 0777))
}

func copyTest(p string) {
	copy("./build/test_test.go", p)
	copy("./build/util_test.go", "test/util_test.go")
}

func runTest(p string) {
	once.Do(buildImage)
	ts := time.Now()
	fmt.Printf("~ test %v: ", p)

	cleanBuild()
	getShader(p)
	copyTest(p)

	if !fLocal {
		// run inside docer to test
		path, _ := os.Getwd()
		c := exec.Command("docker", "run", "-v", filepath.Join(path, "/build")+":/build", "--rm", imageName)
		data, err := c.CombinedOutput()
		if err == nil {
			fmt.Printf("PASS in %v\n", time.Now().Sub(ts))
			return
		}
		fmt.Println("\n")

		fmt.Println(string(data))
		log.Fatalln("test failed")
	} else {
		// run local on this machine (faster?)c
		c := exec.Command("./script/build.sh", "./build/test.comp")
		data, err := c.CombinedOutput()
		if err == nil {
			fmt.Printf("PASS in %v\n", time.Now().Sub(ts))
			return
		}
		fmt.Println("\n")

		fmt.Println(string(data))
		log.Fatalln("test failed")
	}
}

func copy(dst, src string) {
	data, err := ioutil.ReadFile(src)
	ensure(err)
	err = ioutil.WriteFile(dst, data, 0644)
	ensure(err)
}
