#!/bin/bash

MPHPC="./compiler_stack/mphpc"
MPHP="./interpreter_stack/mphp"

cat $1 | $MPHPC  > $1.bc
$MPHP $1.bc

