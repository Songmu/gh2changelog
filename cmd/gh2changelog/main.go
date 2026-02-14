package main

import (
	"context"
	"flag"
	"log"
	"os"
)

func main() {
	log.SetFlags(0)
	
	arguments := os.Args[1:]
	err := run(context.Background(), arguments, os.Stdout, os.Stderr)
	
	if err != nil && err != flag.ErrHelp {
		log.Println(err)
		code := 1
		if exitCoder, ok := err.(interface{ ExitCode() int }); ok {
			code = exitCoder.ExitCode()
		}
		os.Exit(code)
	}
}
