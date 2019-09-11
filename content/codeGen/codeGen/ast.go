package codeGen

import (
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
	Type_Case_ProdType  Type_Case = 1
	Type_Case_ArrayType Type_Case = 2
	Type_Case_RefType   Type_Case = 3
	Type_Case_FunType   Type_Case = 4
	Type_Case_NamedType Type_Case = 5
)

func (ProdType) Case() Type_Case {
	return Type_Case_ProdType
}
func (ArrayType) Case() Type_Case {
	return Type_Case_ArrayType
}
func (RefType) Case() Type_Case {
	return Type_Case_RefType
}
func (FunType) Case() Type_Case {
	return Type_Case_FunType
}
func (NamedType) Case() Type_Case {
	return Type_Case_NamedType
}

type Field struct {
	fieldName string
	fieldType Type
	internal  bool
}

type Method struct {
	methodName string
	methodType *FunType
}

type ProdType struct {
	fields  []Field
	methods []Method
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
	switch x.Case() {
	case Type_Case_ProdType:
		xr := x.(*ProdType)
		HashAccWord(buf, Type_Case_ProdType)
		HashAccWord(buf, Word(len(xr.fields)))
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
		HashAccWord(buf, Word(len(xr.args)))
		for _, arg := range xr.args {
			HashAccString(buf, arg.fieldName)
			HashAccType(buf, arg.fieldType)
		}
		HashAccType(buf, xr.retType)

	case Type_Case_NamedType:
		xr := x.(*NamedType)
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
	concrete  bool
}

func GenGoTypeDecls(decls []TypeDecl) []GoNode {
	ctx := GoGenContext {
		typeMap:  map[string]Type{},
		retDecls: &[]GoNode{},
		retMap:   map[TypeHash]GoNode{},
		concrete: false,
	}

	for _, decl := range decls {
		GenGoTypeDeclAcc(decl.name, decl.type_, ctx)
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

func GenGoTypeDeclAcc(name string, x Type, ctx GoGenContext) GoNode {
	Assert(x != nil)

	switch x.(type) {
	case *ProdType:
		xr := x.(*ProdType)

		ctxConcrete := ctx
		ctxConcrete.concrete = true
	
		interfaceName := name
		interfaceID := GoIdent {name: interfaceName}

		implName := IdToImpl(name)
		implID := GoIdent {name: implName}

		implRefName := IdToImplRef(name)
		implRefID := GoIdent {name: implRefName}

		interfaceFields := []GoField{}
		implFields := []GoField{}
		implRefFields := []GoField{}

		for _, field := range xr.fields {
			implFields = append(implFields, GoField {
				fieldName: field.fieldName,
				fieldType: GenGoTypeAcc(field.fieldType, ctxConcrete),
			})

			interfaceFields = append(interfaceFields, GoField {
				fieldName: field.fieldName,
				fieldType: GoFunType {
					args:     []GoField{},
					retType:  GenGoTypeAcc(field.fieldType, ctx),
				},
			})
		}

		for _, method := range xr.methods {
			interfaceFields = append(interfaceFields, GoField {
				fieldName: method.methodName,
				fieldType: GenGoTypeAcc(method.methodType, ctx),
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
			fieldName: "ContentHash",
			fieldType: GoFunType {
				retType:  GoTypeByteArray(),
				args:     []GoField{},
			},
		})

		implFields = append(implFields, GoField {
			fieldName: "cached_contentHash",
			fieldType: GoTypeByteArray(),
		})

		implRefFields = append(implRefFields, GoField {
			fieldName: "contentHash",
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

		implReceiverVar := GoIdent {name: strings.ToLower(name)}
		implRefReceiverVar := GoIdent {name: strings.ToLower(name)}
		for _, field := range xr.fields {
			baseFieldType := GenGoTypeAcc(field.fieldType, ctx)

			implAccessorDecl := GoFunDecl {
				receiverVar: implReceiverVar,
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
				receiverVar: implRefReceiverVar,
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
			receiverVar: implReceiverVar,
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
			receiverVar: implRefReceiverVar,
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
	case *NamedType:
		panic("TODO")
	default:
		errMsg := fmt.Sprintf("TODO: %v\n", x.Case())
		panic(errMsg)
	}
}

func GenGoTypeAcc(x Type, ctx GoGenContext) GoNode {
	// ret := []GoNode{}
	// idMap := map[Type]string

	if match, ok := ctx.retMap[HashType(x)]; ok {
		return match
	}

	// for _, x := range decls {
	switch x.Case() {
	case Type_Case_ProdType:
		panic("Inline ProdType not supported in this context")

	case Type_Case_ArrayType:
		xr := x.(*ArrayType)
		goElementType := GenGoTypeAcc(xr.elementType, ctx)
		return GoArrayType {
			elementType: goElementType,
		}

	case Type_Case_RefType:
		xr := x.(*RefType)
		goTargetType := GenGoTypeAcc(xr.targetType, ctx)
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
		goRetType := GenGoTypeAcc(xr.retType, ctx)
		goArgs := []GoField{}
		for _, arg := range xr.args {
			goArgType := GenGoTypeAcc(arg.fieldType, ctx)
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

	default:
		panic("TODO")
	}
}

type GoNode interface{}

type GoTypeDecl struct {
	name  string
	type_ GoNode
}

type GoFunDecl struct {
	receiverVar  GoIdent
	receiverType GoNode
	funName      string
	funType      GoFunType
	funArgs      GoNode
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

type GoExprDot struct {
	value     GoNode
	fieldName string
}

type GoExprCall struct {
	f    GoNode
	args []GoNode
}

type GoField struct {
	fieldName string
	fieldType GoNode
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

	case GoFunDecl:
		xr := x.(GoFunDecl)
		var goRecv *ast.FieldList = nil
		if xr.receiverType != nil {
			goRecv = &ast.FieldList {
				List: []*ast.Field {
					&ast.Field {
						Names: []*ast.Ident {
							GenAST(xr.receiverVar).(*ast.Ident),
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

	case GoIdent:
		xr := x.(GoIdent)
		return ast.NewIdent(xr.name)

	default:
		panic("Unknown type for GenAST")
	}
}

const Whitespace = " \t\n"
const Symbols = "(){}[],;|&"

const DebugParser = false

type ParseStream struct {
	rs  io.ReadSeeker
	pos Word
}

func (r *ParseStream) Seek(offset Word) {
	if DebugParser {
		fmt.Printf("Seek: %v\n", offset)
	}
	_, err := r.rs.Seek(int64(offset), io.SeekCurrent)
	CheckErr(err)
	r.pos += offset
}

func ReadCheck(r *ParseStream, lenMax Word) string {
	buf := make([]byte, lenMax)
	n, err := r.rs.Read(buf)
	r.pos += Word(n)
	// fmt.Printf("  n: %v  buf: %v  Pos: %v\n", n, string(buf), r.pos)
	if err == io.EOF {
		Assert(n == 0)
		return string(buf[:n])
	}
	if err != nil {
		if DebugParser {
			fmt.Printf("err: %v\n", err)
		}
		panic("Read error")
	}
	// if Word(n) < lenMin {
	// 	panic("Read error")
	// }
	if Word(n) > lenMax {
		panic("Read error")
	}
	return string(buf[:n])
}

// func ReadExact(r *ParseStream, lenExact Word) string {
// 	return ReadCheck(r, lenExact, lenExact)
// }

// func ParseStringExact(r *ParseStream, x string) {
// 	xCheck := ReadExact(r, Word(len(x)))
// 	if xCheck != x {
// 		var s = "Parse error"
// 		r.err = &s
// 	}
// }

// func ParseStringTokenExact(r *ParseStream, x string) {
// 	StepToken(r)
// 	ParseStringExact(r, x)
// }


func ReadToken(r *ParseStream) (string, error) {
	retToken := []byte{}
	StepToken(r)
	for {
		res := ReadCheck(r, 1)
		if len(res) == 0 {
			if len(retToken) > 0 {
				return string(retToken), nil
			} else {
				return "", io.EOF
			}
		}
		Assert(len(res) == 1)
		if strings.Contains(Whitespace, res) {
			r.Seek(-Word(len(res)))
			return string(retToken), nil
		} else if strings.Contains(Symbols, res) {
			if len(retToken) == 0 {
				retToken = append(retToken, res[0])
			} else {
				r.Seek(-Word(len(res)))
			}
			return string(retToken), nil
			// if DebugParser {
				// fmt.Printf("Read token: %v\n", string(retToken))
			// }
		} else {
			retToken = append(retToken, res...)
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
	ret, err := ReadToken(r)
	if err != nil {
		return "", err
	} else if !IsIdent(ret) {
		msg := fmt.Sprintf("Expected identifier; received: \"%v\"", ret)
		return "", ParseError(msg)
	} else {
		return ret, nil
	}
}

func (err ParseError_S) Error() string {
	return err.msg
}

func ReadTokenCheck(r *ParseStream, validTokens []string) (string, error) {
	tok, err := ReadToken(r)
	if err != nil {
		return "", err
	}
	for _, v := range validTokens {
		if tok == v {
			return tok, nil
		}
	}
	return "", ParseError("Unexpected token")
}

func TryReadTokenCheck(r *ParseStream, validTokens []string) (string, bool) {
	initPos := r.pos
	ret, err := ReadTokenCheck(r, validTokens)
	if err != nil {
		r.Seek(initPos - r.pos)
	}
	return ret, (err == nil)
}

func StepToken(r *ParseStream) {
	for {
		res := ReadCheck(r, 1)
		if len(res) == 1 && strings.Contains(Whitespace, res) {
			continue
		} else {
			r.Seek(-Word(len(res)))
			return
		}
	}
}

func ParseField(r *ParseStream) (*Field, error) {
	fieldName, err := ParseIdent(r)

	if DebugParser {
		fmt.Printf(" >>>>> IDENT: %v %v\n", fieldName, err)
	}

	if err != nil {
		return nil, err
	}
	
	fieldType, err := ParseType(r)
	
	if DebugParser {
		fmt.Printf(" >>>>> ParseField ParseType: %v %v\n", fieldType, err)
	}
	
	if err != nil {
		return nil, err
	}
	return &Field {
		fieldName: fieldName,
		fieldType: fieldType,
		internal: false,
	}, nil
}

func TryParseField(r *ParseStream) (*Field, bool) {
	initPos := r.pos

	if DebugParser {
		fmt.Printf("TryParseField initPos: %v\n", initPos)
	}

	ret, err := ParseField(r)
	if err != nil {
		r.Seek(initPos - r.pos)
	}

	if DebugParser {
		fmt.Printf("TryParseField final pos: (%v, %v): %v\n", ret, err, r.pos)
	}

	return ret, (err == nil)
}

func ParseMethod(r *ParseStream) (*Method, error) {
	methodName, err := ParseIdent(r)

	if DebugParser {
		fmt.Printf("ParseMethod start: name: %v (err=%v)\n", methodName, err)
	}

	if err != nil {
		return nil, err
	}

	_, err = ReadTokenCheck(r, []string{"("})
	if err != nil {
		return nil, err
	}

	methodArgs := []Field{}
	for {
		if field, ok := TryParseField(r); ok {
			methodArgs = append(methodArgs, *field)
		} else {
			_, err := ReadTokenCheck(r, []string{")"})
			if err != nil {
				return nil, ParseError("Error parsing method argument")
			} else {
				break
			}
		}
		if _, ok := TryReadTokenCheck(r, []string{","}); !ok {
			_, err := ReadTokenCheck(r, []string{")"})
			if err != nil {
				return nil, err
			} else {
				break
			}
		}
	}

	methodRetType, err := ParseType(r)
	if err != nil {
		return nil, err
	}

	methodType := FunType {
		args:    methodArgs,
		retType: methodRetType,
	}

	return &Method {
		methodName: methodName,
		methodType: &methodType,
	}, nil
}

func TryParseMethod(r *ParseStream) (*Method, bool) {
	initPos := r.pos
	ret, err := ParseMethod(r)
	if err != nil {
		r.Seek(initPos - r.pos)
	}
	return ret, (err == nil)
}

func ParseType(r *ParseStream) (Type, error) {
	tok, err := ReadToken(r)
	if err != nil {
		return nil, err
	}

	if DebugParser {
		fmt.Printf("ParseType token: \"%v\"\n", tok)
	}

	switch {
	case tok == "struct":
		var fields []Field
		var methods []Method
		_, err := ReadTokenCheck(r, []string{"{"})
		if err != nil {
			return nil, err
		}
		for {
			if field, ok := TryParseField(r); ok {
				fields = append(fields, *field)
			} else {
				if method, ok := TryParseMethod(r); ok {
					methods = append(methods, *method)
				} else {
					break
				}
			}
		}
		_, err = ReadTokenCheck(r, []string{"}"})
		if err != nil {
			return nil, err
		}
		return &ProdType {
			fields: fields,
			methods: methods,
		}, nil

	case tok == "[":
		elementType, err := ParseType(r)
		if err != nil {
			return nil, err
		}
		_, err = ReadTokenCheck(r, []string{"]"})
		if err != nil {
			return nil, err
		}
		return &ArrayType {
			elementType: elementType,
		}, nil

	case tok == "&":
		targetType, err := ParseType(r)
		if err != nil {
			return nil, err
		}
		return &RefType {
			targetType: targetType,
		}, nil

	default:
		if IsIdent(tok) {		
			return &NamedType {name: tok}, nil
		} else {
			return nil, ParseError(fmt.Sprintf("Expected type; received \"%v\"", tok))
		}
	}
}

func ParseTypeDecl(r *ParseStream) (*TypeDecl, error) {
	_, err := ReadTokenCheck(r, []string{"type"})
	if err != nil {
		return nil, err
	}
	declName, err := ReadToken(r)
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
	initPos := r.pos
	ret, err := ParseTypeDecl(r)
	if err != nil {
		r.Seek(initPos - r.pos)
	}
	return ret, (err == nil)
}

func ParseTypeDecls(r *ParseStream) []TypeDecl {
	ret := []TypeDecl{}
	for {
		if decl, ok := TryParseTypeDecl(r); ok {
			ret = append(ret, *decl)
		} else {
			return ret
		}
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

	file, err := os.Open("codeGen/test.mgo")
	CheckErr(err)
	r := ParseStream {
		rs: file,
		pos: 0,
	}
	decls := ParseTypeDecls(&r)
	goDecls := GenGoTypeDecls(decls)
	goMod := GenGoMod(goDecls)
	CheckErr(printer.Fprint(os.Stdout, goMod.astFileSet, goMod.astFile))
}
