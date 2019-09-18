package codeGen

import (
	"runtime/debug"
	"bytes"
	"crypto/sha256"
	"encoding/binary"
	"fmt"
	"go/ast"
	"go/printer"
	"go/token"
	"io"
	"os"
	"strings"
)

type Type interface {
	Case() Type_Case
}

type TypeDecl struct {
	name  string
	type_ Type
}

type Type_Case = Word

const (
	Type_Case_NamedType    Type_Case = 1
	Type_Case_AlgType      Type_Case = 2
	Type_Case_ArrayType    Type_Case = 3
	Type_Case_FunType      Type_Case = 4
	Type_Case_RefType      Type_Case = 5
	Type_Case_OptionType   Type_Case = 6
	Type_Case_MapType      Type_Case = 7
)

func (NamedType) Case() Type_Case {
	return Type_Case_NamedType
}
func (AlgType) Case() Type_Case {
	return Type_Case_AlgType
}
func (ArrayType) Case() Type_Case {
	return Type_Case_ArrayType
}
func (FunType) Case() Type_Case {
	return Type_Case_FunType
}
func (RefType) Case() Type_Case {
	return Type_Case_RefType
}
func (OptionType) Case() Type_Case {
	return Type_Case_OptionType
}
func (MapType) Case() Type_Case {
	return Type_Case_MapType
}

type AlgSort = Word

const (
	AlgSort_Prod AlgSort = 1
	AlgSort_Sum  AlgSort = 2
)

type Field struct {
	fieldName     string
	fieldType     Type
	internal      bool
	attributeList []string
}

type Method struct {
	methodName    string
	methodType    *FunType
	attributeList []string
}

type AlgType struct {
	sort          AlgSort
	fields        []Field
	methods       []Method
	attributeList []string
}

type ArrayType struct {
	elementType Type
}

type RefType struct {
	targetType Type
}

type FunType struct {
	args     []Field
	retType  Type
}

type NamedType struct {
	name string
}

type OptionType struct {
	valueType Type
}

type MapType struct {
	keyType   Type
	valueType Type
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
	HashAccWord(buf, len(s))
	*buf = append(*buf, s...)
}

