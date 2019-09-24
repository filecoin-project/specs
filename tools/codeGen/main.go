package main

import (
	"flag"
	"strings"
	"os"
	codeGen "github.com/filecoin-project/specs/codeGen/main/codeGen"
	util "github.com/filecoin-project/specs/codeGen/main/util"
)

func replaceExt(filePath string, srcExt string, dstExt string) string {
	n := len(filePath)
	util.Assert(n >= len(srcExt))
	util.Assert(filePath[n-len(srcExt):] == srcExt)
	return filePath[:n-len(srcExt)] + dstExt
}

func main() {
	flag.Parse()
	args := flag.Args()
	util.Assert(len(args) == 2)

	inputFilePath := args[0]
	outputFilePath := args[1]

	inputFilePathTokens := strings.Split(inputFilePath, "/")
	util.Assert(len(inputFilePathTokens) >= 2)
	packageName := inputFilePathTokens[len(inputFilePathTokens)-2]

	// fmt.Printf(" ===== Parsing: %v\n\n", filePath)
	inputFile, err := os.Open(inputFilePath)
	util.CheckErr(err)
	goMod := codeGen.GenGoModFromFile(inputFile, packageName)

	outputFile, err := os.Create(outputFilePath)
	util.CheckErr(err)
	codeGen.WriteGoMod(goMod, outputFile)
}
