/* zur verwaltung von symtables */

package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

/* offset zwischen symbolen */
const symOffset = 1

/* Liste von functions */
var Functions = make(map[string]*Function)

type SimpleType int

type PHPType int

func (pType PHPType) IsArray() bool {
	return pType > 10
}

func (pType PHPType) BaseType() PHPType {
	if pType.IsArray() {
		return pType - 10
	} else {
		return pType
	}
}

func (pType PHPType) String() string {
	if pType.IsArray() {
		return SimpleType(pType.BaseType()).String() + "[]"
	} else {
		return SimpleType(pType).String()
	}
}

func ArrayOf(phptype PHPType) PHPType {
	return phptype + 10
}

/* Für debug output sind die grundtypen hier noch aufgelistet */
const (
	Nothing SimpleType = 0
	Int     SimpleType = 1
	Str     SimpleType = 2
)

const (
	NoType  PHPType = 0
	Integer PHPType = 1
	String  PHPType = 2
)

//go:generate stringer -type=SimpleType

type BlockContent interface {
	Print(blockid string, level string)
	Process(block *Block) error
}

type Statement struct {
	command     string
	children    []Statement
	varType     PHPType
	childBlocks []Block
	/* registernummer */
	register int
}

/*
symbols sind identifiziert durch string-name
und speichern ausßerdem ein mphp int (= int64) als adresse des speicherbereichs
*/
type Block struct {
	symbols  map[string]int64
	symTypes map[string]PHPType
	/* mapped symbolNames auf registers */
	symRegisters  map[string]int
	statements    []BlockContent
	parent        *Block /* zeigt auf statischen vorgänger */
	nextSymPos    int64
	usedRegisters int
	functionBlock bool
	/* Bei function-blocks wird hier der typ des return-wertes gespeichert (zur prüfung usw...) */
	returnType PHPType
}

type Function struct {
	ReturnType          PHPType
	FunctionDeclaration *Statement
	NumArguments        int
}

func NewFunction(ReturnType PHPType, DeclaringStatement *Statement, NumArguments int) *Function {
	return &Function{ReturnType, DeclaringStatement, NumArguments}
}

/* Fügt die funktion hinzu und markiert sie, sodass sie später (am ende des codes) ausgegeben wird */
func AddFunction(name string, function *Function) error {
	_, ok := Functions[name]
	if ok {
		return errors.New("Function " + name + " declared multiple times!")
	} else {
		Functions[name] = function
		return nil
	}
}

/* Gibt neue registernumber zurück */
func (block *Block) NewRegister() int {
	if block.RegisterFrame() != block {
		fmt.Fprintf(os.Stderr, "\x1b[0;31mWARNING\x1b[0m: Register allocated in non-function-block\n")
	}
	newReg := block.usedRegisters
	block.usedRegisters++
	return newReg
}

/* Liefert den block, der für den aktuellen symbol frame verantwortlich ist
(ist immer der 1. function-block darüber, oder der programblock - wenn kein functionblock vorhanden)
*/
func (block *Block) RegisterFrame() *Block {
	if block.functionBlock {
		return block
	} else {
		return block.parent.RegisterFrame()
	}
}

/* Gibt den block zurück in dem das symbol definiert ist */
/* z.B.:
$abc = 123
if 1 == 1 {
	$abc = 333
	// hier ist der DefiningBlock der programblock für abc
}
*/
func (block *Block) DefiningBlock(name string) *Block {
	_, ok := block.symbols[name]
	if ok {
		return block
	} else {
		if block.parent != nil {
			return block.parent.DefiningBlock(name)
		} else {
			return nil
		}
	}
}

/* liefert den abstand des blocks in dem das symbol definiert ist zum aktuellen block
z.B. symbol ist im aktuellen block definiert = 0
symbol ist im parent block definiert = 1
symbol ist im parent vom parent definiert = 2
wenn symbol nicht definiert wird -1 zurückgegeben */
func (block *Block) getSymbolDistance(name string) int64 {
	_, ok := block.symbols[name]
	if ok {
		return 0
	} else {
		if block.parent != nil {
			pSymDistance := block.parent.getSymbolDistance(name)
			if pSymDistance != -1 {
				return pSymDistance + 1
			} else {
				return pSymDistance
			}
		} else {
			return -1
		}
	}
}

/* Nimmt an, dass symbol defined ist (gibt 9999999 ansonsten) */
func (block *Block) getSymbolPosition(name string) int64 {
	val, ok := block.symbols[name] /* ok = true wenn vorhanden */
	if ok {
		return val
	} else {
		if block.parent != nil {
			return block.parent.getSymbolPosition(name)
		} else {
			return 9999999
		}
	}
}

func (block *Block) isSymbolDefinedInCurrentScope(name string) bool {
	_, ok := block.symbols[name]
	return ok
}

/* schaut ob symbol in symboltabelle ist */
func (block *Block) isSymbolDefined(name string) bool {
	_, ok := block.symbols[name] /* ok = true wenn vorhanden */
	if ok {
		return true
	} else {
		if block.parent != nil {
			return block.parent.isSymbolDefined(name)
		} else {
			return false
		}
	}
}

func (block *Block) getSymbolType(name string) PHPType {
	_, ok := block.symbols[name] /* ok = true wenn vorhanden */
	if ok {
		return block.symTypes[name]
	} else {
		if block.parent != nil {
			return block.parent.getSymbolType(name)
		} else {
			return NoType
		}
	}
}

/* Das symbol register wird immer vom aktuellen registerframe geholt */
func (block *Block) getSymbolRegister(name string) int {
	val, ok := block.symRegisters[name]
	if ok {
		// Wenn das symbol in diesem block defined ist, wird es zurückgegeben
		return val
	} else {
		// Wenn nicht wird eins drüber probiert
		return block.parent.getSymbolRegister(name)
	}
}

func (block *Block) addSymbol(name string, symType PHPType, symRegister int) {
	if !block.isSymbolDefined(name) {
		debugSymbols("defining symbol " + name + " of type " + symType.String() + " at register: " + strconv.Itoa(symRegister))
		block.symbols[name] = block.nextSymPos
		block.symTypes[name] = symType
		block.symRegisters[name] = symRegister
		block.nextSymPos += symOffset
	} else {
		debugSymbols("symbol " + name + " already defined")
	}
}

func (block *Block) isEmpty() bool {
	return len(block.statements) == 0
}

func appendStatement(statements []Statement, toAppend Statement) []Statement {
	if statements == nil {
		statements = make([]Statement, 0)
	}

	return append(statements, toAppend)
}

func prependStatement(first Statement, rest []Statement) []Statement {
	if rest == nil {
		rest = make([]Statement, 0)
	}

	return append([]Statement{first}, rest...)
}

func newEmptyBlock() Block {
	debugSymbols("creating new empty block")

	var block = Block{make(map[string]int64), make(map[string]PHPType), make(map[string]int), []BlockContent{}, nil, 0, 0, false, NoType}

	return block
}

func newBlock(statements []Statement) Block {
	debugSymbols("creating new block")
	bContent := make([]BlockContent, len(statements))

	for i := range statements {
		bContent[i] = &statements[i]
	}

	var block = Block{make(map[string]int64), make(map[string]PHPType), make(map[string]int), bContent, nil, 0, 0, false, NoType}

	return block
}

func NewStatement(command string, children []Statement, sType PHPType) Statement {
	stmt := Statement{command, children, sType, nil, -1}
	stmt.childBlocks = make([]Block, 0)

	return stmt
}
