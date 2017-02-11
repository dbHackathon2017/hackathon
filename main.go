package main

import (
	"flag"
)

func main() {
	var (
		jesse = flag.Bool("j", true, "Only use if you are jesse")
	)

	flag.Parse()

	if *jesse {
		FILES_PATH = "/home/jesse/Electron/Hackathon/hackathon2"
	}

	ServeFrontEnd(1337)
}
