package main

import (
	"flag"
)

func main() {
	var (
		jesse = flag.Bool("compiled", true, "Decides wheter to use the compiled statics or not. Useful for modifying")
	)

	flag.Parse()

	if *jesse {
		FILES_PATH = "/home/jesse/Electron/Hackathon/hackathon2"
	}

	ServeFrontEnd(1337)
}
