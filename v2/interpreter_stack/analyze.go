package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os"
	"strconv"
)

/* bytecode */
type Command struct {
	commandCode uint8
	p1          int64
	p2          int64
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
	"return":      30,
}

func getCommandName(code uint8) (string, error) {
	for k, v := range commandCodes {
		if v == code {
			return k, nil
		}
	}

	return "", errors.New("command not found for code " + strconv.Itoa(int(code)))
}

func Debug(msg string) {
	fmt.Fprintf(os.Stderr, msg+"\n")

}

func Analyze(file string) error {
	// Erstes byte lesen
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
		Debug("could not read first 8 bytes of program")
		return err
	}

	// Länge des literalspeichers
	litsize := ConvertInt64(intBuffer)
	Debug(strconv.Itoa(int(litsize)) + " bytes of literals")

	reader.Discard(int(litsize))

	// Commands lesen
	cmdBuffer := make([]byte, 17)
	pos := 0
	for {
		_, err = io.ReadFull(reader, cmdBuffer)
		if err != nil {
			if err == io.EOF {
				break
			} else {
				return err
			}
		}

		// Ausgabe
		commandName, err := getCommandName(cmdBuffer[0])
		if err != nil {
			return err
		}
		p1 := ConvertInt64(cmdBuffer[1:9])
		p2 := ConvertInt64(cmdBuffer[9:17])

		Debug(strconv.Itoa(pos) + ": " + commandName + "(" + strconv.Itoa(int(p1)) + ", " + strconv.Itoa(int(p2)) + ")")
		pos++
	}

	return nil
}
