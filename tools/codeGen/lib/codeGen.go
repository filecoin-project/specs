package codeGen

import (
	"fmt"
	"go/ast"
	"go/token"
	"os"
	"strings"
)

func GenGoMod(goDecls []GoNode, packageName string) GoMod {
	var ret GoMod

	var astDecls = []ast.Decl{}

	for _, goDecl := range goDecls {
		switch goDecl.(type) {
		case GoPackageDecl:
			astDecls = append(astDecls, GenAST(goDecl).(ast.Decl))
		}
	}

	importNames := []string{}
	importPaths := []string{}
	for _, goDecl := range goDecls {
		switch goDecl.(type) {
		case GoImportDecl:
			xr := goDecl.(GoImportDecl)
			importNames = append(importNames, xr.name)
			importPaths = append(importPaths, xr.path)
		}
	}

	Assert(len(importNames) == len(importPaths))
	if len(importNames) > 0 {
		goImportDecl := GoImportMultiDecl{
			names: importNames,
			paths: importPaths,
		}
		astDecls = append(astDecls, GenAST(goImportDecl).(ast.Decl))
	}

	for _, goDecl := range goDecls {
		switch goDecl.(type) {
		case GoImportDecl:
			break
		case GoPackageDecl:
			break
		default:
			astDecls = append(astDecls, GenAST(goDecl).(ast.Decl))
		}
	}

	ret.astFileSet = token.NewFileSet()
	ret.astFile = RefGolangASTFile(ast.File{
		Name:  ast.NewIdent(packageName),
		Decls: astDecls,
	})

	return ret
}

type GoGenContext struct {
	retDecls *[]GoNode
	declMap  map[string]GoNode
	typeMap  map[Type]GoNode
	tokens   []string
	usesUtil *[]bool
}

func (ctx GoGenContext) Extend(token string) GoGenContext {
	ret := ctx
	ret.tokens = append(ret.tokens, token)
	return ret
}

func GenGoDecls(topLevelEntries []Entry) []GoNode {
	ctx := GoGenContext{
		retDecls: &[]GoNode{},
		declMap:  map[string]GoNode{},
		typeMap:  map[Type]GoNode{},
		tokens:   []string{},
		usesUtil: &[]bool{false},
	}

	for _, entry := range topLevelEntries {
		switch entry.case_ {
		case Entry_Case_Decl:
			decl := entry.value.(Decl)
			switch decl.Case() {
			case Decl_Case_Type:
				xr := decl.(*TypeDecl)
				GenGoTypeDeclAcc(xr.name, xr.type_, ctx.Extend(xr.name), false)
			case Decl_Case_Import:
				xr := decl.(*ImportDecl)
				GenGoImportDeclAcc(*xr, ctx)
			case Decl_Case_Package:
				xr := decl.(*PackageDecl)
				GenGoPackageDeclAcc(*xr, ctx)
			default:
				panic("Unhandled case")
			}

		case Entry_Case_Empty:
			// skip

		case Entry_Case_Comment:
			// skip

		default:
			panic(fmt.Sprintf("Unhandled case: %v", entry.case_))
		}
	}

	if (*ctx.usesUtil)[0] {
		GenGoImportDeclAcc(ImportDecl{
			name: "util",
			path: "github.com/filecoin-project/specs/util",
		}, ctx)
	}

	return *ctx.retDecls
}

func IdToImpl(name string) string {
	return name + "_I"
}

func IdToImplRef(name string) string {
	return name + "_R"
}

func GoMethodToFieldName(methodName string) string {
	return methodName + "_"
}

func GoTypeToIdent(typeName string) GoIdent {
	return GoIdent{name: strings.ToLower(typeName)[0:1]}
}

func TranslateGoIdent(name string, ctx GoGenContext) GoIdent {
	ret := name
	utilNames := []string{
		"Assert",
		"BigInt",
		"Bytes",
		"Float",
		"Int",
		"Serialization",
		"T",
		"UInt",
		"UInts",
		"UVarint",
	}
	for _, utilName := range utilNames {
		if name == utilName {
			ret = "util." + name
			*ctx.usesUtil = []bool{true}
			break
		}
	}
	return GoIdent{name: ret}
}

