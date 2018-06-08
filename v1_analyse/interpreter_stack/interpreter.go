package main

import (
	"bufio"
	"bytes"
	"encoding/binary"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
	"time"
)

type InstructionInfo struct {
	totalTime int64
	numCalled int64
}

var Stats = make(map[byte]*InstructionInfo)

/* in 8 byte blöcken */
const initialStackSize = 5000

/* in byte */
const initialMemSize = 5000 * 8
const commandSize = 17

type Memory struct {
	memory      []byte
	heapPointer int
}

var breakPoint int

var memory = Memory{make([]byte, initialMemSize), 0}
var stack = make([]int64, initialStackSize)
var stackPointer = 0
var framePointer = 0

var commands []byte

var commandPointer = 0

// Implementierung der eigentlichen kommandos
var commandImpl = []func() (newCommandPointer int, err error){
	/* 0 - invalid */
	func() (int, error) {
		return 0, errors.New("Invalid command: 0 at position " + strconv.Itoa(commandPointer/17))
	},
	/* 1 - push */
	func() (int, error) {
		// 1. arg lesen und auf stack legen
		stack[stackPointer] = P1(commandPointer)
		stackPointer++
		return Advance(commandPointer), nil
	},
	/* 2 - store */
	func() (int, error) {
		// stackinhalt in lokal variable schreiben
		// 1. parameter = position von framepointer aus (wenn 2. variable gesetzt ist, gibt diese die zahl an elternblocks an die besucht werden müssen)

		distance := P2(commandPointer)
		varFramePointer := framePointer
		/* den eigentlichen framepointer der für diese variable zuständig ist herausfinden */
		/* P2 mal dereferenzieren */
		for ; distance > 0; distance-- {
			// der framepointer des übergeordneten blocks liegt immer 1 unter dem aktuellen framepointer
			varFramePointer = int(stack[varFramePointer-1])
		}

		//println("Dereferenced framepointer (store): " + strconv.Itoa(varFramePointer))

		stackPointer--
		stack[varFramePointer+int(P1(commandPointer))] = stack[stackPointer]

		return Advance(commandPointer), nil
	},
	/* 3 - retrieve */
	func() (int, error) {
		// speicherinhalt auf stack laden
		distance := P2(commandPointer)
		varFramePointer := framePointer
		/* den eigentlichen framepointer der für diese variable zuständig ist herausfinden */
		/* P2 mal dereferenzieren */
		for ; distance > 0; distance-- {
			// der framepointer des übergeordneten blocks liegt immer 1 unter dem aktuellen framepointer
			varFramePointer = int(stack[varFramePointer-1])
		}

		//println("Dereferenced framepointer (retr): " + strconv.Itoa(varFramePointer))

		stack[stackPointer] = stack[varFramePointer+int(P1(commandPointer))]
		stackPointer++
		return Advance(commandPointer), nil
	},
	/* 4 - echos (string) */
	func() (int, error) {
		// top of stack ist adresse des auszugebenen string
		// strings sind immer mit dem size attribut als erste 8 byte gespeichert
		stackPointer--
		adr := stack[stackPointer]
		size := memory.ReadInt64At(adr)

		// die ersten 8 byte (die size) wird übersprungen und dann der string gelesen
		fmt.Fprintf(os.Stdout, string(memory.memory[int(adr)+8:int(adr)+8+int(size)]))
		return Advance(commandPointer), nil
	},
	/* 5 - programentry */
	func() (int, error) {
		/* Programmstart = platz für lokal variablen (= globals) am stack reservieren */
		/* 1. parameter des kommandos = anzahl an variablen */
		stackPointer = framePointer + int(P1(commandPointer))
		return Advance(commandPointer), nil
	},
	/* 6 - programexit */
	func() (int, error) {
		/* jetzt noch nichts spannendes: zukünfig soll hier der return wert zurückgegeben werden */
		return Advance(commandPointer), nil
	},
	/* 7 - blockentry */
	func() (int, error) {
		/* der framepointer wird auf den stack gelegt (da ein neuer angelegt wird) */
		stack[stackPointer] = int64(framePointer)
		stackPointer++
		/* der neue framepointer zeigt auf die zelle nach dem alten stackpointer */
		framePointer = stackPointer
		/* der neue stackpointer wird jetzt um die anzahl an lokalen variablen (P1) vorgeschoben */
		stackPointer = framePointer + int(P1(commandPointer))

		/*println("nach blockentry")
		println("sp: " + strconv.Itoa(stackPointer))
		println("fp: " + strconv.Itoa(framePointer))*/

		return Advance(commandPointer), nil
	},
	/* 8 - blockexit */
	func() (int, error) {
		/* der framepointer soll auf dem alten gespeicherten framepointer stehen (in der zelle unter dem aktuellen framepointer gespeichert) */
		/* der stackpointer soll auf die zelle zeigen, in der der alte framepointer gespeichert wurde = top of stack vom vorhergehenden block */
		//stackPointer = stackPointer - int(P1(commandPointer)) /* anzahl lokale vars */ - 1 /* gespeicherter frameppointer */
		//framePointer = int(stack[framePointer-1])

		// Neu ohne P1:
		stackPointer = framePointer - 1
		framePointer = int(stack[framePointer-1])

		return Advance(commandPointer), nil
	},
	/* 9 - add */
	func() (int, error) {
		/* am stack liegen 2 werte - diese werden addiert - die 2 werte gepopped */
		stack[stackPointer-2] = stack[stackPointer-2] + stack[stackPointer-1]

		/* der stack ist jetzt eins weniger */
		stackPointer--

		return Advance(commandPointer), nil
	},
	/* 10 - mult */
	func() (int, error) {
		/* am stack liegen 2 werte - diese werden multipli. - die 2 werte gepopped */
		stack[stackPointer-2] = stack[stackPointer-2] * stack[stackPointer-1]

		/* der stack ist jetzt eins weniger */
		stackPointer--

		return Advance(commandPointer), nil
	},
	/* 11 - subt */
	func() (int, error) {
		/* am stack liegen 2 werte - diese werden subtrahiert - die 2 werte gepopped */
		stack[stackPointer-2] = stack[stackPointer-2] - stack[stackPointer-1]

		/* der stack ist jetzt eins weniger */
		stackPointer--

		return Advance(commandPointer), nil
	},
	/* 12 - div */
	func() (int, error) {
		/* am stack liegen 2 werte - diese werden multipli. - die 2 werte gepopped */
		stack[stackPointer-2] = stack[stackPointer-2] / stack[stackPointer-1]

		/* der stack ist jetzt eins weniger */
		stackPointer--

		return Advance(commandPointer), nil
	},
	/* 13 - echoi (int) */
	func() (int, error) {
		// top of stack ist integer das ausgegeben werden soll
		stackPointer--
		fmt.Fprintf(os.Stdout, strconv.Itoa(int(stack[stackPointer])))

		return Advance(commandPointer), nil
	},
	/* 14 - concat */
	func() (int, error) {
		/* erwartet sich als 2 oberste elemente am stack 2 heap-adressen von strings */
		/* diese werden zusammengehängt und auf dem heap abgelegt */
		/* danach wird die heap-adresse des ergebnisses auf den stack gelegt */

		// gleich auf den nächsten freien platz schreiben
		targetAddress := memory.heapPointer

		// Der 2. string (als oberster am stack)
		stackPointer--
		adr := stack[stackPointer]
		sLen := memory.ReadInt64At(adr)
		string2 := string(memory.memory[int(adr)+8 : int(adr)+8+int(sLen)])

		// Der 1. string (darunter am stack)
		stackPointer--
		adr = stack[stackPointer]
		sLen = memory.ReadInt64At(adr)
		string1 := string(memory.memory[int(adr)+8 : int(adr)+8+int(sLen)])

		// Zusammenhängen
		var buffer bytes.Buffer
		buffer.WriteString(string1)
		buffer.WriteString(string2)

		// Platz reservieren am heap
		memory.heapPointer += buffer.Len() + 8 // das +8 ist weil die größe auch geschrieben wird
		// größe schreiben
		memory.WriteInt64To(int64(targetAddress), int64(buffer.Len()))
		// Schreiben
		buffer.Read(memory.memory[targetAddress+8 : targetAddress+8+buffer.Len()])

		// jetzt noch die string-adresse am stack schreiben
		stack[stackPointer] = int64(targetAddress)
		stackPointer++

		return Advance(commandPointer), nil
	},
	/* 15 - gt */
	func() (int, error) {
		/* am stack liegen 2 werte */
		a := stack[stackPointer-2]
		b := stack[stackPointer-1]
		if a > b {
			stack[stackPointer-2] = 1
		} else {
			stack[stackPointer-2] = 0
		}

		/* der stack ist jetzt eins weniger */
		stackPointer--

		return Advance(commandPointer), nil
	},
	/* 16 - lt */
	func() (int, error) {
		/* am stack liegen 2 werte */
		a := stack[stackPointer-2]
		b := stack[stackPointer-1]
		if a < b {
			stack[stackPointer-2] = 1
		} else {
			stack[stackPointer-2] = 0
		}

		/* der stack ist jetzt eins weniger */
		stackPointer--

		return Advance(commandPointer), nil
	},
	/* 17 - neq */
	func() (int, error) {
		/* am stack liegen 2 werte */
		a := stack[stackPointer-2]
		b := stack[stackPointer-1]
		if a != b {
			stack[stackPointer-2] = 1
		} else {
			stack[stackPointer-2] = 0
		}

		/* der stack ist jetzt eins weniger */
		stackPointer--

		return Advance(commandPointer), nil
	},
	/* 18 - eq */
	func() (int, error) {
		/* am stack liegen 2 werte */
		a := stack[stackPointer-2]
		b := stack[stackPointer-1]
		if a == b {
			stack[stackPointer-2] = 1
		} else {
			stack[stackPointer-2] = 0
		}

		/* der stack ist jetzt eins weniger */
		stackPointer--

		return Advance(commandPointer), nil
	},
	/* 19 - jfalse */
	func() (int, error) {
		/* bei false soll gesprungen werden */
		// = commandpointer bearbeiten
		/* dabei wird gepopped */

		stackPointer--
		if stack[stackPointer] == 0 {
			// = false also jump
			return Jump(commandPointer, int(P1(commandPointer))), nil
		} else {
			return Advance(commandPointer), nil
		}
	},
	/* 20 - j */
	func() (int, error) {
		/* unconditional jump */
		return Jump(commandPointer, int(P1(commandPointer))), nil
	},
	/* 21 - alloc */
	func() (int, error) {
		/* alloc reserviert platz für ein int64 array am heap
		vor alloc liegt die größe des arrays am stack (in int64 blöcken)
		darunter stehen die einzelnen werte
		danach liegt die adresse des arrays am stack */

		stackPointer--
		length := stack[stackPointer]

		// Platz reservieren
		location := memory.heapPointer
		memory.heapPointer += (int(length) * 8) + 8 //das +8 ist weil bei einem array die größe mitgespeichert wird

		// ins erste feld die größe schreiben
		memory.WriteInt64To(int64(location), length)

		// dann von hinten beginnend das array auffüllen
		fillPosition := location + (int(length) * 8)       /* eigentlich loc + 8 + (len * 8) - 8 (+8 weil das length field übersprungen wird, -8 weil ja ins letzt feld geschrieben werden soll und nicht danach) */
		for ; fillPosition > location; fillPosition -= 8 { /* es wird geschrieben bis die fillposition gleich der location ist (d.h. am -1. feld dort wo die length steht) */
			/* wert vom stack holen */
			stackPointer--
			memory.WriteInt64To(int64(fillPosition), stack[stackPointer])
		}

		// jetzt ist array befüllt also wird nur noch die location auf den stack geschrieben
		stack[stackPointer] = int64(location)
		stackPointer++

		return Advance(commandPointer), nil
	},
	/* 22 - storei */
	func() (int, error) {
		/* storei ist ein indirect store:
		der wert wird nicht auf die stack-adresse geschrieben, die im store befehl angegeben ist,
		sondern es wird der wert in der stack-adresse gelesen (inkl. dereferencing und so)
		und dieser wert als heap-adresse aufgefasst in die geschrieben wird.
		storei erwartet sich einen offset auf dem stack (wird als 8byte felder gesehen)
		offset = 2 = schreibe 16 byte weiter vorne
		WICHTIG: Es gibt immer einen standardoffset von 8 byte (weil die ersten 8 byte das length field in arrays sind)*/

		offset := int64(8)
		stackPointer--
		offset += stack[stackPointer] * 8 /* store auf $a[0] ist somit 8 + 0*8 bytes von location weg (= nach length field) */

		// Die distance am stack herausfinden ----------------
		distance := P2(commandPointer)
		varFramePointer := framePointer
		/* den eigentlichen framepointer der für diese variable zuständig ist herausfinden */
		/* P2 mal dereferenzieren */
		for ; distance > 0; distance-- {
			// der framepointer des übergeordneten blocks liegt immer 1 unter dem aktuellen framepointer
			varFramePointer = int(stack[varFramePointer-1])
		}
		//-------------------------------------

		// die location des zu ändernden arrays holen
		location := stack[varFramePointer+int(P1(commandPointer))]

		// länge des arrays holen
		length := memory.ReadInt64At(location)

		// Prüfen ob wir das überhaupt schreiben dürfen
		// Fehler wenn: loc + off < loc + 8 (= erste beschreibbare pos)
		// oder wenn loc + off > (loc + 8 + length * 8 - 8) (= letzte beschreibbare position - 8 vor ende)
		if ((location + offset) < location+8) || (location+offset) > (location+(length*8)) {
			return 0, errors.New("Array index out of bounds (storei)! \nArray location: " + strconv.Itoa(int(location)) + "; Array min: " + strconv.Itoa(int(location)+8) + "; Array max: " + strconv.Itoa(int(location)+(int(length)*8)) + "; Read Position: " + strconv.Itoa(int(location)+int(offset)))
		}

		// stackpointer verringern damit er auf dem eigentlich zu schreibenden wert steht
		stackPointer--

		// und schreiben
		memory.WriteInt64To(location+offset, stack[stackPointer])

		return Advance(commandPointer), nil
	},
	/* 23 - retrievei */
	func() (int, error) {
		/* indirect retrieve:
		es wird nicht der speicherinhalt an der angegebenen stackposition geschrieben,
		sondern der wert in der angegeben stackposition als adresse zu einem array auf dem heap angesehen
		dort wird dann (+ offset auf oberster stackposition) ein int64 gelesen und am stack gelegt */

		offset := int64(8)
		stackPointer--
		offset += stack[stackPointer] * 8 /* retr auf $a[0] ist somit 8 + 0*8 bytes von location weg (= nach length field) */

		// Die distance am stack herausfinden ----------------
		distance := P2(commandPointer)
		varFramePointer := framePointer
		/* den eigentlichen framepointer der für diese variable zuständig ist herausfinden */
		/* P2 mal dereferenzieren */
		for ; distance > 0; distance-- {
			// der framepointer des übergeordneten blocks liegt immer 1 unter dem aktuellen framepointer
			varFramePointer = int(stack[varFramePointer-1])
		}
		//-------------------------------------

		// die location des zu lesenden arrays holen
		location := stack[varFramePointer+int(P1(commandPointer))]

		// länge des arrays holen
		length := memory.ReadInt64At(location)

		// Prüfen ob wir das überhaupt lesen dürfen
		// Fehler wenn: loc + off < loc + 8 (= erste beschreibbare pos)
		// oder wenn loc + off > (loc + 8 + length * 8 - 8) (= letzte beschreibbare position - 8 vor ende)
		if ((location + offset) < location+8) || (location+offset) > (location+(length*8)) {
			return 0, errors.New("Array index out of bounds (retrievei)! \nArray location: " + strconv.Itoa(int(location)) + "; Array min: " + strconv.Itoa(int(location)+8) + "; Array max: " + strconv.Itoa(int(location)+int(length)*8) + "; Read Position: " + strconv.Itoa(int(location)+int(offset)))
		}

		stack[stackPointer] = memory.ReadInt64At(location + offset)
		stackPointer++

		return Advance(commandPointer), nil
	},
	/* 24 - len */
	func() (int, error) {
		/* array length = nur auslesen des wertes im array-feld auf pos 0 */

		// Die distance am stack herausfinden ----------------
		distance := P2(commandPointer)
		varFramePointer := framePointer
		/* den eigentlichen framepointer der für diese variable zuständig ist herausfinden */
		/* P2 mal dereferenzieren */
		for ; distance > 0; distance-- {
			// der framepointer des übergeordneten blocks liegt immer 1 unter dem aktuellen framepointer
			varFramePointer = int(stack[varFramePointer-1])
		}
		//-------------------------------------

		// die location des zu lesenden arrays holen
		location := stack[varFramePointer+int(P1(commandPointer))]

		// länge des arrays holen
		length := memory.ReadInt64At(location)

		stack[stackPointer] = length
		stackPointer++

		return Advance(commandPointer), nil
	},
	/* 25 - disc */
	func() (int, error) {
		/* discard */
		stackPointer--

		return Advance(commandPointer), nil
	},
	/* 26 - call */
	func() (int, error) {
		/* function call */

		// Stackpointer zurücksetzen (um n+1)
		stackPointer = stackPointer - (int(P2(commandPointer))) - 1

		// den return pointer eins unter dem stackpointer speichern
		stack[stackPointer-1] = int64(commandPointer + commandSize) /* es wird der befehl nach dem aktuellen gespeichert */

		//Debug("Saving return-pointer: " + strconv.Itoa(commandPointer+commandSize) + " (" + strconv.Itoa((commandPointer+commandSize)/17) + ")")

		// zu befehl springen
		return Jump(commandPointer, int(P1(commandPointer))), nil
	},
	/* 27 - funcentry */
	func() (int, error) {
		/* function-blockentry */

		/* der framepointer wird auf den stack gelegt (da ein neuer angelegt wird) */
		stack[stackPointer] = int64(framePointer)
		stackPointer++
		/* der neue framepointer zeigt auf die zelle nach dem alten stackpointer */
		framePointer = stackPointer
		/* der neue stackpointer wird jetzt um die anzahl an lokalen variablen (P1) vorgeschoben */
		stackPointer = framePointer + int(P1(commandPointer))

		return Advance(commandPointer), nil
	},
	/* 28 - funcexit */
	func() (int, error) {
		/* function-exit */

		// davor: stack ...(rp)(ofp)[FP](a)(b)(c)(rv)...[SP]

		/*Befehl funcexit n (anz var, anz lokale und rv)
		z.B. funcexit 4
		stack ...(rv)[SP]

		macht
		rp = ([FP]-2)
		SP--
		([FP-2]) = $rv // return value
		SP = SP - 4
		... (rv)[SP](ofp)[FP]
		FP = ofp
		jump rp*/

		retPointer := stack[framePointer-2]

		//Debug("Got return pointer: " + strconv.Itoa(int(retPointer)) + " (" + strconv.Itoa(int(retPointer)/17) + ")")

		//println("retpointer for funcexit", strconv.Itoa(int(retPointer)))

		// jetzt return value schreiben
		// return value steht hinter den funktionsparametern:
		// ...[FP](a)(b)(c)(rv)...[SP]

		// der funcexit befehl kommt mit der anzahl der variablen als parameter + 1 wegen return value
		// d.h. der letzte (n-1) ist der return value
		stack[framePointer-2] = stack[framePointer+int(P2(commandPointer))]

		// jetzt wird der stackpointer auf die stack-position gesetzt wo der alte FP steht
		stackPointer = stackPointer - int(P1(commandPointer)) - 1

		// jetzt wird der alte framepointer wieder geholt
		framePointer = int(stack[stackPointer])

		// stackpointer steht jetzt eins über return value

		// Jetzt kann gejumped werden
		return JumpAbsolute(int(retPointer)), nil
	},
	/* 29 - blockreturn */
	func() (int, error) {
		/* mach mehrere blocks auf einmal zu */
		// anzahl der blocks liegt in P1

		anz := int(P1(commandPointer))

		for ; anz > 0; anz-- {
			stackPointer = framePointer - 1
			framePointer = int(stack[framePointer-1])
		}

		return Advance(commandPointer), nil
	},
}

