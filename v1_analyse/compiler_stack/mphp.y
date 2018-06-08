%{
package main

var programBlock Block

/* literal 1 */
var litOne = NewStatement("litint", []Statement{ NewStatement("1", nil, NoType) }, Integer)
%}

%union {
  value int
  text string
  block Block
  statements []Statement
  statement Statement
}

%token K_PROG_START K_PROG_END K_IF K_ELSE K_FUNC K_FOREACH K_FOR K_AS K_WHILE K_ECHO K_RETURN K_ARRAY INT IDFUNC IDVAR C_GT C_LT C_EQ C_ASSIGN C_NEQ C_INC C_DEC C_CONCAT C_PLUS C_MINUS C_AST C_PARA_L C_PARA_R C_BRAK_L C_BRAK_R C_CURL_L C_CURL_R C_DIV C_SEMIK C_COMMA STRING
%type <block> program block elseblock
%type <statements> statements arglist nearglist varlist nevarlist
%type <statement> statement expr varaccess varassign binop blockstatement conditional loop whileloop unaryop forloop cvarassign foreachloop singvar arrayaccess funccall idfunc funcdec 
%%
program: K_PROG_START statements K_PROG_END
			{
				$$ = newBlock($2)
				programBlock = $$
			}
;

statements: { $$ = nil }
			| statements statement C_SEMIK
			{ 
				$$ = appendStatement($1, $2)
			}
			| statements blockstatement
			{
				$$ = appendStatement($1, $2)
			}
;

statement: 	K_ECHO expr
			{ $$ = NewStatement("echo", appendStatement(nil, $2), NoType) }
			| K_RETURN expr
			{ $$ = NewStatement("return", appendStatement(nil, $2), NoType) }
			| varassign
			{ $$ = $1 }
			| funccall
			{ $$ = NewStatement("proccall", appendStatement(nil, $1), NoType) }
			| 
			{ $$ = Statement{} }/* kann auch leer sein */
;


/* ein statement dass einen block (geschw. klammern) beinhaltet */
blockstatement:	conditional
				{ $$ = $1 }
				| loop
				{ $$ = $1 }
				| funcdec
				{ $$ = $1 }
;

/* function declaration */
funcdec:	K_FUNC idfunc C_PARA_L varlist C_PARA_R block
			{
				$$ = NewStatement("func", prependStatement($2, $4), NoType)
				$$.childBlocks = append($$.childBlocks, $6)
			}
;

idfunc:		IDFUNC
			{
				$$ = NewStatement(yyText(yylex), nil, NoType)
			}
;

funccall:	idfunc C_PARA_L arglist C_PARA_R
			{
				$$ = NewStatement("funccall", prependStatement($1, $3), NoType)
			}
;

loop:		forloop
			{ $$ = $1 }
			| whileloop
			{ $$ = $1 }
			| foreachloop
			{ $$ = $1 }
;

forloop:	K_FOR C_PARA_L cvarassign C_SEMIK expr C_SEMIK cvarassign C_PARA_R block
			{
				$$ = NewStatement("for", []Statement{$3, $5, $7}, NoType)
				$$.childBlocks = append($$.childBlocks, $9)
			}
;

whileloop:	K_WHILE C_PARA_L expr C_PARA_R block
			{ 
				$$ = NewStatement("while", []Statement{$3}, NoType) 
				$$.childBlocks = append($$.childBlocks, $5)
			}
;

foreachloop:	K_FOREACH C_PARA_L singvar K_AS singvar C_PARA_R block
				{
					$$ = NewStatement("foreach", []Statement{$3, $5}, NoType)
					$$.childBlocks = append($$.childBlocks, $7)
				}
;

block:		C_CURL_L statements C_CURL_R
			{ $$ = newBlock($2) }
;

conditional:	K_IF C_PARA_L expr C_PARA_R block elseblock
				{
					$$ = NewStatement("if", []Statement{$3}, NoType)
					$$.childBlocks = append($$.childBlocks, $5)
					$$.childBlocks = append($$.childBlocks, $6)
				}
;

elseblock:	{ $$ = newEmptyBlock() }
			| K_ELSE block
			{ $$ = $2 }
;

