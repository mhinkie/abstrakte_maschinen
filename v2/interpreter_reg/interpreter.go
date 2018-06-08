package main

import (
	"bufio"
	"encoding/binary"
	"fmt"
	"os"
)

/* SYSTEM PARAMTER */
const rspaceSize = 3000
const heapSize = 3000

/* ENDE SYSTEM PARAMETER */

/* SYSTEMSTATUS--------------------------------------------------------------------------- */
var commandPointer = 0

// Registerspace = komplett verfügbarer speicher
var rspace = make([]int, rspaceSize)

// Zeigt auf den start des aktuellen register frames
var rFramePointer = 0

// Registerframe = teil des registerspaces für diese funktion
// Startet mit 3 Werten für evt. Pointer usw. und dann mit den eigentlichen registern
// Ausnahme: Der globalscope (programmscope) startet bei 0 (und nicht -3 wegen den 3 werten)
var registerFrame = rspace[rFramePointer:len(rspace)]

// Der Teil des registerFrames, der die eigentlichen adressierbaren register beinhaltet
var registers = registerFrame[3:len(registerFrame)]

// Aktuelle Auslastung des rspaces (= höchste registeradresse absolut)
//var rspaceMax = 0

var heap = make([]byte, heapSize)

// Zeigt auf oberstes element des heaps
var heapMax = 0

/* ENDE SYSTEMSTATUS----------------------------------------------------------------------- */

/* schreibt einen int wert an angegebene memory pos */
func WriteInt(buffer []byte, value int) {
	binary.PutVarint(buffer, int64(value))
}

func ConvertInt(bytes []byte) int {
	val, _ := binary.Varint(bytes)
	return int(val)
}

func Analyze(fname string) error {
	err := Parse(fname)
	if err != nil {
		return err
	}

	for i, command := range program {
		if command != nil {
			fmt.Fprintf(os.Stderr, fmt.Sprintf("%03d: ", i)+command.String()+"\n")
		}
	}

	return nil
}

func Run(fname string, interpreter func() error) error {
	err := Parse(fname)
	if err != nil {
		return err
	}

	return interpreter()
}

/* Standard interpreter funktion */
func Interpret() error {
	for commandPointer < len(commandList) {
		err := commandImpl[commandList[commandPointer]]()
		if err != nil {
			return err
		}

		commandPointer++
	}

	return nil
}

func PrintStep() {
	fmt.Printf("%d: %s", commandPointer, program[commandPointer].String())
}

/* stepper */
func Step() error {
	for commandPointer < len(commandList) {
		// Warten auf eingabe
		PrintStep()
		reader := bufio.NewReader(os.Stdin)
		fmt.Print("[Enter] for next step")
		reader.ReadString('\n')

		err := commandImpl[commandList[commandPointer]]()

		if err != nil {
			return err
		}

		commandPointer++
	}

	return nil
}

func main() {
	args := os.Args[1:]

	var err error

	for i := 0; i < len(args); i++ {
		if args[i] == "-a" {
			// Analyze
			err = Analyze(args[i+1])
		} else if args[i] == "-s" {
			err = Run(args[i+1], Step)
		} else {
			err = Run(args[i], Interpret)
		} /*else if args[i] == "-d" {
			// Ausführen mit dump danach
			err = Run(args[i+1], interpreter)
			Dump()
		} else if args[i] == "-s" {
			breakPoint = -1
			err = Run(args[i+1], Step)
		} else if args[i] == "-b" {
			breakPoint, err = strconv.Atoi(args[i+1])
			breakPoint = breakPoint * 17
			err = Run(args[i+2], Step)
		} */

	}

	if err != nil {
		fmt.Fprintf(os.Stderr, "Es gab fehler: %s\n", err.Error())
		fmt.Fprintf(os.Stderr, "Aufgetreten bei Command %d\n", commandPointer)
	}

}
