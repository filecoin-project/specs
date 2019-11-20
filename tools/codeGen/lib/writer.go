package codeGen

import (
	"bytes"
	"fmt"
	"go/printer"
	"io"
	"os"
	"strings"
)

func StrFmtLen(s string) IntOption {
	if strings.Contains(s, "\n") {
		return IntOptionNone()
	} else {
		return IntOptionSome(len(s))
	}
}

func WriteDSLTypeFmtLen(type_ Type, ctx WriteDSLContext) IntOption {
	buf := bytes.NewBuffer([]byte{})
	WriteDSLType(buf, type_, ctx)
	return StrFmtLen(buf.String())
}

func WriteDSLDecl(dst io.Writer, decl Decl, ctx WriteDSLContext) {
	switch decl.Case() {
	case Decl_Case_Package:
		xr := decl.(*PackageDecl)
		fmt.Fprintf(dst, "package %s", xr.name)

	case Decl_Case_Import:
		xr := decl.(*ImportDecl)
		fmt.Fprintf(dst, "import %s \"%s\"", xr.name, xr.path)

	case Decl_Case_Type:
		xr := decl.(*TypeDecl)
		fmt.Fprintf(dst, "type %s ", xr.name)
		WriteDSLType(dst, xr.type_, ctx)
	}
}

func WriteDSLContextInit() WriteDSLContext {
	return WriteDSLContext{
		indent:            0,
		useNewlines:       true,
		useFieldAlignment: false,
		alignment:         WriteDSLAlignmentNone(),
		isEnum:            false,
	}
}

func WriteDSLModule(dst io.Writer, mod Module) {
	WriteDSLBlockEntries(dst, mod.entries, WriteDSLContextInit())
}

func WriteGoMod(goMod GoMod, outputFile *os.File) {
	CheckErr(printer.Fprint(outputFile, goMod.astFileSet, goMod.astFile))
}

const (
	ALIGN_IND_SYM  int = 0
	ALIGN_IND_TYPE int = 1

	ALIGN_NUM_INDS int = 2
)

type WriteDSLAlignment = [ALIGN_NUM_INDS]IntOption

func WriteDSLAlignmentNone() WriteDSLAlignment {
	ret := [ALIGN_NUM_INDS]IntOption{}
	for i := 0; i < ALIGN_NUM_INDS; i++ {
		ret[i] = IntOptionNone()
	}
	return ret
}

func WriteDSLAlignmentMin(x, y *WriteDSLAlignment) *WriteDSLAlignment {
	if x == nil {
		return y
	}
	if y == nil {
		return x
	}
	ret := WriteDSLAlignmentNone()
	for i := 0; i < ALIGN_NUM_INDS; i++ {
		ret[i] = IntOptionMin(x[i], y[i])
	}
	return RefWriteDSLAlignment(ret)
}

func WriteDSLAlignmentMax(x, y *WriteDSLAlignment) *WriteDSLAlignment {
	if x == nil {
		return y
	}
	if y == nil {
		return x
	}
	ret := WriteDSLAlignmentNone()
	for i := 0; i < ALIGN_NUM_INDS; i++ {
		ret[i] = IntOptionMax(x[i], y[i])
	}
	return RefWriteDSLAlignment(ret)
}

const ALIGN_METHOD_ARGS_MARGIN int = 12

func WriteDSLAlignmentRes(xMin, xMax *WriteDSLAlignment, n []int) (ret *WriteDSLAlignment) {
	// defer func() {
	// 	fmt.Printf("WriteDSLAlignmentRes: %v, %v, %v -> %v\n", xMin, xMax, n, ret)
	// }()

	if xMin == nil {
		xMin = xMax
	}
	retVal := WriteDSLAlignmentNone()
	ret = RefWriteDSLAlignment(retVal)
	for i := 0; i < ALIGN_NUM_INDS; i++ {
		if xMin[i].IsNone() || xMax[i].IsNone() {
			Assert(xMin[i].IsNone() && xMax[i].IsNone())
			ret[i] = IntOptionNone()
		} else {
			iMin := xMin[i].Get()
			iMax := xMax[i].Get()
			if iMax == 0 {
				ret[i] = IntOptionSome(0)
			} else if n[i] <= 1 || (i+1 < ALIGN_NUM_INDS && n[i+1] <= 1) {
				ret[i] = IntOptionNone()
			} else {
				gap := 1
				if iMax > 2 || (iMax-iMin >= 1) {
					gap = 2
				}
				// if iMax >= 12 {
				// 	gap = 4
				// }
				ret[i] = IntOptionSome(iMax + gap)
			}
		}
	}
	return ret
}

