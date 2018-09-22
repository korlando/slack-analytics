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
  path string
}

func parseFlags() (opt Options) {
  var p, path string
  pDefault := dataPath
  pDesc := "Path to the data folder."
  flag.StringVar(&p, "p", pDefault, pDesc)
  flag.StringVar(&path, "path", pDefault, pDesc)
  flag.Parse()
  opt = Options{
    path: p,
  }
  if path != pDefault {
    opt.path = path
  }
  return
}

func main() {
  opt := parseFlags()
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