/* schreibt einen int64 wert an angegebene memory pos */
func (buf Memory) WriteInt64To(location int64, value int64) {
	binary.PutVarint(buf.memory[int(location):int(location)+8], value)
}

/* holt einen int64 wert von position locatin */
func (buf Memory) ReadInt64At(location int64) int64 {
	return ConvertInt64(buf.memory[int(location) : int(location)+8])
}

func P2(commandPointer int) int64 {
	return ConvertInt64(commands[commandPointer+9 : commandPointer+17])
}

func P1(commandPointer int) int64 {
	return ConvertInt64(commands[commandPointer+1 : commandPointer+9])
}

/* helper */
func Jump(cPointer int, relativePos int) int {
	//println("Jumping to: " + strconv.Itoa((cPointer+(commandSize*relativePos))/17) + " (" + strconv.Itoa(cPointer+(commandSize*relativePos)) + ")")
	return cPointer + (commandSize * relativePos)
}

/* macht eigentlich nichts, ist nur dazu da damit es explizit geschrieben werden kann */
func JumpAbsolute(absolutePos int) int {
	return absolutePos
}

func Advance(cPointer int) int {
	return cPointer + commandSize
}

func Msg(msg string) {
	fmt.Fprintf(os.Stderr, msg+"\n")
}

