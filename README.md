# PHP Interpreters written in GO (and C)
Comparison of stack and register implementations of interpreters for a subset of PHP.
The interpreters were tested using predefined php-files (fib.php, prime.php...).
There are 4 versions: 1 register-interpreter written in GO, 2 stack-interpreters written in GO, and
1 stack interpreter written in C. An overview of the design choices (and errors) can be found in
my [presentation](presentation.pdf).

## prerequesites
- go compiler
- optional (for v3): gcc
- goyacc (https://godoc.org/golang.org/x/tools/cmd/goyacc)
- nex (http://www-cs-students.stanford.edu/∼blynn/nex/)
- stringer (https://godoc.org/golang.org/x/tools/cmd/stringer)

## Directories
The versions directories (v1, v2, v3) each contain scripts to run the interpreter (including the compiler) - e.g.:
```
cd v1
./mphp_stack.sh ../benchmarks/fib.php
```
Note: The programs might have to be recompiled (depending on system). Run `make` in the respective
`compiler_stack`/`compiler_reg` and `interpreter_stack`/`interpreter_reg` directories.

### v1
Version 1 of stack and register interpreters (and compilers) written in GO.

### v1_analyse
Version 1 interpreters with ability to record command runtimes.

### v2
Version 2: New stack interpreter (reduced command set - no block entries or exists) written in GO.

### v3
Version 3: New stack interpreter - this time written in C using direct threaded code. This interpreter
uses the same compiler and command-set as v2.

### benchmarks
The results of the different benchmarks. A detailed (but not very well annotated) overview can be found in [bm.ods](bm.ods).
Also contains the source files for the 3 benchmarks.
