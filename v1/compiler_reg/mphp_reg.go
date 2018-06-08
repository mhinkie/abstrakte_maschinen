/* mphp register language */
/* jedem statement wird ein register zugewiesen (dort steht nach abarbeitung das ergebnis drin) */

package main

import (
	"errors"
	"strconv"
)

/*
Der erste context (symbolContext) gibt immer den wirklichen kontext (im source-code)
mit diesem sollen symbol-prüfungen und ähnliches durchgeführt werden
Der zweite context ist der in dem wirklich die symbole für die ausgabe definiert sind.
großer unterschied: ein if z.B. öffnet zwar einen neuen block (und context) im symbolContext
aber keinen neuen registerFrameContext - der wird nur bei funktionen aufgemacht */
var statementProc map[string]func(statement *Statement, context *Block) error

func init() {
	statementProc = map[string]func(*Statement, *Block) error{
		"assign": func(statement *Statement, context *Block) error {
			if statement.children[0].command == "var" {
				// Links steht eine variable, ist diese noch nicht angelegt wird sie neu angelegt (und bekommt dadurch ein register)
				// ist sie schon vorhanden, wird ein neues register geholt und die variable angelegt.
				// dem rechten kind wird dieses register gegeben und dann die befehle fürs rechte kind ausgegeben.

				// Als erstes das rechte kind ausgeben
				err := statement.children[1].Process(context)
				if err != nil {
					return err
				}

				symName := statement.children[0].children[0].command
				symType := statement.children[1].varType
				if !context.isSymbolDefined(symName) {
					// Neues register vergeben (vergebn wird im richtigen context)
					newreg := context.RegisterFrame().NewRegister()
					// Muss angelegt werden (in beiden contexts)
					context.addSymbol(symName, symType, newreg)
				}

				if context.getSymbolType(symName) != symType {
					return errors.New("Cannot assign value of type " + symType.String() + " to variable of type " + context.getSymbolType(symName).String())
				}

				// Für besseren debug-output wird es als register für die variable gesetzt
				statement.children[0].register = context.getSymbolRegister(symName)

				/* Global ist eine variable die in einem block definiert ist der einen registerframe hat der vom eigenen unterschiedlich ist */
				if context.DefiningBlock(symName).RegisterFrame() != context.RegisterFrame() {
					return errors.New("Handling of globals not implemented!")
				}

				Assert("Symbol "+symName+" has wrong type!", symType == context.getSymbolType(symName))

				// Dann das ergebnis in das register des symbols speichern
				// Das ganze ist eigentlich ein unnötiger move-befehl (der dann später auch wegoptimiert werden sollte)
				// Wenn beide kinder variablen sind, ist der store aber undebindingt notwendig (hier wird auch nur von register nach reigster kopiert, das aber absichtlich)
				if statement.children[0].command == "var" && statement.children[1].command == "access" {
					AppendCommand(NewDoubleValueCommand(OpStore, Register, context.getSymbolRegister(symName), Register, statement.children[1].register, ""))
				} else if statement.children[0].command == "var" && statement.children[1].command == "aalloc" {
					// Bei arrayzugriff kann der zugriff so auch nicht wegoptimiert werden
					AppendCommand(NewDoubleValueCommand(OpStore, Register, context.getSymbolRegister(symName), Register, statement.children[1].register, ""))
				} else {
					AppendCommand(NewDoubleValueCommand(OpStoreUnoptimized, Register, context.getSymbolRegister(symName), Register, statement.children[1].register, ""))
				}

				return nil
			} else if statement.children[0].command == "array" {
				// Array assign ist ein hstore mit dem register in dem der arraypointer steht als output
				// op1 = der wert der assigned wird (register)
				// op2 = offset
				arrStatement := statement.children[0]

				err := statement.children[1].Process(context)
				if err != nil {
					return err
				}
				valueType := statement.children[1].varType

				arrSymbol := arrStatement.children[0].children[0].command

				/* Index statement auswerten */
				err = arrStatement.children[1].Process(context)
				if err != nil {
					return err
				}

				// Array muss schon defined sein
				if !context.isSymbolDefined(arrSymbol) {
					return errors.New("implicit declaration of arrays is not supported!")
				}

				// Typ des arrays muss passen
				if !context.getSymbolType(arrSymbol).IsArray() {
					return errors.New("cannot use array assign on variable of type " + context.getSymbolType(arrSymbol).String())
				}

				// Schaun ob arraytype richtig ist
				if context.getSymbolType(arrSymbol).BaseType() != valueType {
					return errors.New("cannot assing value of type " + valueType.String() + " to array of type " + context.getSymbolType(arrSymbol).BaseType().String())
				}

				// Schaun ob global ist
				/* Global ist eine variable die in einem block definiert ist der einen registerframe hat der vom eigenen unterschiedlich ist */
				if context.DefiningBlock(arrSymbol).RegisterFrame() != context.RegisterFrame() {
					return errors.New("Handling of globals not implemented!")
				}

				// Prüfen ob der index passt
				if arrStatement.children[1].varType != Integer {
					return errors.New("cannot use value of type " + arrStatement.children[1].varType.String() + " as array-index")
				}

				arrRegister := context.getSymbolRegister(arrSymbol)

				// Und jetzt kann endlich zugewiesen werden
				AppendCommand(NewTripleValueCommand(OpHStore, Register, arrRegister, Register, statement.children[1].register, "", Register, arrStatement.children[1].register, ""))

				return nil
			} else {
				return errors.New("Invalid statement type for access: " + statement.children[0].command)
			}
		},
		"access": func(statement *Statement, context *Block) error {

			if statement.children[0].command == "var" {
				// Variablenzugriff
				// Variable muss vorhanden sein (statement type = var-type)
				symName := statement.children[0].children[0].command

				if !context.isSymbolDefined(symName) {
					return errors.New("No symbol of name " + symName + " defined!")
				}

				if context.DefiningBlock(symName).RegisterFrame() != context.RegisterFrame() {
					return errors.New("Handling of globals not implemented!")
				}

				// Hier muss nur das register gesetzt werden (es wird nichts ausgegeben)
				statement.register = context.getSymbolRegister(symName)
				statement.varType = context.getSymbolType(symName)

				return nil
			} else if statement.children[0].command == "array" {
				// Arrayzugriff
				arraySymbol := statement.children[0].children[0].children[0].command

				// Prüfen ob es das symbol überhaupt gibt
				if !context.isSymbolDefined(arraySymbol) {
					return errors.New("No symbol of name " + arraySymbol + " defined!")
				}

				arrayType := context.getSymbolType(arraySymbol)

				// Prüfen ob das auch wirklich ein array ist
				if !arrayType.IsArray() {
					return errors.New("Trying to access array on base-type-variable (" + arraySymbol + " - type: " + arrayType.String() + ")")
				}

				if context.DefiningBlock(arraySymbol).RegisterFrame() != context.RegisterFrame() {
					return errors.New("Handling of globals not implemented!")
				}

				// Beim access auf arrays wird ein register benötigt
				if statement.register < 0 {
					statement.register = context.RegisterFrame().NewRegister()
				}

				// holt den base-type T vom array-type T[]
				statement.varType = arrayType.BaseType()

				// Jetzt wird das register der array-variable geholt
				arrayRegister := context.getSymbolRegister(arraySymbol)

				// Jetzt wird das rechte kind ausgegeben (= index)
				err := statement.children[0].children[1].Process(context)
				if err != nil {
					return err
				}

				// Das rechte kind muss hier unbedingt ein int sein
				if statement.children[0].children[1].varType != Integer {
					return errors.New("Cannot use value of type " + statement.children[0].children[1].varType.String() + " as array index")
				}

				// Und jetzt ein hread ausgeben
				AppendCommand(NewTripleValueCommand(OpHRead, Register, statement.register, Register, arrayRegister, "", Register, statement.children[0].children[1].register, ""))

				return nil
			} else {
				return errors.New("Invalid statement type for access: " + statement.children[0].command)
			}
		},
		"litstring": func(statement *Statement, context *Block) error {
			// ist schon ein register gesetzt?
			if statement.register < 0 {
				// Für das literal wird ein neues register vergeben:
				statement.register = context.RegisterFrame().NewRegister()
			}

			// Und dort hineine wird es mit store gespeichert
			AppendCommand(NewDoubleValueCommand(OpStore, Register, statement.register, LiteralString, 0, statement.children[0].command))

			statement.varType = String

			return nil
		},
		"litint": func(statement *Statement, context *Block) error {
			// int literal einfach store in register
			if statement.register < 0 {
				statement.register = context.RegisterFrame().NewRegister()
			}

			intVal, err := strconv.Atoi(statement.children[0].command)
			if err != nil {
				return err
			}
			AppendCommand(NewDoubleValueCommand(OpStore, Register, statement.register, LiteralInt, intVal, ""))

			statement.varType = Integer

			return nil
		},
		"echo": func(statement *Statement, context *Block) error {
			// Anhand des typen des kinds herausfinden, was ausgegeben wird
			if statement.children[0].command == "litint" {
				intVal, err := strconv.Atoi(statement.children[0].children[0].command)
				if err != nil {
					return err
				}
				// Wenn das kind ein literal ist braucht man kein listring ausgeben (und damit unnötig platz am heap verbrauchen)
				AppendCommand(NewInputOnlyCommand(OpEchoS, LiteralString, intVal, ""))
			} else if statement.children[0].command == "litstring" {
				// Wenn das kind ein literal ist braucht man kein listring ausgeben (und damit unnötig platz am heap verbrauchen)
				AppendCommand(NewInputOnlyCommand(OpEchoS, LiteralString, 0, statement.children[0].children[0].command))
			} else {
				// Kind processen
				err := statement.children[0].Process(context)
				if err != nil {
					return err
				}

				switch statement.children[0].varType {
				case Integer:
					// echoi
					AppendCommand(NewInputOnlyCommand(OpEchoI, Register, statement.children[0].register, ""))
				case String:
					// echos
					AppendCommand(NewInputOnlyCommand(OpEchoS, Register, statement.children[0].register, ""))
				default:
					return errors.New("No echo-command available for type " + statement.children[0].varType.String())
				}
			}

			return nil
		},
		"plus": func(statement *Statement, context *Block) error {
			// = linkes kind - rechtes kind
			err := statement.children[0].Process(context)
			if err != nil {
				return err
			}

			err = statement.children[1].Process(context)
			if err != nil {
				return err
			}

			// Addieren und in ein temp reg speichern
			Assert("Can only add integers!", statement.children[0].varType == Integer && statement.children[1].varType == Integer)

			// int literal einfach store in register
			if statement.register < 0 {
				statement.register = context.RegisterFrame().NewRegister()
			}

			AppendCommand(NewTripleValueCommand(OpAdd, Register, statement.register, Register, statement.children[0].register, "", Register, statement.children[1].register, ""))

			statement.varType = Integer

			return nil
		},
		"minus": func(statement *Statement, context *Block) error {
			// = linkes kind - rechtes kind
			err := statement.children[0].Process(context)
			if err != nil {
				return err
			}

			err = statement.children[1].Process(context)
			if err != nil {
				return err
			}

			// Addieren und in ein temp reg speichern
			Assert("Can only subt integers!", statement.children[0].varType == Integer && statement.children[1].varType == Integer)

			// int literal einfach store in register
			if statement.register < 0 {
				statement.register = context.RegisterFrame().NewRegister()
			}

			AppendCommand(NewTripleValueCommand(OpSubt, Register, statement.register, Register, statement.children[0].register, "", Register, statement.children[1].register, ""))

			statement.varType = Integer

			return nil
		},
		"mult": func(statement *Statement, context *Block) error {
			// = linkes kind - rechtes kind
			err := statement.children[0].Process(context)
			if err != nil {
				return err
			}

			err = statement.children[1].Process(context)
			if err != nil {
				return err
			}

			// Addieren und in ein temp reg speichern
			Assert("Can only mult integers!", statement.children[0].varType == Integer && statement.children[1].varType == Integer)

			// int literal einfach store in register
			if statement.register < 0 {
				statement.register = context.RegisterFrame().NewRegister()
			}

			AppendCommand(NewTripleValueCommand(OpMult, Register, statement.register, Register, statement.children[0].register, "", Register, statement.children[1].register, ""))

			statement.varType = Integer

			return nil
		},
		"div": func(statement *Statement, context *Block) error {
			// = linkes kind - rechtes kind
			err := statement.children[0].Process(context)
			if err != nil {
				return err
			}

			err = statement.children[1].Process(context)
			if err != nil {
				return err
			}

			// Addieren und in ein temp reg speichern
			Assert("Can only div integers!", statement.children[0].varType == Integer && statement.children[1].varType == Integer)

			// int literal einfach store in register
			if statement.register < 0 {
				statement.register = context.RegisterFrame().NewRegister()
			}

			AppendCommand(NewTripleValueCommand(OpDiv, Register, statement.register, Register, statement.children[0].register, "", Register, statement.children[1].register, ""))

			statement.varType = Integer

			return nil
		},
		"concat": func(statement *Statement, context *Block) error {
			// erwarte mir 2 strings
			err := statement.children[0].Process(context)
			if err != nil {
				return err
			}

			err = statement.children[1].Process(context)
			if err != nil {
				return err
			}

			Assert("Can only concat strings", statement.children[0].varType == String && statement.children[1].varType == String)

			statement.varType = String

			// braucht für das ergebnis ein register
			if statement.register < 0 {
				statement.register = context.RegisterFrame().NewRegister()
			}

			AppendCommand(NewTripleValueCommand(OpConcat, Register, statement.register, Register, statement.children[0].register, "", Register, statement.children[1].register, ""))

			return nil
		},
		"gt": func(statement *Statement, context *Block) error {
			err := statement.children[0].Process(context)
			if err != nil {
				return err
			}

			err = statement.children[1].Process(context)
			if err != nil {
				return err
			}

			// Addieren und in ein temp reg speichern
			Assert("Can only compare integer!", statement.children[0].varType == Integer && statement.children[1].varType == Integer)

			if statement.register < 0 {
				statement.register = context.RegisterFrame().NewRegister()
			}

			AppendCommand(NewTripleValueCommand(OpGt, Register, statement.register,
				Register, statement.children[0].register, "", Register,
				statement.children[1].register, ""))

			statement.varType = Integer

			return nil
		},
		"lt": func(statement *Statement, context *Block) error {
			err := statement.children[0].Process(context)
			if err != nil {
				return err
			}

			err = statement.children[1].Process(context)
			if err != nil {
				return err
			}

			// Addieren und in ein temp reg speichern
			Assert("Can only compare integer!", statement.children[0].varType == Integer && statement.children[1].varType == Integer)

			if statement.register < 0 {
				statement.register = context.RegisterFrame().NewRegister()
			}

			AppendCommand(NewTripleValueCommand(OpLt, Register, statement.register,
				Register, statement.children[0].register, "", Register,
				statement.children[1].register, ""))

			statement.varType = Integer

			return nil
		},
		"eq": func(statement *Statement, context *Block) error {
			err := statement.children[0].Process(context)
			if err != nil {
				return err
			}

			err = statement.children[1].Process(context)
			if err != nil {
				return err
			}

			// Addieren und in ein temp reg speichern
			Assert("Can only compare integer!", statement.children[0].varType == Integer && statement.children[1].varType == Integer)

			if statement.register < 0 {
				statement.register = context.RegisterFrame().NewRegister()
			}

			AppendCommand(NewTripleValueCommand(OpEq, Register, statement.register,
				Register, statement.children[0].register, "", Register,
				statement.children[1].register, ""))

			statement.varType = Integer

			return nil
		},
		"neq": func(statement *Statement, context *Block) error {
			err := statement.children[0].Process(context)
			if err != nil {
				return err
			}

			err = statement.children[1].Process(context)
			if err != nil {
				return err
			}

			// Addieren und in ein temp reg speichern
			Assert("Can only compare integer!", statement.children[0].varType == Integer && statement.children[1].varType == Integer)

			if statement.register < 0 {
				statement.register = context.RegisterFrame().NewRegister()
			}

			AppendCommand(NewTripleValueCommand(OpNeq, Register, statement.register,
				Register, statement.children[0].register, "", Register,
				statement.children[1].register, ""))

			statement.varType = Integer

			return nil
		},
		"if": func(statement *Statement, context *Block) error {
			// If-statment: children[0]==prüfung also als erstes auszugeben
			// dann kommt der if-true-block (der ist immer auszugeben)
			// dann der else block (kann auch nicht vorhanden sein)
			// beide gibts als childblocks

			hasElseBlock := len(statement.childBlocks[1].statements) > 0
			// Labels erzeugen für dieses if
			label := NewLabel()
			labelElse := label + "_ELSE"
			labelEnd := label + "_END"

			// Kind ausgeben
			err := statement.children[0].Process(context)
			if err != nil {
				return err
			}

			Assert("Can only make conditional jump if statement has type int", statement.children[0].varType == Integer)

			// Jetzt wird ein jfalse ausgegeben (springt wenn übergebenes register 0 ist)
			// Es wird in den else zweig gesprungen
			jumpCommand := NewInputOnlyCommand(OpJFalse, Register, statement.children[0].register, "")
			jumpCommand.UnresolvedJump = true
			jumpCommand.ResolvingHint = labelElse
			AppendCommand(jumpCommand)

			ifBlock := &statement.childBlocks[0]
			ifBlock.parent = context

			// Jetzt wird der true-block ausgegeben
			for _, stmt := range statement.childBlocks[0].statements {
				err = stmt.Process(ifBlock)
				if err != nil {
					return err
				}
			}

			// Jetzt wird ein jump ans ende des ifs ausgeführt
			// Der wird nur gebraucht, wenn es auch einen else-block gibt
			if hasElseBlock {
				jumpEnd := NewEmptyCommand(OpJump)
				jumpEnd.UnresolvedJump = true
				jumpEnd.ResolvingHint = labelEnd
				AppendCommand(jumpEnd)
			}

			// Jetzt wird das Label das den Else block markiert ausgegeben
			AppendCommand(NewInputOnlyCommand(Label, LiteralString, 0, labelElse))

			if hasElseBlock {
				elseBlock := &statement.childBlocks[1]
				elseBlock.parent = context

				// Und jetzt wird der else block ausgegeben
				for _, stmt := range statement.childBlocks[1].statements {
					err = stmt.Process(elseBlock)
					if err != nil {
						return err
					}
				}

				// Und dann noch ein label das das ende markiert
				AppendCommand(NewInputOnlyCommand(Label, LiteralString, 0, labelEnd))
			}

			return nil
		},
		"while": func(statement *Statement, context *Block) error {
			/* Ausgabe zu while:
			label vor schleife
			statement ausgegeben
			jfalse nach schleife
			block ausgeben
			j vor schleife
			label nach schleife */

			// Labels erzeugen für schleife
			label := NewLabel()
			labelBefore := label + "_BEFORE"
			labelAfter := label + "_AFTER"

			// Jetzt das Label vor der schleife ausgegeben
			AppendCommand(NewInputOnlyCommand(Label, LiteralString, 0, labelBefore))

			// Dann statement ausgegeben
			err := statement.children[0].Process(context)
			if err != nil {
				return err
			}

			// Conditional Jump nach die schleife
			jumpCommand := NewInputOnlyCommand(OpJFalse, Register, statement.children[0].register, "")
			jumpCommand.UnresolvedJump = true
			jumpCommand.ResolvingHint = labelAfter
			AppendCommand(jumpCommand)

			// Danach wird der schleifeninhalt ausgegeben
			whileBlock := &statement.childBlocks[0]
			whileBlock.parent = context

			for _, stmt := range statement.childBlocks[0].statements {
				err = stmt.Process(whileBlock)
				if err != nil {
					return err
				}
			}

			// Dann ein jump vor die schleife
			jumpStart := NewEmptyCommand(OpJump)
			jumpStart.UnresolvedJump = true
			jumpStart.ResolvingHint = labelBefore
			AppendCommand(jumpStart)

			// Dann ein label nach der schleife
			AppendCommand(NewInputOnlyCommand(Label, LiteralString, 0, labelAfter))

			return nil
		},
		"for": func(statement *Statement, context *Block) error {
			/* Ausgabe für für for loop:
			Das zuweisungsstatement wird im context, des loops ausgegeben
			label anfang
			dann das prüfungsstatement ausgegeben
			jfalse schleifenende
			schleifeninhalt
			dann das manipulationsstatement (das dritte) im context des loops
			j anfang
			label schleifenende
			*/

			var err error

			// Labels erzeugen für schleife
			label := NewLabel()
			labelBefore := label + "_BEFORE"
			labelAfter := label + "_AFTER"

			// Schleifenkontext holen
			forBlock := &statement.childBlocks[0]
			forBlock.parent = context

			// zuweisungsstatement (nur wenn nicht leer)
			err = statement.children[0].Process(forBlock)
			if err != nil {
				return err
			}

			// Label anfang
			AppendCommand(NewInputOnlyCommand(Label, LiteralString, 0, labelBefore))

			// Prüfungsstatement (nur wenn nicht leer)
			err = statement.children[1].Process(forBlock)
			if err != nil {
				return err
			}

			// JFalse ans schleifenende
			jumpEnd := NewInputOnlyCommand(OpJFalse, Register, statement.children[1].register, "")
			jumpEnd.UnresolvedJump = true
			jumpEnd.ResolvingHint = labelAfter
			AppendCommand(jumpEnd)

			// Dann schleifeninhalt ausgeben
			for _, stmt := range statement.childBlocks[0].statements {
				err = stmt.Process(forBlock)
				if err != nil {
					return err
				}
			}

			// Dann zuweisungsstatement (ende) ausgeben
			err = statement.children[2].Process(forBlock)
			if err != nil {
				return err
			}

			// Dann ein Jump an den Anfang
			jumpStart := NewEmptyCommand(OpJump)
			jumpStart.UnresolvedJump = true
			jumpStart.ResolvingHint = labelBefore
			AppendCommand(jumpStart)

			// Und am ende das label fürs schleifenende ausgeben
			AppendCommand(NewInputOnlyCommand(Label, LiteralString, 0, labelAfter))

			return nil
		},
		"aalloc": func(statement *Statement, context *Block) error {
			// Gibt mir ein array
			// Prüfungen: alle exprs im aalloc müssen den selben typ haben

			var err error

			len := len(statement.children)
			bsize := 8 /* nur integer kommen am heap */

			if statement.register < 0 {
				statement.register = context.RegisterFrame().NewRegister()
			}

			// Array platz auf heap reservieren
			AppendCommand(NewTripleValueCommand(OpAalloc, Register, statement.register, LiteralInt, len, "", LiteralInt, bsize, ""))

			var itemType PHPType
			itemType = -100 // Unset

			// Jetzt die kind-expressions ausgeben und storen
			for i, stmt := range statement.children {
				err = stmt.Process(context)
				if err != nil {
					return err
				}

				if itemType == -100 {
					itemType = stmt.varType
				} else {
					// Itemtype ist schon gesetzt es muss geprüft werden ob der auch passt
					if itemType != stmt.varType {
						return errors.New("Array containing values of different types: " + itemType.String() + " and " + stmt.varType.String())
					}
				}

				// hstore ausgeben (speichert wert dann auf heap)
				// Output register wird gelesen als store auf heap in array-adressen von (output)
				AppendCommand(NewTripleValueCommand(OpHStore, Register, statement.register, Register, stmt.register, "", LiteralInt, i, ""))
			}

			// Jetzt kann typ des arrays bestimmt werden
			statement.varType = ArrayOf(itemType)

			return nil
		},
		"foreach": func(statement *Statement, context *Block) error {
			/* Ausgabe foreach:
			Anlagen eines neuen symbols (name hint_iterator)
			Anlegen des wert-symbols (wenns dieses noch nicht gibt)
			Anlegen eines len-symbols
			zuweisen von wert 0
			zuweisen von len(array) ans rlenarray
			label start
			lt rhint_iterator rlenarray
			JFalse ende
			wert-symbol den wert arr[hint_iterator] zuweisen
			schleifeninhalt
			hint_iterator erhöhen
			Jump start
			label ende */

			var err error

			foreachBlock := &statement.childBlocks[0]
			foreachBlock.parent = context

			// Iterator-variable
			hint := NewLabel() // Hier werden einfach label-werte verwendet (die sind fix nicht vergeben weil sie nicht mit $ anfangen)
			itSymbol := hint + "_FOREACH_I"
			lenSymbol := hint + "_FOREACH_LEN"
			labelStart := hint + "_BEFORE"
			labelEnd := hint + "_AFTER"

			// Prüfungen
			valueSymbol := statement.children[1].children[0].command

			// Arraytyp
			arraySymbol := statement.children[0].children[0].command
			if !context.isSymbolDefined(arraySymbol) {
				return errors.New("Variable " + arraySymbol + " is not defined!")
			}
			if !context.getSymbolType(arraySymbol).IsArray() {
				return errors.New("Variable " + arraySymbol + " is not an array: cannot use foreach")
			}
			arrayBaseType := context.getSymbolType(arraySymbol).BaseType()
			arrayRegister := context.getSymbolRegister(arraySymbol)

			// Das value-symbol
			if context.isSymbolDefined(valueSymbol) {
				// Wenn schon defined muss typ mit array base type übereinstimmen
				if context.getSymbolType(valueSymbol) != arrayBaseType {
					return errors.New("Cannot iterator " + arraySymbol + " (type: " + arrayBaseType.String() + ") using variable of type " + context.getSymbolType(valueSymbol).String())
				}
			} else {
				// Variable wird jetzt defined (im foreach-scope)
				newreg := foreachBlock.RegisterFrame().NewRegister()
				foreachBlock.addSymbol(valueSymbol, arrayBaseType, newreg)
			}
			valueRegister := foreachBlock.getSymbolRegister(valueSymbol)

			// Jetzt den iterator für die schleife anlegen
			itRegister := foreachBlock.RegisterFrame().NewRegister()
			foreachBlock.addSymbol(itSymbol, Integer, itRegister)

			// Jetzt die len-variable anlegen
			lenRegister := foreachBlock.RegisterFrame().NewRegister()
			foreachBlock.addSymbol(lenSymbol, Integer, lenRegister)

			// Dem iterator-symbol wird jetzt 0 zugewiesen
			// Das kann eigentlich weggelassenen werden, weil go alle array-werte (damit auch die register) auf 0 setzt
			// Is aber besser wenn es explizit passiert (außerdem ist ja nur ein befehl...)
			AppendCommand(NewDoubleValueCommand(OpStore, Register, itRegister, LiteralInt, 0, ""))

			// Jetzt wird die array length in lenRegister gelesen
			AppendCommand(NewDoubleValueCommand(OpLen, Register, lenRegister, Register, arrayRegister, ""))

			// Label start
			AppendCommand(NewInputOnlyCommand(Label, LiteralString, 0, labelStart))

			// jetzt wird ein vergleich ausgebeben (it < len?)...
			compRegister := foreachBlock.RegisterFrame().NewRegister()
			AppendCommand(NewTripleValueCommand(OpLt, Register, compRegister, Register, itRegister, "", Register, lenRegister, ""))

			// ...und gejumped wenn false
			jumpEnd := NewInputOnlyCommand(OpJFalse, Register, compRegister, "")
			jumpEnd.UnresolvedJump = true
			jumpEnd.ResolvingHint = labelEnd
			AppendCommand(jumpEnd)

			// Dem value-symbol wird jetzt der aktuelle wert zugewiesen
			AppendCommand(NewTripleValueCommand(OpHRead, Register, valueRegister, Register, arrayRegister, "", Register, itRegister, ""))

			// Jetzt kann endlich der schleifeninhalt ausgegeben werden
			for _, stmt := range statement.childBlocks[0].statements {
				err = stmt.Process(foreachBlock)
				if err != nil {
					return err
				}
			}

			// nach dem inhalt wird der hint-iterator erhöht
			AppendCommand(NewTripleValueCommand(OpAdd, Register, itRegister, Register, itRegister, "", LiteralInt, 1, ""))

			// Und zum start gejumped
			jumpStart := NewEmptyCommand(OpJump)
			jumpStart.UnresolvedJump = true
			jumpStart.ResolvingHint = labelStart
			AppendCommand(jumpStart)

			// Und das ende-label ausgeben
			AppendCommand(NewInputOnlyCommand(Label, LiteralString, 0, labelEnd))

			return nil
		},
		"func": func(statement *Statement, context *Block) error {
			/* Function declaration: die Funktion wird hinzugefügt und für die spätere verarbeitung markiert */

			// Für den functionblock wird noch der parent richtig gesetzt
			statement.childBlocks[0].parent = context

			function := NewFunction(NoType, statement, len(statement.children)-1) /* len - 1 weil das erste kind der funktionsname ist */

			err := AddFunction(statement.children[0].command, function)
			if err != nil {
				return err
			}

			return nil
		},
		"return": func(statement *Statement, context *Block) error {
			/* gibt einen wert aus der function zurück */
			// wenn return in einer funktion steht, wird ein return statement ausgegeben
			// wenn return im programblock steht, wird ein exit statement ausgegeben
			var err error

			err = statement.children[0].Process(context)
			if err != nil {
				return err
			}

			returnType := statement.children[0].varType

			// Ein programblock ist ein block, dessen registerframe der oberste block ist (also keinen parent hat)
			isProgramBlock := context.RegisterFrame().parent == nil
			if isProgramBlock {
				if returnType != Integer {
					// Es dürfen nur integer werte als exit-code verwendet werden
					return errors.New("Cannot use value of type " + returnType.String() + " as exit code")
				} else {
					AppendCommand(NewInputOnlyCommand(OpExit, Register, statement.children[0].register, ""))
				}
			} else {
				// "normales" return

				// Prüfen bzw. festlegen des return-types
				if context.RegisterFrame().returnType == NoType {
					// Wenn noch kein typ festgelegt ist, wird dieser returntype einfach angenommen
					context.RegisterFrame().returnType = returnType
				} else {
					// ansonsten muss geprüft werden, ob der eigene returntype mit dem schon gesetzten übereinstimmt
					if returnType != context.RegisterFrame().returnType {
						return errors.New("Cannot infer single type as return type. Have " + returnType.String() + " and " + context.RegisterFrame().returnType.String())
					}
				}

				// Jetzt kann return ausgegeben werden
				AppendCommand(NewInputOnlyCommand(OpReturn, Register, statement.children[0].register, ""))
			}

			return nil
		},
		"proccall": func(statement *Statement, context *Block) error {
			/* procedure call wird wie funccall behandelt (ist halt ein register unnötig befüllt... egal) */

			return statement.children[0].Process(context)
		},
		"funccall": func(statement *Statement, context *Block) error {
			/* funktionsaufruf - sieht in etwa so aus:
			TRE: - 6.0 - funccall (3 children)
			TRE: - 6.0.0 - mini (0 children)
			TRE: - 6.0.1 - access (1 children)
			TRE: - 6.0.1.0 - var (1 children)
			TRE: - 6.0.1.0.0 - $a (0 children)
			TRE: - 6.0.2 - access (1 children)
			TRE: - 6.0.2.0 - var (1 children)
			TRE: - 6.0.2.0.0 - $b (0 children)

			1. kind = name
			2. - n. kind = argumente
			*/

			functionName := statement.children[0].command

			callArray := make([]int, 0)
			for i := 1; i < len(statement.children); i++ {
				err := statement.children[i].Process(context)
				if err != nil {
					return err
				}
				callArray = append(callArray, statement.children[i].register)
			}

			// Jetzt erst register reservieren, damit register auf jedem fall nach den parameter-registern kommt
			if statement.register < 0 {
				statement.register = context.RegisterFrame().NewRegister()
			}

			funccall := NewDoubleValueCommand(OpCall, Register, statement.register, RegisterArray, 0, "")
			funccall.Op1.ValueArray = callArray
			funccall.UnresolvedJump = true
			funccall.ResolvingHint = functionName
			AppendCommand(funccall)

			// Problem: was ist das für ein typ?
			statement.varType = Integer

			return nil
		},
	}
}

