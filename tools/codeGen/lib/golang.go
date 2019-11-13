package codeGen

import (
	"fmt"
	"go/ast"
	"go/token"
)

type GoNode interface {
	implements_GoNode()
}

func (GoTypeDecl) implements_GoNode()        {}
func (GoPackageDecl) implements_GoNode()     {}
func (GoImportDecl) implements_GoNode()      {}
func (GoImportMultiDecl) implements_GoNode() {}
func (GoEnumDecl) implements_GoNode()        {}
func (GoFunDecl) implements_GoNode()         {}
func (GoFunType) implements_GoNode()         {}
func (GoPtrType) implements_GoNode()         {}
func (GoTupleType) implements_GoNode()       {}
func (GoStmtReturn) implements_GoNode()      {}
func (GoStmtExpr) implements_GoNode()        {}
func (GoExprDot) implements_GoNode()         {}
func (GoExprEq) implements_GoNode()          {}
func (GoExprCast) implements_GoNode()        {}
func (GoExprCall) implements_GoNode()        {}
func (GoExprStruct) implements_GoNode()      {}
func (GoExprAddrOf) implements_GoNode()      {}
func (GoExprLitNil) implements_GoNode()      {}
func (GoExprLitStr) implements_GoNode()      {}
func (GoField) implements_GoNode()           {}
func (GoProdType) implements_GoNode()        {}
func (GoArrayType) implements_GoNode()       {}
func (GoMapType) implements_GoNode()         {}
func (GoIdent) implements_GoNode()           {}

type GoTypeDecl struct {
	name  string
	type_ GoNode
}

type GoPackageDecl struct {
	name string
}

type GoImportDecl struct {
	name string
	path string
}

type GoImportMultiDecl struct {
	names []string
	paths []string
}

type GoEnumDecl struct {
	name      string
	caseNames []string
}

type GoFunDecl struct {
	receiverVar  *GoIdent
	receiverType GoNode
	funName      string
	funType      GoFunType
	funArgs      []GoNode
	funBody      []GoNode
}

type GoProdTypeCase = int

const (
	GoProdTypeCase_Interface GoProdTypeCase = 0
	GoProdTypeCase_Struct    GoProdTypeCase = 1
)

type GoFunType struct {
	args    []GoField
	retType GoNode
}

type GoPtrType struct {
	targetType GoNode
}

