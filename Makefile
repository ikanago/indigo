indigo: *.go
	go build -o indigo *.go

indigo2: indigo
	./indigo > indigo2.s
	clang -o indigo2 indigo2.s