func HashAccType(buf *[]byte, x Type) {
	switch x.Case() {
	case Type_Case_AlgType:
		xr := x.(*AlgType)
		HashAccWord(buf, Type_Case_AlgType)
		HashAccWord(buf, len(xr.fields))
		for _, field := range xr.fields {
			HashAccString(buf, field.fieldName)
			HashAccType(buf, field.fieldType)
			HashAccBool(buf, field.internal)
		}

	case Type_Case_ArrayType:
		xr := x.(*ArrayType)
		HashAccWord(buf, Type_Case_ArrayType)
		HashAccType(buf, xr.elementType)

	case Type_Case_RefType:
		xr := x.(*RefType)
		HashAccWord(buf, Type_Case_RefType)
		HashAccType(buf, xr.targetType)

	case Type_Case_FunType:
		xr := x.(*FunType)
		HashAccWord(buf, Type_Case_FunType)
		HashAccWord(buf, len(xr.args))
		for _, arg := range xr.args {
			HashAccString(buf, arg.fieldName)
			HashAccType(buf, arg.fieldType)
		}
		HashAccType(buf, xr.retType)

	case Type_Case_NamedType:
		xr := x.(*NamedType)
		HashAccWord(buf, Type_Case_NamedType)
		HashAccString(buf, xr.name)

	case Type_Case_OptionType:
		xr := x.(*OptionType)
		HashAccWord(buf, Type_Case_OptionType)
		HashAccType(buf, xr.valueType)

	case Type_Case_MapType:
		xr := x.(*MapType)
		HashAccWord(buf, Type_Case_MapType)
		HashAccType(buf, xr.keyType)
		HashAccType(buf, xr.valueType)

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
	ret.astFile = &ast.File{
		Name:  ast.NewIdent("fileName"),
		Decls: astDecls,
	}

	return ret
}

type GoGenContext struct {
	typeMap   map[string]Type
	retDecls  *[]GoNode
	retMap    map[TypeHash]GoNode
	tokens    []string
	concrete  bool
}

func (ctx GoGenContext) Extend(token string) GoGenContext {
	ret := ctx
	ret.tokens = append(ret.tokens, token)
	return ret
}

func (ctx GoGenContext) Concrete(concrete bool) GoGenContext {
	ret := ctx
	ret.concrete = concrete
	return ret
}

func GenGoTypeDecls(decls []TypeDecl) []GoNode {
	ctx := GoGenContext {
		typeMap:  map[string]Type{},
		retDecls: &[]GoNode{},
		retMap:   map[TypeHash]GoNode{},
		tokens:   []string{},
		concrete: false,
	}

	for _, decl := range decls {
		GenGoTypeDeclAcc(decl.name, decl.type_, ctx.Extend(decl.name))
	}

	return *ctx.retDecls
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

func IdToImplRef(name string) string {
	return name + "_R"
}

func GoTypeByteArray() GoNode {
	return GoArrayType {
		elementType: GoIdent {
			name: "byte",
		},
	}
}

func GoTypeAny() GoNode {
	return GoProdType {
		typeCase: GoProdTypeCase_Interface,
		fields:   []GoField{},
	}
}

func GenGoIdentityBody(receiverVar GoNode) []GoNode {
	ret := []GoNode {
		GoStmtReturn {
			receiverVar,
		},
	}
	return ret
}

func GenGoDerefCacheBody(receiverVar GoNode) []GoNode {
	// FIXME: look up hash in store if cached_impl is nil
	ret := []GoNode {
		GoStmtReturn {
			GoExprDot {
				value:     receiverVar,
				fieldName: "cached_impl",
			},
		},
	}
	return ret
}

func GenGoAccessorBody(receiverVar GoNode, fieldName string) []GoNode {
	ret := []GoNode {
		GoStmtReturn {
			GoExprDot {
				value:     receiverVar,
				fieldName: fieldName,
			},
		},
	}
	return ret
}

func GenGoDerefAccessorBody(receiverVar GoNode, fieldName string) []GoNode {
	ret := []GoNode {
		GoStmtReturn {
			GoExprDot {
				value: GoExprCall {
					f: GoExprDot { value: receiverVar, fieldName: "Impl", },
					args: []GoNode{},
				},
				fieldName: fieldName,
			},
		},
	}
	return ret
}

func GenGoMethodCall(obj GoNode, methodName string, args []GoNode) GoNode {
	ret := GoExprCall {
		f: GoExprDot {
			value: obj,
			fieldName: methodName,
		},
		args: args,
	}
	return ret
}

func GoTypeToIdent(typeName string) GoIdent {
	return GoIdent { name: strings.ToLower(typeName)[0:1] }
}

func GenGoTypeDeclAcc(name string, x Type, ctx GoGenContext) GoNode {
	Assert(x != nil)

	switch x.(type) {
	case *NamedType:
		xr := x.(*NamedType)
		return GoIdent {name: xr.name}

	case *OptionType:
		xr := x.(*OptionType)
		goValueType := GenGoTypeAcc(xr.valueType, ctx.Extend("OptValue"))
		goValueType.implements_GoNode()
		return GoIdent {name: "_TODO_"}

	case *MapType:
		xr := x.(*MapType)
		goKeyType := GenGoTypeAcc(xr.keyType, ctx.Extend("MapKey"))
		goValueType := GenGoTypeAcc(xr.valueType, ctx.Extend("MapValue"))
		goKeyType.implements_GoNode()
		goValueType.implements_GoNode()
		return GoIdent {name: "_TODO_"}

	case *AlgType:
		xr := x.(*AlgType)
	
		interfaceName := name
		interfaceID := GoIdent {name: interfaceName}

		implName := IdToImpl(name)
		implID := GoIdent {name: implName}

		implRefName := IdToImplRef(name)
		implRefID := GoIdent {name: implRefName}

		interfaceFields := []GoField{}
		implFields := []GoField{}
		implRefFields := []GoField{}

		if xr.sort == AlgSort_Prod {
			for _, field := range xr.fields {
				implFields = append(implFields, GoField {
					fieldName: field.fieldName,
					fieldType: GenGoTypeAcc(field.fieldType, ctx.Extend(field.fieldName).Concrete(true)),
				})

				interfaceFields = append(interfaceFields, GoField {
					fieldName: field.fieldName,
					fieldType: GoFunType {
						args:     []GoField{},
						retType:  GenGoTypeAcc(field.fieldType, ctx.Extend(field.fieldName)),
					},
				})
			}
		}

		for _, method := range xr.methods {
			interfaceFields = append(interfaceFields, GoField {
				fieldName: method.methodName,
				fieldType: GenGoTypeAcc(method.methodType, ctx.Extend(method.methodName)),
			})
		}

		interfaceFields = append(interfaceFields, GoField {
			fieldName: "Impl",
			fieldType: GoFunType {
				retType:  GoPtrType { targetType: implID },
				args:     []GoField{},
			},
		})

		interfaceFields = append(interfaceFields, GoField {
			fieldName: "CID",
			fieldType: GoFunType {
				retType:  GoTypeByteArray(),
				args:     []GoField{},
			},
		})

		var caseTypeID GoNode = nil

		if xr.sort == AlgSort_Sum {
			caseNames := []string{}

			caseTypeName := name + "_Case"
			caseTypeID = GoIdent { caseTypeName }

			*ctx.retDecls = append(*ctx.retDecls, GoTypeDecl {
				name: caseTypeName,
				type_: GoIdent { name: "Word" },
			})

			for _, field := range xr.fields {
				caseNames = append(caseNames, field.fieldName)
				caseWhich := GoIdent { caseTypeName + "_" + field.fieldName }

				caseInterfaceName := name + "_" + field.fieldName
				GenGoTypeDeclAcc(caseInterfaceName, field.fieldType, ctx.Extend(field.fieldName))

				caseInterfaceType := GoIdent { caseInterfaceName }
				caseImplType := GoIdent { IdToImpl(caseInterfaceName) }
				// caseImplPtrType := GoPtrType { targetType: caseImplType }

				interfaceFields = append(interfaceFields, GoField {
					fieldName: "As_" + field.fieldName,
					fieldType: GoFunType {
						retType:   GoIdent { caseInterfaceName },
						args:      []GoField{},
					},
				})

				caseAsDeclArg := GoTypeToIdent(name + "_" + field.fieldName)
				caseAsDeclBody := []GoNode {
					GoStmtExpr {
						expr: GoExprCall {
							f: GoIdent { name: "Assert" },
							args: []GoNode {
								GoExprEq {
									GenGoMethodCall(caseAsDeclArg, "Which", []GoNode{}),
									caseWhich,
								},
							},
						},
					},
					GoStmtReturn {
						value: GoExprCast {
							arg:     caseAsDeclArg,
							resType: caseInterfaceType,
						},
					},
				}
				caseAsDecl := GoFunDecl {
					receiverVar: &caseAsDeclArg,
					receiverType: GoIdent { caseInterfaceName },
					funName: "As_" + field.fieldName,
					funType: GoFunType {
						args: []GoField {},
						retType: caseInterfaceType,
					},
					funArgs: []GoNode{},
					funBody: caseAsDeclBody,
				}

				*ctx.retDecls = append(*ctx.retDecls, caseAsDecl)



				caseNewDeclArg := GoTypeToIdent(name + "_" + field.fieldName)
				caseNewDeclBody := []GoNode {
					GoStmtReturn {
						GoExprAddrOf {
							target: GoExprStruct {
								type_:  caseImplType,
								fields: []GoField {
									GoField {
										fieldName: "cached_cid",
										fieldType: GoExprLitNil {},
									},
									GoField {
										fieldName: "rawValue",
										fieldType: caseNewDeclArg,
									},
									GoField {
										fieldName: "which",
										fieldType: caseWhich,
									},
								},
							},
						},
					},
				}
				caseNewDecl := GoFunDecl {
					receiverVar: nil,
					receiverType: nil,
					funName: name + "_Make_" + field.fieldName,
					funType: GoFunType {
						args: []GoField {
							GoField {
								fieldName: caseNewDeclArg.name,
								fieldType: GoIdent { caseInterfaceName },
							},
						},
						retType: interfaceID,
					},
					funArgs: []GoNode{caseNewDeclArg},
					funBody: caseNewDeclBody,
				}

				*ctx.retDecls = append(*ctx.retDecls, caseNewDecl)
			}

			caseTypeDecl := GoEnumDecl {
				name:      caseTypeName,
				caseNames: caseNames,
			}

			*ctx.retDecls = append(*ctx.retDecls, caseTypeDecl)

			interfaceFields = append(interfaceFields, GoField {
				fieldName: "Which",
				fieldType: GoFunType {
					retType: caseTypeID,
					args:    []GoField{},
				},
			})
		}

		implFields = append(implFields, GoField {
			fieldName: "cached_cid",
			fieldType: GoTypeByteArray(),
		})

		if xr.sort == AlgSort_Sum {
			implFields = append(implFields, GoField {
				fieldName: "rawValue",
				fieldType: GoTypeAny(),
			})

			implFields = append(implFields, GoField {
				fieldName: "which",
				fieldType: caseTypeID,
			})
		}

		implRefFields = append(implRefFields, GoField {
			fieldName: "cid",
			fieldType: GoTypeByteArray(),
		})

		cachedObjectField := GoField {
			fieldName: "cached_impl",
			fieldType: GoPtrType { targetType: implID },
		}
		implRefFields = append(implRefFields, cachedObjectField)

		interfaceDecl := GoTypeDecl{
			name: interfaceName,
			type_: GoProdType{
				typeCase:   GoProdTypeCase_Interface,
				fields:     interfaceFields,
			},
		}

		implDecl := GoTypeDecl{
			name: implName,
			type_: GoProdType{
				typeCase:   GoProdTypeCase_Struct,
				fields:     implFields,
			},
		}

		implRefDecl := GoTypeDecl{
			name: implRefName,
			type_: GoProdType{
				typeCase:   GoProdTypeCase_Struct,
				fields:     implRefFields,
			},
		}

		*ctx.retDecls = append(*ctx.retDecls, interfaceDecl)
		*ctx.retDecls = append(*ctx.retDecls, implDecl)
		*ctx.retDecls = append(*ctx.retDecls, implRefDecl)

		implReceiverVar := GoTypeToIdent(name)
		implRefReceiverVar := GoTypeToIdent(name)
		for _, field := range xr.fields {
			baseFieldType := GenGoTypeAcc(field.fieldType, ctx.Extend(field.fieldName))

			implAccessorDecl := GoFunDecl {
				receiverVar: &implReceiverVar,
				receiverType: GoPtrType {
					targetType: implID,
				},
				funName: field.fieldName,
				funType: GoFunType {
					args:    []GoField{},
					retType: baseFieldType,
				},
				funArgs: []GoNode{},
				funBody: GenGoAccessorBody(implReceiverVar, field.fieldName),
			}

			implRefAccessorDecl := GoFunDecl {
				receiverVar: &implRefReceiverVar,
				receiverType: GoPtrType {
					targetType: implRefID,
				},
				funName: field.fieldName,
				funType: GoFunType {
					args:    []GoField{},
					retType: baseFieldType,
				},
				funArgs: []GoNode{},
				funBody: GenGoDerefAccessorBody(implRefReceiverVar, field.fieldName),
			}

			*ctx.retDecls = append(*ctx.retDecls, implAccessorDecl)
			*ctx.retDecls = append(*ctx.retDecls, implRefAccessorDecl)
		}

		implImplDecl := GoFunDecl {
			receiverVar: &implReceiverVar,
			receiverType: GoPtrType { targetType: implID, },
			funName: "Impl",
			funType: GoFunType {
				args:    []GoField{},
				retType: GoPtrType { targetType: implID },
			},
			funArgs: []GoNode{},
			funBody: GenGoIdentityBody(implReceiverVar),
		}

		implRefImplDecl := GoFunDecl {
			receiverVar: &implRefReceiverVar,
			receiverType: GoPtrType { targetType: implRefID, },
			funName: "Impl",
			funType: GoFunType {
				args:    []GoField{},
				retType: GoPtrType { targetType: implID },
			},
			funArgs: []GoNode{},
			funBody: GenGoDerefCacheBody(implReceiverVar),
		}

		*ctx.retDecls = append(*ctx.retDecls, implImplDecl)
		*ctx.retDecls = append(*ctx.retDecls, implRefImplDecl)

		// TODO: assert not in retMap
		ctx.retMap[HashType(x)] = interfaceID
		return interfaceID

	case *ArrayType:
		panic("TODO")
	case *RefType:
		panic("TODO")
	case *FunType:
		panic("TODO")
	default:
		errMsg := fmt.Sprintf("TODO: %v\n", x.Case())
		panic(errMsg)
	}
}

func GenGoTypeAcc(x Type, ctx GoGenContext) GoNode {
	if match, ok := ctx.retMap[HashType(x)]; ok {
		return match
	}

	switch x.Case() {
	case Type_Case_AlgType:
		typeName := strings.Join(ctx.tokens, "_")
		return GenGoTypeDeclAcc(typeName, x, ctx.Concrete(false))

	case Type_Case_ArrayType:
		xr := x.(*ArrayType)
		goElementType := GenGoTypeAcc(xr.elementType, ctx.Extend("ArrayElement"))
		return GoArrayType {
			elementType: goElementType,
		}

	case Type_Case_RefType:
		xr := x.(*RefType)
		goTargetType := GenGoTypeAcc(xr.targetType, ctx.Extend("RefTarget"))
		switch goTargetType.(type) {
		case GoIdent:
			if ctx.concrete {
				return GoPtrType {
					targetType: GoIdent { name: IdToImplRef(goTargetType.(GoIdent).name), },
				}
			} else {
				return goTargetType
			}
		default:
			panic("General reference types not yet supported")
		}

	case Type_Case_FunType:
		xr := x.(*FunType)
		goRetType := GenGoTypeAcc(xr.retType, ctx.Extend("FunRet"))
		goArgs := []GoField{}
		for i, arg := range xr.args {
			goArgType := GenGoTypeAcc(arg.fieldType, ctx.Extend(fmt.Sprintf("FunArg%v", i)))
			goArg := GoField {
				fieldName: arg.fieldName,
				fieldType: goArgType,
			}
			goArgs = append(goArgs, goArg)
		}
		return GoFunType {
			retType:  goRetType,
			args:     goArgs,
		}

	case Type_Case_NamedType:
		xr := x.(*NamedType)
		return GoIdent {
			name: xr.name,
		}

	case Type_Case_OptionType:
		typeName := strings.Join(ctx.tokens, "_")
		return GenGoTypeDeclAcc(typeName, x, ctx.Concrete(false))

	case Type_Case_MapType:
		typeName := strings.Join(ctx.tokens, "_")
		return GenGoTypeDeclAcc(typeName, x, ctx.Concrete(false))

	default:
		panic("TODO")
	}
}

type GoNode interface {
	implements_GoNode()
}

func (GoTypeDecl) implements_GoNode() {}
func (GoEnumDecl) implements_GoNode() {}
func (GoFunDecl) implements_GoNode() {}
func (GoFunType) implements_GoNode() {}
func (GoPtrType) implements_GoNode() {}
func (GoStmtReturn) implements_GoNode() {}
func (GoStmtExpr) implements_GoNode() {}
func (GoExprDot) implements_GoNode() {}
func (GoExprEq) implements_GoNode() {}
func (GoExprCast) implements_GoNode() {}
func (GoExprCall) implements_GoNode() {}
func (GoExprStruct) implements_GoNode() {}
func (GoExprAddrOf) implements_GoNode() {}
func (GoExprLitNil) implements_GoNode() {}
func (GoField) implements_GoNode() {}
func (GoProdType) implements_GoNode() {}
func (GoArrayType) implements_GoNode() {}
func (GoIdent) implements_GoNode() {}

type GoTypeDecl struct {
	name  string
	type_ GoNode
}

type GoEnumDecl struct {
	name       string
	caseNames  []string
}

type GoFunDecl struct {
	receiverVar  *GoIdent
	receiverType GoNode
	funName      string
	funType      GoFunType
	funArgs      []GoNode
	funBody      []GoNode
}

type GoProdTypeCase = Word

const (
	GoProdTypeCase_Interface GoProdTypeCase = 0
	GoProdTypeCase_Struct    GoProdTypeCase = 1
)

type GoFunType struct {
	args     []GoField
	retType  GoNode
}

type GoPtrType struct {
	targetType GoNode
}

type GoStmtReturn struct {
	value GoNode
}

type GoStmtExpr struct {
	expr GoNode
}

type GoExprDot struct {
	value     GoNode
	fieldName string
}

type GoExprEq struct {
	lhs GoNode
	rhs GoNode
}

type GoExprCast struct {
	arg     GoNode
	resType GoNode
}

type GoExprCall struct {
	f    GoNode
	args []GoNode
}

type GoExprStruct struct {
	type_  GoNode
	fields []GoField
}

type GoExprAddrOf struct {
	target GoNode
}

type GoExprLitNil struct { }

type GoField struct {
	fieldName string
	fieldType GoNode  // TODO: rename
}

type GoProdType struct {
	typeCase  GoProdTypeCase
	fields    []GoField
}

type GoArrayType struct {
	elementType GoNode
}

type GoIdent struct {
	name string
}

func GenAST(x GoNode) ast.Node {
	switch x.(type) {
	case GoTypeDecl:
		xr := x.(GoTypeDecl)
		return &ast.GenDecl{
			Tok: token.TYPE,
			Specs: []ast.Spec{
				&ast.TypeSpec{
					Name: ast.NewIdent(xr.name),
					Type: GenAST(xr.type_).(ast.Expr),
				},
			},
		}

	case GoEnumDecl:
		xr := x.(GoEnumDecl)
		caseTypeNames := []*ast.Ident{}
		caseTypeValues := []ast.Expr{}
		for i, caseName := range xr.caseNames {
			caseTypeNames = append(caseTypeNames, ast.NewIdent(xr.name + "_" + caseName))
			caseTypeValue := &ast.BasicLit {
				Kind:  token.INT,
				Value: fmt.Sprintf("%v", i+1),
			}
			caseTypeValues = append(caseTypeValues, caseTypeValue)
		}
		return &ast.GenDecl {
			Tok: token.CONST,
			Specs: []ast.Spec {
				&ast.ValueSpec {
					Names: caseTypeNames,
					Type: ast.NewIdent(xr.name),
					Values: caseTypeValues,
				},
			},
		}

	case GoFunDecl:
		xr := x.(GoFunDecl)
		var goRecv *ast.FieldList = nil
		if xr.receiverType != nil {
			goRecv = &ast.FieldList {
				List: []*ast.Field {
					&ast.Field {
						Names: []*ast.Ident {
							GenAST(*xr.receiverVar).(*ast.Ident),
						},
						Type: GenAST(xr.receiverType).(ast.Expr),
					},
				},
			}
		}
		goFunType := GenAST(xr.funType).(*ast.FuncType)
		goFunStmts := []ast.Stmt{}
		for _, stmt := range xr.funBody {
			goFunStmt := GenAST(stmt).(ast.Stmt)
			goFunStmts = append(goFunStmts, goFunStmt)
		}
		goFunBody := &ast.BlockStmt {
			List: goFunStmts,
		}
		return &ast.FuncDecl {
			Recv: goRecv,
			Name: ast.NewIdent(xr.funName),
			Type: goFunType,
			Body: goFunBody,
		}

	case GoProdType:
		xr := x.(GoProdType)
		fields := []*ast.Field{}
		for _, field := range xr.fields {
			field := &ast.Field{
				Names: []*ast.Ident{
					ast.NewIdent(field.fieldName),
				},
				Type: GenAST(field.fieldType).(ast.Expr),
			}
			fields = append(fields, field)
		}
		fieldList := &ast.FieldList{
			List: fields,
		}
		switch xr.typeCase {
		case GoProdTypeCase_Interface:
			return &ast.InterfaceType{
				Methods: fieldList,
			}
		case GoProdTypeCase_Struct:
			return &ast.StructType{
				Fields: fieldList,
			}
		default:
			panic("typeCase not recognized")
		}

	case GoArrayType:
		xr := x.(GoArrayType)
		return &ast.ArrayType{
			Elt: GenAST(xr.elementType).(ast.Expr),
		}

	case GoFunType:
		xr := x.(GoFunType)
		goParamFields := []*ast.Field{}
		for _, arg := range xr.args {
			goParamFields = append(goParamFields, &ast.Field {
				Names: []*ast.Ident {ast.NewIdent(arg.fieldName)},
				Type:  GenAST(arg.fieldType).(ast.Expr),
			})
		}

		return &ast.FuncType{
			Params: &ast.FieldList{
				List: goParamFields,
			},
			Results: &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: nil,
						Type:  GenAST(xr.retType).(ast.Expr),
					},
				},
			},
		}

	case GoPtrType:
		xr := x.(GoPtrType)
		return &ast.StarExpr {
			X: GenAST(xr.targetType).(ast.Expr),
		}

	case GoStmtReturn:
		xr := x.(GoStmtReturn)
		return &ast.ReturnStmt {
			Results: []ast.Expr {
				GenAST(xr.value).(ast.Expr),
			},
		}

	case GoExprDot:
		xr := x.(GoExprDot)
		return &ast.SelectorExpr {
			X:   GenAST(xr.value).(ast.Expr),
			Sel: ast.NewIdent(xr.fieldName),
		}

	case GoExprCall:
		xr := x.(GoExprCall)
		goArgs := []ast.Expr{}
		for _, f := range xr.args {
			goArgs = append(goArgs, GenAST(f).(ast.Expr))
		}
		return &ast.CallExpr {
			Fun:  GenAST(xr.f).(ast.Expr),
			Args: goArgs,
		}

	case GoStmtExpr:
		xr := x.(GoStmtExpr)
		return &ast.ExprStmt {
			X: GenAST(xr.expr).(ast.Expr),
		}

	case GoExprLitNil:
		return &ast.BasicLit {
			Kind:  token.STRING,
			Value: "nil",
		}

	case GoExprStruct:
		xr := x.(GoExprStruct)
		goFields := []ast.Expr{}
		for _, field := range xr.fields {
			goField := &ast.KeyValueExpr {
				Key:   ast.NewIdent(field.fieldName),
				Value: GenAST(field.fieldType).(ast.Expr),
			}
			goFields = append(goFields, goField)
		}
		return &ast.CompositeLit {
			Type:  GenAST(xr.type_).(ast.Expr),
			Elts:  goFields,
		}

	case GoExprAddrOf:
		xr := x.(GoExprAddrOf)
		return &ast.UnaryExpr {
			Op: token.AND,
			X:  GenAST(xr.target).(ast.Expr),
		}

	case GoExprEq:
		xr := x.(GoExprEq)
		return &ast.BinaryExpr {
			Op: token.EQL,
			X:  GenAST(xr.lhs).(ast.Expr),
			Y:  GenAST(xr.rhs).(ast.Expr),
		}

	case GoExprCast:
		xr := x.(GoExprCast)
		return &ast.TypeAssertExpr {
			X:    GenAST(xr.arg).(ast.Expr),
			Type: GenAST(xr.resType).(ast.Expr),
		}

	case GoIdent:
		xr := x.(GoIdent)
		return ast.NewIdent(xr.name)

	default:
		fmt.Printf("Unknown type: %T %v\n", x, x)
		panic("Unknown type for GenAST")
	}
}

