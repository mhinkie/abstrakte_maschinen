package main

import "strconv"
import "os"
import "encoding/binary"
import "errors"

var program = make([]*Command, 0)

/* Speichert die labels nachdem sie der optimizer über den code gelaufen ist */
var labels = make(map[string]int)

/* Wird für jedes neue Label erhöht */
var currentLabelIndex = 0

type OperandType uint8
type OpCode uint8

var UnusedOperand = makeOperand(Unused, 0, "")

const (
	Register      OperandType = 0
	LiteralInt    OperandType = 1
	LiteralString OperandType = 2
	RegisterArray OperandType = 3
	Unused        OperandType = 255
)

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

//go:generate stringer -type=OpCode

type Operand struct {
	Type        OperandType
	ValueInt    int
	ValueString string
	ValueArray  []int
}

/* bytecode */
type Command struct {
	CommandCode OpCode
	Output      Operand
	Op1         Operand
	Op2         Operand
	/* Hierzu muss noch das zugehörige label gefunden werden */
	UnresolvedJump bool
	/* Hint zum auflösen des jumps bzw. der function */
	ResolvingHint string
}

func (operand *Operand) String() string {
	switch operand.Type {
	case Register:
		return "\x1b[0;33mr" + strconv.Itoa(int(operand.ValueInt)) + "\x1b[0m"
	case LiteralInt:
		return "\x1b[0;34m## " + strconv.Itoa(int(operand.ValueInt)) + "\x1b[0m"
	case LiteralString:
		strcut := 10
		if len(operand.ValueString) < 10 {
			strcut = len(operand.ValueString)
		}

		return "\x1b[0;92m\"" + operand.ValueString[0:strcut] + ".." + "\"" + "\x1b[0m"
	case RegisterArray:
		if operand.ValueArray != nil {
			// Array von registern
			retString := "{"
			for _, reg := range operand.ValueArray {
				retString += " \x1b[0;33mr" + strconv.Itoa(int(reg)) + "\x1b[0m,"
			}
			retString += " }"
			return retString
		} else {
			return "NULL"
		}
	case Unused:
		return "_"
	default:
		return ""
	}
}

func (cmd *Command) String() string {
	cString := "(" + cmd.CommandCode.String() + ":\t\t"
	cString += cmd.Output.String() + " \t <== "
	cString += cmd.Op1.String() + ", "
	cString += cmd.Op2.String()
	cString += ")"
	if cmd.ResolvingHint != "" {
		cString += " \x1b[0;90m//" + cmd.ResolvingHint + "\x1b[0m"
	}

	return cString
}

/* ausgabe */
func (command *Command) Write() {
	cCodeB := byte(command.CommandCode)
	// 1 byte commandCode
	os.Stdout.Write([]byte{cCodeB})

	command.Output.Write()
	command.Op1.Write()
	command.Op2.Write()
}

func outputInt64(value int64) {
	oByte := make([]byte, 8)
	binary.PutVarint(oByte, value)

	os.Stdout.Write(oByte)
}

func (operand Operand) Write() {
	os.Stdout.Write([]byte{byte(operand.Type)})
	switch operand.Type {
	case Register:
		outputInt64(int64(operand.ValueInt))
	case LiteralInt:
		outputInt64(int64(operand.ValueInt))
	case LiteralString:
		outputInt64(int64(len(operand.ValueString)))
		os.Stdout.Write([]byte(operand.ValueString))
	case RegisterArray:
		// Länge ausgeben
		outputInt64(int64(len(operand.ValueArray)))
		for _, reg := range operand.ValueArray {
			outputInt64(int64(reg))
		}
	}
}

