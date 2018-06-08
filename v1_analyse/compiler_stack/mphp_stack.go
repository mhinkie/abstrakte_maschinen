/* mphp stack language */

package main

import "errors"
import "strconv"

const returnValueSymbolName = "RETVAL"

var statementProc map[string]func(statement *Statement, block *Block) error

func init() {
	statementProc = map[string]func(*Statement, *Block) error{
		"echo": func(statement *Statement, context *Block) error {
			/* hat 1 child */
			err := statement.children[0].Process(context)
			if err != nil {
				return err
			}

			switch statement.children[0].varType {
			case Integer:
				AddCommand(NewCommand("echoi", 0, 0))
			case String:
				AddCommand(NewCommand("echos", 0, 0))
			default:
				Assert("Cannot infer string or int type for echo statement", false)
			}

			// TODO
			return nil
		},
		"access": func(statement *Statement, context *Block) error {
			switch statement.children[0].command {
			case "var":
				/* schaun ob symbol da ist */
				/* struktur: access -> var -> "xyz" - deswegen das 1. kind vom 1. kind holen */
				symName := statement.children[0].children[0].command
				if !context.isSymbolDefined(symName) {
					return errors.New("Symbol " + symName + " is not defined!")
				}
				symType := context.getSymbolType(symName)

				// varType setzen damit übergeordneter befehl weiss was für ein typ kommt
				statement.varType = symType

				command := NewCommand("retrieve", context.getSymbolPosition(symName), context.getSymbolDistance(symName))
				command.resolvingHint = symName /* damit debug output besser lesbar ist */
				AddCommand(command)

				return nil
			case "array":
				var err error
				/* schaun ob symbol da ist */
				/* access -> array -> var -> name */
				symName := statement.children[0].children[0].children[0].command

				if !context.isSymbolDefined(symName) {
					return errors.New("Symbol " + symName + " is not defined!")
				}

				// Was ist das für ein symbol
				arrayType := context.getSymbolType(symName)

				switch arrayType {
				case ArrayString:
					statement.varType = String
				case ArrayInt:
					statement.varType = Integer
				default:
					return errors.New("Invalid array type " + arrayType.String())
				}

				/* vorgehensweise:
				es wird die array-index-expression ausgewertet
				= ergebnis steht am stack
				dann wird retreivei (= indirect retrieve) ausgeführt
				ergebnis von retrievei = wert am stack */

				/* index expr liegt auf: access -> array -> children[1] */
				err = statement.children[0].children[1].Process(context)
				if err != nil {
					return err
				}

				command := NewCommand("retrievei", context.getSymbolPosition(symName), context.getSymbolDistance(symName))
				command.resolvingHint = symName /* nur für output-lesbarkeit */
				AddCommand(command)

				return nil
			default:
				return errors.New("No way to access a " + statement.children[0].command + " variable")
			}

		},
		"assign": func(statement *Statement, context *Block) error {
			switch statement.children[0].command {
			case "var":
				/* rechtes kind = wert */
				err := statement.children[1].Process(context)
				if err != nil {
					return err
				}
				symType := statement.children[1].varType

				/* assign = neues symbol (likes kind) */

				symName := statement.children[0].children[0].command
				if context.isSymbolDefined(symName) && context.getSymbolType(symName) != symType {
					return errors.New("symbol " + symName + " already defined with different type")
				}
				context.addSymbol(symName, symType)

				// Adresse des literals ist jetzt am stack
				// diese wird jetzt in variable gespeichert
				cmd := NewCommand("store", context.getSymbolPosition(symName), context.getSymbolDistance(symName))
				cmd.resolvingHint = symName /* damit debug-output besser lesbar ist */
				AddCommand(cmd)

				return nil
			case "array":
				/* rechtes kind = wert */
				err := statement.children[1].Process(context)
				if err != nil {
					return err
				}
				symType := statement.children[1].varType

				/* assign -> array -> var -> symName */
				symName := statement.children[0].children[0].children[0].command
				var contentType PHPType

				// Type check
				// schaun, wenn das symbol schon vorhanden ist, ob es auch den richtigen type hat
				if context.isSymbolDefined(symName) {
					switch context.getSymbolType(symName) {
					case ArrayString:
						contentType = String
					case ArrayInt:
						contentType = Integer
					default:
						return errors.New("Invalid array-type: " + context.getSymbolType(symName).String())
					}

					if contentType != symType {
						return errors.New("Cannot assign this value to the array (type error) - value-type: " + symType.String() + " - array-type: " + contentType.String())
					}
				}

				// Symbol wird hier nicht hinzugefügt (weil ja nur ein wert des arrays assigned wird und kein komplettes neues)

				// jetzt noch dafür sorgen dass der index-wert am stack liegt
				err = statement.children[0].children[1].Process(context)
				if err != nil {
					return err
				}

				// Schreiben mit storei (indirect)
				cmd := NewCommand("storei", context.getSymbolPosition(symName), context.getSymbolDistance(symName))
				cmd.resolvingHint = symName /* für debug output */
				AddCommand(cmd)

				return nil
			default:
				return errors.New("No way to assign a " + statement.children[0].command + " variable")
			}

		},
		"litstring": func(statement *Statement, context *Block) error {
			/* processing für string literale */
			// = hinzufügen zu literal speicher
			// und position pushen
			// im literalspeicher steht dann an dieser position ein int64 für die länge (in byte) und danach der string

			// Literal anlegen
			pos := AddLiteral(statement.children[0].command)

			debugOutput("added literal " + statement.children[0].command + " at position " + strconv.Itoa(int(pos)))

			AddCommand(NewCommand("push", pos, 0))

			return nil
		},
		"litint": func(statement *Statement, context *Block) error {
			/* processing für integer literale = einfach nur pushen */
			/* im kind steht das eigentliche literal */

			i, err := strconv.Atoi(statement.children[0].command)
			if err != nil {
				return err
			}

			AddCommand(NewCommand("push", int64(i), 0))

			return nil
		},
		"plus": func(statement *Statement, context *Block) error {
			/* processing für plus = kinder processen und dann add befehl ausgeben */
			err := statement.children[0].Process(context)
			if err != nil {
				return err
			}
			err = statement.children[1].Process(context)
			if err != nil {
				return err
			}

			Assert("Operands of + operation have to be Integer", statement.children[0].varType == Integer && statement.children[1].varType == Integer)

			statement.varType = Integer

			AddCommand(NewCommand("add", 0, 0))

			return nil
		},
		"mult": func(statement *Statement, context *Block) error {
			/* processing für mul */
			err := statement.children[0].Process(context)
			if err != nil {
				return err
			}
			err = statement.children[1].Process(context)
			if err != nil {
				return err
			}

			Assert("Operands of * operation have to be Integer", statement.children[0].varType == Integer && statement.children[1].varType == Integer)

			statement.varType = Integer

			AddCommand(NewCommand("mult", 0, 0))

			return nil
		},
		"minus": func(statement *Statement, context *Block) error {
			/* processing für mul */
			err := statement.children[0].Process(context)
			if err != nil {
				return err
			}
			err = statement.children[1].Process(context)
			if err != nil {
				return err
			}

			Assert("Operands of - operation have to be Integer", statement.children[0].varType == Integer && statement.children[1].varType == Integer)

			statement.varType = Integer

			AddCommand(NewCommand("subt", 0, 0))

			return nil
		},
		"div": func(statement *Statement, context *Block) error {
			/* processing für mul */
			err := statement.children[0].Process(context)
			if err != nil {
				return err
			}
			err = statement.children[1].Process(context)
			if err != nil {
				return err
			}

			Assert("Operands of / operation have to be Integer", statement.children[0].varType == Integer && statement.children[1].varType == Integer)

			statement.varType = Integer

			AddCommand(NewCommand("div", 0, 0))

			return nil
		},
		"concat": func(statement *Statement, context *Block) error {
			err := statement.children[0].Process(context)
			if err != nil {
				return err
			}
			err = statement.children[1].Process(context)
			if err != nil {
				return err
			}

			Assert("Operants of string concatenation '.' have to be of type string", statement.children[0].varType == String && statement.children[1].varType == String)

			statement.varType = String

			AddCommand(NewCommand("concat", 0, 0))

			return nil
		},
		"gt": func(statement *Statement, context *Block) error {
			/* processing für mul */
			err := statement.children[0].Process(context)
			if err != nil {
				return err
			}
			err = statement.children[1].Process(context)
			if err != nil {
				return err
			}

			Assert("Operands of > operation have to be Integer", statement.children[0].varType == Integer && statement.children[1].varType == Integer)

			statement.varType = Integer

			AddCommand(NewCommand("gt", 0, 0))

			return nil
		},
		"lt": func(statement *Statement, context *Block) error {
			/* processing für mul */
			err := statement.children[0].Process(context)
			if err != nil {
				return err
			}
			err = statement.children[1].Process(context)
			if err != nil {
				return err
			}

			Assert("Operands of < operation have to be Integer", statement.children[0].varType == Integer && statement.children[1].varType == Integer)

			statement.varType = Integer

			AddCommand(NewCommand("lt", 0, 0))

			return nil
		},
		"neq": func(statement *Statement, context *Block) error {
			/* processing für mul */
			err := statement.children[0].Process(context)
			if err != nil {
				return err
			}
			err = statement.children[1].Process(context)
			if err != nil {
				return err
			}

			Assert("Operands of != operation have to be Integer", statement.children[0].varType == Integer && statement.children[1].varType == Integer)

			statement.varType = Integer

			AddCommand(NewCommand("neq", 0, 0))

			return nil
		},
		"eq": func(statement *Statement, context *Block) error {
			/* processing für mul */
			err := statement.children[0].Process(context)
			if err != nil {
				return err
			}
			err = statement.children[1].Process(context)
			if err != nil {
				return err
			}

			Assert("Operands of == operation have to be Integer", statement.children[0].varType == Integer && statement.children[1].varType == Integer)

			statement.varType = Integer

			AddCommand(NewCommand("eq", 0, 0))

			return nil
		},
		"if": func(statement *Statement, context *Block) error {
			/* if abarbeitung */
			/* als erstes wird die expr ausgewertet - es bleibt ein int im stack über */
			err := statement.children[0].Process(context)
			if err != nil {
				return err
			}

			Assert("Expression in if-statement must be of type Integer", statement.children[0].varType == Integer)

			/* neuen hint-namen erzeugen lassen */
			hint := NewResHint()
			elseHint := hint + "_ELSE"
			afterHint := hint + "_AFTER"

			jumpCommand := NewCommand("jfalse", 0, 0)
			jumpCommand.unresolvedJump = true
			if statement.childBlocks[1].isEmpty() {
				/* gibt es keinen else block wird gleich auf after gejumped */
				jumpCommand.resolvingHint = afterHint
			} else {
				/* bei zero wird auf den else-block gejumped (zero ist false) */
				jumpCommand.resolvingHint = elseHint

			}
			AddCommand(jumpCommand)

			/* dann wird der if-block ausgegeben (wenn true wird dieser ausgeführt) */
			err = statement.childBlocks[0].Process(context) /* parent = der aktuelle block */
			if err != nil {
				return err
			}

			if statement.childBlocks[1].isEmpty() {
				/* gibt es keinen elseblock wird nichts getan */
			} else {
				/* nach dem if-block kommt ein unconditional jump auf ein label nach dem if */
				jumpAfter := NewCommand("j", 0, 0)
				jumpAfter.unresolvedJump = true
				jumpAfter.resolvingHint = afterHint
				AddCommand(jumpAfter)

				/* dann wird das label gespeichert, das den else block markiert */
				SaveLabelPosition(elseHint)

				/* dann wird der else-block ausgegeben */
				err = statement.childBlocks[1].Process(context)
				if err != nil {
					return err
				}
			}

			/* danach wird das label dass das ende des ifs markiert ausgegeben */
			SaveLabelPosition(afterHint)

			return nil
		},
		"while": func(statement *Statement, context *Block) error {
			return ProcessWhile(statement, context)
		},
		"for": func(statement *Statement, context *Block) error {
			return ProcessFor(statement, context)
		},
		"aalloc": func(statement *Statement, context *Block) error {
			/* hat ein array von statements als kinder, die ausgewertet werden
			die statements haben alle einen typ, wenn dieser für alle String ist wird es ein string array,
			wenn dieser für alle int ist wird es ein int array, anonsten fehler */
			var contentType PHPType

			for i := 0; i < len(statement.children); i++ {
				err := statement.children[i].Process(context)
				if err != nil {
					return err
				}

				if i == 0 {
					contentType = statement.children[i].varType
				} else {
					// Prüfen ob type passt
					Assert("All array elements have to be of same type", contentType == statement.children[i].varType)
				}
			}

			/* die elemente dürfen nicht NoType sein */
			Assert("Type NoType is invalid for array elements", contentType != NoType)

			// Die ergebnisse der child-statements sind jetzt am stack
			/* alloc macht aus den letzten x werten am stack ein array und
			pushed einen zeiger darauf auf den stack (egal was das für elemente sind - kann stack ziemlich kaputt machen)
			stack davor e,e,e,l| z.B. für l=2 - stack danach: e|...*/

			AddCommand(NewCommand("push", int64(len(statement.children)), 0))

			AddCommand(NewCommand("alloc", 0, 0))

			// Typ ermitteln
			if contentType == String {
				statement.varType = ArrayString
			} else if contentType == Integer {
				statement.varType = ArrayInt
			}

			return nil
		},
		"foreach": func(statement *Statement, context *Block) error {
			return ProcessForeach(statement, context)
		},
		"proccall": func(statement *Statement, context *Block) error {
			/* ein procedure call ist nur ein function call (das child darunter), bei dem das ergebnis verworfen wird */
			err := statement.children[0].Process(context)
			if err != nil {
				return err
			}

			// Verwirft das oberste element am stack
			AddCommand(NewCommand("disc", 0, 0))

			return nil
		},
		"funccall": func(statement *Statement, context *Block) error {
			var err error
			/* um das ganze einfacher zu machen wird angenommen, dass funktionen immer vom type int sind */
			statement.varType = Integer

			/* als erstes werden 2 plätze am stack reserviert */
			AddCommand(NewCommand("push", 0, 0))
			AddCommand(NewCommand("push", 0, 0))

			/* dann werden die parameter geladen (vom ersten beginnend) */
			/* die parameter beginnen beim statement am index 1 (0 ist der fkt name) */
			for i := 1; i < len(statement.children); i++ {
				err = statement.children[i].Process(context)
				if err != nil {
					return err
				}
			}

			/* danach wird der call befehl ausgegeben - 1. param = funktionsadresse, 2. param = anzahl errechneter variablen */
			cmd := NewCommand("call", 0, int64(len(statement.children)-1))
			cmd.unresolvedFunction = true
			cmd.resolvingHint = statement.children[0].command
			AddCommand(cmd)

			return nil
		},
		"func": func(statement *Statement, context *Block) error {
			var entryCommand *Command
			/* funktionsdefinition - siehe doc/funktionen */

			funcName := statement.children[0].command
			funcBlock := &(statement.childBlocks[0])

			// dem block wird gleich das exit-label gegeben
			funcBlock.funcExitName = funcName + "_FUNC_EXIT"
			println("Setting funcexit label " + funcBlock.funcExitName)

			/* neues hilfslabel, das die funktion überspringt, wenn sie einfach so gelesen wird */
			hint := NewResHint()
			afterFunc := hint + "_" + funcName + "_FUNC_AFTER"

			// Jump command ausgeben, dass die funktion überspringt, wenn sie so erreicht wird
			cmd := NewCommand("j", 0, 0)
			cmd.unresolvedJump = true
			cmd.resolvingHint = afterFunc
			AddCommand(cmd)

			// function adresse als symbol speichern
			SaveFunctionPosition(funcName)

			// entry
			entryCommand = AddCommand(NewCommand("funcentry", 0, 0))

			// Alle parameter als lokale symbols hinzufügen (mit funcBlock als context)
			// dürfen noch nicht defined sein
			for i := 1; i < len(statement.children); i++ {
				symName := statement.children[i].children[0].command
				if funcBlock.isSymbolDefinedInCurrentScope(symName) {
					return errors.New("Function-argument " + symName + " is used multiple times")
				}

				funcBlock.addSymbol(symName, Integer)
			}

			// Jetzt den return value als symbol hinzufügen
			funcBlock.addSymbol(returnValueSymbolName, Integer)

			// Funktionsinhalt ausgeben
			for i := 0; i < len(funcBlock.statements); i++ {
				err := funcBlock.statements[i].Process(funcBlock)
				if err != nil {
					return err //bei fehlern wird abgebrochen
				}
			}

			// Return ausgeben (funcexit)
			SaveLabelPosition(funcBlock.funcExitName)
			AddCommand(NewCommand("funcexit", int64(len(funcBlock.symbols)), int64(funcBlock.getSymbolPosition(returnValueSymbolName))))

			/* jetzt noch die symbols vom entrycommand setzen */
			entryCommand.p1 = int64(len(funcBlock.symbols))

			// labe nach func speichern
			SaveLabelPosition(afterFunc)

			return nil
		},
		"return": func(statement *Statement, context *Block) error {
			/* return macht ein store auf RETVAL und gibt ein funcexit aus */
			if !context.isSymbolDefined(returnValueSymbolName) {
				return errors.New("No return symbol defined, return called outside function?")
			}

			err := statement.children[0].Process(context)
			if err != nil {
				return err
			}

			if statement.children[0].varType != Integer {
				debugTree("Function returns value of type " + statement.children[0].varType.String() + ". This might produce undefined behaviour.")
			}

			// returnwert liegt jetzt am stack und in RETVAL gespeichert
			// Adresse des literals ist jetzt am stack
			// diese wird jetzt in variable gespeichert
			cmd := NewCommand("store", context.getSymbolPosition(returnValueSymbolName), context.getSymbolDistance(returnValueSymbolName))
			cmd.resolvingHint = returnValueSymbolName /* damit debug-output besser lesbar ist */
			AddCommand(cmd)

			/* die distanz zum umschließenden funktionsblock muss gefunden werden */
			distance := 0
			current := context
			for ; !current.isFunctionBlock(); current = current.parent {
				distance++
			}

			/* stellt sicher, dass alle blocks die offen sind (anzahl = distance) geschlossen werden */
			AddCommand(NewCommand("blockreturn", int64(distance), 0))

			// Jetzt noch zum funktionsexit jumpen
			cmd = NewCommand("j", 0, 0)
			cmd.unresolvedJump = true
			cmd.resolvingHint = context.getFuncExitLabel()
			AddCommand(cmd)

			return nil
		},
	}
}

