//line mphp.y:2
package main

import __yyfmt__ "fmt"

//line mphp.y:2
var programBlock Block

/* literal 1 */
var litOne = NewStatement("litint", []Statement{NewStatement("1", nil, NoType)}, Integer)

//line mphp.y:10
type yySymType struct {
	yys        int
	value      int
	text       string
	block      Block
	statements []Statement
	statement  Statement
}

const K_PROG_START = 57346
const K_PROG_END = 57347
const K_IF = 57348
const K_ELSE = 57349
const K_FUNC = 57350
const K_FOREACH = 57351
const K_FOR = 57352
const K_AS = 57353
const K_WHILE = 57354
const K_ECHO = 57355
const K_RETURN = 57356
const K_ARRAY = 57357
const INT = 57358
const IDFUNC = 57359
const IDVAR = 57360
const C_GT = 57361
const C_LT = 57362
const C_EQ = 57363
const C_ASSIGN = 57364
const C_NEQ = 57365
const C_INC = 57366
const C_DEC = 57367
const C_CONCAT = 57368
const C_PLUS = 57369
const C_MINUS = 57370
const C_AST = 57371
const C_PARA_L = 57372
const C_PARA_R = 57373
const C_BRAK_L = 57374
const C_BRAK_R = 57375
const C_CURL_L = 57376
const C_CURL_R = 57377
const C_DIV = 57378
const C_SEMIK = 57379
const C_COMMA = 57380
const STRING = 57381

var yyToknames = [...]string{
	"$end",
	"error",
	"$unk",
	"K_PROG_START",
	"K_PROG_END",
	"K_IF",
	"K_ELSE",
	"K_FUNC",
	"K_FOREACH",
	"K_FOR",
	"K_AS",
	"K_WHILE",
	"K_ECHO",
	"K_RETURN",
	"K_ARRAY",
	"INT",
	"IDFUNC",
	"IDVAR",
	"C_GT",
	"C_LT",
	"C_EQ",
	"C_ASSIGN",
	"C_NEQ",
	"C_INC",
	"C_DEC",
	"C_CONCAT",
	"C_PLUS",
	"C_MINUS",
	"C_AST",
	"C_PARA_L",
	"C_PARA_R",
	"C_BRAK_L",
	"C_BRAK_R",
	"C_CURL_L",
	"C_CURL_R",
	"C_DIV",
	"C_SEMIK",
	"C_COMMA",
	"STRING",
}
var yyStatenames = [...]string{}

const yyEofCode = 1
const yyErrCode = 2
const yyInitialStackSize = 16

//line mphp.y:267

//line yacctab:1
var yyExca = [...]int{
	-1, 1,
	1, -1,
	-2, 0,
}

const yyPrivate = 57344

const yyLast = 207

