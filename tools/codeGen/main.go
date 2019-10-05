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

func main() {
	flag.Parse()
	argsOrig := flag.Args()
	Assert(len(argsOrig) > 1)
	cmd := argsOrig[0]
	args := argsOrig[1:]

	var inputFilePath, outputFilePath string
	var inputFile, outputFile *os.File
	var err error

	if cmd == "gen" || cmd == "fmt" || cmd == "sym" {
		inputFilePath = args[0]
		inputFile, err = os.Open(inputFilePath)
		CheckErr(err)

		if cmd == "gen" || cmd == "fmt" {
			Assert(len(args) == 2)
			outputFilePath = args[1]
			outputFile, err = os.Create(outputFilePath)
			CheckErr(err)
		}
	}

	switch cmd {
	case "gen":
		inputFilePathTokens := strings.Split(inputFilePath, "/")
		Assert(len(inputFilePathTokens) >= 2)
		packageName := inputFilePathTokens[len(inputFilePathTokens)-2]
		goMod := codeGen.GenGoModFromFile(inputFile, packageName)
		codeGen.WriteGoMod(goMod, outputFile)

	case "fmt":
		mod := codeGen.ParseDSLModuleFromFile(inputFile)
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