type WriteDSLContext struct {
	indent            int
	useNewlines       bool
	useFieldAlignment bool
	alignment         WriteDSLAlignment
	isEnum            bool
}

func (ctx WriteDSLContext) Indent() WriteDSLContext {
	ret := ctx
	ret.indent = ctx.indent + 1
	return ret
}

func (ctx WriteDSLContext) Align(alignment WriteDSLAlignment) WriteDSLContext {
	ret := ctx
	ret.alignment = alignment
	return ret
}

func (ctx WriteDSLContext) AlignNone() WriteDSLContext {
	return ctx.Align(WriteDSLAlignmentNone())
}

func (ctx WriteDSLContext) IndentStr() string {
	indentStr := ""
	for i := 0; i < ctx.indent; i++ {
		indentStr += "    "
	}
	return indentStr
}

func WriteDSLAlignGap(xLen IntOption, xAlign IntOption) int {
	gapDefault := 1
	// gapAlignExtra := 2
	if xLen.IsSome() && xLen.Get() == 0 {
		if xAlign.IsNone() {
			return 0
		}
	}
	if xLen.IsNone() || xAlign.IsNone() {
		return gapDefault
	}
	iLen := xLen.Get()
	iAlign := xAlign.Get()
	return IntMax(iAlign-iLen, 0)
}

func WriteDSLFieldSym(dst io.Writer, field Field, ctx WriteDSLContext) {
	if field.fieldName == nil {
		return
	}
	if field.fieldType.Case() == Type_Case_NamedType {
		if field.fieldType.(*NamedType).name == *field.fieldName {
			return
		}
	}
	fmt.Fprintf(dst, "%s", DerefCheckString(field.fieldName))
}

func WriteDSLFieldSymFmtLen(field Field, ctx WriteDSLContext) IntOption {
	buf := bytes.NewBuffer([]byte{})
	WriteDSLFieldSym(buf, field, ctx)
	return StrFmtLen(buf.String())
}

func WriteDSLField(dst io.Writer, field Field, ctx WriteDSLContext) {
	symLen := WriteDSLFieldSymFmtLen(field, ctx)
	symGap := WriteDSLAlignGap(symLen, ctx.alignment[ALIGN_IND_SYM])
	offset := 0
	if symLen == IntOptionSome(0) {
		if ctx.alignment[ALIGN_IND_SYM].IsSome() {
			offset += ctx.alignment[ALIGN_IND_SYM].Get()
		}
		symGap = 0
	}
	typeLen := WriteDSLTypeFmtLen(field.fieldType, ctx)
	typeTarget := IntOptionAdd(ctx.alignment[ALIGN_IND_TYPE], IntOptionSome(offset))
	typeGap := WriteDSLAlignGap(typeLen, typeTarget)

	WriteDSLFieldSym(dst, field, ctx)
	if typeLen != IntOptionSome(0) || len(field.attributeList) > 0 {
		WriteRepeat(dst, " ", symGap)
	}
	WriteDSLType(dst, field.fieldType, ctx)
	if len(field.attributeList) > 0 {
		WriteRepeat(dst, " ", typeGap)
	}
	WriteDSLAttributeList(dst, field.attributeList, ctx)
}

func WriteDSLFieldFmtLen(field Field, ctx WriteDSLContext) IntOption {
	buf := bytes.NewBuffer([]byte{})
	WriteDSLField(buf, field, ctx)
	return StrFmtLen(buf.String())
}

const ENTRIES_BLOCK_INLINE_LEN_MAX int = 60

