package codeGen

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"runtime/debug"
	"strings"
)

const Whitespace = " \t\n"
const Symbols = "(){}[],;|&?://*\""

const DebugParser = false

type ParseFmtInfo struct {
	err *ParseError_S
	// startPos, endPos int
	hitNewlineAny bool
	hitNewlineTop bool
	noopEntries   []Entry // Comment or Empty
}

func ParseFmtInfoInit() ParseFmtInfo {
	return ParseFmtInfo{
		err:           nil,
		hitNewlineAny: false,
		hitNewlineTop: false,
		noopEntries:   []Entry{},
	}
}

func ParseFmtInfoError(err *ParseError_S) ParseFmtInfo {
	return ParseFmtInfo{
		err:           err,
		hitNewlineAny: false,
		hitNewlineTop: false,
		noopEntries:   []Entry{},
	}
}

func (info ParseFmtInfo) Comments() []Comment {
	ret := []Comment{}
	for _, entry := range info.noopEntries {
		if entry.case_ == Entry_Case_Comment {
			ret = append(ret, entry.value.(Comment))
		}
	}
	return ret
}

func (err *ParseError_S) UnifyError(errOther *ParseError_S) *ParseError_S {
	if err == nil {
		return errOther
	}
	if errOther == nil {
		return err
	}
	if errOther.pos > err.pos {
		return errOther
	}
	if errOther.pos < err.pos {
		return err
	}
	subsumed := true
	for _, s := range errOther.msgAbbrev {
		if !SliceContainsString(err.msgAbbrev, s) {
			subsumed = false
		}
	}
	if subsumed {
		return err
	}
	return RefParseError(ParseError_S{
		msgAbbrev: append(err.msgAbbrev, errOther.msgAbbrev...),
		msg:       append(err.msg, errOther.msg...),
		pos:       err.pos,
	})
}

func (info ParseFmtInfo) UnifyFmtInfoCaptureComments(
	r *ParseStream, infoOther ParseFmtInfo) (ParseFmtInfo, []Entry) {

	ret := info.UnifyFmtInfo(r, infoOther)
	retNoopEntries := ret.noopEntries
	ret.noopEntries = []Entry{}
	return ret, retNoopEntries
}

func (info ParseFmtInfo) UnifyFmtInfoRejectCommentsExt(
	r *ParseStream, infoOther ParseFmtInfo, shieldNewlines bool) ParseFmtInfo {

	ret := info.UnifyFmtInfoExt(r, infoOther, shieldNewlines)
	if len(infoOther.Comments()) > 0 {
		minEndPos := infoOther.Comments()[0].endPos + 1
		commentTextAbbrev := ""
		for _, comment := range infoOther.Comments() {
			// fmt.Printf("??? \"%v\"\n", comment.commentText)
			if comment.endPos < minEndPos {
				minEndPos = comment.endPos
				commentTextAbbrev = TextAbbrev(comment.commentText, 16)
			}
		}
		commentErr := r.GenParseErrorExt(
			fmt.Sprintf("Unexpected comment \"%v\"", commentTextAbbrev), minEndPos)
		ret.err = ret.err.UnifyError(commentErr)

		// panic(commentErr)
	}
	return ret
}

func (info ParseFmtInfo) UnifyFmtInfoRejectComments(
	r *ParseStream, infoOther ParseFmtInfo) ParseFmtInfo {
	return info.UnifyFmtInfoRejectCommentsExt(r, infoOther, false)
}

func (info ParseFmtInfo) UnifyFmtInfoExt(
	r *ParseStream, infoOther ParseFmtInfo, shieldNewlines bool) ParseFmtInfo {

	noopEntriesRet := append(info.noopEntries, infoOther.noopEntries...)
	hitNewlineTop := info.hitNewlineTop || (infoOther.hitNewlineTop && !shieldNewlines)
	ret := ParseFmtInfo{
		err: info.err.UnifyError(infoOther.err),
		// startPos: IntMin(info.startPos, infoOther.startPos),
		// endPos: IntMax(info.endPos, infoOther.endPos),
		hitNewlineAny: info.hitNewlineAny || infoOther.hitNewlineAny,
		hitNewlineTop: hitNewlineTop,
		noopEntries:   noopEntriesRet,
	}
	return ret
}

func (info ParseFmtInfo) UnifyFmtInfo(r *ParseStream, infoOther ParseFmtInfo) ParseFmtInfo {
	return info.UnifyFmtInfoExt(r, infoOther, false)
}

type LineInfo struct {
	line int
	col  int
}

type ParseStreamState struct {
	pos int
}

type ParseStream struct {
	rs         io.ReadSeeker
	state      ParseStreamState
	stateStack []ParseStreamState
	lineMap    []LineInfo
	buffer     *bytes.Buffer
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
	r.stateStack = r.stateStack[0 : n-1]
}

func (r *ParseStream) Seek(offset int) {
	// if DebugParser {
	// 	fmt.Printf("Seek: %v\n", offset)
	// }
	_, err := r.rs.Seek(int64(offset), io.SeekCurrent)
	CheckErr(err)
	r.state.pos += offset
}

