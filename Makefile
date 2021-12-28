all: build

build:
	gccgo -c ./options/options.go -o options.o
	gccgo -c ./dispatch/dispatch.go -I. -o dispatch.o
	gccgo -c twig.go -I. -o twig.o
	gccgo -o twig options.o dispatch.o twig.o

clean:
	@rm twig twig.o dispatch.o options.o
