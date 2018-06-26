<?php

function fib($a) {
	if ($a < 2) {
		return 1;
	} else {
		return fib($a - 1) + fib($a - 2);
	}
}

for($i=0;$i<40;$i++) {
	echo fib($i);
	echo " ";
}
?>
