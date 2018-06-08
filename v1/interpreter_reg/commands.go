package main

import (
	"errors"
	"fmt"
	"os"
)

type ArithmeticFunction func(int, int) int

//go:generate stringer -type=OpCode

const (
	OpProgEntry OpCode = 1
	OpProgExit  OpCode = 2
	OpFuncEntry OpCode = 3
	OpFuncExit  OpCode = 4
	OpExit      OpCode = 5
	OpReturn    OpCode = 6
	OpCall      OpCode = 7
	/* 20 aufwärts = alles was mit speichern und laden zu tun hat */
	OpStore  OpCode = 20
	OpAalloc OpCode = 21 /* heap alloc */
	OpHStore OpCode = 22 /* Store auf heap */
	OpHRead  OpCode = 23 /* read auf heap */
	OpLen    OpCode = 24 /* ließt länge eines heap-eintrags */
	/* 50 aufwärts = ausgabe */
	OpEchoI OpCode = 50
	OpEchoS OpCode = 51
	/* 80 aufwärts = kontrollfluss */
	OpJFalse OpCode = 80
	OpJump   OpCode = 81
	/* 100 aufwärts = arithmetik/boolean und andere operatoren */
	OpAdd    OpCode = 100
	OpSubt   OpCode = 101
	OpMult   OpCode = 102
	OpDiv    OpCode = 103
	OpConcat OpCode = 110
	OpGt     OpCode = 120
	OpLt     OpCode = 121
	OpEq     OpCode = 122
	OpNeq    OpCode = 123
	/* 200 aufwärts = alles was nur für debug oder zwischendarstellungen verwendet wird */
	/* OpStoreUnoptimized ist ein store, das im späteren verlauf wegoptimiert werden sollte.
	Im output-prozess wird geprüft, ob es wirklich wegoptimiert wurde, und wenn nicht ein fehler ausgegeben */
	OpStoreUnoptimized OpCode = 200
	Label              OpCode = 201
)

/* Die eigentlichen command-implementations */
// Implementierung der eigentlichen kommandos
var commandImpl = make([]func() error, 256)

func init() {
	// Alle werden als invalid gekennzeichnet
	for i := 0; i < len(commandImpl); i++ {
		commandImpl[i] = InvalidOp
	}

	// Jetzt werden die einzelnen Befehle festgelegt
	commandImpl[OpProgEntry] = ProgEntry
	commandImpl[OpProgExit] = ProgExit
	commandImpl[OpFuncEntry] = FuncEntry
	commandImpl[OpFuncExit] = FuncExit
	commandImpl[OpExit] = Exit
	commandImpl[OpReturn] = Return
	commandImpl[OpCall] = Call

	commandImpl[OpStore] = Store
	commandImpl[OpAalloc] = Aalloc
	commandImpl[OpHStore] = HStore
	commandImpl[OpHRead] = HRead
	commandImpl[OpLen] = Len

	commandImpl[OpEchoI] = EchoI
	commandImpl[OpEchoS] = EchoS

	commandImpl[OpJFalse] = JFalse
	commandImpl[OpJump] = Jump

	commandImpl[OpAdd] = func() error { return IntArith(func(a int, b int) int { return a + b }) }
	commandImpl[OpSubt] = func() error { return IntArith(func(a int, b int) int { return a - b }) }
	commandImpl[OpMult] = func() error { return IntArith(func(a int, b int) int { return a * b }) }
	commandImpl[OpDiv] = func() error { return IntArith(func(a int, b int) int { return a / b }) }

	commandImpl[OpConcat] = Concat

	commandImpl[OpGt] = func() error {
		return IntArith(func(a int, b int) int {
			if a > b {
				return 1
			} else {
				return 0
			}
		})
	}

	commandImpl[OpLt] = func() error {
		return IntArith(func(a int, b int) int {
			if a < b {
				return 1
			} else {
				return 0
			}
		})
	}

	commandImpl[OpEq] = func() error {
		return IntArith(func(a int, b int) int {
			if a == b {
				return 1
			} else {
				return 0
			}
		})
	}

	commandImpl[OpNeq] = func() error {
		return IntArith(func(a int, b int) int {
			if a != b {
				return 1
			} else {
				return 0
			}
		})
	}

}

