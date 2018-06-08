/* Utility funktionen zum arbeiten mit goyacc und nex */

package main

/* Stellt die Text-Funktion vom darunterliegenden lexer zur verfügung */
func yyText(yylex yyLexer) string {
	l := (*Lexer)(yylex.(*Lexer))
	return l.Text()
}