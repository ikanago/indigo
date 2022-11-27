indigo: *.go output
	go build -o output/indigo *.go

indigo2: indigo
	./output/indigo > output/indigo2.s
	clang -o output/indigo2 output/indigo2.s

output:
	mkdir -p output
