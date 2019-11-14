package codeGen

type Decl_Case = int

const (
	Decl_Case_Type    Decl_Case = 1
	Decl_Case_Package Decl_Case = 2
	Decl_Case_Import  Decl_Case = 3
)

type Module struct {
	entries      []Entry
	parseFmtInfo *ParseFmtInfo
}

func (mod Module) Decls() []Decl {
	ret := []Decl{}
	for _, entry := range mod.entries {
		if entry.case_ == Entry_Case_Decl {
			ret = append(ret, entry.value.(Decl))
		}
	}
	return ret
}

type Decl interface {
	Case() Decl_Case
	Name() string
}

type TypeDecl struct {
	name         string
	type_        Type
	parseFmtInfo *ParseFmtInfo
}

type PackageDecl struct {
	name         string
	parseFmtInfo *ParseFmtInfo
}

type ImportDecl struct {
	name         string
	path         string
	parseFmtInfo *ParseFmtInfo
}

func (x TypeDecl) Name() string {
	return x.name
}

func (x PackageDecl) Name() string {
	return x.name
}

func (x ImportDecl) Name() string {
	return x.name
}

func (TypeDecl) Case() Decl_Case {
	return Decl_Case_Type
}

func (PackageDecl) Case() Decl_Case {
	return Decl_Case_Package
}

func (ImportDecl) Case() Decl_Case {
	return Decl_Case_Import
}

type Type interface {
	Case() Type_Case
}

type Type_Case = int

const (
	Type_Case_NamedType  Type_Case = 1
	Type_Case_AlgType    Type_Case = 2
	Type_Case_ArrayType  Type_Case = 3
	Type_Case_FunType    Type_Case = 4
	Type_Case_RefType    Type_Case = 5
	Type_Case_OptionType Type_Case = 6
	Type_Case_MapType    Type_Case = 7
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

type AlgSort = int

const (
	AlgSort_Prod AlgSort = 1
	AlgSort_Sum  AlgSort = 2
)

type Field struct {
	fieldName     *string
	fieldType     Type
	attributeList []string
	parseFmtInfo  *ParseFmtInfo
}

type Method struct {
	methodName    string
	methodArgs    []Entry
	argsFmtInfo   *ParseFmtInfo
	methodRetType Type
	attributeList []string
	parseFmtInfo  *ParseFmtInfo
}

func DSLTrivialStruct() Type {
	return RefAlgType(AlgType{
		sort:          AlgSort_Prod,
		entries:       []Entry{},
		attributeList: []string{},
		parseFmtInfo:  nil,
		isInterface:   false,
		isEnum:        false,
	})
}

func DSLTypeIsTrivialStruct(type_ Type) bool {
	if type_.Case() != Type_Case_AlgType {
		return false
	}
	xr := type_.(*AlgType)
	if xr.sort != AlgSort_Prod {
		return false
	}
	if len(xr.attributeList) > 0 {
		return false
	}
	return len(xr.entries) == 0
}

func DSLTrivialStructField(fieldName *string, info *ParseFmtInfo) *Field {
	return RefField(Field{
		fieldName:     fieldName,
		fieldType:     DSLTrivialStruct(),
		attributeList: []string{},
		parseFmtInfo:  info,
	})
}

func ExtractFieldNames(fields []Field) []string {
	ret := []string{}
	for _, field := range fields {
		Assert(field.fieldName != nil)
		ret = append(ret, *field.fieldName)
	}
	return ret
}

func (method Method) MethodType() *FunType {
	args := []Field{}
	for _, arg := range method.methodArgs {
		switch arg.case_ {
		case Entry_Case_Field:
			args = append(args, arg.value.(Field))
		case Entry_Case_Comment:
			break
		case Entry_Case_Empty:
			break
		default:
			Assert(false)
		}
	}
	retType := method.methodRetType
	return RefFunType(FunType{
		args:    args,
		retType: retType,
	})
}

type Comment struct {
	isInline    bool
	isBlock     bool
	commentText string
	endPos      int
}

type Entry_Case = int

const (
	Entry_Case_Empty   Entry_Case = 1
	Entry_Case_Field   Entry_Case = 2
	Entry_Case_Method  Entry_Case = 3
	Entry_Case_Comment Entry_Case = 4
	Entry_Case_Decl    Entry_Case = 5
)

type Entry struct {
	case_ Entry_Case
	value interface{}
}

func EntryField(field Field) Entry {
	return Entry{
		case_: Entry_Case_Field,
		value: field,
	}
}

func EntryMethod(method Method) Entry {
	return Entry{
		case_: Entry_Case_Method,
		value: method,
	}
}

func EntryComment(comment Comment) Entry {
	return Entry{
		case_: Entry_Case_Comment,
		value: comment,
	}
}

func EntryDecl(decl Decl) Entry {
	return Entry{
		case_: Entry_Case_Decl,
		value: decl,
	}
}

func EntryEmpty() Entry {
	return Entry{
		case_: Entry_Case_Empty,
		value: nil,
	}
}

func EntryIsInlineComment(entry Entry) bool {
	return entry.case_ == Entry_Case_Comment && entry.value.(Comment).isInline
}

type AlgType struct {
	sort           AlgSort
	entries        []Entry
	entriesFmtInfo *ParseFmtInfo
	attributeList  []string
	parseFmtInfo   *ParseFmtInfo
	isInterface    bool
	isEnum         bool
}

func (x *AlgType) Methods() []Method {
	ret := []Method{}
	for _, entry := range x.entries {
		if entry.case_ == Entry_Case_Method {
			ret = append(ret, entry.value.(Method))
		}
	}
	return ret
}

func (x *AlgType) Fields() []Field {
	ret := []Field{}
	for _, entry := range x.entries {
		if entry.case_ == Entry_Case_Field {
			ret = append(ret, entry.value.(Field))
		}
	}
	return ret
}

type ArrayType struct {
	elementType  Type
	parseFmtInfo *ParseFmtInfo
}

type RefType struct {
	targetType   Type
	parseFmtInfo *ParseFmtInfo
}

type FunType struct {
	args         []Field
	retType      Type
	parseFmtInfo *ParseFmtInfo
}

type NamedType struct {
	name         string
	parseFmtInfo *ParseFmtInfo
}

type OptionType struct {
	valueType    Type
	parseFmtInfo *ParseFmtInfo
}

type MapType struct {
	keyType      Type
	valueType    Type
	parseFmtInfo *ParseFmtInfo
}