/* block processing für foreach blocks */
func ProcessForeach(foreachStatement *Statement, parent *Block) error {
	/* folgender output
	01: blockentry x + 2 (x sind alle anderen symbole + 1x internal symbol + 1x var-symbol)
	02: anlegen von zählervariable (push 0, store ...) - es wird ein symbol ohne $ vorne (NewInternalSymbol) angelegt
	03: label start
	04: retrieve internalSymbol
	05: len arr
	06: lt
	07: jfalse end
	08: retrieve internalSymbol
	09: retrievei arr
	10: store val
	11: blockinhalt
	12: push 1
	13: retrieve internalSymbol
	14: add
	15: store internalSymbol
	16: jump start
	17: label end
	18: blockexit x + 2 */

	var entryCommand *Command
	var intSym string /* interner zähler */
	var intSymType PHPType
	var arrSym string /* array */
	var valSym string /* zugriffsvariable */
	var valSymType PHPType

	foreachBlock := &(foreachStatement.childBlocks[0])

	/* hints für labels */
	hint := NewResHint()
	startHint := hint + "_START"
	afterHint := hint + "_AFTER"

	foreachBlock.parent = parent
	foreachBlock.debugHelper = "foreachBlock"

	/* die zählervariable als symbol anlegen */
	intSym = NewInternalSymbol()
	intSymType = Integer
	if foreachBlock.isSymbolDefined(intSym) && foreachBlock.getSymbolType(intSym) != intSymType {
		return errors.New("symbol " + intSym + " (foreach counter) already defined with different type")
	}
	foreachBlock.addSymbol(intSym, intSymType)

	/* Prüfen ob das array überhaupt da ist */
	arrSym = foreachStatement.children[0].children[0].command
	if !foreachBlock.isSymbolDefined(arrSym) {
		return errors.New("symbol " + arrSym + " (foreach array) is not defined!")
	}

	/* jetzt muss ich herausfinden, welchen typ das array und damit die zugriffsvariable hat */
	switch foreachBlock.getSymbolType(arrSym) {
	case ArrayInt:
		valSymType = Integer
	case ArrayString:
		valSymType = String
	default:
		return errors.New("Invalid Type for foreach array (only ArrayInt and ArrayString is allowed)")
	}

	/* jetzt wird die zugriffsvariable (val) als symbol definiert */
	valSym = foreachStatement.children[1].children[0].command
	if foreachBlock.isSymbolDefined(valSym) {
		return errors.New("Foreach variable " + valSym + " is already defined")
	}
	foreachBlock.addSymbol(valSym, valSymType)

	/* wird dann auf x + 2 locals gesetzt */
	/* 01 */
	entryCommand = AddCommand(NewCommand("blockentry", 0, 0))

	/* 02 */
	/* jetzt wird symbol für zählervariable anglegt
	Das symbol wird nicht im parent block angelegt, sondern im foreach block */

	// Jetzt wird die Zählervariable mit 0 intialisiert
	AddCommand(NewCommand("push", 0, 0))
	// diese wird jetzt in variable gespeichert
	cmd := NewCommand("store", foreachBlock.getSymbolPosition(intSym), foreachBlock.getSymbolDistance(intSym))
	cmd.resolvingHint = intSym /* damit debug-output besser lesbar ist */
	AddCommand(cmd)

	/* 03 */
	// label für schleifenstart */
	SaveLabelPosition(startHint)

	/* 04 - wert von internal symbol holen */
	cmd = NewCommand("retrieve", foreachBlock.getSymbolPosition(intSym), foreachBlock.getSymbolDistance(intSym))
	cmd.resolvingHint = intSym
	AddCommand(cmd)

	/* 05 - länge des arrays holen */
	cmd = NewCommand("len", foreachBlock.getSymbolPosition(arrSym), foreachBlock.getSymbolDistance(arrSym))
	cmd.resolvingHint = arrSym
	AddCommand(cmd)

	/* 06 - lt vergleich */
	AddCommand(NewCommand("lt", 0, 0))

	/* 07 - wenn nicht kleiner (jfalse) wird die schleife verlassen */
	jumpToEnd := NewCommand("jfalse", 0, 0)
	jumpToEnd.unresolvedJump = true
	jumpToEnd.resolvingHint = afterHint
	AddCommand(jumpToEnd)

	/* 08 - jetzt wird mithilfe vom internal-symbol der aktuelle arrayinhalt in val geschrieben */
	/* retrieve */
	cmd = NewCommand("retrieve", foreachBlock.getSymbolPosition(intSym), foreachBlock.getSymbolDistance(intSym))
	cmd.resolvingHint = intSym
	AddCommand(cmd)

	/* 09 - retrieve indirect auf array */
	cmd = NewCommand("retrievei", foreachBlock.getSymbolPosition(arrSym), foreachBlock.getSymbolDistance(arrSym))
	cmd.resolvingHint = arrSym
	AddCommand(cmd)

	/* 10 - ergebnis in val speichern */
	cmd = NewCommand("store", foreachBlock.getSymbolPosition(valSym), foreachBlock.getSymbolDistance(valSym))
	cmd.resolvingHint = valSym
	AddCommand(cmd)

	/* 11 - blockinhalt */
	/* block inhalt */
	for i := 0; i < len(foreachBlock.statements); i++ {
		err := foreachBlock.statements[i].Process(foreachBlock)
		if err != nil {
			return err //bei fehlern wird abgebrochen
		}
	}

	/* 12 - push 1 um zu erhöhen */
	AddCommand(NewCommand("push", 1, 0))

	/* 13 - retrieve internal sym */
	cmd = NewCommand("retrieve", foreachBlock.getSymbolPosition(intSym), foreachBlock.getSymbolDistance(intSym))
	cmd.resolvingHint = intSym
	AddCommand(cmd)

	/* 14 - add */
	AddCommand(NewCommand("add", 0, 0))

	/* 15 - intsym store */
	cmd = NewCommand("store", foreachBlock.getSymbolPosition(intSym), foreachBlock.getSymbolDistance(intSym))
	cmd.resolvingHint = intSym
	AddCommand(cmd)

	/* 16 - jump start */
	jumptToStart := NewCommand("j", 0, 0)
	jumptToStart.unresolvedJump = true
	jumptToStart.resolvingHint = startHint
	AddCommand(jumptToStart)

	/* 17 - label end */
	SaveLabelPosition(afterHint)

	/* 18 - blockexit */
	/* hier weiss ich schon wieviele symbols der block hat */
	AddCommand(NewCommand("blockexit", int64(len(foreachBlock.symbols)), 0))

	/* jetzt noch die symbols vom entrycommand setzen */
	entryCommand.p1 = int64(len(foreachBlock.symbols))

	return nil
}

