<?php
/* expected stack vm:
push <adr von asdfasdf> - pushed die adresse auf den stack
store aStringVar - speichert diese adresse als aStringVar
load aStringVar - holt sich die adresse wieder
echo - output als string
*/	

$aStringVar = "asdfasdf";

echo $aStringVar;
?>
