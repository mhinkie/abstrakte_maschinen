<?php

function fib($a) {
	if ($a < 2) {
		return 1;
	} else {
		return fib($a - 1) + fib($a - 2);
	}
}

echo fib(2);
?>