const Whitespace = " \t\n"
const Symbols = "(){}[],;|&?://*"

const DebugParser = false

type LineInfo struct {
	line  Word
	col   Word
}

type ParseStreamState struct {
	pos Word

}

type ParseStream struct {
	rs          io.ReadSeeker
	state       ParseStreamState
	stateStack  []ParseStreamState
	lineMap     []LineInfo
	buffer      *bytes.Buffer
}

func (r *ParseStream) Push() {
	if DebugParser {
		fmt.Printf("Push():  %v\n", r.PosDebug())
	}

	r.stateStack = append(r.stateStack, r.state)
}

func (r *ParseStream) Pop(restore bool) {
	if DebugParser {
		fmt.Printf("Pop(%v):  %v\n", restore, r.PosDebug())
	}

	n := len(r.stateStack)
	Assert(n > 0)
	if restore {
		r.Seek(r.stateStack[n-1].pos - r.state.pos)
		r.state = r.stateStack[n-1]
	}
	r.stateStack = r.stateStack[0:n-1]
}

func (r *ParseStream) Seek(offset Word) {
	// if DebugParser {
	// 	fmt.Printf("Seek: %v\n", offset)
	// }
	_, err := r.rs.Seek(int64(offset), io.SeekCurrent)
	CheckErr(err)
	r.state.pos += offset
}

