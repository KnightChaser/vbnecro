package main

import (
	"log"
	"os"
)

func init() {
	// Configure log output and flags.
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Ldate | log.Ltime | log.Lshortfile)
}
