GC=go

main: threadring.go
	$(GC) build -o threadring.out