func Run(file string, interpreter func() error) error {
	f, err := os.Open(file)
	if err != nil {
		return err
	}

	defer f.Close()

	//Buffer in dem int64 werte stehen
	intBuffer := make([]byte, 8)
	reader := bufio.NewReader(f)
	_, err = io.ReadFull(reader, intBuffer)
	if err != nil {
		Msg("could not read first 8 bytes of program")
		return err
	}

	// Länge des literalspeichers
	litsize := ConvertInt64(intBuffer)
	// Literalspeicher anlegen
	literals := memory.memory[0:litsize]

	// Literalspeicher befüllen
	_, err = io.ReadFull(reader, literals)
	if err != nil {
		Msg("error copying literals")
		return err
	}

	// Heappointer = nach literals
	memory.heapPointer = int(litsize)

	// Commands kopieren
	fileinfo, err := f.Stat()
	if err != nil {
		Msg("error optaining file info")
		return err
	}
	commandBlockSize := int(fileinfo.Size()) - 8 /* erstes byte */ - int(litsize)

	if commandBlockSize%17 != 0 {
		return errors.New("incomplete commands in command-block")
	}

	/* speicher für kommandos anlegen */
	commands = make([]byte, commandBlockSize)

	// und kopieren
	_, err = io.ReadFull(reader, commands)
	if err != nil {
		Msg("could not load all commands")
		return err
	}

	return interpreter()
}

