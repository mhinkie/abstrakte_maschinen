UTIL := parser.go commands.go

all: interpreter.go $(UTIL) *_string.go
	go build -o mphp interpreter.go $(UTIL) *_string.go

*_string.go: interpreter.go $(UTIL)
	go generate

clean:
	rm mphp
