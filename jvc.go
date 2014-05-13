// Copyright 2014 Fredy Wijaya
//
// Permission is hereby granted, free of charge, to any person obtaining
// a copy of this software and associated documentation files (the
// "Software"), to deal in the Software without restriction, including
// without limitation the rights to use, copy, modify, merge, publish,
// distribute, sublicense, and/or sell copies of the Software, and to
// permit persons to whom the Software is furnished to do so, subject to
// the following conditions:
//
// The above copyright notice and this permission notice shall be
// included in all copies or substantial portions of the Software.
//
// THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND,
// EXPRESS OR IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF
// MERCHANTABILITY, FITNESS FOR A PARTICULAR PURPOSE AND
// NONINFRINGEMENT. IN NO EVENT SHALL THE AUTHORS OR COPYRIGHT HOLDERS BE
// LIABLE FOR ANY CLAIM, DAMAGES OR OTHER LIABILITY, WHETHER IN AN ACTION
// OF CONTRACT, TORT OR OTHERWISE, ARISING FROM, OUT OF OR IN CONNECTION
// WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE SOFTWARE.

package main

import (
	"archive/zip"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"strings"
)

type JavaClass struct {
	Magic uint32
	Minor uint16
	Major uint16
}

func ReadJarFile(jarFile string) (string, uint16, uint16, error) {
	r, e := zip.OpenReader(jarFile)
	if e != nil {
		return "", 0, 0, errors.New("Unable to read " + jarFile)
	}
	defer r.Close()

	for _, f := range r.File {
		if strings.HasSuffix(f.Name, ".class") {
			rc, e := f.Open()
			if e != nil {
				continue
			}
			defer rc.Close()
			return read(rc, f.Name)
		}
	}
	return "", 0, 0, errors.New("Unable to find a class file in " + jarFile)
}

func ReadClassFile(classFile string) (string, uint16, uint16, error) {
	r, e := os.Open(classFile)
	if e != nil {
		return "", 0, 0, errors.New("Unable to read " + classFile)
	}
	defer r.Close()

	return read(r, classFile)
}

func read(r io.Reader, classFile string) (string, uint16, uint16, error) {
	var jc JavaClass
	binary.Read(r, binary.BigEndian, &jc)
	if jc.Magic != 0xcafebabe {
		return "", 0, 0, errors.New("Invalid class file: " + classFile)
	}
	return determineJavaVersion(jc.Major, jc.Minor), jc.Major, jc.Minor, nil
}

func determineJavaVersion(major, minor uint16) string {
	if major == 45 && minor == 3 {
		return "1.1"
	} else if major == 46 && minor == 0 {
		return "1.2"
	} else if major == 47 && minor == 0 {
		return "1.3"
	} else if major == 48 && minor == 0 {
		return "1.4"
	} else if major == 49 && minor == 0 {
		return "1.5"
	} else if major == 50 && minor == 0 {
		return "1.6"
	} else if major == 51 && minor == 0 {
		return "1.7"
	} else {
		return "1.8"
	}
}

func validateArgs() {
	if len(os.Args) != 2 {
		printUsage()
		os.Exit(1)
	}
}

func printUsage() {
	fmt.Println("Usage:", os.Args[0], "<class_file|jar_file>")
}

func main() {
	validateArgs()

	var ver string
	var major uint16
	var minor uint16
	var e error
	if strings.HasSuffix(os.Args[1], ".class") {
		ver, major, minor, e = ReadClassFile(os.Args[1])
		if e != nil {
			fmt.Println("Error:", e)
			os.Exit(1)
		}
	} else if strings.HasSuffix(os.Args[1], ".jar") {
		ver, major, minor, e = ReadJarFile(os.Args[1])
		if e != nil {
			fmt.Println("Error:", e)
			os.Exit(1)
		}
	} else {
		fmt.Println("Error: Only .class or .jar file is supported")
		os.Exit(1)
	}
	fmt.Printf("Class version : %d.%d\n", major, minor)
	fmt.Printf("Java version  : %s\n", ver)
}