func (r *ParseStream) PosDebug() string {
	var i Word = r.state.pos - 1
	var lineInfo LineInfo
	if i < 0 {
		lineInfo = LineInfo { line: 0, col: 0 }
	} else {
		lineInfo = r.lineMap[i]
	}
	return fmt.Sprintf("line %v, column %v", lineInfo.line + 1, lineInfo.col + 1)
}

func (r *ParseStream) GenParseError(errMsg string) ParseError_S {
	var i Word = r.state.pos - 1
	if i < 0 {
		i = 0
	}
	errMsgNew := fmt.Sprintf("Parse error (%v)\n\n", r.PosDebug())

	for {
		if i > 0 && r.buffer.Bytes()[i] != '\n' {
			i--
		} else {
			break
		}
	}
	for i++; i < len(r.buffer.Bytes()); i++ {
		c := string(r.buffer.Bytes()[i:i+1])
		errMsgNew = errMsgNew + c
		if c == "\n" {
			break
		}
	}
	if i == len(r.buffer.Bytes()) {
		errMsgNew += "\n"
	}
	errMsgNew += "\n"
	errMsgNew += errMsg
	errMsgNew += "\n\n"
	errMsgNew += string(debug.Stack())
	return ParseError(errMsgNew)
}

func (r *ParseStream) Get(lenMax Word, advance bool) string {
	buf := make([]byte, lenMax)

	n, err := r.rs.Read(buf)
	if !advance {
		defer r.Seek(-n)
	}

	ret := string(buf[:n])
	for i, c := range buf[:n] {
		p := r.state.pos + i
		if p >= len(r.lineMap) {
			Assert(p == len(r.lineMap))
			prevLineInfo := LineInfo { line: 0, col: 0 }
			if p >= 1 {
				prevLineInfo = r.lineMap[p-1]
			}
			newLineInfo := LineInfo {
				line: prevLineInfo.line,
				col:  prevLineInfo.col + 1,
			}
			if c == '\n' {
				newLineInfo = LineInfo {
					line: prevLineInfo.line + 1,
					col:  0,
				}
			}
			r.lineMap = append(r.lineMap, newLineInfo)
			r.buffer.WriteByte(c)
		}
	}
	r.state.pos += n
	if err == io.EOF {
		Assert(n == 0)
		return ret
	}
	if err != nil {
		if DebugParser {
			fmt.Printf("err == nil: %v\n", err == nil)
		}
		panic("Read error")
	}
	if n > lenMax {
		panic("Read error")
	}
	return ret
}


