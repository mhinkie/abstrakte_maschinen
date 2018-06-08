package main

import (
	"errors"
	"fmt"
	"os"
	"strconv"
)

/* DEBUG CONFIGURATION */
const OUTPUT_LEX = true
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

func printLiterals() {
	debugOutput("Literals:")

	curPos := 0

	for _, literal := range literals {
		debugOutput(strconv.Itoa(curPos) + ": " + string(literal.value) + " (length: " + strconv.Itoa(int(literal.length)) + ")")
		curPos += 8 + int(literal.length)
	}
}

/* gibt programmcodes aus */
func printProgram() {
	debugOutput("Program: (cCode, p1, p2, unresFunc, unresJump, resHint)")

	for i, command := range program {
		debugOutput(fmt.Sprintf("%03d: ", i) + command.String())
	}
}

/* gibt zusammenfassung über symbole usw. nach tree-processing aus */
func printTreeProcessSummary(block Block) {
	debugTree("Tree processing summary:")
	debugTree("symbol tables:")

	printSymTable(block, "-")
}

func printSymTable(block Block, level string) {
	debugTree(level + " got " + strconv.Itoa(len(block.symbols)) + " symbols:")
	for symbol, position := range block.symbols {
		debugTree(level + " " + symbol + ": " + strconv.Itoa(int(position)))
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

func getCommandName(code uint8) (string, error) {
	for k, v := range commandCodes {
		if v == code {
			return k, nil
		}
	}

	return "", errors.New("command not found for code " + strconv.Itoa(int(code)))
}

func (block Block) Print(blockid string, level string) {
	blockid = blockid + "-"
	if block.functionBlock {
		fmt.Fprintf(os.Stderr, "TRE: %s = functionblock with %d statements/childblocks:\n", blockid, len(block.statements))
	} else {
		fmt.Fprintf(os.Stderr, "TRE: %s = block with %d statements/childblocks:\n", blockid, len(block.statements))
	}

	for i := 0; i < len(block.statements); i++ {
		block.statements[i].Print(blockid, strconv.Itoa(i))
	}

}

func (statement Statement) Print(blockid string, level string) {
	fmt.Fprintf(os.Stderr, "TRE: %s %s - %s (%d children)\n", blockid, level, string(statement.command), len(statement.children))

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

func printLabels() {
	fmt.Fprintf(os.Stderr, "TRE: Labels:\n")
	for key, value := range labels {
		fmt.Fprintf(os.Stderr, "TRE: %s: pos %d\n", key, value)
	}
}

func printFunctions() {
	fmt.Fprintf(os.Stderr, "TRE: Functions:\n")
	for key, value := range functions {
		fmt.Fprintf(os.Stderr, "TRE: %s: pos %d\n", key, value)
	}
}