type GoTupleType struct {
	elementTypes []GoNode
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

type GoExprLitNil struct{}

type GoExprLitStr struct {
	str string
}

type GoField struct {
	fieldName *string
	fieldType GoNode // TODO: rename
}

type GoProdType struct {
	typeCase GoProdTypeCase
	fields   []GoField
}

type GoArrayType struct {
	elementType GoNode
}

type GoMapType struct {
	keyType   GoNode
	valueType GoNode
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

	case GoImportDecl:
		xr := x.(GoImportDecl)
		return &ast.GenDecl{
			Tok: token.IMPORT,
			Specs: []ast.Spec{
				&ast.ImportSpec{
					Name: ast.NewIdent(xr.name),
					Path: &ast.BasicLit{
						Kind:  token.STRING,
						Value: xr.path,
					},
				},
			},
		}

	case GoImportMultiDecl:
		xr := x.(GoImportMultiDecl)
		importSpecs := []ast.Spec{}
		n := len(xr.names)
		Assert(len(xr.paths) == n)
		for i, name := range xr.names {
			path := xr.paths[i]
			importSpecs = append(importSpecs, RefGolangASTImportSpec(ast.ImportSpec{
				Name: ast.NewIdent(name),
				Path: RefGolangASTBasicLit(ast.BasicLit{
					Kind:  token.STRING,
					Value: path,
				}),
			}))
		}
		return &ast.GenDecl{
			Tok:   token.IMPORT,
			Specs: importSpecs,
		}

	case GoEnumDecl:
		xr := x.(GoEnumDecl)
		caseTypeNames := []*ast.Ident{}
		caseTypeValues := []ast.Expr{}
		for i, caseName := range xr.caseNames {
			caseTypeNames = append(caseTypeNames, ast.NewIdent(xr.name+"_"+caseName))
			caseTypeValue := RefGolangASTBasicLit(ast.BasicLit{
				Kind:  token.INT,
				Value: fmt.Sprintf("%v", i+1),
			})
			caseTypeValues = append(caseTypeValues, caseTypeValue)
		}
		return &ast.GenDecl{
			Tok: token.CONST,
			Specs: []ast.Spec{
				&ast.ValueSpec{
					Names:  caseTypeNames,
					Type:   ast.NewIdent(xr.name),
					Values: caseTypeValues,
				},
			},
		}

	case GoFunDecl:
		xr := x.(GoFunDecl)
		var goRecv *ast.FieldList = nil
		if xr.receiverType != nil {
			goRecv = &ast.FieldList{
				List: []*ast.Field{
					&ast.Field{
						Names: []*ast.Ident{
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
		goFunBody := &ast.BlockStmt{
			List: goFunStmts,
		}
		return &ast.FuncDecl{
			Recv: goRecv,
			Name: ast.NewIdent(xr.funName),
			Type: goFunType,
			Body: goFunBody,
		}

	case GoProdType:
		xr := x.(GoProdType)
		fields := []*ast.Field{}
		for _, field := range xr.fields {
			field := RefGolangASTField(ast.Field{
				Names: []*ast.Ident{
					ast.NewIdent(DerefCheckString(field.fieldName)),
				},
				Type: GenAST(field.fieldType).(ast.Expr),
			})
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

	case GoMapType:
		xr := x.(GoMapType)
		return &ast.MapType{
			Key:   GenAST(xr.keyType).(ast.Expr),
			Value: GenAST(xr.valueType).(ast.Expr),
		}

	case GoFunType:
		xr := x.(GoFunType)
		goParamFields := []*ast.Field{}
		for _, arg := range xr.args {
			goNames := []*ast.Ident{}
			if arg.fieldName != nil {
				goNames = append(goNames, ast.NewIdent(DerefCheckString(arg.fieldName)))
			}
			goParamFields = append(goParamFields, RefGolangASTField(ast.Field{
				Names: goNames,
				Type:  GenAST(arg.fieldType).(ast.Expr),
			}))
		}

		goResultFields := []*ast.Field{}
		switch xr.retType.(type) {
		case GoTupleType:
			rr := xr.retType.(GoTupleType)
			for _, ri := range rr.elementTypes {
				goResultFields = append(goResultFields, &ast.Field{
					Names: nil,
					Type:  GenAST(ri).(ast.Expr),
				})
			}
		default:
			goResultFields = append(goResultFields, &ast.Field{
				Names: nil,
				Type:  GenAST(xr.retType).(ast.Expr),
			})
		}

		return &ast.FuncType{
			Params: &ast.FieldList{
				List: goParamFields,
			},
			Results: &ast.FieldList{
				List: goResultFields,
			},
		}

	case GoPtrType:
		xr := x.(GoPtrType)
		return &ast.StarExpr{
			X: GenAST(xr.targetType).(ast.Expr),
		}

	case GoTupleType:
		panic("tuple type only permitted as function return")

	case GoStmtReturn:
		xr := x.(GoStmtReturn)
		return &ast.ReturnStmt{
			Results: []ast.Expr{
				GenAST(xr.value).(ast.Expr),
			},
		}

	case GoExprDot:
		xr := x.(GoExprDot)
		return &ast.SelectorExpr{
			X:   GenAST(xr.value).(ast.Expr),
			Sel: ast.NewIdent(xr.fieldName),
		}

	case GoExprCall:
		xr := x.(GoExprCall)
		goArgs := []ast.Expr{}
		for _, f := range xr.args {
			goArgs = append(goArgs, GenAST(f).(ast.Expr))
		}
		return &ast.CallExpr{
			Fun:  GenAST(xr.f).(ast.Expr),
			Args: goArgs,
		}

	case GoStmtExpr:
		xr := x.(GoStmtExpr)
		return &ast.ExprStmt{
			X: GenAST(xr.expr).(ast.Expr),
		}

	case GoExprLitNil:
		return &ast.BasicLit{
			Kind:  token.STRING,
			Value: "nil",
		}

	case GoExprLitStr:
		xr := x.(GoExprLitStr)
		return &ast.BasicLit{
			Kind:  token.STRING,
			Value: fmt.Sprintf("\"%v\"", xr.str),
		}

	case GoExprStruct:
		xr := x.(GoExprStruct)
		goFields := []ast.Expr{}
		for _, field := range xr.fields {
			goField := RefGolangASTKeyValueExpr(ast.KeyValueExpr{
				Key:   ast.NewIdent(DerefCheckString(field.fieldName)),
				Value: GenAST(field.fieldType).(ast.Expr),
			})
			goFields = append(goFields, goField)
		}
		return &ast.CompositeLit{
			Type: GenAST(xr.type_).(ast.Expr),
			Elts: goFields,
		}

	case GoExprAddrOf:
		xr := x.(GoExprAddrOf)
		return &ast.UnaryExpr{
			Op: token.AND,
			X:  GenAST(xr.target).(ast.Expr),
		}

	case GoExprEq:
		xr := x.(GoExprEq)
		return &ast.BinaryExpr{
			Op: token.EQL,
			X:  GenAST(xr.lhs).(ast.Expr),
			Y:  GenAST(xr.rhs).(ast.Expr),
		}

	case GoExprCast:
		xr := x.(GoExprCast)
		return &ast.TypeAssertExpr{
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

func GoTypeByteArray() GoNode {
	return GoArrayType{
		elementType: GoIdent{
			name: "byte",
		},
	}
}

func GoTypeAny() GoNode {
	return GoProdType{
		typeCase: GoProdTypeCase_Interface,
		fields:   []GoField{},
	}
}

func GenGoIdentityBody(receiverVar GoNode) []GoNode {
	ret := []GoNode{
		GoStmtReturn{
			receiverVar,
		},
	}
	return ret
}

func GenGoDerefCacheBody(receiverVar GoNode) []GoNode {
	// FIXME: look up hash in store if cached_impl is nil
	ret := []GoNode{
		GoStmtReturn{
			GoExprDot{
				value:     receiverVar,
				fieldName: "cached_impl",
			},
		},
	}
	return ret
}

func GenGoAccessorBody(receiverVar GoNode, fieldName string) []GoNode {
	ret := []GoNode{
		GoStmtReturn{
			GoExprDot{
				value:     receiverVar,
				fieldName: fieldName,
			},
		},
	}
	return ret
}

func GenGoDerefAccessorBody(receiverVar GoNode, fieldName string) []GoNode {
	ret := []GoNode{
		GoStmtReturn{
			GoExprDot{
				value: GoExprCall{
					f:    GoExprDot{value: receiverVar, fieldName: "Impl"},
					args: []GoNode{},
				},
				fieldName: fieldName,
			},
		},
	}
	return ret
}

func GenGoMethodCall(obj GoNode, methodName string, args []GoNode) GoNode {
	ret := GoExprCall{
		f: GoExprDot{
			value:     obj,
			fieldName: methodName,
		},
		args: args,
	}
	return ret
}

func GenGoPanicTodoBody() []GoNode {
	ret := []GoNode{
		GoStmtExpr{
			expr: GoExprCall{
				f:    GoIdent{name: "panic"},
				args: []GoNode{GoExprLitStr{str: "TODO"}},
			},
		},
	}
	return ret
}

type GoMod struct {
	astFileSet *token.FileSet
	astFile    *ast.File
}