func EntriesReqNewlines(parentInfo *ParseFmtInfo, entries []Entry, ctx WriteDSLContext) (ret bool) {
	// defer func() {
	// 	fmt.Printf("???? ReqNewlines: \n  parentInfo:%#v\n  entries:%#v\n  ret:%#v\n", parentInfo, entries, ret)
	// }()

	if parentInfo != nil && parentInfo.hitNewlineTop {
		return true
	}

	itemsLen := 0

	for _, entry := range entries {
		switch entry.case_ {
		case Entry_Case_Comment:
			xr := entry.value.(Comment)
			if !xr.isBlock {
				return true
			}
			itemsLen += len(xr.commentText) + 4

		case Entry_Case_Field:
			xr := entry.value.(Field)
			if xr.parseFmtInfo != nil && xr.parseFmtInfo.hitNewlineTop {
				return true
			}
			currLen := WriteDSLFieldFmtLen(xr, ctx)
			if currLen.IsSome() {
				itemsLen += currLen.Get()
			} else {
				return true
			}

		case Entry_Case_Method:
			xr := entry.value.(Method)
			if xr.parseFmtInfo != nil && xr.parseFmtInfo.hitNewlineTop {
				return true
			}
			currLen := WriteDSLMethodFmtLen(xr, ctx)
			if currLen.IsSome() {
				itemsLen += currLen.Get()
			} else {
				return true
			}

		case Entry_Case_Empty:
			// pass

		default:
			Assert(false)
		}
	}

	return (itemsLen > ENTRIES_BLOCK_INLINE_LEN_MAX)
}

func AlignDSLBlockEntries(entries []Entry, ctx WriteDSLContext) []WriteDSLAlignment {
	alignActive := true
	alignActiveCount := 0
	alignActiveCountNontrivial := make([]int, ALIGN_NUM_INDS)
	alignLengths := []WriteDSLAlignment{}
	var alignLengthsMin, alignLengthsMax, alignLengthsRes *WriteDSLAlignment
	alignLengthsMax = RefWriteDSLAlignment(WriteDSLAlignment{IntOptionSome(0), IntOptionSome(0)})

	alignFlush := func() {
		for i := 0; i < alignActiveCount; i++ {
			if alignActiveCount > 1 {
				alignLengths = append(alignLengths, *alignLengthsRes)
			} else {
				alignLengths = append(alignLengths, WriteDSLAlignment{IntOptionNone(), IntOptionNone()})
			}
		}
		alignActive = false
		alignActiveCount = 0
		alignActiveCountNontrivial = make([]int, ALIGN_NUM_INDS)
		alignLengthsMin = nil
		alignLengthsMax = RefWriteDSLAlignment(WriteDSLAlignment{IntOptionSome(0), IntOptionSome(0)})
		alignLengthsRes = nil
	}

	if !ctx.useNewlines {
		alignActive = false
	}

	for _, entry := range entries {
		// fmt.Printf(" >>> proc entry %#v  %#v, %#v\n", entry, alignActive, alignActiveCount)
		// fmt.Printf("    >>> %#v\n", alignLengthsMax)

		if !alignActive {
			alignLengths = append(alignLengths, WriteDSLAlignment{IntOptionNone(), IntOptionNone()})
			if entry.case_ == Entry_Case_Empty && ctx.useNewlines {
				alignActive = true
			}
			continue
		}

		alignLengthsCurr := RefWriteDSLAlignment(WriteDSLAlignment{IntOptionNone(), IntOptionNone()})
		entryTrivial := false

		switch entry.case_ {
		case Entry_Case_Field:
			field := entry.value.(Field)
			symLen := WriteDSLFieldSymFmtLen(field, ctx)
			typeLen := IntOptionNone()
			if !ctx.isEnum {
				typeLen = WriteDSLTypeFmtLen(field.fieldType, ctx)
			}
			if symLen.IsNone() || typeLen.IsNone() {
				alignFlush()
				break
			}
			if symLen.Get() > 0 {
				alignLengthsCurr = RefWriteDSLAlignment(WriteDSLAlignment{symLen, typeLen})
			} else {
				alignLengthsCurr = RefWriteDSLAlignment(WriteDSLAlignment{typeLen, IntOptionSome(0)})
			}

		case Entry_Case_Method:
			method := entry.value.(Method)
			symLen := WriteDSLMethodSymFmtLen(method, ctx)
			typeLen := WriteDSLTypeFmtLen(method.methodRetType, ctx)
			if DSLTypeIsTrivialStruct(method.methodRetType) {
				typeLen = IntOptionSome(0)
			}
			if symLen.IsNone() || typeLen.IsNone() {
				alignFlush()
				break
			}
			if len(method.methodArgs) > 0 &&
				alignLengthsMax[ALIGN_IND_SYM].IsSome() &&
				symLen.Get() > alignLengthsMax[ALIGN_IND_SYM].Get()+ALIGN_METHOD_ARGS_MARGIN {
				alignFlush()
				break
			}
			if symLen.Get() > 0 {
				alignLengthsCurr = RefWriteDSLAlignment(WriteDSLAlignment{symLen, typeLen})
			} else {
				alignLengthsCurr = RefWriteDSLAlignment(WriteDSLAlignment{typeLen, IntOptionSome(0)})
			}

		case Entry_Case_Comment:
			entryTrivial = true

		case Entry_Case_Empty:
			entryTrivial = true

		case Entry_Case_Decl:
			Assert(false)
		}

		if alignActive {
			if !entryTrivial {
				alignLengthsMin = WriteDSLAlignmentMin(alignLengthsMin, alignLengthsCurr)
				alignLengthsMax = WriteDSLAlignmentMax(alignLengthsMax, alignLengthsCurr)

				for i := 0; i < ALIGN_NUM_INDS; i++ {
					if alignLengthsCurr[i].IsNone() || alignLengthsCurr[i].Get() == 0 {
						break
					}
					alignActiveCountNontrivial[i]++
				}
			}
			alignActiveCount += 1
			alignLengthsRes = WriteDSLAlignmentRes(alignLengthsMin, alignLengthsMax, alignActiveCountNontrivial)
		} else {
			alignLengths = append(alignLengths, *alignLengthsCurr)
		}
	}

	if alignActive {
		alignFlush()
	}

	Assert(len(alignLengths) == len(entries))
	// Assert(len(fmtLengths) == len(entries))

	return alignLengths
}

