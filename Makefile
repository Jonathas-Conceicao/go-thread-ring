GC=go
TRAGET=threadring

build: $(TRAGET).go
	$(GC) build -o $(TRAGET).out

run: build
	./$(TRAGET).out 1000

tests: build
	./runTests.sh

cleanAll: cleanTests cleanExecutable

cleanTests:
	rm -f ./tests/*

cleanExecutable:
	rm -f  $(TRAGET).out
