CLANG = clang
FLAGS = -O3 -Wall -Wextra -fPIE
AFLAGS = $(FLAGS) -c
FLTO = -Wall # TODO: ensure we can activate -flto again

.PHONY: clean compile transpile chec_comp build all test_build wipe

# by default also clean all intermediate files except results
all: build clean

build: check_comp copy_runtime transpile compile test_build

check_comp:
	glslangValidator data/*.comp
	mkdir -p data/build
	mkdir -p data/output

copy_runtime:
	cp -rf src/runtime data/build/runtime 

transpile:
	cd data && go run /src/gl2cl/*.go -buildpath build -outpath output -headerpath build/runtime *.comp
	cd data && cp build/runtime/shared.h output/shared.h

compile:
	$(CLANG) data/build/runtime/runtime.c -o data/build/runtime.o $(AFLAGS) $(FLTO)
	$(CLANG) data/build/kernel.cl -o data/build/kernel.o -Xclang -finclude-default-header $(AFLAGS) $(FLTO)
	ar rc data/output/kernel.a data/build/kernel.o data/build/runtime.o

test_build:
	cd data/output && go build .
	cd data/output && go test .

clean:
	rm -rf data/build/**

wipe:
	rm -rf data/build/** data/output/** data/*.comp