func WriteDSLBlockEntries(
	dst io.Writer,
	entries []Entry,
	ctx WriteDSLContext) {

	alignLengths := []WriteDSLAlignment{}
	if ctx.useFieldAlignment {
		alignLengths = AlignDSLBlockEntries(entries, ctx)
	} else {
		for i := 0; i < len(entries); i++ {
			alignLengths = append(alignLengths, WriteDSLAlignmentNone())
		}
	}

	holdEntryNewline := (len(entries) > 0 && EntryIsInlineComment(entries[0]))

	iMin := 0
	for i := 0; i <= len(entries); i++ {
		iMin = i
		if i >= len(entries) || entries[i].case_ != Entry_Case_Empty {
			break
		}
	}
	iMax := len(entries) - 1
	for i := len(entries) - 1; i >= -1; i-- {
		iMax = i
		if i < 0 || entries[i].case_ != Entry_Case_Empty {
			break
		}
	}
	if iMin == len(entries) {
		Assert(iMax == -1)
		iMin = 0
		iMax = len(entries) - 1
	}

	for i, entry := range entries {
		if i < iMin || i > iMax {
			Assert(entries[i].case_ == Entry_Case_Empty)
			continue
		}
		if entries[i].case_ == Entry_Case_Empty {
			if i > 0 && entries[i-1].case_ == Entry_Case_Empty {
				continue
			}
		}
		if ctx.useNewlines && !holdEntryNewline && entry.case_ != Entry_Case_Empty {
			fmt.Fprint(dst, ctx.IndentStr())
		}
		ctxSub := ctx.Align(alignLengths[i])
		switch entry.case_ {
		case Entry_Case_Field:
			field := entry.value.(Field)
			if ctx.isEnum {
				fmt.Fprintf(dst, "%s", DerefCheckString(field.fieldName))
			} else {
				WriteDSLField(dst, field, ctxSub)
			}

		case Entry_Case_Method:
			method := entry.value.(Method)
			WriteDSLMethod(dst, method, ctxSub)

		case Entry_Case_Comment:
			comment := entry.value.(Comment)
			if comment.isInline && i > 0 {
				fmt.Fprintf(dst, "  ")
			}
			WriteDSLComment(dst, comment, ctxSub)

		case Entry_Case_Empty:
			// defer to newline below

		case Entry_Case_Decl:
			Assert(!ctx.useFieldAlignment)
			decl := entry.value.(Decl)
			WriteDSLDecl(dst, decl, ctxSub)

		default:
			panic("Entry case not supported")
		}

		holdEntryNewline = (i+1 < len(entries) && EntryIsInlineComment(entries[i+1]))

		if ctx.useNewlines {
			if !holdEntryNewline {
				fmt.Fprintf(dst, "\n")
			}
		} else {
			if i < len(entries)-1 {
				fmt.Fprintf(dst, ", ")
			}
		}
	}
}

