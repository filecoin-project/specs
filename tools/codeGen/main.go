package main

import (
	"flag"
	"fmt"
	codeGen "github.com/filecoin-project/specs/codeGen/lib"
	util "github.com/filecoin-project/specs/codeGen/util"
	"os"
	"strings"
)

var Assert = util.Assert
var CheckErr = util.CheckErr

func replaceExt(filePath string, srcExt string, dstExt string) string {
	n := len(filePath)
	Assert(n >= len(srcExt))
	Assert(filePath[n-len(srcExt):] == srcExt)
	return filePath[:n-len(srcExt)] + dstExt
}

const USAGE = `SYNOPSIS
	%[1]s <command> src.id [out.go]

COMMANDS
	gen <idsrc> <goout>     parse contents of <idsrc>, compile, and output to <goout>
	fmt <idsrc> [<idsrc2>]  parse <idsrc>, and write formatted output to <idsrc2> (or <idsrc>)
	sym <idsrc>             parse contents of <idsrc>, and write symbol table to STDOUT

EXAMPLES
	# compile file.id to file.gen.go
	%[1]s gen a/b/file.id a/b/file.gen.go

	# format file.id
	%[1]s fmt a/b/file.id

	# format file.id to file2.id
	%[1]s fmt a/b/file.id a/b/file2.id

	# output symbol table of file.id
	%[1]s sym a/b/file.id
`

func main() {
	flag.Usage = func() {
		fmt.Printf(USAGE, os.Args[0])
		os.Exit(0)
	}

	flag.Parse()
	argsOrig := flag.Args()
	Assert(len(argsOrig) > 1)
	cmd := argsOrig[0]
	args := argsOrig[1:]

	var inputFilePath, outputFilePath string
	var inputFile, outputFile *os.File
	var err error

	// first argument
	if cmd == "gen" || cmd == "fmt" || cmd == "sym" {
		inputFilePath = args[0]
		inputFile, err = os.Open(inputFilePath)
		CheckErr(err)
	}

	// second argument
	if cmd == "gen" {
		Assert(len(args) == 2)
		outputFilePath = args[1]
	} else if cmd == "fmt" {
		outputFilePath = args[0] // replace file
		if len(args) == 2 {
			outputFilePath = args[1]
		}
	}
	// open files last
	// defer opening outputFile until the end (after parsing)
	// so that fmt can output to the input filename

	switch cmd {
	case "gen":
		inputFilePathTokens := strings.Split(inputFilePath, "/")
		Assert(len(inputFilePathTokens) >= 2)
		packageName := inputFilePathTokens[len(inputFilePathTokens)-2]
		goMod := codeGen.GenGoModFromFile(inputFile, packageName)
		outputFile, err = os.Create(outputFilePath)
		CheckErr(err)
		codeGen.WriteGoMod(goMod, outputFile)

	case "fmt":
		mod := codeGen.ParseDSLModuleFromFile(inputFile)
		outputFile, err = os.Create(outputFilePath)
		CheckErr(err)
		codeGen.WriteDSLModule(outputFile, mod)

	case "sym":
		Assert(len(args) >= 2)
		mod := codeGen.ParseDSLModuleFromFile(inputFile)
		decls := mod.Decls()
		declsMap := map[string]codeGen.Decl{}
		for _, decl := range decls {
			declsMap[decl.Name()] = decl
		}
		declsPrint := []codeGen.Entry{}
		for i, sym := range args[1:] {
			if i > 0 {
				declsPrint = append(declsPrint, codeGen.EntryEmpty())
			}
			decl, ok := declsMap[sym]
			if !ok {
				panic(fmt.Sprintf("Error: symbol not found: %v\n", sym))
			}
			declsPrint = append(declsPrint, codeGen.EntryDecl(decl))
		}
		codeGen.WriteDSLBlockEntries(os.Stdout, declsPrint, codeGen.WriteDSLContextInit())

	default:
		Assert(false)
	}
}
