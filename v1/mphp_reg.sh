#!/bin/bash

MPHPC="./compiler_reg/mphpc"
MPHP="./interpreter_reg/mphp"

cat $1 | $MPHPC  > $1.bc
$MPHP $1.bc

