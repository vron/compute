package main

import (
	"bytes"
	"log"
	"regexp"
	"strings"
)

type qualifier struct {
	name  string
	value string
}

type layoutSpec struct {
	// the ones inside the parenthesis
	qualifiers []qualifier
	// definitions
	definitions []string
	// final name, e.g. the name of the variables
	name string
}

// this file implements handling of layout specifications.
func handleLayout(buff []byte, info *Info) {
	spec := parseLayoutSpec(buff)
	arg := argFromLayout(spec)
	info.Args = append(info.Args, arg)
}

func argFromLayout(s layoutSpec) (a Argument) {
	a.Name = s.name

	// find the type:
	if has(s.definitions, "image2D") {
		if hasq(s.qualifiers, "rgba32f") {
			a.Type.GL = "image2D"
			a.Type.CL = "__global float*"
			return
		}
	}
	log.Fatalln("unknown layout spec")
	return
}

func has(a []string, v string) bool {
	for _, vv := range a {
		if vv == v {
			return true
		}
	}
	return false
}

func hasq(a []qualifier, v string) bool {
	for _, vv := range a {
		if vv.name == v {
			return true
		}
	}
	return false
}

func parseLayoutSpec(buf []byte) (l layoutSpec) {
	// hacky as is, but split on parenthesis, then regexp out what we want
	l.qualifiers = parseQualifiers(buf)
	l.definitions = parseDefinitions(buf)
	l.name = l.definitions[len(l.definitions)-1]
	l.definitions = l.definitions[:len(l.definitions)-1]
	return
}

func parseQualifiers(buf []byte) (qs []qualifier) {
	re := regexp.MustCompile(`\(.*\)`)
	buf = re.Find(buf)
	buf = buf[1 : len(buf)-1]
	args := bytes.Split(buf, []byte(","))
	for _, arg := range args {
		if bytes.Contains(arg, []byte("=")) {
			parts := bytes.Split(arg, []byte("="))
			qs = append(qs, qualifier{strings.TrimSpace(string(parts[0])), strings.TrimSpace(string(parts[1]))})

		} else {
			qs = append(qs, qualifier{strings.TrimSpace(string(arg)), ""})
		}
	}
	return
}

func parseDefinitions(buf []byte) (ds []string) {
	buf = bytes.Split(buf, []byte(")"))[1]
	args := bytes.Split(buf, []byte(" "))
	for _, arg := range args {
		ds = append(ds, strings.TrimSpace(string(arg)))
	}
	return
}