func Interpret() error {
	var err error
	var newCP int /* commandpointer zwischenspeicher */
	var start, end int64
	var iinfo *InstructionInfo
	var ok bool
	for commandPointer < len(commands) {
		/* kommando aus dem implmentierungsarray holen, das dem kommandocode an der position commandpointer entspricht */
		if commands[commandPointer] >= byte(len(commandImpl)) {
			return errors.New("no command with code: " + strconv.Itoa(int(commands[commandPointer])) + " at position " + strconv.Itoa(commandPointer/17))
		}
		start = time.Now().UnixNano()
		newCP, err = commandImpl[commands[commandPointer]]()
		if err != nil {
			return err
		}
		end = time.Now().UnixNano()

		// Infos speichern
		iinfo, ok = Stats[commands[commandPointer]]
		if ok {
			// Hinzufügen
			iinfo.numCalled += 1
			iinfo.totalTime += (end - start)
		} else {
			Stats[commands[commandPointer]] = &InstructionInfo{(end - start), 1}
		}
		commandPointer = newCP
	}

	for cCode, iinfo := range Stats {
		cName, _ := getCommandName(cCode)
		fmt.Fprintf(os.Stderr, "%s: TotalTime: %d, MeanTime: %d, #Calls: %d\n", cName, iinfo.totalTime, int(iinfo.totalTime/iinfo.numCalled), iinfo.numCalled)
	}

	return nil
}

