package main

import "fmt"
import "strconv"
import "os"

/* DEBUG CONFIGURATION */
const OUTPUT_LEX = false
const OUTPUT_PAR = true
const OUTPUT_SYM = true
const OUTPUT_OUT = true

/*---------------------*/

func debugLex(text ...interface{}) {
	if OUTPUT_LEX {
		var output = "LEX: "
		for _, elem := range text {
			output += string(elem.(string)) + " "
		}
		fmt.Fprintf(os.Stderr, output+"\n")
	}
}

func debugParse(text ...interface{}) {
	if OUTPUT_PAR {
		var output = "PAR: "
		for _, elem := range text {
			output += string(elem.(string)) + " "
		}
		fmt.Fprintf(os.Stderr, output+"\n")
	}
}

func debugSymbols(text ...interface{}) {
	if OUTPUT_SYM {
		var output = "SYM: "
		for _, elem := range text {
			output += string(elem.(string)) + " "
		}
		fmt.Fprintf(os.Stderr, output+"\n")
	}
}

func debugOutput(text ...interface{}) {
	if OUTPUT_OUT {
		var output = "OUT: "
		for _, elem := range text {
			output += string(elem.(string)) + " "
		}
		fmt.Fprintf(os.Stderr, output+"\n")
	}
}

func debugTree(text ...interface{}) {
	if OUTPUT_OUT {
		var output = "TRE: "
		for _, elem := range text {
			output += string(elem.(string)) + " "
		}
		fmt.Fprintf(os.Stderr, output+"\n")
	}
}

func printTree(block Block) {
	block.Print("", "")
}

/* gibt zusammenfassung über symbole usw. nach tree-processing aus */
func printTreeProcessSummary(block Block) {
	debugTree("Tree processing summary:")
	debugTree("symbol tables:")

	printSymTable(block, "-")
}

func printLabels() {
	debugOutput("Labels:")
	for label, position := range labels {
		debugOutput(label + " at " + strconv.Itoa(position))
	}
}

func printFunctions() {
	debugOutput("Functions: ")
	for name, function := range Functions {
		debugOutput(fmt.Sprintf("%s %s (%d arguments)", name, function.ReturnType.String(), function.NumArguments))
	}
}

func printProgram() {
	commands := 0
	for i, command := range program {
		if command != nil {
			debugOutput(fmt.Sprintf("%03d: ", i) + command.String())
			commands++
		}
	}

	debugOutput("Program length: " + strconv.Itoa(commands))
}

func printSymTable(block Block, level string) {
	debugTree(level + " got " + strconv.Itoa(len(block.symbols)) + " symbols:")
	for symbol, _ := range block.symbols {
		debugTree(level + " " + symbol + " (" + block.getSymbolType(symbol).String() + "): " + strconv.Itoa(block.getSymbolRegister(symbol)))
	}
	for i := 0; i < len(block.statements); i++ {
		switch block.statements[i].(type) {
		case *Statement:
			stmt := (block.statements[i]).(*Statement)
			for _, cb := range stmt.childBlocks {
				printSymTable(cb, level+"-")
			}
		case *Block:
			printSymTable(*(block.statements[i].(*Block)), level+"-")
		}
	}
}

func (block Block) Print(blockid string, level string) {
	blocktype := "block"
	if block.functionBlock {
		blocktype = "functionblock"
	}
	blockid = blockid + "-"
	fmt.Fprintf(os.Stderr, "TRE: %s = %s with %d statements/childblocks:\n", blockid, blocktype, len(block.statements))

	for i := 0; i < len(block.statements); i++ {
		block.statements[i].Print(blockid, strconv.Itoa(i))
	}

}

func (statement Statement) Print(blockid string, level string) {
	var pReg string
	if statement.register == -1 {
		pReg = ""
	} else {
		pReg = " [" + strconv.Itoa(statement.register) + "]"
	}
	fmt.Fprintf(os.Stderr, "TRE: %s %s - %s%s (%d children)\n", blockid, level, string(statement.command), pReg, len(statement.children))

	for i := 0; i < len(statement.children); i++ {
		statement.children[i].Print(blockid, level+"."+strconv.Itoa(i))
	}

	if statement.childBlocks != nil {
		for i := 0; i < len(statement.childBlocks); i++ {
			statement.childBlocks[i].Print(blockid, level+"."+strconv.Itoa(i))
		}
	}
}

func Assert(message string, condition bool) {
	if !condition {
		fmt.Fprintf(os.Stderr, message+"\n")
		os.Exit(1)
	}
}
