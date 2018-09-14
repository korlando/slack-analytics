package main

import (
  "flag"
  "fmt"
  "log"
  "strconv"
  "strings"

  sa "github.com/korlando/slackanalytics"
)

const (
  dataPath = "../data"
)

var (
  commonWords = []string{"i", "you", "he", "she", "it", "we", "they", "me", "him", "her", "us", "them", "what", "who", "whom", "this", "that", "these", "those", "the", "to", "a", "is", "of", "and", "in", "on", "for", "not", "like", "have", "my"}
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
  ws := sa.GetWordStats(messages)
  wordCounts := sa.GetSortedWords(ws)
  fmt.Println("Total words: " + strconv.Itoa(ws.TotalWords))
  fmt.Println("Avg word length: " + strconv.FormatFloat(ws.AvgLength, 'f', 3, 64))
  i := 0
  j := 0
  for j < 20 {
    wc := wordCounts[i]
    i += 1
    w := wc.Word
    isCommon := false
    for _, c := range commonWords {
      if c == strings.ToLower(w) {
        isCommon = true
        continue
      }
    }
    if isCommon {
      continue
    }
    fmt.Println(w + " " + strconv.Itoa(wc.Count))
    j += 1
  }
}
