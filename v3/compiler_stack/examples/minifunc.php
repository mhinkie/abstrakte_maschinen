<?php
$a = 3;
$b = 4;

/* comment */

noparam();

function mini($x, $y) {
	return $x + $y;
}

function noparam() {
	echo "noparam";
}

echo mini($a, $b);

noparam();

?>