/* block processing für for blocks */
func ProcessFor(forStatement *Statement, parent *Block) error {
	/* folgender output:
	blockentry
	auswertung declaration/assign statement children[0] - wenn vorhanden
	label startschleife
	auswertung expr children[1]
	jfalse nach_der_schleife (vor blockexit)
	blockinhalt
	auswertung assign statement children[2] - wenn vorhanden
	j startschleife
	label nach_der_schleife
	blockexit
	*/

	var entryCommand *Command
	forBlock := &(forStatement.childBlocks[0])

	/* hints für labels */
	hint := NewResHint()
	startHint := hint + "_START"
	afterHint := hint + "_AFTER"

	forBlock.parent = parent
	forBlock.debugHelper = "forBlock"

	// jetzt gehts los
	entryCommand = AddCommand(NewCommand("blockentry", 0, 0)) /* wird noch vervollständigt */

	// Auswertung declaration (hat parent forBlock)
	err := forStatement.children[0].Process(forBlock)
	if err != nil {
		return err
	}

	// Label startschleife
	SaveLabelPosition(startHint)

	// expr auswerten
	err = forStatement.children[1].Process(forBlock)
	if err != nil {
		return err
	}
	Assert("Expression in for-statement must be of type Integer", forStatement.children[1].varType == Integer)

	/* jump wenn false (nach die schleife) */
	jumpToEnd := NewCommand("jfalse", 0, 0)
	jumpToEnd.unresolvedJump = true
	jumpToEnd.resolvingHint = afterHint
	AddCommand(jumpToEnd)

	/* block inhalt */
	for i := 0; i < len(forBlock.statements); i++ {
		err := forBlock.statements[i].Process(forBlock)
		if err != nil {
			return err //bei fehlern wird abgebrochen
		}
	}

	/* assign statement (children[2]) */
	err = forStatement.children[2].Process(forBlock)
	if err != nil {
		return err
	}

	/* nach schleifeninhalt wird wieder nach oben gejumped */
	jumpToStart := NewCommand("j", 0, 0)
	jumpToStart.unresolvedJump = true
	jumpToStart.resolvingHint = startHint
	AddCommand(jumpToStart)

	/* label nach der schleife */
	SaveLabelPosition(afterHint)

	/* blockexit */
	/* hier weiss ich schon wieviele symbols der block hat */
	AddCommand(NewCommand("blockexit", int64(len(forBlock.symbols)), 0))

	/* jetzt noch die symbols vom entrycommand setzen */
	entryCommand.p1 = int64(len(forBlock.symbols))

	return nil
}