func GenGoImportDeclAcc(decl ImportDecl, ctx GoGenContext) GoNode {
	goImportDecl := GoImportDecl{
		name: decl.name,
		path: "\"" + decl.path + "\"",
	}
	*ctx.retDecls = append(*ctx.retDecls, goImportDecl)
	return goImportDecl
}

func GenGoPackageDeclAcc(decl PackageDecl, ctx GoGenContext) GoNode {
	goPackageDecl := GoPackageDecl{
		name: decl.name,
	}
	*ctx.retDecls = append(*ctx.retDecls, goPackageDecl)
	return goPackageDecl
}

func GenGoTypeDeclAcc(name string, x Type, ctx GoGenContext, declAlias bool) (ret GoNode) {
	Assert(x != nil)

	if t, ok := ctx.declMap[name]; ok {
		ret = t
		return
	}

	if x.Case() == Type_Case_FunType {
		panic("TODO")
	}

	ret = GenGoTypeAcc(x, ctx)
	retDecl := GoTypeDecl{
		name:      name,
		type_:     ret,
		declAlias: declAlias,
	}

	if t, ok := ctx.declMap[name]; ok {
		ret = t
	} else {
		ctx.declMap[name] = ret
		*ctx.retDecls = append(*ctx.retDecls, retDecl)
	}
	return
}

func GenGoAlgTypeSerializers(x Type, ctx GoGenContext, name string, interfaceID GoNode) {
	serializeDecl := GoFunDecl{
		receiverVar:  nil,
		receiverType: nil,
		funName:      "Serialize_" + name,
		funType: GoFunType{
			args: []GoField{
				GoField{
					fieldName: nil,
					fieldType: interfaceID,
				},
			},
			retType: TranslateGoIdent("Serialization", ctx),
		},
		funArgs: []GoNode{GoIdent{name: "x"}},
		funBody: GenGoPanicTodoBody(),
	}

	deserializeDecl := GoFunDecl{
		receiverVar:  nil,
		receiverType: nil,
		funName:      "Deserialize_" + name,
		funType: GoFunType{
			args: []GoField{
				GoField{
					fieldName: nil,
					fieldType: TranslateGoIdent("Serialization", ctx),
				},
			},
			retType: GoTupleType{
				elementTypes: []GoNode{
					interfaceID,
					GoIdent{name: "error"},
				},
			},
		},
		funArgs: []GoNode{GoIdent{name: "x"}},
		funBody: GenGoPanicTodoBody(),
	}

	*ctx.retDecls = append(*ctx.retDecls, serializeDecl)
	*ctx.retDecls = append(*ctx.retDecls, deserializeDecl)
}