/* Zu diesem Zeitpunkt ist parent noch nicht gesetzt */
func (block *Block) Process(parent *Block) error {
	var err error

	if parent == nil {
		// ist ein program block (wie ein function block)
		var entryCommand *Command
		entryCommand = NewInputOnlyCommand(OpProgEntry, LiteralInt, 0 /* wird noch gesetzt */, "")
		AppendCommand(entryCommand)

		for i := 0; i < len(block.statements); i++ {
			err = block.statements[i].Process(block)
			if err != nil {
				return err
			}
		}

		exitCommand := NewInputOnlyCommand(OpProgExit, LiteralInt, block.usedRegisters, "")
		AppendCommand(exitCommand)

		entryCommand.Op1.ValueInt = block.usedRegisters

		// Nachdem der programblock ausgegeben wurde, werden die einzelnen functions ausgegeben
		for _, function := range Functions {
			err = function.Process(parent)
			if err != nil {
				return err
			}
		}
	}
	return nil
}

func (statement *Statement) Process(parent *Block) error {
	if statement.command == "" {
		return nil // Ist ein leeres statement (wird nicht ausgegeben)
	} else {
		procFunction := statementProc[statement.command]
		if procFunction != nil {
			return procFunction(statement, parent)
		} else {
			debugTree("No processor found for " + statement.command)
			return nil
		}
	}
}

