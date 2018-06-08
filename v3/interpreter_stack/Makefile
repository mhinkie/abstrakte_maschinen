all: interpreter.c
		gcc -o mphp interpreter.c 

test: all
	./mphp examples/prime.php.bc

clean:
		rm mphp
