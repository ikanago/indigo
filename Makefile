indigo: *.go output
	go build -o output/indigo *.go

indigo2: indigo
	./output/indigo > output/indigo2.s
	clang -o output/indigo2 output/indigo2.s

.PHONY: output
output:
	mkdir -p output

test: indigo test/**/*.go test/**/expected.txt
	./run_tests.sh

unittest: indigo *_test.go
	go test -v

testall: test unittest

.PHONY: check
check:
	go fmt *.go
	go vet .