func (r *ParseStream) GetLineInfo(pos int) LineInfo {
	var ret LineInfo
	if pos < 0 {
		ret = LineInfo{line: 0, col: 0}
	} else {
		ret = r.lineMap[pos]
	}
	return ret
}

func (r *ParseStream) PosDebugExt(pos int) string {
	lineInfo := r.GetLineInfo(pos)
	return fmt.Sprintf("line %v, column %v", lineInfo.line+1, lineInfo.col+1)
}

func (r *ParseStream) PosDebug() string {
	return r.PosDebugExt(r.state.pos - 1)
}

func (r *ParseStream) BufferToNextNewline(pos int) {
	r.Push()
	defer r.Pop(true)

	posInit := r.state.pos
	r.Seek(pos - posInit)
	for {
		s := r.Read(1)
		if s == "" || s == "\n" {
			return
		}
	}
}

func (r *ParseStream) GenExcerptExt(pos int) string {
	Assert(pos < len(r.buffer.Bytes()))

	lineInfo := r.GetLineInfo(pos)

	r.BufferToNextNewline(IntMax(pos, 0))
	ret := ""
	i := pos
	for {
		if i > -1 && r.buffer.Bytes()[i] != '\n' {
			i--
		} else {
			break
		}
	}
	for i++; i < len(r.buffer.Bytes()); i++ {
		c := string(r.buffer.Bytes()[i : i+1])
		ret += c
		if c == "\n" {
			break
		}
	}
	if i == len(r.buffer.Bytes()) {
		ret += "\n"
	}
	if lineInfo.col == 0 {
		ret += "\u2196\n"
	} else {
		Assert(lineInfo.col > 0)
		ret += WriteRepeatString(" ", lineInfo.col-1) + "\u2191\n"
	}
	return ret
}

func (r *ParseStream) GenExcerpt() string {
	return r.GenExcerptExt(r.state.pos - 1)
}

func (r *ParseStream) GenParseErrorExt(errMsg string, pos int) *ParseError_S {
	errMsgFull := fmt.Sprintf("Parse error (%v)\n\n", r.PosDebugExt(pos))
	errMsgFull += r.GenExcerptExt(pos)
	errMsgFull += "\n"
	errMsgFull += errMsg
	errMsgFull += "\n\n"
	errMsgFull += string(debug.Stack())
	ret := ParseError(errMsg, errMsgFull, pos)
	return RefParseError(ret)
}

func (r *ParseStream) GenParseError(errMsg string) *ParseError_S {
	var i int = r.state.pos - 1
	return r.GenParseErrorExt(errMsg, i)
}

