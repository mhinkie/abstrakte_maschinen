all: interpreter.c
		gcc -O3 -o mphp interpreter.c

test: all
	./mphp examples/prime.php.bc

clean:
		rm mphp