func (r *ParseStream) Read(lenMax Word) string {
	return r.Get(lenMax, true)
}

func (r *ParseStream) Peek(lenMax Word) string {
	return r.Get(lenMax, false)
}


func (r *ParseStream) ParseComments() (hitComment bool, err error, hitNewline bool) {
	hitComment = false

	if DebugParser {
		fmt.Printf(" ** ParseComments body: %v \n", r.PosDebug())
	}

	depth := 0
	singleLineComment := false
	for {
		if _, ok := r.PeekExact(1); !ok {
			return
		}

		if res, ok := r.PeekExact(1); strings.Contains(Whitespace, res) && ok {
			if res == "\n" {
				singleLineComment = false
				if depth == 0 {
					hitNewline = true
				}
			}
			r.Seek(1)
			continue
		}

		if depth == 0 {
			if res, ok := r.PeekExact(2); res == "//" && ok {
				r.Seek(2)
				singleLineComment = true
			}
		}

		if res, ok := r.PeekExact(2); res == "/*" && ok {
			r.Seek(2)
			depth += 1
		} else if res, ok := r.PeekExact(2); res == "*/" && ok {
			r.Seek(2)
			depth -= 1
		} else {
			if depth > 0 || singleLineComment {
				r.Seek(1)
			} else {
				return
			}
		}
	}
}


