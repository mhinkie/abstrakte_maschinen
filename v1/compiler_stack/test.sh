#!/bin/bash

SUMMARY=""

FILES=("examples/stringlit.php" "examples/intarith.php" "examples/stringlit_long.php" "examples/boolean.php" "examples/conditional.php" "examples/conditional_bvars.php" "examples/unaryops.php" "examples/whileloop.php" "examples/forloop.php" "examples/array.php" "examples/foreach.php" "examples/minifunc.php" "examples/grossefunc.php" "examples/prime.php")

for file in ${FILES[*]}
do
	echo "-----------------TEST: $file ----------------"
	cat $file | ./mphpc > $file.bc
	if [ $? -eq 0 ]
	then
		SUMMARY="${SUMMARY}${file}:\033[0;32mSUCCESS\033[0m\n"
		echo -e "---------------------\033[0;32mSUCCESS\033[0m----------------------"
	else 
		SUMMARY="${SUMMARY}${file}:\033[0;31mFAILURE\033[0m\n"
		echo -e "---------------------\033[0;31mFAILURE\033[0m----------------------"
	fi
done

echo "SUMMARY:"
echo -e $SUMMARY
