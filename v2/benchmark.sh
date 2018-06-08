#!/bin/bash
FILE=$1

START=`date +%s%N`
./mphp_stack.sh $1
END=`date +%s%N`

DUR=$(($END-$START))
DURS=`bc <<< "scale=6; ${DUR}/1000000000"`

echo ""
echo ""
echo "MPHP_STACK: Start ${START}; End ${END}; Duration ${DUR} ns (${DURS} s)"

START=`date +%s%N`
./mphp_reg.sh $1
END=`date +%s%N`

DUR=$(($END-$START))
DURS=`bc <<< "scale=6; ${DUR}/1000000000"`

echo ""
echo ""
echo "MPHP_REG: Start ${START}; End ${END}; Duration ${DUR} ns (${DURS} s)"


START=`date +%s%N`
php $1
END=`date +%s%N`

DUR=$(($END-$START))
DURS=`bc <<< "scale=6; ${DUR}/1000000000"`

echo ""
echo ""
echo "PHP: Start ${START}; End ${END}; Duration ${DUR} ns (${DURS} s)"
