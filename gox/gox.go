package main

import (
	"fmt"
	"go/scanner"
	"go/token"
	"io/ioutil"
	"os"
	"path/filepath"
)

func main() {

	// fmt.Println("gox running with:", os.Args, wd, err)
	files := findFiles()
	for _, f := range files {
		err := transform(f)
		if err != nil {
			fmt.Println("transform error:", err)
		}
	}
}

func transform(file string) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		return err
	}

	fset := token.NewFileSet() // positions are relative to fset

	// src is the input that we want to tokenize.
	// src := []byte("cos(x) + 1i*sin(x) // Euler")
	src := data

	// Initialize the scanner.
	var s scanner.Scanner
	fset2 := token.NewFileSet()                       // positions are relative to fset
	file2 := fset.AddFile("", fset2.Base(), len(src)) // register input "file"
	s.Init(file2, src, nil /* no error handler */, scanner.ScanComments)

	// Repeated calls to Scan yield the token sequence found in the input.
	var t1, t2, t3 token.Token
	for {
		pos, tok, lit := s.Scan()
		if tok == token.EOF {
			break
		}
		// if level > 0 {
		// 	fmt.Print(strings.Repeat("    ", level))
		// }
		if t1 == token.LSS {
			if tok == token.IDENT {
				fmt.Println(">>> Possible start token!")
			}
			if tok == token.NOT {
				t2 = tok
			}
		} else {
			t1 = tok
		}
		_ = t2
		_ = t3
		fmt.Printf("%s\t%s\t%q\n", fset.Position(pos), tok, lit)
	}

	return nil
}

func findFiles() []string {
	wd, err := os.Getwd()
	if err != nil {
		fmt.Println("could't get the working directory:", err)
		os.Exit(-1)
	}
	dir, err := os.Open(wd)
	if err != nil {
		fmt.Println("could't open the directory:", err)
		os.Exit(-1)
	}
	names, err := dir.Readdirnames(-1)
	if err != nil {
		fmt.Println("could't read the directory:", err)
		os.Exit(-1)
	}
	var files []string
	for _, n := range names {
		if filepath.Ext(n) == ".go" {
			files = append(files, filepath.Join(wd, n))
		}
	}
	return files
}
