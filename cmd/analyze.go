package main

import (
	"flag"
	"log"

	sa "github.com/korlando/slackanalytics"
)

const (
	dataPath = "../data"
)

type Options struct {
	path    string
	msgFile bool
}

func parseFlags() (opt Options) {
	var p, path string
	// m is true if the path points to a JSON file containing
	// messages alone, vs. a folder with a Slack dump
	var m bool
	pDefault := dataPath
	pDesc := "Path to the data folder."
	flag.StringVar(&p, "p", pDefault, pDesc)
	flag.StringVar(&path, "path", pDefault, pDesc)
	flag.BoolVar(&m, "m", false, "Path points to a JSON file containing a messages array.")
	flag.Parse()
	opt = Options{
		path:    p,
		msgFile: m,
	}
	if path != pDefault {
		opt.path = path
	}
	return
}

func main() {
	opt := parseFlags()
	if opt.msgFile {
		messages, err := sa.ReadMessagesFromFile(opt.path)
		if err != nil {
			log.Fatal(err)
		}
		sa.ExportMessageAnalysis(messages)
		return
	}
	users, err := sa.GetUsers(opt.path)
	if err != nil {
		log.Fatal(err)
	}
	channels, err := sa.GetChannels(opt.path)
	if err != nil {
		log.Fatal(err)
	}
	sa.GetAndPrintStats(users, channels)
}
