package main

import (
	"codeberg.org/HackTheOxidation/twig/dispatch"
	"codeberg.org/HackTheOxidation/twig/options"
)

func main() {
	opts := options.GetOptions()
	dispatch.DispatchAndExecute(&opts)
}