/* Block processing für while-blöcke */
func ProcessWhile(whileStatement *Statement, parent *Block) error {
	/* das while statement kommt mit einem kind statement (= expr, die bedingung)
	die bedingung wird in den block als erstes statement eingebaut
	vor dieses statement kommt ein label (= dort wird nach dem druchlauf hingesprungen)
	nach der auswertung des statements kann entweder nach die schleife gesprungen werden (= auch ein label)
	oder garnicht gesprungen werden (= ausführung der schleife)
	am ende des blocks ist wieder ein sprungstatement.
	also:
	blockentry
	label startschleife
	auswertung expr
	jfalse nach_der_schleife (vor blockexit)
	blockinhalt
	j startschleife
	label nach_der_schleife
	blockexit*/
	var entryCommand *Command
	whileBlock := &(whileStatement.childBlocks[0])

	/* hints für labels */
	hint := NewResHint()
	startHint := hint + "_START"
	afterHint := hint + "_AFTER"

	whileBlock.parent = parent
	whileBlock.debugHelper = "whileblock"

	/* wird später um die anzahl der lokalen var. ergänzt */
	entryCommand = AddCommand(NewCommand("blockentry", 0, 0))

	/* label für schleifenstart */
	SaveLabelPosition(startHint)

	/* expr auswerten - weil diese ja schon im block ist, wird whileBlock als parent mitgegeben und nicht der eigene parent */
	err := whileStatement.children[0].Process(whileBlock)
	if err != nil {
		return err
	}
	Assert("Expression in while-statement must be of type Integer", whileStatement.children[0].varType == Integer)

	/* jump wenn false (nach die schleife) */
	jumpToEnd := NewCommand("jfalse", 0, 0)
	jumpToEnd.unresolvedJump = true
	jumpToEnd.resolvingHint = afterHint
	AddCommand(jumpToEnd)

	/* block inhalt */
	for i := 0; i < len(whileBlock.statements); i++ {
		err := whileBlock.statements[i].Process(whileBlock)
		if err != nil {
			return err //bei fehlern wird abgebrochen
		}
	}

	/* nach schleifeninhalt wird wieder nach oben gejumped */
	jumpToStart := NewCommand("j", 0, 0)
	jumpToStart.unresolvedJump = true
	jumpToStart.resolvingHint = startHint
	AddCommand(jumpToStart)

	/* label nach der schleife */
	SaveLabelPosition(afterHint)

	/* blockexit */
	/* hier weiss ich schon wieviele symbols der block hat */
	AddCommand(NewCommand("blockexit", int64(len(whileBlock.symbols)), 0))

	/* jetzt noch die symbols vom entrycommand setzen */
	entryCommand.p1 = int64(len(whileBlock.symbols))

	return nil
}