func ReadToken(r *ParseStream) (ret string, hitNewline bool, err error) {
	if DebugParser {
		defer func() {
			fmt.Printf("ReadToken: \"%v\", %v, %v\n", ret, hitNewline, err == nil)
		}()
	}

	retToken := []byte{}
	_, err, hitNewline = r.ParseComments()
	if err != nil {
		ret = ""
		return
	}

	for {
		c, ok := r.PeekExact(1)
		if !ok {
			if len(retToken) > 0 {
				ret = string(retToken)
				return
			} else {
				ret = ""
				err = io.EOF
				return
			}
		}

		Assert(len(c) == 1)

		if strings.Contains(Whitespace, c) {
			ret = string(retToken)
			return

		} else if strings.Contains(Symbols, c) {
			if len(retToken) == 0 {
				retToken = append(retToken, c[0])
				r.Seek(1)
			}
			ret = string(retToken)
			return

		} else {
			retToken = append(retToken, c...)
			r.Seek(1)
			continue
		}
	}
}

type ParseError_S struct {
	msg string
}

func ParseError(msg string) ParseError_S {
	return ParseError_S {msg: msg}
}

func IsLower(c byte) bool {
	return (c >= 'a' && c <= 'z')
}

func IsUpper(c byte) bool {
	return (c >= 'A' && c <= 'Z')
}

func IsAlpha(c byte) bool {
	return IsLower(c) || IsUpper(c)
}

func IsDigit(c byte) bool {
	return (c >= '0' && c <= '9')
}

func IsIdent(s string) bool {
	for i := 0; i < len(s); i++ {
		c := s[i]
		if IsAlpha(c) || c == '_' {
			continue
		}
		if IsDigit(c) && i > 0 {
			continue
		}
		return false
	}
	return true
}

func ParseIdent(r *ParseStream) (string, error) {
	ret, _, err := ReadToken(r)
	if err != nil {
		return "", err
	} else if !IsIdent(ret) {
		msg := fmt.Sprintf("Expected identifier; received: \"%v\"", ret)
		return "", r.GenParseError(msg)
	} else {
		return ret, nil
	}
}

func (err ParseError_S) Error() string {
	return err.msg
}

func ReadTokenCheck(r *ParseStream, validTokens []string) (string, error) {
	for _, v := range validTokens {
		if v == "\n" {
			if tok, ok := PeekToken(r, true); ok && (tok == "\n") {
				if DebugParser {
					fmt.Printf("ReadTokenCheck read newline (special case): %v\n", r.PosDebug())
				}
				return tok, nil
			}
		}
	}

	tok, _, err := ReadToken(r)
	if err != nil {
		return "", err
	}
	for _, v := range validTokens {
		if tok == v {
			return tok, nil
		}
	}

	validTokensDisp := []string{}
	for _, tok := range validTokens {
		validTokenDisp := "\"" + tok + "\""
		validTokensDisp = append(validTokensDisp, validTokenDisp)
	}
	errMsg := fmt.Sprintf("Unexpected token \"%v\" (expected: %v)", tok, strings.Join(validTokensDisp, ", "))
	return "", r.GenParseError(errMsg)
}

func (r *ParseStream) ReadTokenSequenceCheck(tokenSeq []string) (string, error) {
	ret := ""
	for _, token := range tokenSeq {
		currRet, err := ReadTokenCheck(r, []string{token})
		ret += currRet
		if err != nil {
			return ret, err
		}
	}
	return ret, nil
}

func TryReadTokenCheck(r *ParseStream, validTokens []string) (ret string, err error) {
	r.Push()
	defer func() { r.Pop(err != nil) }()

	ret, err = ReadTokenCheck(r, validTokens)
	return
}

func (r *ParseStream) PeekExact(len_ Word) (string, bool) {
	ret := r.Peek(len_)
	return ret, (len(ret) == len_)
}

func PeekToken(r *ParseStream, matchNewline bool) (string, bool) {
	r.Push()
	defer func() { r.Pop(true) }()

	ret, hitNewline, err := ReadToken(r)
	if hitNewline && matchNewline {
		return "\n", true
	}
	if err != nil {
		Assert(err == io.EOF)
		return "", false
	} else {
		return ret, true
	}
}

func ParseAttributeList(r *ParseStream) ([]string, error) {
	if tok, ok := PeekToken(r, false); ok && (tok == "@") {
		_, _ = ReadTokenCheck(r, []string{"@"})
		fTryParse := func (rr *ParseStream) (interface{}, error) {
			ret, _, err := ReadToken(rr)
			if err != nil {
				return nil, err
			}
			return ret, nil
		}
		ret := []string{}
		fAppend := func (x interface{}) {
			ret = append(ret, x.(string))
		}
		err := ParseDelimitedList(r, "(", []string{","}, ")", fTryParse, fAppend, false)
		if err != nil {
			return nil, err
		}
		return ret, nil
	} else {
		return []string{}, nil
	}
}

