package main

import (
	"fmt"
	"os"
)

func main() {
	args := os.Args[1:]
	if len(args) < 1 {
		fmt.Fprintln(os.Stderr, "Expected source file name")
		os.Exit(1)
	}

	fileName := args[0]
	file, err := os.Open(fileName)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot open the source file: %s\n", err.Error())
		os.Exit(1)
	}

	fileData := make([]byte, 1024*1024)
	count, err := file.Read(fileData)
	if err != nil {
		fmt.Fprintf(os.Stderr, "Cannot read the source file: %s\n", err.Error())
		os.Exit(1)
	}

	source := string(fileData[:count])
	tokenStream, err := Tokenize(source)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	ast, err := Parse(tokenStream)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	if err := ast.InferType(); err != nil {
		fmt.Fprintf(os.Stderr, "%s\n", err.Error())
		os.Exit(1)
	}

	Generate(ast)
}