/* Zu diesem Zeitpunkt ist parent noch nicht gesetzt */
func (block *Block) Process(parent *Block) error {

	var entryCommand *Command
	/* Plan: statements eins nach dem anderen durchgehen
	- wenn ein string- oder array-literal gefunden wird, in den speicher schreiben
	- wenn eine variablendefinition gefunden wird, in die symboltabelle schreiben
	- wenn ein variablenzugriff gefunden wird, prüfen ob in symboltabelle vorhanden
	*/
	block.parent = parent
	block.debugHelper = "regularBlock"
	/* block entry */
	if block.parent == nil {
		block.debugHelper = "programblock"
		// program block
		// das blockentry kommando wird später noch mit der anzahl der lokalen variablen ergänzt
		entryCommand = AddCommand(NewCommand("progentry", 0, 0))
	} else {
		// normaler block
		entryCommand = AddCommand(NewCommand("blockentry", 0, 0))
	}

	for i := 0; i < len(block.statements); i++ {
		err := block.statements[i].Process(block)
		if err != nil {
			return err //bei fehlern wird abgebrochen
		}
	}

	if block.parent == nil {
		AddCommand(NewCommand("progexit", int64(len(block.symbols)), 0))
	} else {
		AddCommand(NewCommand("blockexit", int64(len(block.symbols)), 0))
	}

	entryCommand.p1 = int64(len(block.symbols))

	return nil
}

func (statement *Statement) Process(parent *Block) error {
	//debugTree("processing Statement: " + statement.command)

	if statement.command == "" {
		// Das statement hat kein command (wird also nicht benötigt)
		return nil
	}

	procFunction := statementProc[statement.command]
	if procFunction != nil {
		return procFunction(statement, parent)
	} else {
		debugTree("No processor found for " + statement.command)
		return nil
	}
}