func InvalidOp() error {
	return errors.New("Invalid Operation!")
}

func ProgEntry() error {
	// Passiert mal nichts
	//rspaceMax += program[commandPointer].Op1.ValueInt
	return nil
}

func ProgExit() error {
	//rspaceMax -= program[commandPointer].Op1.ValueInt
	commandPointer = len(commandList) // Damit das program beendet wird
	return nil
}

func Store() error {
	// Output von store ist ein register
	// Input ist entweder ein register oder ein literal
	switch program[commandPointer].Op1.Type {
	case Register:
		// Kopieren von einem ins andere
		registers[program[commandPointer].Output.ValueInt] = registers[program[commandPointer].Op1.ValueInt]
	case LiteralInt:
		// Integer Literal
		registers[program[commandPointer].Output.ValueInt] = program[commandPointer].Op1.ValueInt
	case LiteralString:
		// String Literal == Platz am Heap reservieren und dann dort rein schreiben
		// Ablage am heap: -8 byte von pointer aus liegt size in byte (als int64) und dann die bytes
		sLength := len(program[commandPointer].Op1.ValueString)
		WriteInt(heap[heapMax:heapMax+8], sLength)
		heapMax += 8
		stringLocation := heapMax
		// String kopieren (dst, src)
		copy(heap[heapMax:heapMax+sLength], program[commandPointer].Op1.ValueString)
		// Adresse in das angegebene register schreiben
		registers[program[commandPointer].Output.ValueInt] = stringLocation
		heapMax += sLength
	}
	return nil
}

func EchoI() error {
	// Will register oder integer literal
	switch program[commandPointer].Op1.Type {
	case Register:
		fmt.Print(registers[program[commandPointer].Op1.ValueInt])
	case LiteralInt:
		fmt.Print(program[commandPointer].Op1.ValueInt)
	}
	return nil
}

func EchoS() error {
	// Gibt einen string aus
	// Operand kann entweder Register (mit adresse einen strings) oder literal string sein
	switch program[commandPointer].Op1.Type {
	case Register:
		// Aus register lesen
		sAddress := registers[program[commandPointer].Op1.ValueInt]
		// Länge lesen
		sLength := ConvertInt(heap[sAddress-8 : sAddress])
		fmt.Print(string(heap[sAddress : sAddress+sLength]))
	case LiteralString:
		// Direkt ausgeben
		fmt.Print(program[commandPointer].Op1.ValueString)
	}

	return nil
}

/* template für alle arithmetik-funktionen */
func IntArith(fnc ArithmeticFunction) error {
	// Ich erwarte mir hier immer entweder LitInteger oder Register
	var op1 int
	var op2 int

	switch program[commandPointer].Op1.Type {
	case Register:
		op1 = registers[program[commandPointer].Op1.ValueInt]
	case LiteralInt:
		op1 = program[commandPointer].Op1.ValueInt
	}

	switch program[commandPointer].Op2.Type {
	case Register:
		op2 = registers[program[commandPointer].Op2.ValueInt]
	case LiteralInt:
		op2 = program[commandPointer].Op2.ValueInt
	}

	registers[program[commandPointer].Output.ValueInt] = fnc(op1, op2)

	return nil
}