func GenGoTypeAcc(x Type, ctx GoGenContext) (ret GoNode) {
	if match, ok := ctx.typeMap[x]; ok {
		return match
	}
	defer func() { ctx.typeMap[x] = ret }()

	switch x.Case() {
	case Type_Case_AlgType:
		name := strings.Join(ctx.tokens, "_")

		xr := x.(*AlgType)

		interfaceName := name
		interfaceID := GoIdent{name: interfaceName}

		var caseTypeName string
		var caseTypeID GoNode = nil

		if xr.sort == AlgSort_Sum {
			caseNames := []string{}

			if xr.isEnum {
				caseTypeName = name
			} else {
				caseTypeName = name + "_Case"
			}
			caseTypeID = GoIdent{caseTypeName}

			*ctx.retDecls = append(*ctx.retDecls, GoTypeDecl{
				name:      caseTypeName,
				type_:     TranslateGoIdent("UVarint", ctx),
				declAlias: false,
			})

			for _, field := range xr.Fields() {
				Assert(field.fieldName != nil)
				fieldName := *field.fieldName
				caseNames = append(caseNames, fieldName)
			}

			caseTypeDecl := GoEnumDecl{
				name:      caseTypeName,
				caseNames: caseNames,
			}

			*ctx.retDecls = append(*ctx.retDecls, caseTypeDecl)
		}

		GenGoAlgTypeSerializers(x, ctx, name, interfaceID)

		if xr.isEnum {
			Assert(xr.sort == AlgSort_Sum)
			ctx.declMap[name] = interfaceID
			ret = interfaceID
			break
		}

		implName := IdToImpl(name)
		implID := GoIdent{name: implName}

		implRefName := IdToImplRef(name)
		implRefID := GoIdent{name: implRefName}

		interfaceFields := []GoField{}
		implFields := []GoField{}
		implRefFields := []GoField{}

		if xr.sort == AlgSort_Prod {
			for _, field := range xr.Fields() {
				fieldName := DerefCheckString(field.fieldName)

				implFields = append(implFields, GoField{
					fieldName: RefString(GoMethodToFieldName(fieldName)),
					fieldType: GenGoTypeAcc(field.fieldType, ctx.Extend(fieldName)),
				})

				interfaceFields = append(interfaceFields, GoField{
					fieldName: RefString(fieldName),
					fieldType: GoFunType{
						args:    []GoField{},
						retType: GenGoTypeAcc(field.fieldType, ctx.Extend(fieldName)),
					},
				})
			}
		}

		for _, method := range xr.Methods() {
			interfaceFields = append(interfaceFields, GoField{
				fieldName: RefString(method.methodName),
				fieldType: GenGoTypeAcc(method.MethodType(), ctx.Extend(method.methodName)),
			})
		}

		if !xr.isInterface {
			interfaceFields = append(interfaceFields, GoField{
				fieldName: RefString("Impl"),
				fieldType: GoFunType{
					retType: GoPtrType{targetType: implID},
					args:    []GoField{},
				},
			})
		}

		// TODO: re-enable
		//
		// interfaceFields = append(interfaceFields, GoField {
		// 	fieldName: "CID",
		// 	fieldType: GoFunType {
		// 		retType:  GoTypeByteArray(),
		// 		args:     []GoField{},
		// 	},
		// })

		if xr.sort == AlgSort_Sum {
			for _, field := range xr.Fields() {
				Assert(field.fieldName != nil)
				fieldName := *field.fieldName

				caseWhich := GoIdent{caseTypeName + "_" + fieldName}

				caseInterfaceName := name + "_" + fieldName
				GenGoTypeDeclAcc(caseInterfaceName, field.fieldType, ctx.Extend(fieldName), true)

				caseInterfaceType := GoIdent{caseInterfaceName}
				// caseImplType := GoIdent { IdToImpl(caseInterfaceName) }
				// caseImplPtrType := GoPtrType { targetType: caseImplType }

				interfaceFields = append(interfaceFields, GoField{
					fieldName: RefString("As_" + fieldName),
					fieldType: GoFunType{
						retType: GoIdent{caseInterfaceName},
						args:    []GoField{},
					},
				})

				caseAsDeclArg := GoTypeToIdent(name + "_" + fieldName)
				caseAsDeclBody := []GoNode{
					GoStmtExpr{
						expr: GoExprCall{
							f: TranslateGoIdent("Assert", ctx),
							args: []GoNode{
								GoExprEq{
									GenGoMethodCall(caseAsDeclArg, "Which", []GoNode{}),
									caseWhich,
								},
							},
						},
					},
					GoStmtReturn{
						value: GoExprCast{
							arg:     GoExprDot{value: caseAsDeclArg, fieldName: "rawValue"},
							resType: caseInterfaceType,
						},
					},
				}
				caseAsDecl := GoFunDecl{
					receiverVar:  RefGoIdent(caseAsDeclArg),
					receiverType: GoPtrType{targetType: implID},
					funName:      "As_" + fieldName,
					funType: GoFunType{
						args:    []GoField{},
						retType: caseInterfaceType,
					},
					funArgs: []GoNode{},
					funBody: caseAsDeclBody,
				}

				*ctx.retDecls = append(*ctx.retDecls, caseAsDecl)

				caseNewDeclArg := GoTypeToIdent(name + "_" + fieldName)
				caseNewDeclBody := []GoNode{
					GoStmtReturn{
						GoExprAddrOf{
							target: GoExprStruct{
								type_: implID,
								fields: []GoField{
									GoField{
										fieldName: RefString("cached_cid"),
										fieldType: GoExprLitNil{},
									},
									GoField{
										fieldName: RefString("rawValue"),
										fieldType: caseNewDeclArg,
									},
									GoField{
										fieldName: RefString("which"),
										fieldType: caseWhich,
									},
								},
							},
						},
					},
				}
				caseNewDecl := GoFunDecl{
					receiverVar:  nil,
					receiverType: nil,
					funName:      name + "_Make_" + fieldName,
					funType: GoFunType{
						args: []GoField{
							GoField{
								fieldName: RefString(caseNewDeclArg.name),
								fieldType: GoIdent{caseInterfaceName},
							},
						},
						retType: interfaceID,
					},
					funArgs: []GoNode{caseNewDeclArg},
					funBody: caseNewDeclBody,
				}

				*ctx.retDecls = append(*ctx.retDecls, caseNewDecl)
			}

			interfaceFields = append(interfaceFields, GoField{
				fieldName: RefString("Which"),
				fieldType: GoFunType{
					retType: caseTypeID,
					args:    []GoField{},
				},
			})
		}

		implFields = append(implFields, GoField{
			fieldName: RefString("cached_cid"),
			fieldType: GoTypeByteArray(),
		})

		if xr.sort == AlgSort_Sum {
			implFields = append(implFields, GoField{
				fieldName: RefString("rawValue"),
				fieldType: GoTypeAny(),
			})

			implFields = append(implFields, GoField{
				fieldName: RefString("which"),
				fieldType: caseTypeID,
			})

			whichDeclArg := GoTypeToIdent(name)
			whichDeclBody := []GoNode{
				GoStmtReturn{
					value: GoExprDot{
						value:     whichDeclArg,
						fieldName: "which",
					},
				},
			}
			whichDecl := GoFunDecl{
				receiverVar:  RefGoIdent(whichDeclArg),
				receiverType: GoPtrType{targetType: implID},
				funName:      "Which",
				funType: GoFunType{
					args:    []GoField{},
					retType: caseTypeID,
				},
				funArgs: []GoNode{},
				funBody: whichDeclBody,
			}

			*ctx.retDecls = append(*ctx.retDecls, whichDecl)
		}

		implRefFields = append(implRefFields, GoField{
			fieldName: RefString("cid"),
			fieldType: GoTypeByteArray(),
		})

		cachedObjectField := GoField{
			fieldName: RefString("cached_impl"),
			fieldType: GoPtrType{targetType: implID},
		}
		implRefFields = append(implRefFields, cachedObjectField)

		interfaceDecl := GoTypeDecl{
			name: interfaceName,
			type_: GoProdType{
				typeCase: GoProdTypeCase_Interface,
				fields:   interfaceFields,
			},
			declAlias: false,
		}

		implDecl := GoTypeDecl{
			name: implName,
			type_: GoProdType{
				typeCase: GoProdTypeCase_Struct,
				fields:   implFields,
			},
			declAlias: false,
		}

		implRefDecl := GoTypeDecl{
			name: implRefName,
			type_: GoProdType{
				typeCase: GoProdTypeCase_Struct,
				fields:   implRefFields,
			},
			declAlias: false,
		}

		*ctx.retDecls = append(*ctx.retDecls, interfaceDecl)

		if !xr.isInterface {
			*ctx.retDecls = append(*ctx.retDecls, implDecl)
			*ctx.retDecls = append(*ctx.retDecls, implRefDecl)
		}

		implReceiverVar := GoTypeToIdent(name)
		implRefReceiverVar := GoTypeToIdent(name)

		if xr.sort == AlgSort_Prod && !xr.isInterface {
			for _, field := range xr.Fields() {
				Assert(field.fieldName != nil)
				fieldName := *field.fieldName

				baseFieldType := GenGoTypeAcc(field.fieldType, ctx.Extend(fieldName))

				implAccessorDecl := GoFunDecl{
					receiverVar: RefGoIdent(implReceiverVar),
					receiverType: GoPtrType{
						targetType: implID,
					},
					funName: fieldName,
					funType: GoFunType{
						args:    []GoField{},
						retType: baseFieldType,
					},
					funArgs: []GoNode{},
					funBody: GenGoAccessorBody(implReceiverVar, GoMethodToFieldName(fieldName)),
				}

				implRefAccessorDecl := GoFunDecl{
					receiverVar: RefGoIdent(implRefReceiverVar),
					receiverType: GoPtrType{
						targetType: implRefID,
					},
					funName: fieldName,
					funType: GoFunType{
						args:    []GoField{},
						retType: baseFieldType,
					},
					funArgs: []GoNode{},
					funBody: GenGoDerefAccessorBody(implRefReceiverVar, GoMethodToFieldName(fieldName)),
				}

				*ctx.retDecls = append(*ctx.retDecls, implAccessorDecl)
				*ctx.retDecls = append(*ctx.retDecls, implRefAccessorDecl)
			}
		}

		implImplDecl := GoFunDecl{
			receiverVar:  RefGoIdent(implReceiverVar),
			receiverType: GoPtrType{targetType: implID},
			funName:      "Impl",
			funType: GoFunType{
				args:    []GoField{},
				retType: GoPtrType{targetType: implID},
			},
			funArgs: []GoNode{},
			funBody: GenGoIdentityBody(implReceiverVar),
		}

		implRefImplDecl := GoFunDecl{
			receiverVar:  RefGoIdent(implRefReceiverVar),
			receiverType: GoPtrType{targetType: implRefID},
			funName:      "Impl",
			funType: GoFunType{
				args:    []GoField{},
				retType: GoPtrType{targetType: implID},
			},
			funArgs: []GoNode{},
			funBody: GenGoDerefCacheBody(implReceiverVar),
		}

		if !xr.isInterface {
			*ctx.retDecls = append(*ctx.retDecls, implImplDecl)
			*ctx.retDecls = append(*ctx.retDecls, implRefImplDecl)
		}

		ctx.declMap[name] = interfaceID

		ret = interfaceID

	case Type_Case_ArrayType:
		xr := x.(*ArrayType)
		goElementType := GenGoTypeAcc(xr.elementType, ctx.Extend("ArrayElement"))
		ret = GoArrayType{
			elementType: goElementType,
		}

	case Type_Case_RefType:
		xr := x.(*RefType)
		// TODO: check that target type is hashable
		// goTargetType := GenGoTypeAcc(xr.targetType, ctx.Extend("RefTarget"))
		goTargetType := GenGoTypeAcc(xr.targetType, ctx)
		ret = goTargetType

	case Type_Case_FunType:
		xr := x.(*FunType)
		goRetType := GenGoTypeAcc(xr.retType, ctx.Extend("FunRet"))
		goArgs := []GoField{}
		for i, arg := range xr.args {
			goArgType := GenGoTypeAcc(arg.fieldType, ctx.Extend(fmt.Sprintf("FunArg%v", i)))
			goArg := GoField{
				fieldName: arg.fieldName,
				fieldType: goArgType,
			}
			goArgs = append(goArgs, goArg)
		}
		ret = GoFunType{
			retType: goRetType,
			args:    goArgs,
		}

	case Type_Case_NamedType:
		xr := x.(*NamedType)
		ret = TranslateGoIdent(xr.name, ctx)

	case Type_Case_OptionType:
		xr := x.(*OptionType)
		ret = GenGoTypeAcc(RefAlgType(AlgType{
			sort: AlgSort_Sum,
			entries: []Entry{
				EntryField(Field{
					fieldName:     RefString("Some"),
					fieldType:     xr.valueType,
					attributeList: []string{},
				}),
				EntryField(Field{
					fieldName: RefString("None"),
					fieldType: RefAlgType(AlgType{
						sort:          AlgSort_Prod,
						entries:       []Entry{},
						attributeList: []string{},
						parseFmtInfo:  nil,
						isInterface:   false,
						isEnum:        false,
					}),
					attributeList: []string{},
				}),
			},
			attributeList: []string{},
			parseFmtInfo:  nil,
			isInterface:   false,
		}), ctx)

	case Type_Case_MapType:
		xr := x.(*MapType)
		goKeyType := GenGoTypeAcc(xr.keyType, ctx.Extend("MapKey"))
		goValueType := GenGoTypeAcc(xr.valueType, ctx.Extend("MapValue"))
		ret = GoMapType{keyType: goKeyType, valueType: goValueType}

	default:
		panic("TODO")
	}

	return
}

func GenGoModFromFile(file *os.File, packageName string) GoMod {
	mod := ParseDSLModuleFromFile(file)
	goDecls := GenGoDecls(mod.entries)
	goMod := GenGoMod(goDecls, packageName)
	return goMod
}
