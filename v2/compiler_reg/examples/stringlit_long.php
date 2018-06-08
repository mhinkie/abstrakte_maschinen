<?php
$s1 = "hohoho";
$s2 = "asdfasdf";
$s3 = "asdfasdfasdfasdf";
$s4 = "s4";

echo "asdf";
echo "hoho";
echo $s1 . " anhang";
echo $s3 . $s2;
echo $s1 . $s1 . " aosdoa";
echo "hoahoasdfhasodf";

$s7 = $s3 . $s2;
echo $s7;
echo $s7 . $s3;

/* erwartet:
alex@hihilenovo:~/git/am/interpreter_stack (master)$ php ../compiler_stack/examples/stringlit_long.php
asdfhohohohoho anhangasdfasdfasdfasdfasdfasdfhohohohohoho aosdoahoahoasdfhasodfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdf
*/
/* habe:
alex@hihilenovo:~/git/am/interpreter_stack (master)$ ./mphp examples/stringlit_long.php.bc 
asdfhohohohoho anhangasdfasdfasdfasdfasdfasdfhohohohohoho aosdoahoahoasdfhasodfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdf

und 

alex@hihilenovo:~/git/am/interpreter_reg (master)$ ./mphp examples/stringlit_long.php.bc 
asdfhohohohoho anhangasdfasdfasdfasdfasdfasdfhohohohohoho aosdoahoahoasdfhasodfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdfasdf
*/

?>