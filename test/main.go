// Command test is used to run all the end-to-end tests for project compute.
package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"os"
	"os/exec"
	"path/filepath"
	"regexp"
	"strings"
	"sync"
	"time"
)

var (
	fLocal    bool
	fBench    bool
	once      sync.Once
	imageName string
)

func init() {
	flag.BoolVar(&fLocal, "local", true, "run on local machine instead of container")
	flag.BoolVar(&fBench, "bench", false, "run benchmars instead of tests")
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
	benchFiles, err := filepath.Glob("./test/benchmarks/*.go")
	ensure(err)
	if !fBench {
		for _, testFile := range testFiles {
			if !re.MatchString(testFile) {
				continue
			}
			runTest(testFile)
		}
	}
	for _, testFile := range benchFiles {
		if !re.MatchString(testFile) {
			continue
		}
		runTest(testFile)
	}

	fmt.Printf("\ntests ok\n")
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
	if !fBench {
		fmt.Printf("%30s: ", filepath.Base(p))
	}

	cleanBuild()
	getShader(p)
	copyTest(p)

	path, _ := os.Getwd()
	if !fLocal {
		if fBench {
			panic("benchmars inside DOC not yet suppored")
		}
		// run inside docer to test
		c := exec.Command("docker", "run", "-v", filepath.Join(path, "/build")+":/build", "--rm", imageName)
		data, err := c.CombinedOutput()
		if err == nil {
			fmt.Printf("ok\n")
			return
		}
		fmt.Println("\n")

		fmt.Println(string(data))
		log.Fatalln("test failed")
	} else {
		// run local on this machine
		c := exec.Command("./script/build.sh", "./build/test.comp")
		data, err := c.CombinedOutput()
		if err == nil {
			if !fBench {
				fmt.Printf("ok\n")
				return
			}

			// run the benchmar and log the output
			c := exec.Command("go", "test", "-run", "xxxxxx", "-benchtime", "2s", "-bench", ".")
			c.Dir = filepath.Join(path, "build", "go")
			c.Stderr = os.Stderr
			c.Stdout = resultsOnly(os.Stdout)
			err := c.Run()
			if err != nil {
				log.Fatalln("bench failed", err)
			}
			return
		}
		fmt.Println("\n")

		fmt.Println(string(data))
		log.Fatalln("test failed")
	}
}

// return a filtering writer that will only write lines that
// are benchmar results to the output such that we get a nice
// printing of it all..
func resultsOnly(w io.Writer) io.Writer {
	pr, pw := io.Pipe()

	go func() {
		scanner := bufio.NewScanner(pr)
		for scanner.Scan() {
			line := scanner.Text()
			if strings.HasPrefix(line, "Benchmar") {
				io.WriteString(w, line+"\n")
			}
		}
		if err := scanner.Err(); err != nil {
			log.Println("reading input:", err)
		}
	}()
	return pw
}

func copy(dst, src string) {
	data, err := ioutil.ReadFile(src)
	ensure(err)
	err = ioutil.WriteFile(dst, data, 0644)
	ensure(err)
}