var yyAct = [...]int{

	94, 36, 31, 9, 67, 10, 14, 97, 3, 21,
	89, 62, 64, 28, 86, 49, 50, 52, 95, 51,
	29, 37, 53, 54, 55, 56, 40, 91, 112, 49,
	50, 52, 57, 51, 84, 107, 53, 54, 55, 56,
	96, 66, 38, 58, 41, 42, 57, 106, 92, 14,
	68, 60, 40, 61, 83, 59, 48, 65, 70, 47,
	69, 46, 71, 72, 73, 74, 75, 76, 77, 78,
	79, 81, 44, 49, 50, 52, 88, 51, 43, 15,
	53, 54, 55, 56, 27, 90, 23, 102, 2, 13,
	57, 99, 22, 19, 17, 39, 93, 104, 18, 12,
	45, 100, 98, 108, 103, 10, 14, 88, 111, 14,
	68, 110, 105, 113, 49, 50, 52, 11, 51, 6,
	35, 53, 54, 55, 56, 5, 85, 87, 49, 50,
	52, 57, 51, 63, 101, 53, 54, 55, 56, 49,
	50, 52, 82, 51, 1, 57, 53, 54, 55, 56,
	0, 80, 0, 0, 0, 16, 57, 20, 26, 24,
	0, 25, 7, 8, 0, 0, 23, 27, 34, 32,
	23, 27, 0, 0, 0, 0, 0, 0, 49, 50,
	52, 0, 51, 30, 109, 53, 54, 55, 56, 0,
	0, 0, 33, 4, 16, 57, 20, 26, 24, 0,
	25, 7, 8, 0, 0, 23, 27,
}
var yyPact = [...]int{

	84, -1000, -1000, 188, -1000, -24, -1000, 153, 153, -1000,
	-1000, -1000, -1000, -1000, 20, 48, 42, -1000, -1000, -1000,
	69, -1000, -1000, -1000, 31, 29, 26, -1000, -1000, 159,
	153, -6, -1000, -1000, 25, -1000, -1000, 159, 153, -1000,
	153, -1000, -1000, 153, 153, 11, 66, 153, 66, 153,
	153, 153, 153, 153, 153, 153, 153, 153, 120, 153,
	159, 109, 23, -1000, -4, 95, 66, -27, -1000, 54,
	16, 159, 159, 159, 159, 159, 159, 159, 159, 159,
	-1000, 17, -1000, -1000, 153, -16, 9, -1000, -31, 153,
	-16, 66, -1000, -1000, 80, -1000, -16, 66, 10, -1000,
	4, -1000, -16, 149, -1000, -1000, 66, -16, -1000, -1000,
	-3, -1000, -16, -1000,
}
var yyPgo = [...]int{

	0, 144, 0, 134, 8, 11, 133, 14, 127, 125,
	12, 2, 3, 120, 119, 117, 99, 98, 95, 94,
	4, 93, 9, 92, 1, 79, 89,
}
var yyR1 = [...]int{

	0, 1, 4, 4, 4, 9, 9, 9, 9, 9,
	14, 14, 14, 26, 25, 24, 16, 16, 16, 19,
	17, 21, 2, 15, 3, 3, 10, 10, 10, 10,
	10, 10, 10, 5, 5, 6, 6, 7, 7, 8,
	8, 13, 13, 13, 13, 13, 13, 13, 13, 13,
	20, 20, 12, 12, 18, 18, 11, 11, 22, 23,
}
var yyR2 = [...]int{

	0, 3, 0, 3, 2, 2, 2, 1, 1, 0,
	1, 1, 1, 6, 1, 4, 1, 1, 1, 9,
	5, 7, 3, 6, 0, 2, 3, 1, 1, 1,
	4, 1, 1, 0, 1, 1, 3, 0, 1, 1,
	3, 3, 3, 3, 3, 3, 3, 3, 3, 3,
	0, 1, 3, 2, 1, 1, 1, 1, 1, 4,
}
var yyChk = [...]int{

	-1000, -1, 4, -4, 5, -9, -14, 13, 14, -12,
	-24, -15, -16, -26, -11, -25, 6, -19, -17, -21,
	8, -22, -23, 17, 10, 12, 9, 18, 37, -10,
	30, -11, 16, 39, 15, -13, -24, -10, 22, -18,
	32, 24, 25, 30, 30, -25, 30, 30, 30, 19,
	20, 23, 21, 26, 27, 28, 29, 36, -10, 30,
	-10, -10, -5, -6, -10, -10, 30, -20, -12, -10,
	-22, -10, -10, -10, -10, -10, -10, -10, -10, -10,
	31, -5, 33, 31, 38, 31, -7, -8, -22, 37,
	31, 11, 31, -5, -2, 34, 31, 38, -10, -2,
	-22, -3, 7, -4, -2, -7, 37, 31, -2, 35,
	-20, -2, 31, -2,
}
var yyDef = [...]int{

	0, -2, 2, 9, 1, 0, 4, 0, 0, 7,
	8, 10, 11, 12, 0, 0, 0, 16, 17, 18,
	0, 56, 57, 14, 0, 0, 0, 58, 3, 5,
	0, 27, 28, 29, 0, 31, 32, 6, 0, 53,
	0, 54, 55, 33, 0, 0, 50, 0, 0, 0,
	0, 0, 0, 0, 0, 0, 0, 0, 0, 33,
	52, 0, 0, 34, 35, 0, 37, 0, 51, 0,
	0, 41, 42, 43, 44, 45, 46, 47, 48, 49,
	26, 0, 59, 15, 33, 0, 0, 38, 39, 0,
	0, 0, 30, 36, 24, 2, 0, 37, 0, 20,
	0, 23, 0, 9, 13, 40, 50, 0, 25, 22,
	0, 21, 0, 19,
}
var yyTok1 = [...]int{

	1,
}
var yyTok2 = [...]int{

	2, 3, 4, 5, 6, 7, 8, 9, 10, 11,
	12, 13, 14, 15, 16, 17, 18, 19, 20, 21,
	22, 23, 24, 25, 26, 27, 28, 29, 30, 31,
	32, 33, 34, 35, 36, 37, 38, 39,
}
var yyTok3 = [...]int{
	0,
}