/* OUTPUT */
/* 1. mal durchgehen = register verdichten und optimieren */
func optimizeProgram() error {
	newProgram := make([]*Command, 0)
	newProgramSize := 0 /* die größe des ausgabeprogramms - alle weggelassenen befehle werden hier nicht mitgezählt */
	for _, command := range program {
		switch command.CommandCode {
		case OpStoreUnoptimized:
			// Das ist ein unoptimized store, das wird nicht ausgegeben, sondern beim befehl davor das destination register
			// auf das destination register dieses stores gesetzt
			if newProgram[newProgramSize-1].Output.Type == Register {
				// Register überehmen
				newProgram[newProgramSize-1].Output.ValueInt = command.Output.ValueInt
			} else {
				debugOutput("Cannot optimize store because destination type is not register")
			}
		case Label:
			// Die Position des labels im neuen programm wird gespeichert (das label selbst aber nicht ausgegeben)
			labels[command.Op1.ValueString] = newProgramSize
		default:
			// Befehl wird einfach hinzugefügt
			newProgram = append(newProgram, command)
			newProgramSize++
		}
	}

	// Das optimierte programm wird jetzt verwendet
	program = newProgram

	return nil
}

/* 2. durchgang = output */
func outputProgram() error {
	for _, command := range program {
		if command != nil {
			command.Resolve()
			command.Write()
		}
	}

	return nil
}

/* Löst alle unfertigen referenzen auf (labels usw) */
func (cmd *Command) Resolve() error {
	if cmd.UnresolvedJump {
		// Ein jump muss aufgelöst werden
		adr, ok := labels[cmd.ResolvingHint]

		if !ok {
			return errors.New("Label " + cmd.ResolvingHint + " not found during output!")
		}

		cmd.Op2.Type = LiteralInt
		cmd.Op2.ValueInt = adr
	}
	return nil
}

/* Labelnamen */
func NewLabel() string {
	labelName := "_L" + strconv.Itoa(currentLabelIndex)
	currentLabelIndex++
	return labelName
}

/* Programm */
func AppendCommand(command *Command) {
	program = append(program, command)
}

func ProgramLength() int {
	return len(program)
}

/* Operanden */

func makeOperand(Type OperandType, ValueInt int, ValueString string) Operand {
	return Operand{Type, ValueInt, ValueString, nil}
}

/* erzeugt commands:
NewEmptyCommand macht ein command ohne operanden
NewSingleValueCommand macht ein command mit dem Output-operanden gesetzt
NewDoubleValueCommand hat output und op1
NewTripleValueCommand hat output, op1 und op2
*/
func NewCommand(CommandCode OpCode, Output Operand, Op1 Operand, Op2 Operand) *Command {
	c := NewEmptyCommand(CommandCode)
	c.Output = Output
	c.Op1 = Op1
	c.Op2 = Op2

	return c
}

func NewEmptyCommand(CommandCode OpCode) *Command {
	c := Command{}
	c.CommandCode = CommandCode
	c.Output = UnusedOperand
	c.Op1 = UnusedOperand
	c.Op2 = UnusedOperand

	return &c
}

func NewSingleValueCommand(CommandCode OpCode, TypeOutput OperandType, ValueIntOutput int) *Command {
	c := NewEmptyCommand(CommandCode)
	c.Output = makeOperand(TypeOutput, ValueIntOutput, "")

	return c
}

func NewInputOnlyCommand(CommandCode OpCode, TypeOp1 OperandType, ValueIntOp1 int, ValueStringOp1 string) *Command {
	c := NewEmptyCommand(CommandCode)
	c.Op1 = makeOperand(TypeOp1, ValueIntOp1, ValueStringOp1)

	return c
}

func NewDoubleValueCommand(CommandCode OpCode, TypeOutput OperandType, ValueIntOutput int, TypeOp1 OperandType, ValueIntOp1 int, ValueStringOp1 string) *Command {
	c := NewSingleValueCommand(CommandCode, TypeOutput, ValueIntOutput)
	c.Op1 = makeOperand(TypeOp1, ValueIntOp1, ValueStringOp1)

	return c
}

func NewTripleValueCommand(CommandCode OpCode, TypeOutput OperandType, ValueIntOutput int, TypeOp1 OperandType, ValueIntOp1 int, ValueStringOp1 string, TypeOp2 OperandType, ValueIntOp2 int, ValueStringOp2 string) *Command {
	c := NewDoubleValueCommand(CommandCode, TypeOutput, ValueIntOutput, TypeOp1, ValueIntOp1, ValueStringOp1)
	c.Op2 = makeOperand(TypeOp2, ValueIntOp2, ValueStringOp2)

	return c
}
