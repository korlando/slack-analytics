package main

import (
  "fmt"
  "log"
  "sort"
  "strconv"
  "strings"

  sa "github.com/korlando/slackanalytics"
)

const (
  path = "../data"
)

type WordMap struct {
  TotalWords uint
  Words      map[string]uint
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

func main() {
  messages, err := sa.ReadAllMessages(path)
  if err != nil {
    log.Fatal(err)
  }
  wm := WordMap{0, make(map[string]uint)}
  for _, m := range messages {
    if m.Text == "" {
      continue
    }
    words := strings.Fields(m.Text)
    for _, w := range words {
      if w == "" {
        continue
      }
      wm.TotalWords += 1
      count, ok := wm.Words[w]
      if !ok {
        wm.Words[w] = 1
      } else {
        wm.Words[w] = count + 1
      }
    }
  }
  wordCounts := make([]WordCount, len(wm.Words))
  i := 0
  for word, count := range wm.Words {
    wordCounts[i] = WordCount{word, count}
    i += 1
  }
  sort.Sort(sortByCount(wordCounts))
  fmt.Println("Total words: " + strconv.Itoa(int(wm.TotalWords)))
  for i := 0; i < 20; i++ {
    wc := wordCounts[i]
    fmt.Println(wc.Word + " " + strconv.Itoa(int(wc.Count)))
  }
}
