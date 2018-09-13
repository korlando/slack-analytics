package main

import (
  "encoding/json"
  "fmt"
  "io/ioutil"
  "log"
  "os"
  "sort"
  "strconv"
  "strings"
)

const (
  path = "../data"
)

type Message struct {
  User      string `json:"user"`
  Type      string `json:"type"`
  SubType   string `json:"subtype"`
  Text      string `json:"text"`
  TimeStamp string `json:"ts"`
}

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
  fileInfos, err := ioutil.ReadDir(path)
  if err != nil {
    log.Fatal(err)
  }
  wm := WordMap{0, make(map[string]uint)}
  for _, f := range fileInfos {
    // only channels are dirs
    if !f.IsDir() {
      continue
    }
    // look at each json file in channel (1 per day)
    channelPath := path + "/" + f.Name()
    jsonFiles, err := ioutil.ReadDir(channelPath)
    if err != nil {
      log.Fatal(err)
    }
    for _, j := range jsonFiles {
      file, err := os.Open(channelPath + "/" + j.Name())
      if err != nil {
        log.Fatal(err)
      }
      defer file.Close()
      jsonBytes, err := ioutil.ReadAll(file)
      if err != nil {
        log.Fatal(err)
      }
      var messages []Message
      json.Unmarshal(jsonBytes, &messages)
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
