<?php

function fib($a) {
	if ($a < 2) {
		return 1;
	} else {
		return fib($a - 1) + fib($a - 2);
	}
}

for($i=0;$i<10;$i++) {
	echo fib($i);
	echo " ";
}

$x = 3;
$y = 4;
if($x == 3) {
	$y = 5;
	$z = 2;
} else {
	$z = 3;
}
?>
