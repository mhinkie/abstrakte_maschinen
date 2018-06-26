all: interpreter.c
		gcc -ggdb -o mphp interpreter.c 

test: all
	./mphp examples/prime.php.bc

clean:
		rm mphp
