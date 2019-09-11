package codeGen

import (
	"crypto/sha256"
	"encoding/binary"
	"go/ast"
	"go/printer"
	"go/token"
	"os"
)

type Type interface {}

type TypeDecl struct {
	name  string
	type_ Type
}

type Type_Case = Word
const (
	Type_Case_ProdType  Type_Case = 0
	Type_Case_ArrayType Type_Case = 1
	Type_Case_RefType   Type_Case = 2
	Type_Case_FunType   Type_Case = 3
	Type_Case_NamedType Type_Case = 4
)

type Field struct {
	fieldName  string
	fieldType  Type
	internal    bool
}

type ProdType struct {
	fields []Field
}

type ArrayType struct {
	elementType Type
}

type RefType struct {
	targetType Type
}

type FunType struct {
	argTypes []Type
	retType  Type
}

type NamedType struct {
	name string
}

type TypeHash string

func HashAccWord(buf *[]byte, w Word) {
	bufAcc := make([]byte, 8)
	binary.LittleEndian.PutUint64(bufAcc, uint64(w))
	*buf = append(*buf, bufAcc...)
}

func HashAccBool(buf *[]byte, b bool) {
	var c byte = 0
	if b {
		c = 1
	}
	*buf = append(*buf, c)
}

func HashAccString(buf *[]byte, s string) {
	HashAccWord(buf, Word(len(s)))
	*buf = append(*buf, s...)
}

func HashAccType(buf *[]byte, x Type) {
	switch x.(type) {
	case ProdType:
		xr := x.(ProdType)
		HashAccWord(buf, Type_Case_ProdType)
		HashAccWord(buf, Word(len(xr.fields)))
		for _, field := range xr.fields {
			HashAccString(buf, field.fieldName)
			HashAccType(buf, field.fieldType)
			HashAccBool(buf, field.internal)
		}

	case NamedType:
		xr := x.(NamedType)
		HashAccWord(buf, Type_Case_NamedType)
		HashAccString(buf, xr.name)

	default:
		panic("TODO")
	}
}

func HashSHA256(x []byte) []byte {
	h := sha256.New()
	h.Write(x)
	return h.Sum(nil)
}

func HashType(x Type) TypeHash {
	buf := &[]byte{}
	HashAccType(buf, x)
	// fmt.Printf("?? %v", *buf)
	return TypeHash(HashSHA256(*buf))
}

type GoMod struct {
	astFileSet *token.FileSet
	astFile    *ast.File
}

func GenGoMod(goDecls []GoNode) GoMod {
	var ret GoMod

	var astDecls = []ast.Decl{}
	for _, goDecl := range goDecls {
		astDecls = append(astDecls, GenAST(goDecl).(ast.Decl))
	}

	ret.astFileSet = token.NewFileSet()
	ret.astFile = &ast.File {
		Name: ast.NewIdent("fileName"),
		Decls: astDecls,
	}

	return ret
}

func GenGoTypeDecls(decls []TypeDecl) []GoNode {
	typeMap := map[string]Type{}
	retDecls := &[]GoNode{}
	retMap := map[TypeHash]GoNode{}

	for _, decl := range decls {
		GenGoTypeDeclAcc(decl.name, decl.type_, typeMap, retDecls, retMap)
	}

	return *retDecls
}

func ExtractFieldNames(fields []Field) []string {
	ret := []string{}
	for _, field := range fields {
		ret = append(ret, field.fieldName)
	}
	return ret
}

func IdToImpl(name string) string {
	return name + "_I"
}

func GenGoTypeDeclAcc(name string, x Type, typeMap map[string]Type, retDecls *[]GoNode, retMap map[TypeHash]GoNode) GoNode {
	switch x.(type) {
	case ProdType:
		xr := x.(ProdType)

		interfaceName := name
		interfaceID := GoNamedType { name: interfaceName, }

		fieldNames := ExtractFieldNames(xr.fields)
		interfaceFieldTypes := []GoNode{}
		implFieldTypes := []GoNode{}

		for _, field := range xr.fields {
			implFieldType := GenGoTypeAcc(field.fieldType, typeMap, retDecls, retMap)
			interfaceFieldType := GoFunType {
				argTypes: []GoNode{},
				retType:  implFieldType,
			}

			implFieldTypes = append(implFieldTypes, implFieldType)
			interfaceFieldTypes = append(interfaceFieldTypes, interfaceFieldType)
		}

		interfaceDecl := GoTypeDecl {
			name: interfaceName,
			type_: GoProdType {
				typeCase:   GoProdTypeCase_Interface,
				fieldNames: fieldNames,
				fieldTypes: interfaceFieldTypes,
			},
		}

		implName := IdToImpl(name)

		implDecl := GoTypeDecl {
			name:  implName,
			type_: GoProdType {
				typeCase:   GoProdTypeCase_Struct,
				fieldNames: fieldNames,
				fieldTypes: implFieldTypes,
			},
		}

		*retDecls = append(*retDecls, interfaceDecl)
		*retDecls = append(*retDecls, implDecl)

		// TODO: assert not in retMap
		retMap[HashType(x)] = interfaceID
		return interfaceID

	default:
		panic("TODO")
	}
}

