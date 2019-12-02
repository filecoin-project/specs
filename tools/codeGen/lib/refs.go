package codeGen

import (
	"go/ast"
)

func RefParseFmtInfo(x ParseFmtInfo) *ParseFmtInfo {
	return &x
}

func RefWriteDSLAlignment(x WriteDSLAlignment) *WriteDSLAlignment {
	return &x
}

func RefEntry(x Entry) *Entry {
	return &x
}

func RefAlgType(x AlgType) *AlgType {
	return &x
}

func RefField(x Field) *Field {
	return &x
}

func RefFunType(x FunType) *FunType {
	return &x
}

func RefArrayType(x ArrayType) *ArrayType {
	return &x
}

func RefOptionType(x OptionType) *OptionType {
	return &x
}

func RefMapType(x MapType) *MapType {
	return &x
}

func RefNamedType(x NamedType) *NamedType {
	return &x
}

func RefDSLRefType(x RefType) *RefType {
	return &x
}

func RefDSLType(x Type) *Type {
	return &x
}

func RefPackageDecl(x PackageDecl) *PackageDecl {
	return &x
}

func RefImportDecl(x ImportDecl) *ImportDecl {
	return &x
}

func RefTypeDecl(x TypeDecl) *TypeDecl {
	return &x
}

func RefGoIdent(x GoIdent) *GoIdent {
	return &x
}

func RefParseStream(x ParseStream) *ParseStream {
	return &x
}

func RefParseError(x ParseError_S) *ParseError_S {
	return &x
}

func RefGolangASTFile(x ast.File) *ast.File {
	return &x
}

func RefGolangASTImportSpec(x ast.ImportSpec) *ast.ImportSpec {
	return &x
}

func RefGolangASTField(x ast.Field) *ast.Field {
	return &x
}

func RefGolangASTKeyValueExpr(x ast.KeyValueExpr) *ast.KeyValueExpr {
	return &x
}

func RefGolangASTBasicLit(x ast.BasicLit) *ast.BasicLit {
	return &x
}

func GoNode_Ref(x GoNode) *GoNode {
	return &x
}
