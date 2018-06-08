UTIL := analyze.go

all: interpreter.go $(UTIL)
	go build -o mphp interpreter.go $(UTIL)

clean:
	rm mphp