func WriteDSLBlock(
	dst io.Writer,
	entries []Entry,
	ctx WriteDSLContext,
	leftDelim string,
	rightDelim string) {

	if len(entries) == 0 {
		fmt.Fprintf(dst, "%s%s", leftDelim, rightDelim)
		return
	}

	fmt.Fprint(dst, leftDelim)
	ctxSub := ctx
	if ctx.useNewlines {
		ctxSub = ctxSub.Indent()
	}

	holdEntryNewline := (len(entries) > 0 && EntryIsInlineComment(entries[0]))

	if ctx.useNewlines && !holdEntryNewline {
		fmt.Fprintf(dst, "\n")
	}

	WriteDSLBlockEntries(dst, entries, ctxSub)

	if ctx.useNewlines {
		fmt.Fprint(dst, ctx.IndentStr())
	}

	fmt.Fprint(dst, rightDelim)
}

func WriteDSLMethodSym(dst io.Writer, method Method, ctx WriteDSLContext) {
	fmt.Fprintf(dst, "%s", method.methodName)
	ctxSub := WriteDSLContext{
		indent:            ctx.indent,
		useNewlines:       EntriesReqNewlines(method.argsFmtInfo, method.methodArgs, ctx),
		useFieldAlignment: true,
		alignment:         WriteDSLAlignmentNone(),
		isEnum:            false,
	}
	WriteDSLBlock(dst, method.methodArgs, ctxSub, "(", ")")
}

func WriteDSLMethodSymFmtLen(method Method, ctx WriteDSLContext) IntOption {
	buf := bytes.NewBuffer([]byte{})
	WriteDSLMethodSym(buf, method, ctx)
	return StrFmtLen(buf.String())
}

func WriteDSLMethod(dst io.Writer, method Method, ctx WriteDSLContext) {
	typeTrivial := DSLTypeIsTrivialStruct(method.methodRetType)

	symLen := WriteDSLMethodSymFmtLen(method, ctx)
	symGap := WriteDSLAlignGap(symLen, ctx.alignment[ALIGN_IND_SYM])
	if symLen == IntOptionSome(0) {
		Assert(false)
	}

	typeLen := WriteDSLTypeFmtLen(method.methodRetType, ctx)
	if typeTrivial {
		typeLen = IntOptionSome(0)
	}
	typeGap := WriteDSLAlignGap(typeLen, ctx.alignment[ALIGN_IND_TYPE])
	WriteDSLMethodSym(dst, method, ctx)
	if typeLen != IntOptionSome(0) || len(method.attributeList) > 0 {
		WriteRepeat(dst, " ", symGap)
	}
	if !typeTrivial {
		WriteDSLType(dst, method.methodRetType, ctx)
	}
	if len(method.attributeList) > 0 {
		WriteRepeat(dst, " ", typeGap)
	}
	WriteDSLAttributeList(dst, method.attributeList, ctx)
}

