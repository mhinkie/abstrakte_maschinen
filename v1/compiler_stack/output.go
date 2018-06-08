package main

import "strconv"
import "os"
import "encoding/binary"

var program = make([]*Command, 0)

var literals = make([]Literal, 0)

/* aktuelle position (länge) des literalspeichers (in byte) */
/* pro literal wird ein int64 (8 byte) + die länge in byte geschrieben */
var litpos = 0

/* dieser index wird bei jeder verwendung eines hints erhöht */
var hintIndex = 0
var intSymbolIndex = 0

/* mapped die label-namen auf die position im code */
var labels = make(map[string]int, 0)

var functions = make(map[string]int, 0)

type Literal struct {
	length int64
	value  []byte
}

var commandCodes = map[string]uint8{
	"push":        1,
	"store":       2,
	"retrieve":    3,
	"echos":       4,
	"progentry":   5,
	"progexit":    6,
	"blockentry":  7,
	"blockexit":   8,
	"add":         9,
	"mult":        10,
	"subt":        11,
	"div":         12,
	"echoi":       13,
	"concat":      14,
	"gt":          15,
	"lt":          16,
	"neq":         17,
	"eq":          18,
	"jfalse":      19,
	"j":           20,
	"alloc":       21,
	"storei":      22,
	"retrievei":   23,
	"len":         24,
	"disc":        25,
	"call":        26,
	"funcentry":   27,
	"funcexit":    28,
	"blockreturn": 29,
}

/* bytecode */
type Command struct {
	commandCode uint8
	p1          int64
	p2          int64
	/* speichert ob irgendwelche referenzen noch nicht resolved sind */
	unresolvedFunction bool
	unresolvedJump     bool
	/* speichert einen namen der zum resolven dieser referenz verwendet wird (name der funktion/des jumps) */
	resolvingHint string
}

func outputInt64(value int64) {
	oByte := make([]byte, 8)
	binary.PutVarint(oByte, value)

	os.Stdout.Write(oByte)
}

func (command *Command) Write() {
	// 1 byte commandCode
	os.Stdout.Write([]byte{command.commandCode})

	// 8 byte p1
	outputInt64(command.p1)
	// 8 byte p2
	outputInt64(command.p2)

	// = in summe 17 byte
}

/* Löst alle noch offenen referenzen auf */
func (command *Command) Resolve(currentPos int) {
	if command.unresolvedJump {
		/* in der jump-label-tabelle nachschaun, wo der jump ist */
		command.p1 = int64(labels[command.resolvingHint] - currentPos)
		command.unresolvedJump = false
	}
	if command.unresolvedFunction {
		/* in der function-table nachschaun, wo die func ist */
		command.p1 = int64(functions[command.resolvingHint] - currentPos)
		command.unresolvedFunction = false
	}
}

/* Gibt das Programm aus und löst dabei noch nicht aufgelöste referenzen auf */
func outputProgram() {
	//uint64 werden als little endian ausgegeben
	// Erst literals ausgeben

	// 1. 8 byte = länge des literal speichers
	outputInt64(int64(litpos))

	// Dann literale ausgeben
	for _, literal := range literals {
		// länge ausgeben
		outputInt64(literal.length)

		// value ausgeben
		os.Stdout.Write(literal.value)
	}

	// Danach code ausgeben
	for i := range program {
		program[i].Resolve(i)
		program[i].Write()
	}
}

/* fügt string literal hinzu und gibt position zurück */
func AddLiteral(value string) int64 {
	lit := Literal{int64(len(value)), []byte(value)}
	curLitPos := litpos

	// nächste literal position
	litpos += 8 + len(value)

	// hinzufügen
	literals = append(literals, lit)

	return int64(curLitPos)
}

/* für output */
func (command Command) String() string {
	return "(" + strconv.Itoa(int(command.commandCode)) + ", " + strconv.Itoa(int(command.p1)) + ", " + strconv.Itoa(int(command.p2)) + ", " + strconv.FormatBool(command.unresolvedFunction) + ", " + strconv.FormatBool(command.unresolvedJump) + ", " + command.resolvingHint + ")"
}

func AddCommand(command Command) *Command {
	program = append(program, &command)
	return &command
}

/* erstellt ein neues command */
func NewCommand(commandName string, p1 int64, p2 int64) Command {
	commandCode := commandCodes[commandName]

	return Command{commandCode, p1, p2, false, false, ""}
}

func NewResHint() string {
	hintIndex++
	return "HINT" + strconv.Itoa(hintIndex)
}

func NewInternalSymbol() string {
	intSymbolIndex++
	return "INTSYM" + strconv.Itoa(intSymbolIndex)
}

/* speichert die aktuelle position im code als die position des labels ab
beim output wird die position dann für die jumps verwendet */
func SaveLabelPosition(labelName string) {
	// aktuelle länge ist die nächste position im kommando
	labels[labelName] = len(program)
}

func SaveFunctionPosition(funcName string) {
	functions[funcName] = len(program)
}