var yyErrorMessages = [...]struct {
	state int
	token int
	msg   string
}{}

//line yaccpar:1

/*	parser for yacc output	*/

var (
	yyDebug        = 0
	yyErrorVerbose = false
)

type yyLexer interface {
	Lex(lval *yySymType) int
	Error(s string)
}

type yyParser interface {
	Parse(yyLexer) int
	Lookahead() int
}

type yyParserImpl struct {
	lval  yySymType
	stack [yyInitialStackSize]yySymType
	char  int
}

func (p *yyParserImpl) Lookahead() int {
	return p.char
}

func yyNewParser() yyParser {
	return &yyParserImpl{}
}

const yyFlag = -1000

func yyTokname(c int) string {
	if c >= 1 && c-1 < len(yyToknames) {
		if yyToknames[c-1] != "" {
			return yyToknames[c-1]
		}
	}
	return __yyfmt__.Sprintf("tok-%v", c)
}

func yyStatname(s int) string {
	if s >= 0 && s < len(yyStatenames) {
		if yyStatenames[s] != "" {
			return yyStatenames[s]
		}
	}
	return __yyfmt__.Sprintf("state-%v", s)
}

func yyErrorMessage(state, lookAhead int) string {
	const TOKSTART = 4

	if !yyErrorVerbose {
		return "syntax error"
	}

	for _, e := range yyErrorMessages {
		if e.state == state && e.token == lookAhead {
			return "syntax error: " + e.msg
		}
	}

	res := "syntax error: unexpected " + yyTokname(lookAhead)

	// To match Bison, suggest at most four expected tokens.
	expected := make([]int, 0, 4)

	// Look for shiftable tokens.
	base := yyPact[state]
	for tok := TOKSTART; tok-1 < len(yyToknames); tok++ {
		if n := base + tok; n >= 0 && n < yyLast && yyChk[yyAct[n]] == tok {
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}
	}

	if yyDef[state] == -2 {
		i := 0
		for yyExca[i] != -1 || yyExca[i+1] != state {
			i += 2
		}

		// Look for tokens that we accept or reduce.
		for i += 2; yyExca[i] >= 0; i += 2 {
			tok := yyExca[i]
			if tok < TOKSTART || yyExca[i+1] == 0 {
				continue
			}
			if len(expected) == cap(expected) {
				return res
			}
			expected = append(expected, tok)
		}

		// If the default action is to accept or reduce, give up.
		if yyExca[i+1] != 0 {
			return res
		}
	}

	for i, tok := range expected {
		if i == 0 {
			res += ", expecting "
		} else {
			res += " or "
		}
		res += yyTokname(tok)
	}
	return res
}