func WriteDSLMethodFmtLen(method Method, ctx WriteDSLContext) IntOption {
	buf := bytes.NewBuffer([]byte{})
	WriteDSLMethod(buf, method, ctx)
	return StrFmtLen(buf.String())
}

func WriteDSLComment(
	dst io.Writer, comment Comment, ctx WriteDSLContext) {

	if comment.isBlock {
		Assert(!comment.isInline)
		fmt.Fprintf(dst, "/*%s*/", comment.commentText)
	} else {
		fmt.Fprintf(dst, "//%s", comment.commentText)
	}
}

func WriteDSLAttributeList(dst io.Writer, attributeList []string, ctx WriteDSLContext) {
	if len(attributeList) == 0 {
		return
	}
	fmt.Fprintf(dst, "@(")
	for i, attr := range attributeList {
		if i > 0 {
			fmt.Fprintf(dst, ", ")
		}
		fmt.Fprintf(dst, "%s", attr)
	}
	fmt.Fprintf(dst, ")")
}

func WriteDSLType(dst io.Writer, type_ Type, ctx WriteDSLContext) {
	switch type_.Case() {
	case Type_Case_NamedType:
		xr := type_.(*NamedType)
		fmt.Fprintf(dst, "%s", xr.name)

	case Type_Case_AlgType:
		xr := type_.(*AlgType)
		if xr.isTuple {
			Assert(xr.sort == AlgSort_Prod)
			Assert(len(xr.attributeList) == 0)
			ctxSub := WriteDSLContext{
				indent:            ctx.indent,
				useNewlines:       EntriesReqNewlines(xr.entriesFmtInfo, xr.entries, ctx),
				useFieldAlignment: true,
				alignment:         WriteDSLAlignmentNone(),
				isEnum:            false,
			}
			WriteDSLBlock(dst, xr.entries, ctxSub, "(", ")")
			break
		}

		isEnum := false
		if xr.sort == AlgSort_Prod {
			if xr.isInterface {
				fmt.Fprintf(dst, "interface")
			} else {
				fmt.Fprintf(dst, "struct")
			}
		} else if xr.sort == AlgSort_Sum {
			if len(xr.Methods()) == 0 {
				isEnum = true
				for _, field := range xr.Fields() {
					if !DSLTypeIsTrivialStruct(field.fieldType) || len(field.attributeList) > 0 {
						isEnum = false
						break
					}
				}
			}
			if isEnum {
				fmt.Fprintf(dst, "enum")
			} else {
				fmt.Fprintf(dst, "union")
			}
		} else {
			Assert(false)
		}

		fmt.Fprintf(dst, " ")
		WriteDSLAttributeList(dst, xr.attributeList, ctx)
		if len(xr.attributeList) > 0 {
			fmt.Fprintf(dst, " ")
		}

		ctxSub := WriteDSLContext{
			indent:            ctx.indent,
			useNewlines:       EntriesReqNewlines(xr.entriesFmtInfo, xr.entries, ctx),
			useFieldAlignment: true,
			alignment:         WriteDSLAlignmentNone(),
			isEnum:            isEnum,
		}
		WriteDSLBlock(dst, xr.entries, ctxSub, "{", "}")

	case Type_Case_ArrayType:
		xr := type_.(*ArrayType)
		fmt.Fprintf(dst, "[")
		WriteDSLType(dst, xr.elementType, ctx)
		fmt.Fprintf(dst, "]")

	case Type_Case_OptionType:
		xr := type_.(*OptionType)
		WriteDSLType(dst, xr.valueType, ctx)
		fmt.Fprintf(dst, "?")

	case Type_Case_MapType:
		xr := type_.(*MapType)
		fmt.Fprintf(dst, "{")
		WriteDSLType(dst, xr.keyType, ctx)
		fmt.Fprintf(dst, ": ")
		WriteDSLType(dst, xr.valueType, ctx)
		fmt.Fprintf(dst, "}")

	case Type_Case_RefType:
		xr := type_.(*RefType)
		fmt.Fprintf(dst, "&")
		WriteDSLType(dst, xr.targetType, ctx)

	default:
		fmt.Printf("Unhandled case: %v\n", type_.Case())
		panic("TODO")
	}
}
