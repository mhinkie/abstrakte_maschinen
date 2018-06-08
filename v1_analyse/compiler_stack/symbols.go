/* zur verwaltung von symtables */

package main

import "strconv"

/* offset zwischen symbolen */
const symOffset = 1

type PHPType int

func (pType PHPType) String() string {
	return strconv.Itoa(int(pType))
}

const (
	NoType      PHPType = 0
	Integer     PHPType = 1
	String      PHPType = 2
	ArrayInt    PHPType = 3
	ArrayString PHPType = 4
)

type BlockContent interface {
	Print(blockid string, level string)
	Process(block *Block) error
}

type Statement struct {
	command     string
	children    []Statement
	varType     PHPType
	childBlocks []Block
}

/*
symbols sind identifiziert durch string-name
und speichern ausßerdem ein mphp int (= int64) als adresse des speicherbereichs
*/
type Block struct {
	symbols      map[string]int64
	symTypes     map[string]PHPType
	statements   []BlockContent
	parent       *Block /* zeigt auf statischen vorgänger */
	nextSymPos   int64
	funcExitName string
	debugHelper  string
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

func (block *Block) isFunctionBlock() bool {
	return block.funcExitName != ""
}

/* Ein funktionsblock hat ein label wo die exit-operationen ausgeführt werden */
func (block *Block) getFuncExitLabel() string {
	if block.funcExitName == "" {
		if block.parent != nil {
			return block.parent.getFuncExitLabel()
		} else {
			return "INVALID LABEL"
		}
	} else {
		return block.funcExitName
	}
}

func (block *Block) addSymbol(name string, symType PHPType) {
	if !block.isSymbolDefined(name) {
		debugSymbols("defining symbol " + name + " of type " + symType.String())
		block.symbols[name] = block.nextSymPos
		block.symTypes[name] = symType
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

	var block = Block{make(map[string]int64), make(map[string]PHPType), []BlockContent{}, nil, 0, "", ""}

	return block
}

func newBlock(statements []Statement) Block {
	debugSymbols("creating new block")
	bContent := make([]BlockContent, len(statements))

	for i := range statements {
		bContent[i] = &statements[i]
	}

	var block = Block{make(map[string]int64), make(map[string]PHPType), bContent, nil, 0, "", ""}

	return block
}

func NewStatement(command string, children []Statement, sType PHPType) Statement {
	stmt := Statement{command, children, sType, nil}
	stmt.childBlocks = make([]Block, 0)

	return stmt
}