func PrintStep() {
	Debug("\033[2J")

	Dump()

	// Commands lesen
	var cmdBuffer []byte
	pos := 0
	for i := 0; i < len(commands); i += 17 {
		cmdBuffer = commands[i : i+17]

		// Ausgabe
		commandName, err := getCommandName(cmdBuffer[0])
		if err != nil {
			return
		}
		p1 := ConvertInt64(cmdBuffer[1:9])
		p2 := ConvertInt64(cmdBuffer[9:17])

		if i == commandPointer {
			Debug("* " + strconv.Itoa(pos) + ": " + commandName + "(" + strconv.Itoa(int(p1)) + ", " + strconv.Itoa(int(p2)) + ")")
		} else {
			Debug(strconv.Itoa(pos) + ": " + commandName + "(" + strconv.Itoa(int(p1)) + ", " + strconv.Itoa(int(p2)) + ")")
		}
		pos++
	}
}

func Step() error {
	var err error
	var newCP int /* commandpointer zwischenspeicher */

	broken := false
	if breakPoint == -1 {
		broken = true
	}

	for commandPointer < len(commands) {
		if commandPointer == breakPoint {
			broken = true
		}
		if broken {
			// Warten auf eingabe
			reader := bufio.NewReader(os.Stdin)
			fmt.Print("[Enter] for next step")
			reader.ReadString('\n')

			PrintStep()
		}

		/* kommando aus dem implmentierungsarray holen, das dem kommandocode an der position commandpointer entspricht */
		if commands[commandPointer] >= byte(len(commandImpl)) {
			return errors.New("no command with code: " + strconv.Itoa(int(commands[commandPointer])) + " at position " + strconv.Itoa(commandPointer/17))
		}
		newCP, err = commandImpl[commands[commandPointer]]()
		if err != nil {
			return err
		}
		commandPointer = newCP
	}

	return nil
}