func (function *Function) Process(parent *Block) error {
	var err error

	functionBlock := &function.FunctionDeclaration.childBlocks[0]
	functionStatement := function.FunctionDeclaration

	// Label ausgeben, damit die function dann gefunden werden kann
	label := functionStatement.children[0].command // = function name
	AppendCommand(NewInputOnlyCommand(Label, LiteralString, 0, label))

	// Entry-command ausgeben - die größe des entry commands wird noch auf die anzahl der block-register gesetzt
	var entryCommand *Command
	entryCommand = NewInputOnlyCommand(OpFuncEntry, LiteralInt, 0 /* wird noch gesetzt */, "")
	AppendCommand(entryCommand)

	// Für alle funktions-argumente im functionblock symbole reservieren
	for i := 1; i < len(functionStatement.children); i++ { /* fängt bei 1 an, weil das erste kind der name ist */
		functionBlock.addSymbol(functionStatement.children[i].children[0].command, Integer, functionBlock.NewRegister())
	}

	// Jetzt den Inhalt ausgeben
	for _, stmt := range functionBlock.statements {
		err = stmt.Process(functionBlock)
		if err != nil {
			return err
		}
	}

	// Exit command ausgeben
	exitCommand := NewInputOnlyCommand(OpFuncExit, LiteralInt, functionBlock.usedRegisters, "")
	AppendCommand(exitCommand)

	entryCommand.Op1.ValueInt = functionBlock.usedRegisters

	return nil
}
