<?php
function one($x) {
	return $x > 10;
}

function two($x) {
	$y = $x;
	$z = $y + 2;
	if(one($z)) {
		echo "gt";
	} else {
		echo "lt";
	}
}

two(11);
?>