expr:		C_PARA_L expr C_PARA_R
			{
				$$ = $2
			}
			| varaccess
			{
				$$ = NewStatement("access", []Statement{$1}, NoType)
			}
			| INT 
			{ $$ = NewStatement("litint", []Statement{ NewStatement(yyText(yylex), nil, NoType) }, Integer) }
			| STRING
			{ $$ = NewStatement("litstring", []Statement{ NewStatement(yyText(yylex)[1:len(yyText(yylex))-1], nil, NoType) }, String) }
			| K_ARRAY C_PARA_L arglist C_PARA_R
			{ $$ = NewStatement("aalloc", $3, NoType) }
			| binop 
			{ $$ = $1 }
			| funccall
			{ $$ = $1 }
;

arglist:	{ $$ = make([]Statement, 0)}
			| nearglist
			{ $$ = $1 }
;

nearglist:	expr
			{
				$$ = appendStatement(nil, $1)
			}
			| expr C_COMMA arglist
			{
				$$ = prependStatement($1, $3)
			}
;

varlist:	{ $$ = make([]Statement, 0) }
			| nevarlist
			{ $$ = $1 }
;

nevarlist:	singvar
			{ $$ = appendStatement(nil, $1) }
			| singvar C_COMMA varlist
			{ $$ = prependStatement($1, $3) }
;	

binop:		expr C_GT expr
			{ 
				$$ = NewStatement("gt", []Statement{ $1, $3 }, NoType)
			}
			| expr C_LT expr
			{ 
				$$ = NewStatement("lt", []Statement{ $1, $3 }, NoType)
			}
			| expr C_NEQ expr
			{ 
				$$ = NewStatement("neq", []Statement{ $1, $3 }, NoType)
			}
			| expr C_EQ expr
			{ 
				$$ = NewStatement("eq", []Statement{ $1, $3 }, NoType)
			}
			| expr C_CONCAT expr
			{ 
				$$ = NewStatement("concat", []Statement{ $1, $3 }, NoType)
			}
			| expr C_PLUS expr
			{ 
				$$ = NewStatement("plus", []Statement{ $1, $3 }, NoType)
			}
			| expr C_MINUS expr
			{ 
				$$ = NewStatement("minus", []Statement{ $1, $3 }, NoType) 
			}
			| expr C_AST expr
			{ 
				$$ = NewStatement("mult", []Statement{ $1, $3 }, NoType) 
			}
			| expr C_DIV expr
			{ 
				$$ = NewStatement("div", []Statement{ $1, $3 }, NoType)
			}
;

cvarassign:	{
				$$ = Statement{}
			}
			| varassign
			{
				$$ = $1
			}
;

varassign:	varaccess C_ASSIGN expr
			{ $$ = NewStatement("assign", []Statement{$1, $3}, NoType) }
			| varaccess unaryop
			{ 
				/* dem unary-op wird das linke kind = variablen-access gesetzt */
				/* der varaccess hat als variable die var-deklaration vom assign */
				$2.children[0] = NewStatement("access", []Statement{ $1 }, NoType)
				$$ = NewStatement("assign", []Statement{$1, $2}, NoType)
				/* aus $a++; wird also:
				$a = $a + 1
				Im gegensatz zu anderen sprachen kann $a++; daher nur als eigenständiges statement verwendet werden.
				*/ 
			}
;

unaryop:	C_INC 
			{
				/* hier wird ein halbvolles x + 1 zurückgegeben (x muss noch ausgefüllt werden) */
				$$ = NewStatement("plus", []Statement{Statement{}, litOne}, NoType)
			}
			| C_DEC
			{
				$$ = NewStatement("minus", []Statement{Statement{}, litOne}, NoType)
			}
;

varaccess:	singvar
			{
				$$ = $1
			}
			| arrayaccess
			{
				$$ = $1
			}
;	 	

singvar:	IDVAR
			{
				$$ = NewStatement("var", appendStatement(nil, NewStatement(yyText(yylex), nil, NoType)), NoType)
			}
;

arrayaccess:	varaccess C_BRAK_L expr C_BRAK_R
				{
					$$ = NewStatement("array", []Statement{$1, $3}, NoType)
				}
;
%%
