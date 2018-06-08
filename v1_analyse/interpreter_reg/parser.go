/* Verantwortlich für parsend des bytecodes */

package main

import (
	"bufio"
	"encoding/binary"
	"errors"
	"io"
	"os"
	"strconv"
)

func Parse(filename string) error {
	f, err := os.Open(filename)
	if err != nil {
		return err
	}

	defer f.Close()
	reader := bufio.NewReader(f)

	comm, err := ParseCommand(reader)
	for err != io.EOF {
		if err != nil && err != io.EOF {
			return err
		}
		if comm != nil {
			program = append(program, comm)
			commandList = append(commandList, comm.CommandCode)
		}
		comm, err = ParseCommand(reader)
	}

	return nil
}

func ParseCommand(f *bufio.Reader) (*Command, error) {
	// 1 byte commandCode
	cCode := make([]byte, 1)
	_, err := io.ReadFull(f, cCode)
	if err != nil {
		return nil, err
	}

	opOut, err := ParseOperand(f)
	if err != nil {
		return nil, err
	}

	op1, err := ParseOperand(f)
	if err != nil {
		return nil, err
	}

	op2, err := ParseOperand(f)
	if err != nil {
		return nil, err
	}

	return NewCommand(OpCode(uint8(cCode[0])), *opOut, *op1, *op2), nil
}

func ParseOperand(reader *bufio.Reader) (*Operand, error) {
	/* Typ lesen = 1 byte */
	cType := make([]byte, 1)
	_, err := io.ReadFull(reader, cType)
	if err != nil {
		return nil, err
	}

	switch OperandType(cType[0]) {
	case Register:
		// Int 64 lesen
		intBuffer := make([]byte, 8)
		_, err = io.ReadFull(reader, intBuffer)
		if err != nil {
			return nil, err
		}
		op := makeOperand(Register, ConvertInt(intBuffer), "")
		return &op, nil
	case LiteralInt:
		// Int 64 lesen
		intBuffer := make([]byte, 8)
		_, err = io.ReadFull(reader, intBuffer)
		if err != nil {
			return nil, err
		}
		op := makeOperand(LiteralInt, ConvertInt(intBuffer), "")
		return &op, nil
	case LiteralString:
		intBuffer := make([]byte, 8)
		_, err = io.ReadFull(reader, intBuffer)
		if err != nil {
			return nil, err
		}

		size := ConvertInt(intBuffer)

		bString := make([]byte, size)
		_, err = io.ReadFull(reader, bString)
		if err != nil {
			return nil, err
		}

		op := makeOperand(LiteralString, 0, string(bString))
		return &op, nil
	case RegisterArray:
		intBuffer := make([]byte, 8)
		_, err = io.ReadFull(reader, intBuffer)
		if err != nil {
			return nil, err
		}

		size := ConvertInt(intBuffer)

		regArray := make([]int, size)
		for i := 0; i < size; i++ {
			_, err = io.ReadFull(reader, intBuffer)
			if err != nil {
				return nil, err
			}
			regArray[i] = ConvertInt(intBuffer)
		}

		op := makeOperand(RegisterArray, 0, "")
		op.ValueArray = regArray
		return &op, nil
	case Unused:
		// Nichts passiert
		op := makeOperand(Unused, 0, "")
		return &op, nil
	default:
		return nil, errors.New("Invalid operandtype: " + strconv.Itoa(int(cType[0])))
	}
}

/*





 */
/* Übernommen aus output.go */
var program = make([]*Command, 0)

/* Hier sind nur die opcodes gespeichert für einfacheren zugriff */
var commandList = make([]OpCode, 0)

var UnusedOperand = makeOperand(Unused, 0, "")

type OperandType uint8
type OpCode uint8

const (
	Register      OperandType = 0
	LiteralInt    OperandType = 1
	LiteralString OperandType = 2
	RegisterArray OperandType = 3
	Unused        OperandType = 255
)

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

func makeOperand(Type OperandType, ValueInt int, ValueString string) Operand {
	return Operand{Type, ValueInt, ValueString, nil}
}

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