func ParseField(r *ParseStream, allowEmptyFieldTypes bool, delimSet []string) (*Field, error) {
	fieldName, err := ParseIdent(r)

	if DebugParser {
		fmt.Printf(" >>>>> IDENT: %v %v\n", fieldName, err == nil)
	}

	if err != nil {
		return nil, err
	}

	if allowEmptyFieldTypes {
		matchNewline := false
		for _, delim := range delimSet {
			if delim == "\n" {
				matchNewline = true
			}
		}
		if tok, ok := PeekToken(r, matchNewline); ok {
			for _, delim := range delimSet {
				if tok == delim {
					return &Field {
						fieldName: fieldName,
						fieldType: &AlgType {
							sort: AlgSort_Prod,
							fields: []Field{},
							methods: []Method{},
							attributeList: []string{},
						},
						internal: false,
						attributeList: []string{},
					}, nil
				}
			}
		}
	}

	fieldType, err := ParseType(r)
	if DebugParser {
		fmt.Printf(" >>>>> ParseField ParseType: %v %v\n", fieldType, err == nil)
	}
	if err != nil {
		return nil, err
	}

	attributeList, err := ParseAttributeList(r)
	if DebugParser {
		fmt.Printf(" >>>>> ParseField ParseAttributeList: %v %v\n", attributeList, err == nil)
	}
	if err != nil {
		return nil, err
	}
	
	return &Field {
		fieldName: fieldName,
		fieldType: fieldType,
		internal: false,
		attributeList: attributeList,
	}, nil
}

func TryParseField(r *ParseStream, allowEmptyFieldTypes bool, delimSet []string) (ret *Field, err error) {
	r.Push()
	defer func(){ r.Pop(err != nil) }()

	ret, err = ParseField(r, allowEmptyFieldTypes, delimSet)
	return
}

func ParseDelimitedList(
	r *ParseStream,
	start string,
	validDelims []string,
	end string,
	fTryParse func(*ParseStream)(interface{}, error),
	fAppend func(interface{})(),
	allowInitDelim bool) error {

	_, err := ReadTokenCheck(r, []string{start})
	if err != nil {
		return err
	}

	if DebugParser {
		fmt.Printf("ParseDelimitedList read start: %v, %v\n", start, r.PosDebug())
	}

	if allowInitDelim {
		ret, err := TryReadTokenCheck(r, validDelims)
		if DebugParser {
			fmt.Printf("ParseDelimitedList read init delim: %v, %v, %v\n", ret, err == nil, r.PosDebug())
		}
	}

	for {
		if x, errTryParse := fTryParse(r); errTryParse == nil {
			if DebugParser {
				fmt.Printf("ParseDelimitedList read item: %T %v\n", x, x)
			}
			fAppend(x)
		} else {
			_, err := ReadTokenCheck(r, []string{end})
			if err != nil {
				return errTryParse
			} else {
				return nil
			}
		}
		if _, errTok := TryReadTokenCheck(r, validDelims); errTok != nil {
			_, err := ReadTokenCheck(r, []string{end})
			if err != nil {
				return err
			} else {
				return nil
			}
		}
	}
}

func ParseFieldList(
	r *ParseStream,
	start string,
	validDelims []string,
	end string,
	allowMethods bool,
	allowInitDelim bool,
	allowEmptyFieldTypes bool,
	) ([]Field, []Method, error) {

	retFields  := []Field{}
	retMethods := []Method{}
	fTryParse := func (r *ParseStream) (interface{}, error) {
		retField, errField := TryParseField(r, allowEmptyFieldTypes, append(validDelims, end))
		if errField == nil {
			return retField, nil
		} else {
			if allowMethods {
				retMethod, errMethod := TryParseMethod(r)
				if errMethod == nil {
					return retMethod, nil
				} else {
					return nil, errField
				}
			} else {
				return nil, errField
			}
		}
	}
	fAppend := func (x interface{}) {
		switch x.(type) {
		case *Field:
			retFields = append(retFields, *(x.(*Field)))
		case *Method:
			retMethods = append(retMethods, *(x.(*Method)))
		default:
			panic("Case not supported in ParseFieldList -> fAppend")
		}
	}
	err := ParseDelimitedList(r, start, validDelims, end, fTryParse, fAppend, allowInitDelim)
	if err != nil {
		return nil, nil, err
	} else {
		return retFields, retMethods, nil
	}
}

func ParseMethod(r *ParseStream) (*Method, error) {
	methodName, err := ParseIdent(r)

	if DebugParser {
		fmt.Printf("ParseMethod start: name: %v (err=%v)\n", methodName, err)
	}

	if err != nil {
		return nil, err
	}

	methodArgs, _, err := ParseFieldList(r, "(", []string{","}, ")", false, false, false)
	if err != nil {
		return nil, err
	}

	methodRetType, err := ParseType(r)
	if err != nil {
		return nil, err
	}

	methodType := FunType {
		args:    methodArgs,
		retType: methodRetType,
	}

	attributeList, err := ParseAttributeList(r)
	if DebugParser {
		fmt.Printf(" >>>>> ParseMethod ParseAttributeList: %v %v\n", attributeList, err)
	}
	if err != nil {
		return nil, err
	}

	return &Method {
		methodName:    methodName,
		methodType:    &methodType,
		attributeList: attributeList,
	}, nil
}

func TryParseMethod(r *ParseStream) (ret *Method, err error) {
	r.Push()
	defer func(){ r.Pop(err != nil) }()

	ret, err = ParseMethod(r)
	return
}