func Concat() error {
	var op1 string
	var op2 string

	switch program[commandPointer].Op1.Type {
	case Register:
		// Aus register lesen
		sAddress := registers[program[commandPointer].Op1.ValueInt]
		// Länge lesen
		sLength := ConvertInt(heap[sAddress-8 : sAddress])
		op1 = string(heap[sAddress : sAddress+sLength])
	case LiteralString:
		// Direkt ausgeben
		op1 = program[commandPointer].Op1.ValueString
	}

	switch program[commandPointer].Op2.Type {
	case Register:
		// Aus register lesen
		sAddress := registers[program[commandPointer].Op2.ValueInt]
		// Länge lesen
		sLength := ConvertInt(heap[sAddress-8 : sAddress])
		op2 = string(heap[sAddress : sAddress+sLength])
	case LiteralString:
		// Direkt ausgeben
		op2 = program[commandPointer].Op2.ValueString
	}

	output := op1 + op2

	// Ablage
	// String Literal == Platz am Heap reservieren und dann dort rein schreiben
	// Ablage am heap: -8 byte von pointer aus liegt size in byte (als int64) und dann die bytes
	sLength := len(output)
	WriteInt(heap[heapMax:heapMax+8], sLength)
	heapMax += 8
	stringLocation := heapMax
	// String kopieren (dst, src)
	copy(heap[heapMax:heapMax+sLength], output)
	// Adresse in das angegebene register schreiben
	registers[program[commandPointer].Output.ValueInt] = stringLocation
	heapMax += sLength

	return nil
}

/* Jump wenn wert in register 0 ist */
func JFalse() error {
	/* Im Op1 steht das register das ausgewertet wird, im op2 steht die zieladresse */
	if registers[program[commandPointer].Op1.ValueInt] == 0 {
		// Jump: Nach dem command wird der commandpointer immer um 1 erhöht
		// deswegen wird die zieladresse -1 gesetzt
		commandPointer = program[commandPointer].Op2.ValueInt - 1
	}
	// Ansonsten wird nichts gemacht
	return nil
}

/* Jump immer (unconditional) */
func Jump() error {
	/* der commandPointer wird 1 vor die zieladresse gesetzt */
	/* zieladresse steht immer in op2 */
	commandPointer = program[commandPointer].Op2.ValueInt - 1

	return nil
}

/* Heap sachen */
func Aalloc() error {
	/* 8 + (len * bsize) auf heap reservieren */
	/* und length schreiben */
	/* in output register den pointer (nach der length) schreiben */

	len := program[commandPointer].Op1.ValueInt
	bsize := program[commandPointer].Op2.ValueInt

	// Length schreiben
	WriteInt(heap[heapMax:heapMax+8], len)
	heapMax += 8
	registers[program[commandPointer].Output.ValueInt] = heapMax
	heapMax += (len * bsize)

	return nil
}

func HStore() error {
	/* Output-register = register des arraypointers */
	/* Input 1 = register oder literal das gespeichert wird */
	/* Input 2 = register oder literal für den offset */
	var input int
	var offset int
	var arrayPointer int

	switch program[commandPointer].Op1.Type {
	case Register:
		input = registers[program[commandPointer].Op1.ValueInt]
	case LiteralInt:
		input = program[commandPointer].Op1.ValueInt
	}

	switch program[commandPointer].Op2.Type {
	case Register:
		offset = registers[program[commandPointer].Op2.ValueInt]
	case LiteralInt:
		offset = program[commandPointer].Op2.ValueInt
	}

	arrayPointer = registers[program[commandPointer].Output.ValueInt]

	/* der offset kommt immer im index-werten (muss also mal bsize multipliziert werden)
	array speichert immer 8 byte - also *8 */
	//println("Writing " + strconv.Itoa(input) + " to " + strconv.Itoa(arrayPointer+(offset*8)) + ":" + strconv.Itoa(arrayPointer+(offset*8)+8) + " arraypointer: " + strconv.Itoa(arrayPointer) + " offset: " + strconv.Itoa(offset))

	WriteInt(heap[arrayPointer+(offset*8):arrayPointer+(offset*8)+8], input)

	return nil
}

func HRead() error {
	/* Ließt aus heap
	ergebnis kommt in output register
	input 1 = array pointer (immer register)
	input 2 = offset (register oder literal)
	*/

	var offset int
	arrayPointer := registers[program[commandPointer].Op1.ValueInt]

	switch program[commandPointer].Op2.Type {
	case Register:
		offset = registers[program[commandPointer].Op2.ValueInt]
	case LiteralInt:
		offset = program[commandPointer].Op2.ValueInt
	}

	//println("Reading from " + strconv.Itoa(arrayPointer) + " offset " + strconv.Itoa(offset) + " = " + strconv.Itoa(arrayPointer+(offset*8)) + ":" + strconv.Itoa(arrayPointer+(offset*8)+8))
	registers[program[commandPointer].Output.ValueInt] = ConvertInt(heap[arrayPointer+(offset*8) : arrayPointer+(offset*8)+8])

	return nil
}