func ConvertInt64(bytes []byte) int64 {
	val, _ := binary.Varint(bytes)
	return val
}

func Dump() {
	println()
	println()

	Debug("Framepointer: " + strconv.Itoa(framePointer))
	Debug("Stackpointer: " + strconv.Itoa(stackPointer))
	Debug("Commandpointer: " + strconv.Itoa(commandPointer/17) + " (" + strconv.Itoa(commandPointer) + ")")

	println()

	Debug("Memory:")
	Debug("Heappointer: " + strconv.Itoa(memory.heapPointer))
	Debug("Notnull fields in heap (might not be full heap): ")
	countz := 0
	for i := 0; i < len(memory.memory); i++ {
		if memory.memory[i] == 0 {
			countz++
		} else {
			countz = 0
		}
		if countz > 20 {
			break
		}
		if i%16 == 0 {
			if i >= 16 {
				fmt.Printf("  %s", string(memory.memory[i-16:i-1]))
			}
			fmt.Printf("\n%04d:", i)
		}
		fmt.Printf(" %02x", memory.memory[i])
	}
	println()
	println()
	Debug("Heap as int64: ")
	countz = 0
	for i := 0; i < len(memory.memory); i += 8 {
		val := ConvertInt64(memory.memory[i : i+8])
		if val == 0 {
			countz++
		} else {
			countz = 0
		}
		if countz > 8 {
			break
		}
		fmt.Printf("\n%04d: %016x", i, val)
	}
	println()
	println()

	Debug("Stack:")
	countz = 0
	for i := 0; i < len(stack); i++ {
		if stack[i] == 0 {
			countz++
		}
		if countz > 5 {
			break
		}
		fmt.Printf("\n%04d: %016x", i, stack[i])
		if i == stackPointer {
			fmt.Printf(" [SP]")
		}
		if i == framePointer {
			fmt.Printf(" [FP]")
		}
	}

	println()
}

func main() {
	args := os.Args[1:]

	var err error

	var interpreter func() error

	interpreter = Interpret

	for i := 0; i < len(args); i++ {
		if args[i] == "-a" {
			// Analyze
			err = Analyze(args[i+1])
			i++
		} else if args[i] == "-d" {
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
		} else {
			err = Run(args[i], interpreter)
		}
	}

	if err != nil {
		fmt.Fprintf(os.Stderr, err.Error()+"\n")
		fmt.Fprintf(os.Stderr, "Command-Position: %d (%d)\n", commandPointer/17, commandPointer)
		os.Exit(1)
	}
}
