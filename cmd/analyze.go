package main

import (
  "flag"
  "fmt"
  "log"
  "sort"
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

type WordStat struct {
  TotalWords   uint
  WordCountMap map[string]uint
}

type WordCount struct {
  Word  string
  Count uint
}

type sortByCount []WordCount

func (s sortByCount) Len() int {
  return len(s)
}

func (s sortByCount) Less(i, j int) bool {
  return s[i].Count > s[j].Count
}

func (s sortByCount) Swap(i, j int) {
  s[i], s[j] = s[j], s[i]
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
  ws := WordStat{0, make(map[string]uint)}
  for _, m := range messages {
    if m.Text == "" {
      continue
    }
    words := strings.Fields(m.Text)
    for _, w := range words {
      if w == "" {
        continue
      }
      ws.TotalWords += 1
      count, ok := ws.WordCountMap[w]
      if !ok {
        ws.WordCountMap[w] = 1
      } else {
        ws.WordCountMap[w] = count + 1
      }
    }
  }
  wordCounts := make([]WordCount, len(ws.WordCountMap))
  i := 0
  for word, count := range ws.WordCountMap {
    wordCounts[i] = WordCount{word, count}
    i += 1
  }
  sort.Sort(sortByCount(wordCounts))
  fmt.Println("Total words: " + strconv.Itoa(int(ws.TotalWords)))
  i = 0
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
    fmt.Println(w + " " + strconv.Itoa(int(wc.Count)))
    j += 1
  }
}