func Len() error {
	// OP1 = register mit arraypointer
	// OP2 = output register
	arrayPointer := registers[program[commandPointer].Op1.ValueInt]
	registers[program[commandPointer].Output.ValueInt] = ConvertInt(heap[arrayPointer-8 : arrayPointer])

	return nil
}

func Call() error {
	// Op1 = ein register array mit den registern in denen die parameter stehen
	// Op2 = der commandpointer zu dem gesprungen werden soll (behandeln wie bei jump)
	// Output = Das register in dem danach der return wert stehen soll

	// Vorgehensweise: Ausgehend vom return-wert register (kurz rrw) wird der neue register-frame für die funktion beginnen
	// es wird folgendes "geschrieben"
	// rrw bleibt unverändert
	// rrw + 1 = saved frame pointer (rFramePointer)
	// rrw + 2 = return instruction pointer (= zeigt eigentlich auf die call funktion weil der rip immer nach dem befehl erhöht wird)
	// rrw + 3 = start des neuen registers - arrays = ab hier liegen dann die werte die als funktionsparameter kommen
	rrw := program[commandPointer].Output.ValueInt
	registers[rrw+1] = rFramePointer
	registers[rrw+2] = commandPointer
	//fmt.Printf("saved commandpointer %d\n", commandPointer)

	// Parameter übernehmen
	for i, paramReg := range program[commandPointer].Op1.ValueArray {
		registers[rrw+3+i] = registers[paramReg]
		//fmt.Printf("copying value")
	}

	// Jetzt werden die register-frame slices angepasst (stack overflow wird einfach mal nicht geprüft)
	rFramePointer = rFramePointer + 3 + rrw /* alter framepointer + offset zu registerbeginn + rrw */
	registerFrame = rspace[rFramePointer:len(rspace)]
	registers = registerFrame[3:len(registerFrame)]

	//fmt.Printf("saved commandpointer (after manipulations): %d\n", registerFrame[2])

	// und jetzt kommt der eigentliche jump in die funktion
	commandPointer = program[commandPointer].Op2.ValueInt - 1 /* -1 weil nach fkt erhöht wird */

	// fertig
	return nil
}

var regused = 0

func FuncEntry() error {
	/* Funktionseintritt:
	OUT: 008: (OpFuncEntry:		_ 	 <== ## 3, _)
	Hier ist eigentlich nichts mehr zu tun (rframep, registerfram, registers... ist alles schon gesetzt)
	*/
	//fmt.Printf("+using at least %d\n", rFramePointer)
	return nil
}

func FuncExit() error {

	originalFramePointer := registerFrame[1]
	originalCommandPointer := registerFrame[2]

	// Wird die funktion verlassen (hier ohne return value),
	// wird als erstes der alte register frame wiederhergestellt:
	rFramePointer = originalFramePointer
	registerFrame = rspace[rFramePointer:len(rspace)]
	registers = registerFrame[3:len(registerFrame)]

	// Und dann der commandpointer wieder gesetzt
	commandPointer = originalCommandPointer /* zeigt 1 vor die zieladresse */
	//fmt.Printf("restored commandpointer: %d\n", commandPointer)

	//fmt.Printf("-using at least %d\n", rFramePointer)

	return nil
}

func Exit() error {
	/* exit aus dem Programm */
	/*
		OUT: 004: (OpStore:		r5 	 <== ## 10, _)
		OUT: 005: (OpExit:		_ 	 <== r5, _) */

	os.Exit(registers[program[commandPointer].Op1.ValueInt])

	return nil
}

func Return() error {
	/* OUT: 010: (OpReturn:		_ 	 <== r2, _) */
	// Schreibt den return-wert in das return-register und macht ein FuncExit
	registerFrame[0] = registers[program[commandPointer].Op1.ValueInt]

	return FuncExit()
}
