UTIL := debug.go symbols.go util.go output.go

all: mphp.nn.go y.go $(UTIL)
	go build -o mphpc mphp.nn.go y.go mphp_stack.go $(UTIL)

mphp.nn.go: mphp.nex
	nex mphp.nex

y.go: mphp.y
	goyacc mphp.y

clean:
	rm y.go y.output mphp.nn.go
	rm mphpc

test: all
	./test.sh
