package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"os/signal"
	"syscall"
)

var (
	MAKE_TRANS bool = false
	USE_DB     bool = false
	FULL_CACHE bool = false
)

func main() {
	var (
		jesse     = flag.Bool("j", false, "Only use if you are jesse")
		makeTrans = flag.Bool("t", false, "Turn this on if you want factom transactions made at bootup")
		db        = flag.Bool("db", false, "Turn this on if you want to use db cache")
		full      = flag.Bool("f", false, "Turn this on if you want to cache factom entries into db")
		robin     = flag.Bool("r", false, "Only use if you are robin")
		debug     = flag.Bool("d", false, "Turn on debugging")
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

	if *db {
		USE_DB = true
		if *full {
			FULL_CACHE = true
		}
	}

	if !(*debug) {
		log.SetOutput(ioutil.Discard)
	}

	MAKE_TRANS = *makeTrans

	c := make(chan os.Signal, 2)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	go func() {
		<-c
		if MainCompany != nil {
			MainCompany.Save(GetCacheList(), FULL_CACHE)
		}
		os.Exit(1)
	}()
	fmt.Printf("Path: %s\n", FILES_PATH)
	ServeFrontEnd(1337)
}
