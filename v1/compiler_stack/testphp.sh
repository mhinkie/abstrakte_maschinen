#!/bin/bash

SUMMARY=""

for file in examples/*.php
do
	echo "-----------------TEST: $file ----------------"
	cat $file | php
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