func yylex1(lex yyLexer, lval *yySymType) (char, token int) {
	token = 0
	char = lex.Lex(lval)
	if char <= 0 {
		token = yyTok1[0]
		goto out
	}
	if char < len(yyTok1) {
		token = yyTok1[char]
		goto out
	}
	if char >= yyPrivate {
		if char < yyPrivate+len(yyTok2) {
			token = yyTok2[char-yyPrivate]
			goto out
		}
	}
	for i := 0; i < len(yyTok3); i += 2 {
		token = yyTok3[i+0]
		if token == char {
			token = yyTok3[i+1]
			goto out
		}
	}

out:
	if token == 0 {
		token = yyTok2[1] /* unknown char */
	}
	if yyDebug >= 3 {
		__yyfmt__.Printf("lex %s(%d)\n", yyTokname(token), uint(char))
	}
	return char, token
}

func yyParse(yylex yyLexer) int {
	return yyNewParser().Parse(yylex)
}

func (yyrcvr *yyParserImpl) Parse(yylex yyLexer) int {
	var yyn int
	var yyVAL yySymType
	var yyDollar []yySymType
	_ = yyDollar // silence set and not used
	yyS := yyrcvr.stack[:]

	Nerrs := 0   /* number of errors */
	Errflag := 0 /* error recovery flag */
	yystate := 0
	yyrcvr.char = -1
	yytoken := -1 // yyrcvr.char translated into internal numbering
	defer func() {
		// Make sure we report no lookahead when not parsing.
		yystate = -1
		yyrcvr.char = -1
		yytoken = -1
	}()
	yyp := -1
	goto yystack

ret0:
	return 0

ret1:
	return 1

yystack:
	/* put a state and value onto the stack */
	if yyDebug >= 4 {
		__yyfmt__.Printf("char %v in %v\n", yyTokname(yytoken), yyStatname(yystate))
	}

	yyp++
	if yyp >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyS[yyp] = yyVAL
	yyS[yyp].yys = yystate

yynewstate:
	yyn = yyPact[yystate]
	if yyn <= yyFlag {
		goto yydefault /* simple state */
	}
	if yyrcvr.char < 0 {
		yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
	}
	yyn += yytoken
	if yyn < 0 || yyn >= yyLast {
		goto yydefault
	}
	yyn = yyAct[yyn]
	if yyChk[yyn] == yytoken { /* valid shift */
		yyrcvr.char = -1
		yytoken = -1
		yyVAL = yyrcvr.lval
		yystate = yyn
		if Errflag > 0 {
			Errflag--
		}
		goto yystack
	}

yydefault:
	/* default state action */
	yyn = yyDef[yystate]
	if yyn == -2 {
		if yyrcvr.char < 0 {
			yyrcvr.char, yytoken = yylex1(yylex, &yyrcvr.lval)
		}

		/* look through exception table */
		xi := 0
		for {
			if yyExca[xi+0] == -1 && yyExca[xi+1] == yystate {
				break
			}
			xi += 2
		}
		for xi += 2; ; xi += 2 {
			yyn = yyExca[xi+0]
			if yyn < 0 || yyn == yytoken {
				break
			}
		}
		yyn = yyExca[xi+1]
		if yyn < 0 {
			goto ret0
		}
	}
	if yyn == 0 {
		/* error ... attempt to resume parsing */
		switch Errflag {
		case 0: /* brand new error */
			yylex.Error(yyErrorMessage(yystate, yytoken))
			Nerrs++
			if yyDebug >= 1 {
				__yyfmt__.Printf("%s", yyStatname(yystate))
				__yyfmt__.Printf(" saw %s\n", yyTokname(yytoken))
			}
			fallthrough

		case 1, 2: /* incompletely recovered error ... try again */
			Errflag = 3

			/* find a state where "error" is a legal shift action */
			for yyp >= 0 {
				yyn = yyPact[yyS[yyp].yys] + yyErrCode
				if yyn >= 0 && yyn < yyLast {
					yystate = yyAct[yyn] /* simulate a shift of "error" */
					if yyChk[yystate] == yyErrCode {
						goto yystack
					}
				}

				/* the current p has no shift on "error", pop stack */
				if yyDebug >= 2 {
					__yyfmt__.Printf("error recovery pops state %d\n", yyS[yyp].yys)
				}
				yyp--
			}
			/* there is no state on the stack with an error shift ... abort */
			goto ret1

		case 3: /* no shift yet; clobber input char */
			if yyDebug >= 2 {
				__yyfmt__.Printf("error recovery discards %s\n", yyTokname(yytoken))
			}
			if yytoken == yyEofCode {
				goto ret1
			}
			yyrcvr.char = -1
			yytoken = -1
			goto yynewstate /* try again in the same state */
		}
	}

	/* reduction by production yyn */
	if yyDebug >= 2 {
		__yyfmt__.Printf("reduce %v in:\n\t%v\n", yyn, yyStatname(yystate))
	}

	yynt := yyn
	yypt := yyp
	_ = yypt // guard against "declared and not used"

	yyp -= yyR2[yyn]
	// yyp is now the index of $0. Perform the default action. Iff the
	// reduced production is ε, $1 is possibly out of range.
	if yyp+1 >= len(yyS) {
		nyys := make([]yySymType, len(yyS)*2)
		copy(nyys, yyS)
		yyS = nyys
	}
	yyVAL = yyS[yyp+1]

	/* consult goto table to find next state */
	yyn = yyR1[yyn]
	yyg := yyPgo[yyn]
	yyj := yyg + yyS[yyp].yys + 1

	if yyj >= yyLast {
		yystate = yyAct[yyg]
	} else {
		yystate = yyAct[yyj]
		if yyChk[yystate] != -yyn {
			yystate = yyAct[yyg]
		}
	}
	// dummy call; replaced with literal code
	switch yynt {

	case 1:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line mphp.y:24
		{
			yyVAL.block = newBlock(yyDollar[2].statements)
			yyVAL.block.functionBlock = true
			programBlock = yyVAL.block
		}
	case 2:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line mphp.y:31
		{
			yyVAL.statements = nil
		}
	case 3:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line mphp.y:33
		{
			yyVAL.statements = appendStatement(yyDollar[1].statements, yyDollar[2].statement)
		}
	case 4:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line mphp.y:37
		{
			yyVAL.statements = appendStatement(yyDollar[1].statements, yyDollar[2].statement)
		}
	case 5:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line mphp.y:43
		{
			yyVAL.statement = NewStatement("echo", appendStatement(nil, yyDollar[2].statement), NoType)
		}
	case 6:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line mphp.y:45
		{
			yyVAL.statement = NewStatement("return", appendStatement(nil, yyDollar[2].statement), NoType)
		}
	case 7:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line mphp.y:47
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 8:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line mphp.y:49
		{
			yyVAL.statement = NewStatement("proccall", appendStatement(nil, yyDollar[1].statement), NoType)
		}
	case 9:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line mphp.y:51
		{
			yyVAL.statement = Statement{}
		}
	case 10:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line mphp.y:57
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 11:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line mphp.y:59
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 12:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line mphp.y:61
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 13:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line mphp.y:66
		{
			yyVAL.statement = NewStatement("func", prependStatement(yyDollar[2].statement, yyDollar[4].statements), NoType)
			yyDollar[6].block.functionBlock = true
			yyVAL.statement.childBlocks = append(yyVAL.statement.childBlocks, yyDollar[6].block)
		}
	case 14:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line mphp.y:74
		{
			yyVAL.statement = NewStatement(yyText(yylex), nil, NoType)
		}
	case 15:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line mphp.y:80
		{
			yyVAL.statement = NewStatement("funccall", prependStatement(yyDollar[1].statement, yyDollar[3].statements), NoType)
		}
	case 16:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line mphp.y:86
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 17:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line mphp.y:88
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 18:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line mphp.y:90
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 19:
		yyDollar = yyS[yypt-9 : yypt+1]
		//line mphp.y:94
		{
			yyVAL.statement = NewStatement("for", []Statement{yyDollar[3].statement, yyDollar[5].statement, yyDollar[7].statement}, NoType)
			yyVAL.statement.childBlocks = append(yyVAL.statement.childBlocks, yyDollar[9].block)
		}
	case 20:
		yyDollar = yyS[yypt-5 : yypt+1]
		//line mphp.y:101
		{
			yyVAL.statement = NewStatement("while", []Statement{yyDollar[3].statement}, NoType)
			yyVAL.statement.childBlocks = append(yyVAL.statement.childBlocks, yyDollar[5].block)
		}
	case 21:
		yyDollar = yyS[yypt-7 : yypt+1]
		//line mphp.y:108
		{
			yyVAL.statement = NewStatement("foreach", []Statement{yyDollar[3].statement, yyDollar[5].statement}, NoType)
			yyVAL.statement.childBlocks = append(yyVAL.statement.childBlocks, yyDollar[7].block)
		}
	case 22:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line mphp.y:115
		{
			yyVAL.block = newBlock(yyDollar[2].statements)
		}
	case 23:
		yyDollar = yyS[yypt-6 : yypt+1]
		//line mphp.y:119
		{
			yyVAL.statement = NewStatement("if", []Statement{yyDollar[3].statement}, NoType)
			yyVAL.statement.childBlocks = append(yyVAL.statement.childBlocks, yyDollar[5].block)
			yyVAL.statement.childBlocks = append(yyVAL.statement.childBlocks, yyDollar[6].block)
		}
	case 24:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line mphp.y:126
		{
			yyVAL.block = newEmptyBlock()
		}
	case 25:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line mphp.y:128
		{
			yyVAL.block = yyDollar[2].block
		}
	case 26:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line mphp.y:132
		{
			yyVAL.statement = yyDollar[2].statement
		}
	case 27:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line mphp.y:136
		{
			yyVAL.statement = NewStatement("access", []Statement{yyDollar[1].statement}, NoType)
		}
	case 28:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line mphp.y:140
		{
			yyVAL.statement = NewStatement("litint", []Statement{NewStatement(yyText(yylex), nil, NoType)}, Integer)
		}
	case 29:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line mphp.y:142
		{
			yyVAL.statement = NewStatement("litstring", []Statement{NewStatement(yyText(yylex)[1:len(yyText(yylex))-1], nil, NoType)}, String)
		}
	case 30:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line mphp.y:144
		{
			yyVAL.statement = NewStatement("aalloc", yyDollar[3].statements, NoType)
		}
	case 31:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line mphp.y:146
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 32:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line mphp.y:148
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 33:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line mphp.y:151
		{
			yyVAL.statements = make([]Statement, 0)
		}
	case 34:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line mphp.y:153
		{
			yyVAL.statements = yyDollar[1].statements
		}
	case 35:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line mphp.y:157
		{
			yyVAL.statements = appendStatement(nil, yyDollar[1].statement)
		}
	case 36:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line mphp.y:161
		{
			yyVAL.statements = prependStatement(yyDollar[1].statement, yyDollar[3].statements)
		}
	case 37:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line mphp.y:166
		{
			yyVAL.statements = make([]Statement, 0)
		}
	case 38:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line mphp.y:168
		{
			yyVAL.statements = yyDollar[1].statements
		}
	case 39:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line mphp.y:172
		{
			yyVAL.statements = appendStatement(nil, yyDollar[1].statement)
		}
	case 40:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line mphp.y:174
		{
			yyVAL.statements = prependStatement(yyDollar[1].statement, yyDollar[3].statements)
		}
	case 41:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line mphp.y:178
		{
			yyVAL.statement = NewStatement("gt", []Statement{yyDollar[1].statement, yyDollar[3].statement}, NoType)
		}
	case 42:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line mphp.y:182
		{
			yyVAL.statement = NewStatement("lt", []Statement{yyDollar[1].statement, yyDollar[3].statement}, NoType)
		}
	case 43:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line mphp.y:186
		{
			yyVAL.statement = NewStatement("neq", []Statement{yyDollar[1].statement, yyDollar[3].statement}, NoType)
		}
	case 44:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line mphp.y:190
		{
			yyVAL.statement = NewStatement("eq", []Statement{yyDollar[1].statement, yyDollar[3].statement}, NoType)
		}
	case 45:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line mphp.y:194
		{
			yyVAL.statement = NewStatement("concat", []Statement{yyDollar[1].statement, yyDollar[3].statement}, NoType)
		}
	case 46:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line mphp.y:198
		{
			yyVAL.statement = NewStatement("plus", []Statement{yyDollar[1].statement, yyDollar[3].statement}, NoType)
		}
	case 47:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line mphp.y:202
		{
			yyVAL.statement = NewStatement("minus", []Statement{yyDollar[1].statement, yyDollar[3].statement}, NoType)
		}
	case 48:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line mphp.y:206
		{
			yyVAL.statement = NewStatement("mult", []Statement{yyDollar[1].statement, yyDollar[3].statement}, NoType)
		}
	case 49:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line mphp.y:210
		{
			yyVAL.statement = NewStatement("div", []Statement{yyDollar[1].statement, yyDollar[3].statement}, NoType)
		}
	case 50:
		yyDollar = yyS[yypt-0 : yypt+1]
		//line mphp.y:215
		{
			yyVAL.statement = Statement{}
		}
	case 51:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line mphp.y:219
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 52:
		yyDollar = yyS[yypt-3 : yypt+1]
		//line mphp.y:225
		{
			yyVAL.statement = NewStatement("assign", []Statement{yyDollar[1].statement, yyDollar[3].statement}, NoType)
		}
	case 53:
		yyDollar = yyS[yypt-2 : yypt+1]
		//line mphp.y:227
		{
			/* dem unary-op wird das linke kind = variablen-access gesetzt */
			/* der varaccess hat als variable die var-deklaration vom assign */
			yyDollar[2].statement.children[0] = NewStatement("access", []Statement{yyDollar[1].statement}, NoType)
			yyVAL.statement = NewStatement("assign", []Statement{yyDollar[1].statement, yyDollar[2].statement}, NoType)
			/* aus $a++; wird also:
			$a = $a + 1
			Im gegensatz zu anderen sprachen kann $a++; daher nur als eigenständiges statement verwendet werden.
			*/
		}
	case 54:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line mphp.y:237
		{
			/* hier wird ein halbvolles x + 1 zurückgegeben (x muss noch ausgefüllt werden) */
			yyVAL.statement = NewStatement("plus", []Statement{Statement{}, litOne}, NoType)
		}
	case 55:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line mphp.y:241
		{
			yyVAL.statement = NewStatement("minus", []Statement{Statement{}, litOne}, NoType)
		}
	case 56:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line mphp.y:247
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 57:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line mphp.y:251
		{
			yyVAL.statement = yyDollar[1].statement
		}
	case 58:
		yyDollar = yyS[yypt-1 : yypt+1]
		//line mphp.y:257
		{
			yyVAL.statement = NewStatement("var", appendStatement(nil, NewStatement(yyText(yylex), nil, NoType)), NoType)
		}
	case 59:
		yyDollar = yyS[yypt-4 : yypt+1]
		//line mphp.y:263
		{
			yyVAL.statement = NewStatement("array", []Statement{yyDollar[1].statement, yyDollar[3].statement}, NoType)
		}
	}
	goto yystack /* stack new state and value */
}
