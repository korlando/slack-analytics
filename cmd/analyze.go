package main

import (
  "flag"
  "fmt"
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
  messages, err := sa.ReadAllMessages(opt.path)
  if err != nil {
    log.Fatal(err)
  }
  sa.GetAndPrintStats(messages)
  users, err := sa.GetUsers(opt.path)
  if err != nil {
    log.Fatal(err)
  }
  // get word stats for individual users
  for _, u := range users {
    userMessages := sa.FilterMessagesByUser(messages, u.Id)
    if len(userMessages) == 0 {
      continue
    }
    name := u.Profile.RealName
    if u.Profile.DisplayName != "" {
      name = u.Profile.DisplayName
    }
    fmt.Println(name + "\n")
    sa.GetAndPrintStats(userMessages)
    fmt.Println()
  }
}
