package main

import (
	"log"
	"os"
)

func main() {
	err := os.RemoveAll("./build/")
	if err != nil {
		log.Fatalln(err)
	}
}
