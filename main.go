package main

import (
	"flag"
	"fmt"
	crunner "github.com/rabugopsl/crunner/console_runner"
	"log"
)

type Flags struct {
	Command     *string
	CommandArgs []string
	Debug       *bool
	Input       *string
}

func main() {

	flags := processFlags()
	params, _ := crunner.LoadInputParams(*flags.Input)
	if err := crunner.RunCommand(*flags.Debug, params, *flags.Command, flags.CommandArgs...); err != nil {
		fmt.Printf("%s\n", err)
		log.Fatal(err)
	}
}

func processFlags() *Flags {
	var input = flag.String("input", "data.json", "Specifies the user input to pass down to the underlying executable.\n")

	var debug = flag.Bool("debug", false, "Specifies whether or not the program should deliver debug data")

	flag.Parse()

	flagstruct := new(Flags)
	flagstruct.Input = input
	flagstruct.Debug = debug

	if len(flag.Args()) < 1 {
		log.Fatal("The application requires at least the executable command to run")
	}
	var command string = flag.Args()[0]
	flagstruct.Command = &command

	if len(flag.Args()) > 1 {
		flagstruct.CommandArgs = flag.Args()[1:]
	}

	return flagstruct
}
