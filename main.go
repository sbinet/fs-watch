package main

import (
	"flag"
	"log"
	"os"
	"os/exec"
	"path/filepath"

	"github.com/fsnotify/fsnotify"
)

func main() {
	log.SetPrefix("fs-watch: ")
	log.SetFlags(0)

	cmdFlag := flag.String("cmd", "", "command to run when fs modifications happen")

	flag.Parse()

	fs, err := fsnotify.NewWatcher()
	if err != nil {
		log.Fatalf("could not create fs-watcher: %v", err)
	}
	defer fs.Close()

	switch {
	case flag.NArg() > 0:
		for _, dir := range flag.Args() {
			abs, err := filepath.Abs(dir)
			if err != nil {
				log.Fatal(err)
			}
			err = fs.Add(abs)
			if err != nil {
				log.Fatalf("could not add directory [%s] to the watch list: %v", dir, err)
			}
		}
	default:
		dir, err := os.Getwd()
		if err != nil {
			log.Fatalf("could not get working directory: %v", err)
		}
		err = fs.Add(dir)
		if err != nil {
			log.Fatalf("could not add directory [%s] to the watch list: %v", dir, err)
		}
	}

	for {
		select {
		case evt := <-fs.Events:
			log.Printf("event: %v", evt)
			if evt.Op&fsnotify.Write != fsnotify.Write {
				continue
			}
			cmd := exec.Command(*cmdFlag)
			cmd.Stdout = os.Stdout
			cmd.Stderr = os.Stderr
			err = cmd.Run()
			if err != nil {
				log.Printf("error running watch command: %v", err)
			}
		case err := <-fs.Errors:
			log.Printf("error: %v", err)
		}
	}

}
