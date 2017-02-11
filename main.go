package main

import (
	"flag"
	"fmt"
)

var (
	MAKE_TRANS bool = false
)

func main() {
	var (
		jesse     = flag.Bool("j", false, "Only use if you are jesse")
		makeTrans = flag.Bool("t", false, "Turn this off if you don't want factom transactions made at bootup")
	)

	flag.Parse()

	if *jesse {
		fmt.Println("Using Jesse's path")
		FILES_PATH = "/home/jesse/Electron/Hackathon/hackathon2"
	}

	MAKE_TRANS = *makeTrans

	ServeFrontEnd(1337)
}
