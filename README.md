# PHP Interpreters written in GO (and C)
Comparison of stack and register implementations of interpreters for a subset of PHP.
The interpreter were tested using predefined php-files (fib.php, prime.php...)

## prerequesites
- goyacc (https://godoc.org/golang.org/x/tools/cmd/goyacc)
- nex (http://www-cs-students.stanford.edu/∼blynn/nex/)
- stringer (https://godoc.org/golang.org/x/tools/cmd/stringer)

## Directories

### v1
Version 1 of stack and register interpreters (and compilers) written in GO.

### v1_analyse
Version 1 interpreters with ability to record command runtimes.

### v2
Version 2: New stack interpreter (reduced command set - no block entries or exists) written in GO.

### v3
Version 3: New stack interpreter - this time written in C using direct threaded code.

### benchmarks
The results of the different benchmarks. A detailed (but not very well annotated) overview can be found in bm.ods.
