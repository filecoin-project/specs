package main

import (
	"flag"
	"main/codeGen"
	"os"
)

func replaceExt(filePath string, srcExt string, dstExt string) string {
	n := len(filePath)
	codeGen.Assert(n >= len(srcExt))
	codeGen.Assert(filePath[n-len(srcExt):] == srcExt)
	return filePath[:n-len(srcExt)] + dstExt
}

func main() {
	flag.Parse()
	args := flag.Args()
	codeGen.Assert(len(args) == 1)

	inputFilePath := args[0]
	outputFilePath := replaceExt(inputFilePath, ".id", ".gen.go")

	// fmt.Printf(" ===== Parsing: %v\n\n", filePath)
	inputFile, err := os.Open(inputFilePath)
	codeGen.CheckErr(err)
	goMod := codeGen.GenGoModFromFile(inputFile)

	outputFile, err := os.Create(outputFilePath)
	codeGen.CheckErr(err)
	codeGen.WriteGoMod(goMod, outputFile)
}