func (r *ParseStream) Get(lenMax int, advance bool) string {
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
			prevLineInfo := LineInfo{line: 0, col: 0}
			if p >= 1 {
				prevLineInfo = r.lineMap[p-1]
			}
			newLineInfo := LineInfo{
				line: prevLineInfo.line,
				col:  prevLineInfo.col + 1,
			}
			if c == '\n' {
				newLineInfo = LineInfo{
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

func (r *ParseStream) Read(lenMax int) string {
	return r.Get(lenMax, true)
}

func (r *ParseStream) Peek(lenMax int) string {
	return r.Get(lenMax, false)
}

func (r *ParseStream) ParseComments() (info ParseFmtInfo) {
	info.noopEntries = []Entry{}

	if DebugParser {
		fmt.Printf(" ** ParseComments body: %v \n", r.PosDebug())
	}

	commentText := ""
	depth := 0
	singleLineComment := false

	whitespaceOnlyLine := true
	for i := r.state.pos - 1; i >= 0; i-- {
		if !strings.Contains(Whitespace, string(r.buffer.Bytes()[i:i+1])) {
			whitespaceOnlyLine = false
			break
		}
		if r.buffer.Bytes()[i] == '\n' {
			break
		}
	}

	for {
		if _, ok := r.PeekExact(1); !ok {
			return
		}

		initPeekRes, initPeekOk := r.PeekExact(1)
		if strings.Contains(Whitespace, initPeekRes) && initPeekOk {
			if initPeekRes == "\n" {
				if whitespaceOnlyLine && depth == 0 && !singleLineComment {
					info.noopEntries = append(info.noopEntries, EntryEmpty())
				}
				if singleLineComment {
					info.noopEntries = append(info.noopEntries, EntryComment(Comment{
						commentText: commentText,
						isBlock:     false,
						isInline:    !info.hitNewlineAny,
						endPos:      r.state.pos,
					}))
					commentText = ""
				}
				singleLineComment = false
				if depth == 0 {
					info.hitNewlineAny = true
					info.hitNewlineTop = true
				}
				whitespaceOnlyLine = true
			}
			if depth > 0 || singleLineComment {
				commentText += initPeekRes
			}
			r.Seek(1)
			continue
		}

		whitespaceOnlyLine = false

		if depth == 0 && !singleLineComment {
			if res, ok := r.PeekExact(2); res == "//" && ok {
				r.Seek(2)
				singleLineComment = true
				continue
			}
		}

		if res, ok := r.PeekExact(2); res == "/*" && ok && !singleLineComment {
			if depth > 0 {
				commentText += res
			}
			r.Seek(2)
			depth += 1
		} else if res, ok := r.PeekExact(2); res == "*/" && ok && !singleLineComment {
			r.Seek(2)
			depth -= 1
			if depth == 0 {
				info.noopEntries = append(info.noopEntries, EntryComment(Comment{
					commentText: commentText,
					isBlock:     true,
					isInline:    false,
					endPos:      r.state.pos,
				}))
				commentText = ""
			} else {
				commentText += res
			}
		} else {
			if depth > 0 || singleLineComment {
				commentText += initPeekRes
				r.Seek(1)
			} else {
				return
			}
		}
	}
}

func ReadToken(r *ParseStream) (ret string, info ParseFmtInfo) {
	var infoSub ParseFmtInfo

	defer func() {
		if len(ret) == 0 && info.err == nil {
			fmt.Printf("ReadToken: %v\n", r.GenExcerpt())
			Assert(false)
		}
	}()

	if DebugParser {
		defer func() {
			fmt.Printf("ReadToken: \"%v\", %v, %v\n", TextAbbrev(ret, 16), info.hitNewlineAny, info.err == nil)
		}()
	}

	retToken := []byte{}
	infoSub = r.ParseComments()
	info = info.UnifyFmtInfo(r, infoSub)
	if info.err != nil {
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
				info.err = r.GenParseError("Reached end-of-file")
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
	msgAbbrev []string
	msg       []string
	pos       int
}

func ParseError(msgAbbrev string, msg string, pos int) ParseError_S {
	return ParseError_S{msgAbbrev: []string{msgAbbrev}, msg: []string{msg}, pos: pos}
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
		if (IsDigit(c) || c == '.') && i > 0 {
			continue
		}
		return false
	}
	return true
}

func ParseIdent(r *ParseStream) (ret string, info ParseFmtInfo) {
	var infoSub ParseFmtInfo

	ret, infoSub = ReadToken(r)
	info = info.UnifyFmtInfo(r, infoSub)

	if info.err != nil {
		ret = ""
		return
	} else if !IsIdent(ret) {
		ret = ""
		msg := fmt.Sprintf("Expected identifier; received: \"%v\"", ret)
		info.err = r.GenParseError(msg)
		return
	} else {
		return
	}
}

func (err ParseError_S) Error() string {
	if len(err.msg) == 1 {
		return err.msg[0]
	} else {
		Assert(len(err.msg) > 1)
		ret := ""
		for i, mi := range err.msg {
			if i > 0 {
				ret += "\n"
			}
			ret += fmt.Sprintf("Attempt %v:\n    %v\n", i+1, mi)
		}
		return ret
	}
}

func ReadTokenCheck(r *ParseStream, validTokens []string) (ret string, info ParseFmtInfo) {
	var tok string
	var infoSub ParseFmtInfo

	for _, v := range validTokens {
		if v == "\n" {
			if tok, ok := PeekToken(r, true); ok && (tok == "\n") {
				if DebugParser {
					fmt.Printf("ReadTokenCheck read newline (special case): %v\n", r.PosDebug())
				}
				ret = tok
				return
			}
		}
	}

	tok, infoSub = ReadToken(r)
	info = info.UnifyFmtInfo(r, infoSub)
	if info.err != nil {
		ret = ""
		return
	}
	for _, v := range validTokens {
		if tok == v {
			ret = tok
			return
		}
	}

	ret = ""
	validTokensDisp := []string{}
	for _, tok := range validTokens {
		validTokenDisp := "\"" + tok + "\""
		validTokensDisp = append(validTokensDisp, validTokenDisp)
	}
	errMsg := fmt.Sprintf("Unexpected token \"%v\" (expected: %v)", tok, strings.Join(validTokensDisp, ", "))
	info.err = r.GenParseError(errMsg)
	return
}

func (r *ParseStream) ReadTokenSequenceCheck(tokenSeq []string) (ret string, info ParseFmtInfo) {
	var infoSub ParseFmtInfo
	var retSub string

	for _, token := range tokenSeq {
		retSub, infoSub = ReadTokenCheck(r, []string{token})
		info = info.UnifyFmtInfo(r, infoSub)
		ret += retSub
		if info.err != nil {
			ret = ""
			return
		}
	}
	return
}

func TryReadTokenCheck(r *ParseStream, validTokens []string) (ret string, info ParseFmtInfo) {
	r.Push()
	defer func() { r.Pop(info.err != nil) }()

	ret, info = ReadTokenCheck(r, validTokens)
	return
}

func (r *ParseStream) PeekExact(len_ int) (string, bool) {
	ret := r.Peek(len_)
	return ret, (len(ret) == len_)
}

func PeekToken(r *ParseStream, matchNewline bool) (string, bool) {
	r.Push()
	defer func() { r.Pop(true) }()

	ret, info := ReadToken(r)
	if info.hitNewlineAny && matchNewline {
		return "\n", true
	}
	if info.err != nil {
		return "", false
	} else {
		return ret, true
	}
}

func ParseAttributeList(r *ParseStream) (ret []string, info ParseFmtInfo) {
	ret = []string{}

	var infoSub ParseFmtInfo

	if tok, ok := PeekToken(r, false); ok && (tok == "@") {
		_, _ = ReadTokenCheck(r, []string{"@"})
		fTryParse := func(rr *ParseStream) (interface{}, ParseFmtInfo) {
			tok, info := ReadToken(rr)
			return tok, info
		}
		fAppend := func(x interface{}) {
			ret = append(ret, x.(string))
		}
		infoSub = ParseDelimitedList(r, "(", []string{","}, ")", fTryParse, fAppend, false)
		info = info.UnifyFmtInfo(r, infoSub)
		if info.err != nil {
			ret = nil
			return
		}
		return
	} else {
		return
	}
}

type EntryParseSpec int

const (
	EntryParseSpec_Case      EntryParseSpec = 1
	EntryParseSpec_DataField EntryParseSpec = 2
	EntryParseSpec_ArgField  EntryParseSpec = 3
	EntryParseSpec_Method    EntryParseSpec = 4
)

func EntryParseSpec_MayOmitName(spec EntryParseSpec) bool {
	switch spec {
	case EntryParseSpec_Case:
		return false
	case EntryParseSpec_DataField:
		return true
	case EntryParseSpec_ArgField:
		return true
	case EntryParseSpec_Method:
		return false
	default:
		panic("Unhandled case")
	}
}

func EntryParseSpec_MustOmitName(spec EntryParseSpec) bool {
	switch spec {
	case EntryParseSpec_Case:
		return false
	case EntryParseSpec_DataField:
		return false
	case EntryParseSpec_ArgField:
		return false
	case EntryParseSpec_Method:
		return false
	default:
		panic("Unhandled case")
	}
}

func EntryParseSpec_MayOmitType(spec EntryParseSpec) bool {
	switch spec {
	case EntryParseSpec_Case:
		return true
	case EntryParseSpec_DataField:
		return false
	case EntryParseSpec_ArgField:
		return false
	case EntryParseSpec_Method:
		return true
	default:
		panic("Unhandled case")
	}
}

func EntryParseSpec_MustOmitType(spec EntryParseSpec) bool {
	switch spec {
	case EntryParseSpec_Case:
		return true
	case EntryParseSpec_DataField:
		return false
	case EntryParseSpec_ArgField:
		return false
	case EntryParseSpec_Method:
		return false
	default:
		panic("Unhandled case")
	}
}

func EntryParseSpec_HasArgs(spec EntryParseSpec) bool {
	switch spec {
	case EntryParseSpec_Case:
		return false
	case EntryParseSpec_DataField:
		return false
	case EntryParseSpec_ArgField:
		return false
	case EntryParseSpec_Method:
		return true
	default:
		panic("Unhandled case")
	}
}

func EntryParseSpec_AllowsUnnamedTypes(spec EntryParseSpec) bool {
	switch spec {
	case EntryParseSpec_Case:
		return false
	case EntryParseSpec_DataField:
		return false
	case EntryParseSpec_ArgField:
		return true
	case EntryParseSpec_Method:
		return false
	default:
		panic("Unhandled case")
	}
}

func EntryParseSpec_Consistent(spec EntryParseSpec, omitName bool, omitType bool) bool {
	if omitName && omitType {
		return false
	}
	if omitName && !EntryParseSpec_MayOmitName(spec) {
		return false
	}
	if !omitName && EntryParseSpec_MustOmitName(spec) {
		return false
	}
	if omitType && !EntryParseSpec_MayOmitType(spec) {
		return false
	}
	if !omitType && EntryParseSpec_MustOmitType(spec) {
		return false
	}
	return true
}

func PeekCheckDelim(r *ParseStream, delimSet []string) bool {
	matchNewline := SliceContainsString(delimSet, "\n")
	if tok, ok := PeekToken(r, matchNewline); ok {
		// fmt.Printf("    >>> PeekCheckDelim: tok = \"%v\"\n", tok)

		for _, delim := range delimSet {
			if tok == delim {
				return true
			}
		}
	}
	return false
}

func ParseEntry(
	r *ParseStream,
	spec EntryParseSpec,
	delimSet []string) (ret *Entry, info ParseFmtInfo) {

	// defer func() {
	// 	if ret != nil {
	// 		fmt.Printf(" >>>>> ParseEntry return: %#v\n", ret)
	// 	}
	// }()

	var entryNameStr string
	var entryName *string
	var entryRetType Type
	var attributeList []string

	var argsEntriesSub []Entry
	var infoSub ParseFmtInfo

	var errAcc *ParseError_S

	// fmt.Printf(" >>>>> ParseEntry start\n%v\n\n\n", r.GenExcerpt())

	r.Push()
	defer func() {
		// fmt.Printf(" >>>>> ParseEntry end: ok=%v, hitNewlineTop=%v\n", info.err == nil, info.hitNewlineTop)
		r.Pop(info.err != nil)
	}()

	for _, omitName := range []bool{false, true} {
		for _, omitType := range []bool{false, true} {
			if !EntryParseSpec_Consistent(spec, omitName, omitType) {
				continue
			}

			r.Pop(true)
			r.Push()

			// if true {
			// 	msgAbbrev := "<nil>"
			// 	if errAcc != nil {
			// 		msgAbbrev = errAcc.msgAbbrev[0]
			// 	}
			// 	fmt.Printf(" >>>>> ParseEntry try: omitName=%v, omitType=%v, delimSet=%v\n errAcc: %v\n\n",
			// 		omitName, omitType, delimSet, msgAbbrev)
			// }

			info = ParseFmtInfoInit()

			if omitName {
				entryName = nil
			} else {
				entryNameStr, infoSub = ParseIdent(r)
				info = info.UnifyFmtInfoRejectComments(r, infoSub)
				entryName = RefString(entryNameStr)

				if DebugParser {
					fmt.Printf(" >>>>> IDENT: %v %v\n", entryName, info.err == nil)
				}

				if info.err != nil {
					errAcc = errAcc.UnifyError(info.err)
					continue
				}
			}

			var argsFmtInfo *ParseFmtInfo = nil
			if EntryParseSpec_HasArgs(spec) {
				specsSub := []EntryParseSpec{EntryParseSpec_ArgField}
				argsEntriesSub, infoSub = ParseEntryList(r, "(", []string{",", "\n"}, ")", false, specsSub)
				argsFmtInfo = RefParseFmtInfo(infoSub)

				info = info.UnifyFmtInfoRejectCommentsExt(r, infoSub, true)

				if info.err != nil {
					errAcc = errAcc.UnifyError(info.err)
					continue
				}
			}

			if omitType {
				entryRetType = DSLTrivialStruct()
			} else {
				entryRetType, infoSub = ParseType(r)
				info = info.UnifyFmtInfoRejectComments(r, infoSub)

				if infoSub.hitNewlineTop && !omitName {
					errAcc = errAcc.UnifyError(
						r.GenParseError(
							fmt.Sprintf(
								"Newline not permitted between name and type (omitName=%v, omitType=%v)",
								omitName, omitType)))
					continue
				}

				if DebugParser {
					fmt.Printf(" >>>>> ParseEntry ParseType: %v %v\n", entryRetType, info.err == nil)
				}
				if info.err != nil {
					errAcc = errAcc.UnifyError(info.err)
					continue
				}
			}

			attributeList, infoSub = ParseAttributeList(r)
			info = info.UnifyFmtInfoRejectComments(r, infoSub)

			if DebugParser {
				fmt.Printf(" >>>>> ParseEntry ParseAttributeList: %v %v\n", attributeList, info.err == nil)
			}
			if info.err != nil {
				errAcc = errAcc.UnifyError(info.err)
				continue
			}

			if PeekCheckDelim(r, delimSet) {
				if entryName == nil {
					Assert(entryRetType != nil)
					if entryRetType.Case() == Type_Case_NamedType {
						entryNameStr = entryRetType.(*NamedType).name
						entryName = RefString(entryNameStr)
					} else {
						if EntryParseSpec_AllowsUnnamedTypes(spec) {
							entryName = nil
						} else {
							errAcc = errAcc.UnifyError(
								r.GenParseError("Unnamed fields must have named types"))
							continue
						}
					}
				}
			} else {
				errAcc = errAcc.UnifyError(
					r.GenParseError("Expected field/method delimiter"))
				continue
			}

			Assert(info.err == nil)

			switch {
			case spec == EntryParseSpec_Method:
				Assert(argsEntriesSub != nil)
				Assert(argsFmtInfo != nil)

				retEntry := EntryMethod(Method{
					methodName:    *entryName,
					methodArgs:    argsEntriesSub,
					argsFmtInfo:   argsFmtInfo,
					methodRetType: entryRetType,
					attributeList: attributeList,
					parseFmtInfo:  RefParseFmtInfo(info),
				})
				ret = RefEntry(retEntry)

			case spec == EntryParseSpec_DataField || spec == EntryParseSpec_ArgField || spec == EntryParseSpec_Case:
				retEntry := EntryField(Field{
					fieldName:     entryName,
					fieldType:     entryRetType,
					attributeList: attributeList,
					parseFmtInfo:  RefParseFmtInfo(info),
				})
				ret = RefEntry(retEntry)

			default:
				panic("Unhandled case")
			}

			return
		}
	}

	Assert(errAcc != nil)
	info = ParseFmtInfoError(errAcc)
	ret = nil
	return
}

func TryParseEntry(r *ParseStream, spec EntryParseSpec, delimSet []string) (ret *Entry, info ParseFmtInfo) {
	r.Push()
	defer func() { r.Pop(info.err != nil) }()

	ret, info = ParseEntry(r, spec, delimSet)
	return
}

func ParseDelimitedList(
	r *ParseStream,
	start string,
	validDelims []string,
	end string,
	fTryParse func(*ParseStream) (interface{}, ParseFmtInfo),
	fAppend func(interface{}),
	allowInitDelim bool) (info ParseFmtInfo) {

	var infoSub ParseFmtInfo
	var noopEntriesSub []Entry
	var xSub interface{}

	_, infoSub = ReadTokenCheck(r, []string{start})
	info = info.UnifyFmtInfoRejectComments(r, infoSub)

	if info.err != nil {
		return
	}

	if DebugParser {
		fmt.Printf("ParseDelimitedList read start: %v, %v\n", start, r.PosDebug())
	}

	Assert(!SliceContainsString(validDelims, ""))

	if allowInitDelim {
		ret, infoSub := TryReadTokenCheck(r, validDelims)
		if infoSub.err == nil {
			info, noopEntriesSub = info.UnifyFmtInfoCaptureComments(r, infoSub)
			for _, noopEntrySub := range noopEntriesSub {
				fAppend(RefEntry(noopEntrySub))
			}
			if DebugParser {
				fmt.Printf("ParseDelimitedList read init delim: %v, %v, %v\n", ret, info.err == nil, r.PosDebug())
			}
		}
	}

	for {
		infoSub = r.ParseComments()
		if infoSub.err == nil {
			info, noopEntriesSub = info.UnifyFmtInfoCaptureComments(r, infoSub)
			for _, noopEntrySub := range noopEntriesSub {
				fAppend(RefEntry(noopEntrySub))
			}
		} else {
			panic(infoSub.err)
		}

		xSub, infoSub = fTryParse(r)
		var errSubAcc *ParseError_S

		if infoSub.err == nil {

			// hitNewlineTopPrev := info.hitNewlineTop
			info, noopEntriesSub = info.UnifyFmtInfoCaptureComments(r, infoSub)
			// info.hitNewlineTop = hitNewlineTopPrev

			for _, noopEntrySub := range noopEntriesSub {
				fAppend(RefEntry(noopEntrySub))
			}
			if DebugParser {
				fmt.Printf("ParseDelimitedList read item: %T %v\n", xSub, xSub)
			}
			fAppend(xSub)
		} else {
			errSubAcc = errSubAcc.UnifyError(infoSub.err)

			_, infoSub = ReadTokenCheck(r, []string{end})
			if DebugParser {
				fmt.Printf("ParseDelimitedList looking for end: %v (%v)\n", end, infoSub.err == nil)
			}

			if infoSub.err == nil {
				info, noopEntriesSub = info.UnifyFmtInfoCaptureComments(r, infoSub)
				for _, noopEntrySub := range noopEntriesSub {
					fAppend(RefEntry(noopEntrySub))
				}
				return
			} else {
				errSubAcc = errSubAcc.UnifyError(infoSub.err)
				info.err = errSubAcc
				return
			}
		}

		_, infoSub = TryReadTokenCheck(r, validDelims)
		if infoSub.err == nil {
			info, noopEntriesSub = info.UnifyFmtInfoCaptureComments(r, infoSub)
			for _, noopEntrySub := range noopEntriesSub {
				fAppend(RefEntry(noopEntrySub))
			}
			continue
		} else {
			_, infoSub = ReadTokenCheck(r, []string{end})
			info, noopEntriesSub = info.UnifyFmtInfoCaptureComments(r, infoSub)
			for _, noopEntrySub := range noopEntriesSub {
				fAppend(RefEntry(noopEntrySub))
			}
			return
		}
	}
}

func ParseEntryList(
	r *ParseStream,
	start string,
	validDelims []string,
	end string,
	allowInitDelim bool,
	specs []EntryParseSpec,
) (retEntries []Entry, info ParseFmtInfo) {

	retEntries = []Entry{}

	fTryParse := func(r *ParseStream) (ret interface{}, info ParseFmtInfo) {
		var retSub *Entry
		var infoSub ParseFmtInfo
		var errAcc *ParseError_S = nil

		for _, spec := range specs {
			retSub, infoSub = TryParseEntry(r, spec, append(validDelims, end))
			if infoSub.err == nil {
				ret = retSub
				info = infoSub
				return
			} else {
				errAcc = errAcc.UnifyError(infoSub.err)
			}
		}

		ret = nil
		Assert(errAcc != nil)
		info = ParseFmtInfoError(errAcc)
		return
	}

	fAppend := func(x interface{}) {
		retEntries = append(retEntries, *x.(*Entry))
	}

	infoSub := ParseDelimitedList(r, start, validDelims, end, fTryParse, fAppend, allowInitDelim)
	info = info.UnifyFmtInfoRejectComments(r, infoSub)
	return
}

func ParseType(r *ParseStream) (ret Type, info ParseFmtInfo) {
	var tok string
	var infoSub ParseFmtInfo
	var attributeList []string
	var entriesSub []Entry
	var elementType, targetType, keyType, valueType Type

	// fmt.Printf(" >>>>> ParseType start\n%v\n", r.GenExcerpt())
	// defer func() {
	// 	fmt.Printf(" >>>>> ParseType end  %#v\n%v\n%v\n", ret, r.state.pos, r.GenExcerpt())
	// }()

	tok, infoSub = ReadToken(r)
	info = info.UnifyFmtInfoRejectComments(r, infoSub)

	if info.err != nil {
		ret = nil
		return
	}

	if DebugParser {
		fmt.Printf("ParseType token: \"%v\"\n", tok)
	}

	switch {
	case tok == "interface" || tok == "struct" || tok == "union" || tok == "enum":
		algSort := AlgSort_Prod
		validDelims := []string{"\n", ","}
		specs := []EntryParseSpec{EntryParseSpec_DataField, EntryParseSpec_Method}

		if tok == "union" || tok == "enum" {
			algSort = AlgSort_Sum
			validDelims = append(validDelims, "|")
			if tok == "enum" {
				specs = []EntryParseSpec{EntryParseSpec_Case}
			}
		}

		attributeList, infoSub = ParseAttributeList(r)
		info = info.UnifyFmtInfoRejectComments(r, infoSub)

		if info.err != nil {
			ret = nil
			return
		}

		entriesSub, infoSub = ParseEntryList(r, "{", validDelims, "}", true, specs)
		entriesFmtInfo := infoSub
		info = info.UnifyFmtInfoRejectCommentsExt(r, infoSub, true)

		if info.err != nil {
			ret = nil
			return
		}

		ret = RefAlgType(AlgType{
			sort:           algSort,
			entries:        entriesSub,
			entriesFmtInfo: RefParseFmtInfo(entriesFmtInfo),
			attributeList:  attributeList,
			parseFmtInfo:   RefParseFmtInfo(info),
			isInterface:    tok == "interface",
			isEnum:         tok == "enum",
		})

	case tok == "[":
		elementType, infoSub = ParseType(r)
		info = info.UnifyFmtInfo(r, infoSub)
		if info.err != nil {
			ret = nil
			return
		}

		_, infoSub = ReadTokenCheck(r, []string{"]"})
		info = info.UnifyFmtInfo(r, infoSub)
		if info.err != nil {
			ret = nil
			return
		}
		ret = RefArrayType(ArrayType{
			elementType:  elementType,
			parseFmtInfo: RefParseFmtInfo(info),
		})

	case tok == "&":
		targetType, infoSub = ParseType(r)
		info = info.UnifyFmtInfo(r, infoSub)
		if info.err != nil {
			ret = nil
			return
		}
		ret = RefDSLRefType(RefType{
			targetType:   targetType,
			parseFmtInfo: RefParseFmtInfo(info),
		})

	case tok == "{":
		keyType, infoSub = ParseType(r)
		info = info.UnifyFmtInfo(r, infoSub)
		if info.err != nil {
			ret = nil
			return
		}
		_, infoSub = ReadTokenCheck(r, []string{":"})
		info = info.UnifyFmtInfo(r, infoSub)
		if info.err != nil {
			ret = nil
			return
		}
		valueType, infoSub = ParseType(r)
		info = info.UnifyFmtInfo(r, infoSub)
		if info.err != nil {
			ret = nil
			return
		}
		_, infoSub = ReadTokenCheck(r, []string{"}"})
		info = info.UnifyFmtInfo(r, infoSub)
		if info.err != nil {
			ret = nil
			return
		}
		ret = RefMapType(MapType{
			keyType:      keyType,
			valueType:    valueType,
			parseFmtInfo: RefParseFmtInfo(info),
		})

	default:
		if IsIdent(tok) {
			ret = RefNamedType(NamedType{name: tok})
		} else {
			ret = nil
			info.err = r.GenParseError(fmt.Sprintf("Expected type; received \"%v\"", tok))
			return
		}
	}

	for {
		if tok, ok := PeekToken(r, false); ok && (tok == "?") {
			_, infoSub = ReadTokenCheck(r, []string{"?"})
			info = info.UnifyFmtInfo(r, infoSub)
			ret = RefOptionType(OptionType{
				valueType:    ret,
				parseFmtInfo: RefParseFmtInfo(info),
			})
		} else {
			break
		}
	}

	return
}

func ParseTypeDecl(r *ParseStream) (ret *TypeDecl, info ParseFmtInfo) {
	var infoSub ParseFmtInfo
	var declName string
	var declType Type

	_, infoSub = ReadTokenCheck(r, []string{"type"})
	info = info.UnifyFmtInfo(r, infoSub)
	if info.err != nil {
		ret = nil
		return
	}

	declName, infoSub = ReadToken(r)
	info = info.UnifyFmtInfo(r, infoSub)
	if info.err != nil {
		ret = nil
		return
	}

	declType, infoSub = ParseType(r)
	info = info.UnifyFmtInfo(r, infoSub)
	if info.err != nil {
		ret = nil
		return
	} else {
		ret = RefTypeDecl(TypeDecl{
			name:  declName,
			type_: declType,
		})
		return
	}
}

func TryParseTypeDecl(r *ParseStream) (ret *TypeDecl, info ParseFmtInfo) {
	r.Push()
	defer func() { r.Pop(info.err != nil) }()

	ret, info = ParseTypeDecl(r)
	return
}

func ParsePackageDecl(r *ParseStream) (ret *PackageDecl, info ParseFmtInfo) {
	var infoSub ParseFmtInfo
	var packageName string

	_, infoSub = ReadTokenCheck(r, []string{"package"})
	info = info.UnifyFmtInfo(r, infoSub)
	if info.err != nil {
		ret = nil
		return
	}

	packageName, infoSub = ReadToken(r)
	info = info.UnifyFmtInfo(r, infoSub)
	ret = RefPackageDecl(PackageDecl{
		name: packageName,
	})
	return
}

func TryParsePackageDecl(r *ParseStream) (ret *PackageDecl, info ParseFmtInfo) {
	r.Push()
	defer func() { r.Pop(info.err != nil) }()

	ret, info = ParsePackageDecl(r)
	return
}

func ReadStringLiteral(r *ParseStream) (ret string, info ParseFmtInfo) {
	var infoSub ParseFmtInfo

	ret = ""

	_, infoSub = ReadTokenCheck(r, []string{"\""})
	info = info.UnifyFmtInfo(r, infoSub)
	if info.err != nil {
		ret = ""
		return
	}

	for {
		c, ok := r.PeekExact(1)
		if !ok {
			ret = ""
			info.err = r.GenParseError("Error parsing string literal")
			return
		}
		r.Seek(1)
		if c == "\"" {
			return
		} else {
			ret += c
		}
	}
}

func ParseImportDecl(r *ParseStream) (ret *ImportDecl, info ParseFmtInfo) {
	var infoSub ParseFmtInfo
	var importName string
	var importPath string

	_, infoSub = ReadTokenCheck(r, []string{"import"})
	info = info.UnifyFmtInfo(r, infoSub)
	if info.err != nil {
		ret = nil
		return
	}

	importName, infoSub = ReadToken(r)
	info = info.UnifyFmtInfo(r, infoSub)
	if info.err != nil {
		ret = nil
		return
	}

	importPath, infoSub = ReadStringLiteral(r)
	info = info.UnifyFmtInfo(r, infoSub)
	if info.err != nil {
		ret = nil
		return
	}

	ret = RefImportDecl(ImportDecl{
		name: importName,
		path: importPath,
	})
	return
}

func TryParseImportDecl(r *ParseStream) (ret *ImportDecl, info ParseFmtInfo) {
	r.Push()
	defer func() { r.Pop(info.err != nil) }()

	ret, info = ParseImportDecl(r)
	return
}

func ParseDSLModuleFromStream(r *ParseStream) Module {
	ret := []Entry{}
	var decl Decl
	var info, infoSub ParseFmtInfo
	var noopEntriesSub []Entry
	var errAcc *ParseError_S = nil

	for {
		infoSub = r.ParseComments()
		if infoSub.err == nil {
			info, noopEntriesSub = info.UnifyFmtInfoCaptureComments(r, infoSub)
			for _, noopEntrySub := range noopEntriesSub {
				ret = append(ret, noopEntrySub)
			}
		} else {
			panic(infoSub.err)
		}

		if _, ok := PeekToken(r, false); !ok {
			Assert(info.err == nil)
			return Module{entries: ret}
		}

		decl, infoSub = TryParseTypeDecl(r)
		if infoSub.err == nil {
			info = info.UnifyFmtInfoRejectComments(r, infoSub)
			if info.err != nil {
				panic(info.err)
			}
			ret = append(ret, EntryDecl(decl))
			continue
		} else {
			errAcc = errAcc.UnifyError(infoSub.err)
		}

		decl, infoSub = TryParsePackageDecl(r)
		if infoSub.err == nil {
			info = info.UnifyFmtInfoRejectComments(r, infoSub)
			if info.err != nil {
				panic(info.err)
			}
			ret = append(ret, EntryDecl(decl))
			continue
		} else {
			errAcc = errAcc.UnifyError(infoSub.err)
		}

		decl, infoSub = TryParseImportDecl(r)
		if infoSub.err == nil {
			info = info.UnifyFmtInfoRejectComments(r, infoSub)
			if info.err != nil {
				panic(info.err)
			}
			ret = append(ret, EntryDecl(decl))
			continue
		} else {
			errAcc = errAcc.UnifyError(infoSub.err)
		}

		Assert(errAcc != nil)
		panic(errAcc)
	}
}

func ParseDSLModuleFromFile(file *os.File) Module {
	r := RefParseStream(ParseStream{
		rs: file,
		state: ParseStreamState{
			pos: 0,
		},
		lineMap: []LineInfo{},
		buffer:  bytes.NewBuffer([]byte{}),
	})
	return ParseDSLModuleFromStream(r)
}
