<?php
/* gibt $a % $b */
function mod($a, $b) {
	while($a > $b) {
		$a = $a - $b;
	}
	if($a == $b) {
		return 0;
	}
	return $a;
}

/* liefert immer nächstes int über sqrt
also $x mit 0 < ($x - sqrt($n)) < 1*/
function sqrtA($n) {
	$sqrt = 0;
	for(; ($sqrt * $sqrt) < $n; $sqrt++) {
	}
	return $sqrt;
}

/* liefert 1 wenn prime */
function isPrime($p) {
	$checkUntil = sqrtA($p);
	for($i = 2; $i < $checkUntil; $i++) {
		if(mod($p, $i) == 0) {
			return 0;
		}
	}
	return 1;
}


for($i=2;$i<100000;$i++) {
	if(isPrime($i)) {
		echo $i;
		echo " ";
	}
}

/*echo "mod: ";
echo mod(5,4);
echo " ";
echo mod(10,2);
echo " ";
echo mod(7,3);
echo " ";

echo "sqrtA: ";
echo sqrtA(10mal10);
echo " ";
echo sqrtA(1234);
echo " ";
echo sqrtA(22);
echo " ";
echo sqrtA(4);
echo " ";

echo "isPrime: ";
echo isPrime(10);
echo " ";
echo isPrime(20);
echo " ";
echo isPrime(11);
echo " ";
echo isPrime(12341214mal2);
echo " ";
echo isPrime(33);*/


?>
