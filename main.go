package main

import (
	"flag"
	"fmt"
	"os"
	"os/signal"
	"syscall"
)

var (
	MAKE_TRANS bool = false
	USE_DB     bool = false
)

func main() {
	var (
		jesse     = flag.Bool("j", false, "Only use if you are jesse")
		makeTrans = flag.Bool("t", false, "Turn this on if you want factom transactions made at bootup")
		db        = flag.Bool("db", false, "Turn this on if you want to use db cache")
		robin     = flag.Bool("r", false, "Only use if you are robin")
		steven    = flag.Bool("s", false, "Only use if you are steven")
	)

	flag.Parse()

	if *jesse {
		fmt.Println("Using Jesse's path")
		FILES_PATH = "/home/jesse/Electron/Hackathon/hackathon2"
	}

	if *robin {
		fmt.Println("Using Robin's path")
		FILES_PATH = "/Users/robin/dev/dbh/hackathon2"
	}

	if *steven {
		fmt.Println("Using Robin's path")
		FILES_PATH = "/Users/stevenmasley/go/src/github.com/dbHackathon2017/hackathon2"
	}

	if *db {
		USE_DB = true
	}

	MAKE_TRANS = *makeTrans

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		if MainCompany != nil {
			MainCompany.Save()
		}
		os.Exit(1)
	}()
	fmt.Printf("Path: %s\n", FILES_PATH)
	ServeFrontEnd(1337)
}