func GenGoTypeAcc(x Type, typeMap map[string]Type, retDecls *[]GoNode, retMap map[TypeHash]GoNode) GoNode {
	// ret := []GoNode{}
	// idMap := map[Type]string

	if match, ok := retMap[HashType(x)]; ok {
		return match
	}

	// for _, x := range decls {
	switch x.(type) {
	case NamedType:
		xr := x.(NamedType)
		return GoNamedType {
			name: xr.name,
		}

	case ProdType:
		Assert(false)

	default:
		Assert(false)
	}
	panic("")
}

type GoNode interface {}

type GoTypeDecl struct {
	name  string
	type_ GoNode
}

type GoProdTypeCase = Word

const (
	GoProdTypeCase_Interface GoProdTypeCase = 0
	GoProdTypeCase_Struct    GoProdTypeCase = 1
)

type GoFunType struct {
	argTypes []GoNode
	retType  GoNode
}

type GoProdType struct {
	typeCase   GoProdTypeCase
	fieldNames []string
	fieldTypes []GoNode
}

type GoArrayType struct {
	elementType GoNode
}

type GoNamedType struct {
	name string
}

func GenAST(x GoNode) ast.Node {
	switch x.(type) {
	case GoTypeDecl:
		xr := x.(GoTypeDecl)
		return &ast.GenDecl {
			Tok: token.TYPE,
			Specs: []ast.Spec {
				&ast.TypeSpec {
					Name: ast.NewIdent(xr.name),
					Type: GenAST(xr.type_).(ast.Expr),
				},
			},
		}

	case GoProdType:
		xr := x.(GoProdType)
		Assert(len(xr.fieldTypes) == len(xr.fieldNames))
		fields := []*ast.Field{}
		for i, fieldName := range xr.fieldNames {
			field := &ast.Field {
				Names: []*ast.Ident {
					ast.NewIdent(fieldName),
				},
				Type: GenAST(xr.fieldTypes[i]).(ast.Expr),
			}
			fields = append(fields, field)
		}
		fieldList := &ast.FieldList {
			List: fields,
		}
		switch xr.typeCase {
		case GoProdTypeCase_Interface:
			return &ast.InterfaceType {
				Methods: fieldList,
			}
		case GoProdTypeCase_Struct:
			return &ast.StructType {
				Fields: fieldList,
			}
		default:
			panic("typeCase not recognized")
		}

	case GoArrayType:
		xr := x.(GoArrayType)
		return &ast.ArrayType {
			Elt: GenAST(xr.elementType).(ast.Expr),
		}

	case GoFunType:
		xr := x.(GoFunType)
		goParamFields := []*ast.Field {}

		return &ast.FuncType {
			Params:  &ast.FieldList {
				List: goParamFields,
			},
			Results: &ast.FieldList {
				List: []*ast.Field {
					&ast.Field {
						Names: nil,
						Type:  GenAST(xr.retType).(ast.Expr),
					},
				},
			},
		}

	case GoNamedType:
		xr := x.(GoNamedType)
		return ast.NewIdent(xr.name)

	default:
		panic("Unknown type for GenAST")
	}
}

func TestAST() {
	// fs := token.NewFileSet()
	// f := &ast.File {
	// 	Name: ast.NewIdent("fileName"),
	// 	Decls: []ast.Decl {
	// 		GenAST(GoTypeDecl {
	// 			"typeName",
	// 			GoNamedType {
	// 				"int64",
	// 			},
	// 		}).(ast.Decl),
	// 	},
	// }
	// CheckErr(printer.Fprint(os.Stdout, fs, f))

	decls := []TypeDecl {
		TypeDecl {
			name:  "Block",
			type_: ProdType {
				fields: []Field {
					Field {
						fieldName: "Field1",
						fieldType: NamedType {"Word"},
						internal: true,
					},
					Field {
						fieldName: "Field2",
						fieldType: NamedType {"byte"},
						internal: false,
					},
				},
			},
		},
	}

	goDecls := GenGoTypeDecls(decls)
	goMod := GenGoMod(goDecls)
	CheckErr(printer.Fprint(os.Stdout, goMod.astFileSet, goMod.astFile))
}