func ParseType(r *ParseStream) (Type, error) {
	tok, _, err := ReadToken(r)
	if err != nil {
		return nil, err
	}

	if DebugParser {
		fmt.Printf("ParseType token: \"%v\"\n", tok)
	}

	var ret Type

	switch {
	case tok == "struct" || tok == "union" || tok == "enum":
		algSort := AlgSort_Prod
		validDelims := []string{"\n", ","}
		if tok == "union" || tok == "enum" {
			algSort = AlgSort_Sum
			validDelims = append(validDelims, "|")
		}

		attributeList, err := ParseAttributeList(r)
		if err != nil {
			return nil, err
		}
		
		allowEmptyFieldTypes := (algSort == AlgSort_Sum)
		fields, methods, err := ParseFieldList(r, "{", validDelims, "}", true, true, allowEmptyFieldTypes)

		if err != nil {
			return nil, err
		}

		ret = &AlgType {
			sort:          algSort,
			fields:        fields,
			methods:       methods,
			attributeList: attributeList,
		}

	case tok == "[":
		elementType, err := ParseType(r)
		if err != nil {
			return nil, err
		}
		_, err = ReadTokenCheck(r, []string{"]"})
		if err != nil {
			return nil, err
		}
		ret = &ArrayType {
			elementType: elementType,
		}

	case tok == "&":
		targetType, err := ParseType(r)
		if err != nil {
			return nil, err
		}
		ret = &RefType {
			targetType: targetType,
		}

	case tok == "{":
		keyType, err := ParseType(r)
		if err != nil {
			return nil, err
		}
		_, err = ReadTokenCheck(r, []string{":"})
		if err != nil {
			return nil, err
		}
		valueType, err := ParseType(r)
		if err != nil {
			return nil, err
		}
		_, err = ReadTokenCheck(r, []string{"}"})
		if err != nil {
			return nil, err
		}
		ret = &MapType {
			keyType:   keyType,
			valueType: valueType,
		}

	default:
		if IsIdent(tok) {		
			ret = &NamedType {name: tok}
		} else {
			return nil, r.GenParseError(fmt.Sprintf("Expected type; received \"%v\"", tok))
		}
	}

	for {
		if tok, ok := PeekToken(r, false); ok && (tok == "?") {
			_, _ = ReadTokenCheck(r, []string{"?"})
			ret = &OptionType {
				valueType: ret,
			}
		} else {
			break
		}
	}

	return ret, nil
}

func ParseTypeDecl(r *ParseStream) (*TypeDecl, error) {
	_, err := ReadTokenCheck(r, []string{"type"})
	if err != nil {
		return nil, err
	}
	declName, _, err := ReadToken(r)
	if err != nil {
		return nil, err
	}
	if declType, err := ParseType(r); err != nil {
		return nil, err
	} else {
		return &TypeDecl {
			name:  declName,
			type_: declType,
		}, nil
	}
}

func TryParseTypeDecl(r *ParseStream) (*TypeDecl, bool) {
	initPos := r.state.pos
	ret, err := ParseTypeDecl(r)
	if err != nil {
		r.Seek(initPos - r.state.pos)
	}
	return ret, (err == nil)
}

func ParseFile(r *ParseStream) []TypeDecl {
	ret := []TypeDecl{}
	for {
		if _, ok := PeekToken(r, false); !ok {
			return ret
		}
		decl, err := ParseTypeDecl(r)
		if err != nil {
			panic(err)
		}
		ret = append(ret, *decl)
	}
}

func TestAST() {
	// testSrc := ``
	// fs := token.NewFileSet()
	// f, err := parser.ParseFile(fs, "test.go", testSrc, 0)
	// if err != nil {
	// 	panic(err)
	// }
	// ast.Print(fs, f)

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

	// decls := []TypeDecl{
	// 	TypeDecl{
	// 		name: "Block",
	// 		type_: ProdType{
	// 			fields: []Field{
	// 				Field{
	// 					fieldName: "MinerAddress",
	// 					fieldType: NamedType{"Address"},
	// 					internal:  false,
	// 				},
	// 				Field{
	// 					fieldName: "TestField1",
	// 					fieldType: NamedType{"byte"},
	// 					internal:  true,
	// 				},
	// 			},
	// 			methods: []Method {
	// 				Method {
	// 					methodName: "SerializeSigned",
	// 					methodType: FunType {
	// 						retType: ArrayType {
	// 							elementType: NamedType{"byte"},
	// 						},
	// 						argTypes: []Type{},
	// 					},
	// 				},
	// 			},
	// 		},
	// 	},
	// }

	// goDecls := GenGoTypeDecls(origDecls)
	// goMod := GenGoMod(goDecls)
	// CheckErr(printer.Fprint(os.Stdout, goMod.astFileSet, goMod.astFile))

	// file, err := os.Open("codeGen/test.mgo")

	pathRoot := "../../../src/subsystems"
	fileSubPaths := []string {
		// "blockchain/components/block_receiver.id",
		// "blockchain/components/block_producer.id",
		// "blockchain/components/block_propagator.id",
		// "blockchain/elections.id",
		// "blockchain/statetree.id",
		// "blockchain/blockchain.id",
		// "blockchain/block.id",

		"clock/clock.id",

		// "datatransfer/datatransfer.id",

		// "libp2p/libp2p.id",

		// "messagepool/message.id",
		// "messagepool/messagepool.id",

		// "payments/payments.id",

		// "repository/repository.id",

		// "retrieval_market/retrieval_market.id",

		// "sectorstorage/components/sector_sealer.id",
		// "sectorstorage/components/sector_store.id",
		// "sectorstorage/components/sector_builder.id",
		// "sectorstorage/components/sectorstoragesubsystem.id",

		// "storage_market/storage_market.id",

		// "storagemining/components/mining_scheduler.id",
		// "storagemining/components/post_prover.id",
		// "storagemining/actors/storagemineractor.id",
		// "storagemining/mining.id",

		// "storagepowerconsensus/components/expectedconsensus.id",
		// "storagepowerconsensus/actors/storagepoweractor.id",

		"vm/components/actor.id",
		"vm/components/vm_interpreter.id",
		"vm/actors/account_actor.id",
		"vm/actors/init_actor.id",
		"vm/vm.id",
	}

	for _, fileSubPath := range fileSubPaths {
		filePath := pathRoot + "/" + fileSubPath
		fmt.Printf(" ===== Parsing: %v\n\n", filePath)
		file, err := os.Open(filePath)
		CheckErr(err)
		r := ParseStream {
			rs: file,
			state: ParseStreamState {
				pos: 0,
			},
			lineMap: []LineInfo{},
			buffer: bytes.NewBuffer([]byte{}),
		}
		decls := ParseFile(&r)
		goDecls := GenGoTypeDecls(decls)
		goMod := GenGoMod(goDecls)
		CheckErr(printer.Fprint(os.Stdout, goMod.astFileSet, goMod.astFile))
	}
}
